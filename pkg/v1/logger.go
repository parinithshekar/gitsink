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

package v1

import "io"

// Methods required for a logger.
type Logger interface {
	Fatalf(message string, args ...interface{})
	Panicf(message string, args ...interface{})
	Debugf(message string, args ...interface{})
	Errorf(message string, args ...interface{})
	Infof(message string, args ...interface{})
	LogLevel() string
	SetLevel(string)
	Tracef(message string, args ...interface{})
	Warningf(message string, args ...interface{})

	WithError(error) Logger
	WithField(string, interface{}) Logger
	WithFields(map[string]interface{}) Logger
	AutoClearFields(bool)

	SetOutput(io.Writer)
}
