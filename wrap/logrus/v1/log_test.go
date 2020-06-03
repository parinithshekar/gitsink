// Copyright 2019 Cisco Systems, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package v1_test

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"testing"

	pkg "github.com/parinithshekar/gitsink/pkg/v1"
	logger "github.com/parinithshekar/gitsink/wrap/logrus/v1"
	require "github.com/stretchr/testify/require"
)

// Enforce interface implementation.
func TestInterface(t *testing.T) {
	var _ pkg.Logger = &logger.Logger{}
}

// Test Error without any fields
func TestLogError(t *testing.T) {
	buf := &bytes.Buffer{}
	log := customLogger(buf)
	err := errors.New("This is an error message.")
	log.Errorf("%s", err.Error())
	// "time="2019-09-10T10:51:43+05:30" level=error msg="This is an error message." file="log_test.go:43"
	require.Contains(t, buf.String(), fmt.Sprintf("level=%s msg=\"%s\" file=\"%s:%d\"", "error", "This is an error message.", "log_test.go", getLineNumber()-2))
}

// Test Waringf
func TestLogWarningf(t *testing.T) {
	buf := &bytes.Buffer{}
	log := customLogger(buf)
	err := errors.New("This is an warning.")
	log.Warningf("%s", err.Error())
	// "time="2019-09-10T10:51:43+05:30" level=error msg="This is an error message." file="log_test.go:43"
	require.Contains(t, buf.String(), fmt.Sprintf("level=%s msg=\"%s\" file=\"%s:%d\"", "warning", "This is an warning.", "log_test.go", getLineNumber()-2))
}

// Test logger with WithError method.
func TestWithError(t *testing.T) {
	buf := &bytes.Buffer{}
	log := customLogger(buf)

	// Check Errorf().
	checkErrorf(t, log, buf)

	// This shouldn't have any effect on WithError.
	log.AutoClearFields(false)

	checkErrorf(t, log, buf)
}

func checkErrorf(t *testing.T, log pkg.Logger, buf *bytes.Buffer) {
	err := errors.New("This is a custom error.")
	log.WithError(err).Errorf("Encountered an error.")
	// time="2019-09-10T10:56:17+05:30" level=error msg="Encountered an error." file="log_test.go:53" error="This is a custom error."
	require.Contains(t, buf.String(), fmt.Sprintf("level=%s msg=\"%s\" file=\"%s:%d\" error=\"%s\"", "error", "Encountered an error.", "log_test.go", getLineNumber()-2, "This is a custom error."))

	// Check if error field is cleared.
	log.Errorf("Encountered an error.")
	// time="2019-09-10T10:56:17+05:30" level=error msg="Encountered an error." file="log_test.go:53" error="This is a custom error."
	require.Contains(t, buf.String(), fmt.Sprintf("level=%s msg=\"%s\" file=\"%s:%d\"", "error", "Encountered an error.", "log_test.go", getLineNumber()-2))
}

// Test logger with WithField method.
func TestWithField(t *testing.T) {
	buf := &bytes.Buffer{}
	log := customLogger(buf)
	log.WithField("key", "value").WithField("key2", "value2").Errorf("Errors with custom field.")
	// time="2019-09-10T11:17:19+05:30" level=error msg="Errors with custom field." file="log_test.go:63" key=value key2=value2
	require.Contains(t, buf.String(), fmt.Sprintf("level=%s msg=\"%s\" file=\"%s:%d\" key=value key2=value2", "error", "Errors with custom field.", "log_test.go", getLineNumber()-2))
	log.WithField("key3", "value3").Errorf("Errors with custom field again.")
	// time="2019-09-10T11:17:19+05:30" level=error msg="Errors with custom field." file="log_test.go:63" key=value key2=value2
	require.Contains(t, buf.String(), fmt.Sprintf("level=%s msg=\"%s\" file=\"%s:%d\" key3=value3", "error", "Errors with custom field again.", "log_test.go", getLineNumber()-2))
}

// Test concurrent modifications to fields.
func TestConcurrentMods(t *testing.T) {
	buf := &bytes.Buffer{}
	log := customLogger(buf)
	log.SetLevel("debug")
	for i := 0; i < 15000; i++ {
		go func() {
			log.WithError(pkg.ErrNoMatch).WithField("key", "value").Debugf("")
			require.Contains(t, buf.String(), fmt.Sprintf("level=%s file=\"%s:%d\" error=\"No alert matched in alert config.\" key=value", "debug", "log_test.go", getLineNumber()-1))
		}()
	}
}

// Test logger with AutoClear disabled.
func TestAutoClearFieldsDisabled(t *testing.T) {
	buf := &bytes.Buffer{}
	log := customLogger(buf)
	log.AutoClearFields(false)
	log.WithField("key", "value").WithField("key2", "value2").Errorf("Errors with custom field.")
	// time="2019-09-10T11:17:19+05:30" level=error msg="Errors with custom field." file="log_test.go:63" key=value key2=value2
	require.Contains(t, buf.String(), fmt.Sprintf("level=%s msg=\"%s\" file=\"%s:%d\" key=value key2=value2", "error", "Errors with custom field.", "log_test.go", getLineNumber()-2))

	log.WithField("key3", "value3").Errorf("Errors with custom field again.")
	require.Contains(t, buf.String(), fmt.Sprintf("level=%s msg=\"%s\" file=\"%s:%d\" key=value key2=value2 key3=value3", "error", "Errors with custom field again.", "log_test.go", getLineNumber()-1))
}

func TestWithFields(t *testing.T) {
	buf := &bytes.Buffer{}
	log := customLogger(buf)
	log.SetLevel("debug")
	log.WithFields(map[string]interface{}{
		"key1": "val1",
		"key2": "val2",
	}).
		WithField("key3", "val3").
		Debugf("Errors with custom fields.")
	// time="2019-09-10T11:21:40+05:30" level=debug msg="Errors with custom fields." file="log_test.go:77" key1=val1 key2=val2 key3=val3
	require.Contains(t, buf.String(), fmt.Sprintf("level=%s msg=\"%s\" file=\"%s:%d\" key1=val1 key2=val2 key3=val3", "debug", "Errors with custom fields.", "log_test.go", getLineNumber()-2))
}

// Test debug log level enabled
func TestWithDebugEnabled(t *testing.T) {
	buf := &bytes.Buffer{}
	log := customLogger(buf)
	log.SetLevel("deBug")
	log.Debugf("Debug log enabled.")
	// time="2019-09-10T11:25:46+05:30" level=debug msg="Debug log enabled." file="log_test.go:88
	require.Contains(t, buf.String(), fmt.Sprintf("level=%s msg=\"%s\" file=\"%s:%d\"", "debug", "Debug log enabled.", "log_test.go", getLineNumber()-2))
}

// Test debug log level disabled
func TestWithDebugDisabled(t *testing.T) {
	buf := &bytes.Buffer{}
	log := customLogger(buf)
	log.SetLevel("deBug")
	log.Debugf("Debug log enabled.")

	// "time="2019-09-09 19:12:11" level="DEBUG" tag="test.with.fields" location="log_test.go:90" msg="Debug log enabled."
	require.Contains(t, buf.String(), "")
}

// Custom logger that writes to a buffer for testing, instead of os.Stderr
func customLogger(output io.Writer) pkg.Logger {
	log := logger.New()
	log.SetOutput(output)
	return log
}

// Test log.Infof without any fields
func TestLogInfo(t *testing.T) {
	buf := &bytes.Buffer{}
	log := customLogger(buf)
	msg := "This is an info message."
	log.Infof(msg)
	// "time="2019-09-10T10:51:43+05:30" level=info msg="This is an info message." file="log_test.go:118"
	require.Contains(t, buf.String(), fmt.Sprintf("level=%s msg=\"%s\" file=\"%s:%d\"", "info", msg, "log_test.go", getLineNumber()-2))
}

// Test log.Tracef without any fields
func TestLogTrace(t *testing.T) {
	buf := &bytes.Buffer{}
	log := customLogger(buf)
	log.SetLevel("trace")
	msg := "This is an trace message."
	log.Tracef(msg)
	// "time="2019-09-10T10:51:43+05:30" level=info msg="This is an info message." file="log_test.go:118"
	require.Contains(t, buf.String(), fmt.Sprintf("level=%s msg=\"%s\" file=\"%s:%d\"", "trace", msg, "log_test.go", getLineNumber()-2))
}

// Test log.Debugf without any fields
func TestLogDebug(t *testing.T) {
	buf := &bytes.Buffer{}
	log := customLogger(buf)
	log.SetLevel("debug")
	msg := "This is an debug message."
	log.Debugf(msg)
	// "time="2019-09-10T10:51:43+05:30" level=info msg="This is an info message." file="log_test.go:118"
	require.Contains(t, buf.String(), fmt.Sprintf("level=%s msg=\"%s\" file=\"%s:%d\"", "debug", msg, "log_test.go", getLineNumber()-2))
}

// Test log.Panicf without any fields
func TestLogPanic(t *testing.T) {
	buf := &bytes.Buffer{}
	log := customLogger(buf)
	msg := "This is an panic message."
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()

	log.Panicf(msg)
}

// Test log.Fatalf without any fields
func TestFatalF(t *testing.T) {
	// log.Fatalf calls os.Exit, so executing it as another process.
	msg := "This is an fatal message."
	if os.Getenv("BE_CRASHER") == "1" {
		log := logger.New()
		log.Fatalf("This is an fatal message.")
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestFatalF")
	cmd.Env = append(os.Environ(), "BE_CRASHER=1")
	output, err := cmd.CombinedOutput()
	require.Error(t, err)
	//time="2019-09-10T12:41:33+05:30" level=fatal msg="This is an fatal message." file="log_test.go:128"
	require.Contains(t, string(output), fmt.Sprintf("level=%s msg=\"%s\" file=\"%s:%d\"", "fatal", msg, "log_test.go", getLineNumber()-9))
}

// Test log.Fatalf without any fields
func TestFatalFSkip(t *testing.T) {
	// log.Fatalf calls os.Exit, so executing it as another process.
	buf := &bytes.Buffer{}
	log := customLogger(buf)
	log.SetLevel("panic")
	log.Fatalf("This is an fatal message.")

	//time="2019-09-10T12:41:33+05:30" level=fatal msg="This is an fatal message." file="log_test.go:128"
	require.Contains(t, buf.String(), "")
}

// Test setting log level.
func TestGetLevel(t *testing.T) {
	// log.Fatalf calls os.Exit, so executing it as another process.
	log := logger.New()
	log.SetLevel("debug")
	log.LogLevel()

	require.EqualValues(t, "debug", log.LogLevel())
}

// Get line number of caller.
func getLineNumber() int {
	_, _, line, _ := runtime.Caller(1)
	return line
}
