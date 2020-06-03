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
//
// This work incorporates works covered by the following notices:
//
// Dave Cheney <dave@cheney.net>
// Copyright (c) 2013 Dave Cheney. All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are
// met:
//
//    * Redistributions of source code must retain the above copyright
// notice, this list of conditions and the following disclaimer.
//    * Redistributions in binary form must reproduce the above
// copyright notice, this list of conditions and the following disclaimer
// in the documentation and/or other materials provided with the
// distribution.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
// "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
// LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
// A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
// OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
// SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
// LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
// THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

package v1

import (
	"os"

	profile "github.com/pkg/profile"
)

// noop is a struct that we use to explicitly do nothing.
type noop struct{}

// Stop is a noop that we use to defer stopping the profiler.
// ref: https://stackoverflow.com/questions/29775836/no-op-explicitly-do-nothing-in-go
func (n *noop) Stop() {}

// Profile parses the PROFILE_MODE environment variable and executes the proper profiling task.
func Start() interface {
	Stop()
} {
	switch os.Getenv("PROFILER_MODE") {
	case "cpu":
		return profile.Start(profile.CPUProfile, profile.ProfilePath("."), profile.NoShutdownHook)
	case "mem":
		return profile.Start(profile.MemProfile, profile.ProfilePath("."), profile.NoShutdownHook)
	case "mutex":
		return profile.Start(profile.MutexProfile, profile.ProfilePath("."), profile.NoShutdownHook)
	case "block":
		return profile.Start(profile.BlockProfile, profile.ProfilePath("."), profile.NoShutdownHook)
	}
	return new(noop)
}
