/*
Copyright (c) 2014-2020 CGCL Labs
Container_Migrate is licensed under Mulan PSL v2.
You can use this software according to the terms and conditions of the Mulan PSL v2.
You may obtain a copy of Mulan PSL v2 at:
         http://license.coscl.org.cn/MulanPSL2
THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND,
EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT,
MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
See the Mulan PSL v2 for more details.
*/
package libcontainer

import "io"

// ErrorCode is the API error code type.
type ErrorCode int

// API error codes.
const (
	// Factory errors
	IdInUse ErrorCode = iota
	InvalidIdFormat

	// Container errors
	ContainerNotExists
	ContainerPaused
	ContainerNotStopped
	ContainerNotRunning
	ContainerNotPaused

	// Process errors
	NoProcessOps

	// Common errors
	ConfigInvalid
	ConsoleExists
	SystemError
)

func (c ErrorCode) String() string {
	switch c {
	case IdInUse:
		return "Id already in use"
	case InvalidIdFormat:
		return "Invalid format"
	case ContainerPaused:
		return "Container paused"
	case ConfigInvalid:
		return "Invalid configuration"
	case SystemError:
		return "System error"
	case ContainerNotExists:
		return "Container does not exist"
	case ContainerNotStopped:
		return "Container is not stopped"
	case ContainerNotRunning:
		return "Container is not running"
	case ConsoleExists:
		return "Console exists for process"
	case ContainerNotPaused:
		return "Container is not paused"
	case NoProcessOps:
		return "No process operations"
	default:
		return "Unknown error"
	}
}

// Error is the API error type.
type Error interface {
	error

	// Returns an error if it failed to write the detail of the Error to w.
	// The detail of the Error may include the error message and a
	// representation of the stack trace.
	Detail(w io.Writer) error

	// Returns the error code for this error.
	Code() ErrorCode
}
