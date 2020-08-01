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
	"strings"
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
	expectedLog := fmt.Sprintf("{\"file\":\"%s:%d\",\"func\":\"%s\",\"level\":\"%s\",\"msg\":\"%s\",\"time\":\"", "log_test.go", getLineNumber()-2, getFuncName(), "error", "This is an error message.")
	require.Contains(t, buf.String(), expectedLog)
}

// Test Waringf
func TestLogWarningf(t *testing.T) {
	buf := &bytes.Buffer{}
	log := customLogger(buf)
	err := errors.New("This is an warning.")
	log.Warningf("%s", err.Error())
	// "time="2019-09-10T10:51:43+05:30" level=error msg="This is an error message." file="log_test.go:43"
	expectedLog := fmt.Sprintf("{\"file\":\"%s:%d\",\"func\":\"%s\",\"level\":\"%s\",\"msg\":\"%s\",\"time\":\"", "log_test.go", getLineNumber()-2, getFuncName(), "warning", "This is an warning.")
	require.Contains(t, buf.String(), expectedLog)
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
	expectedLog := fmt.Sprintf("{\"error\":\"%s\",\"file\":\"%s:%d\",\"func\":\"%s\",\"level\":\"%s\",\"msg\":\"%s\",\"time\":\"", "This is a custom error.", "log_test.go", getLineNumber()-2, getFuncName(), "error", "Encountered an error.")
	require.Contains(t, buf.String(), expectedLog)

	// Check if error field is cleared.
	log.Errorf("Encountered an error.")
	// time="2019-09-10T10:56:17+05:30" level=error msg="Encountered an error." file="log_test.go:53" error="This is a custom error."
	expectedLog = fmt.Sprintf("{\"file\":\"%s:%d\",\"func\":\"%s\",\"level\":\"%s\",\"msg\":\"%s\",\"time\":\"", "log_test.go", getLineNumber()-2, getFuncName(), "error", "Encountered an error.")
	require.Contains(t, buf.String(), expectedLog)
}

// Test logger with WithField method.
func TestWithField(t *testing.T) {
	buf := &bytes.Buffer{}
	log := customLogger(buf)
	log.WithField("key", "value").WithField("key2", "value2").Errorf("Errors with custom field.")
	// time="2019-09-10T11:17:19+05:30" level=error msg="Errors with custom field." file="log_test.go:63" key=value key2=value2
	expectedLog := fmt.Sprintf("{\"file\":\"%s:%d\",\"func\":\"%s\",\"key\":\"value\",\"key2\":\"value2\",\"level\":\"%s\",\"msg\":\"%s\",\"time\":\"", "log_test.go", getLineNumber()-2, getFuncName(), "error", "Errors with custom field.")
	require.Contains(t, buf.String(), expectedLog)
	log.WithField("key3", "value3").Errorf("Errors with custom field again.")
	// time="2019-09-10T11:17:19+05:30" level=error msg="Errors with custom field." file="log_test.go:63" key=value key2=value2
	expectedLog = fmt.Sprintf("{\"file\":\"%s:%d\",\"func\":\"%s\",\"key3\":\"value3\",\"level\":\"%s\",\"msg\":\"%s\",\"time\":\"", "log_test.go", getLineNumber()-2, getFuncName(), "error", "Errors with custom field again.")
	require.Contains(t, buf.String(), expectedLog)
}

// Test concurrent modifications to fields.
func TestConcurrentMods(t *testing.T) {
	buf := &bytes.Buffer{}
	log := customLogger(buf)
	log.SetLevel("debug")
	for i := 0; i < 15000; i++ {
		go func() {
			log.WithError(pkg.ErrNoMatch).WithField("key", "value").Debugf("")
			expectedLog := fmt.Sprintf("{\"error\":\"No alert matched in alert config.\",\"file\":\"%s:%d\",\"func\":\"%s\",\"key\":\"value\",\"level\":\"%s\",\"msg\":\"\",\"time\":\"", "log_test.go", getLineNumber()-1, getFuncName(), "debug")
			require.Contains(t, buf.String(), expectedLog)
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
	expectedLog := fmt.Sprintf("{\"file\":\"%s:%d\",\"func\":\"%s\",\"key\":\"value\",\"key2\":\"value2\",\"level\":\"%s\",\"msg\":\"%s\",\"time\":\"", "log_test.go", getLineNumber()-2, getFuncName(), "error", "Errors with custom field.")
	require.Contains(t, buf.String(), expectedLog)

	log.WithField("key3", "value3").Errorf("Errors with custom field again.")
	expectedLog = fmt.Sprintf("{\"file\":\"%s:%d\",\"func\":\"%s\",\"key\":\"value\",\"key2\":\"value2\",\"key3\":\"value3\",\"level\":\"%s\",\"msg\":\"%s\",\"time\":\"", "log_test.go", getLineNumber()-1, getFuncName(), "error", "Errors with custom field again.")
	require.Contains(t, buf.String(), expectedLog)
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
	expectedLog := fmt.Sprintf("{\"file\":\"%s:%d\",\"func\":\"%s\",\"key1\":\"val1\",\"key2\":\"val2\",\"key3\":\"val3\",\"level\":\"%s\",\"msg\":\"%s\",\"time\":\"", "log_test.go", getLineNumber()-2, getFuncName(), "debug", "Errors with custom fields.")
	require.Contains(t, buf.String(), expectedLog)
}

// Test debug log level enabled
func TestWithDebugEnabled(t *testing.T) {
	buf := &bytes.Buffer{}
	log := customLogger(buf)
	log.SetLevel("deBug")
	log.Debugf("Debug log enabled.")
	// time="2019-09-10T11:25:46+05:30" level=debug msg="Debug log enabled." file="log_test.go:88
	expectedLog := fmt.Sprintf("{\"file\":\"%s:%d\",\"func\":\"%s\",\"level\":\"%s\",\"msg\":\"%s\",\"time\":\"", "log_test.go", getLineNumber()-2, getFuncName(), "debug", "Debug log enabled.")
	require.Contains(t, buf.String(), expectedLog)
}

// Test debug log level disabled
func TestWithDebugDisabled(t *testing.T) {
	buf := &bytes.Buffer{}
	log := customLogger(buf)
	log.SetLevel("info")
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
	expectedLog := fmt.Sprintf("{\"file\":\"%s:%d\",\"func\":\"%s\",\"level\":\"%s\",\"msg\":\"%s\",\"time\":\"", "log_test.go", getLineNumber()-2, getFuncName(), "info", msg)
	require.Contains(t, buf.String(), expectedLog)
}

// Test log.Tracef without any fields
func TestLogTrace(t *testing.T) {
	buf := &bytes.Buffer{}
	log := customLogger(buf)
	log.SetLevel("trace")
	msg := "This is an trace message."
	log.Tracef(msg)
	// "time="2019-09-10T10:51:43+05:30" level=info msg="This is an info message." file="log_test.go:118"
	expectedLog := fmt.Sprintf("{\"file\":\"%s:%d\",\"func\":\"%s\",\"level\":\"%s\",\"msg\":\"%s\",\"time\":\"", "log_test.go", getLineNumber()-2, getFuncName(), "trace", msg)
	require.Contains(t, buf.String(), expectedLog)
}

// Test log.Debugf without any fields
func TestLogDebug(t *testing.T) {
	buf := &bytes.Buffer{}
	log := customLogger(buf)
	log.SetLevel("debug")
	msg := "This is an debug message."
	log.Debugf(msg)
	// "time="2019-09-10T10:51:43+05:30" level=info msg="This is an info message." file="log_test.go:118"
	expectedLog := fmt.Sprintf("{\"file\":\"%s:%d\",\"func\":\"%s\",\"level\":\"%s\",\"msg\":\"%s\",\"time\":\"", "log_test.go", getLineNumber()-2, getFuncName(), "debug", msg)
	require.Contains(t, buf.String(), expectedLog)
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
	expectedLog := fmt.Sprintf("{\"file\":\"%s:%d\",\"func\":\"%s\",\"level\":\"%s\",\"msg\":\"%s\",\"time\":\"", "log_test.go", getLineNumber()-2, getFuncName(), "panic", msg)
	require.Contains(t, buf.String(), expectedLog)
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
	expectedLog := fmt.Sprintf("{\"file\":\"%s:%d\",\"func\":\"%s\",\"level\":\"%s\",\"msg\":\"%s\",\"time\":\"", "log_test.go", getLineNumber()-9, getFuncName(), "fatal", msg)
	require.Contains(t, string(output), expectedLog)
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

func getFuncName() string {
	pc, _, _, _ := runtime.Caller(1)
	function := runtime.FuncForPC(pc)

	functionSplit := strings.Split(function.Name(), ".")
	functionName := functionSplit[len(functionSplit)-1]
	return functionName
}
