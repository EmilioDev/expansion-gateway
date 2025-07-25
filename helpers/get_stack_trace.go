package helpers

import (
	"fmt"
	"runtime"
)

// GetStackTrace returns a slice of strings where each entry represents a stack frame.
// skipFrames determines how many initial stack frames to skip (e.g., 0 = current function).
func GetStackTrace(skipFrames int) []string {
	const maxDepth = 32
	var pcs [maxDepth]uintptr

	n := runtime.Callers(skipFrames+2, pcs[:]) // +2 to skip GetStackTrace and its caller
	frames := runtime.CallersFrames(pcs[:n])

	var stack []string

	for {
		frame, more := frames.Next()
		stack = append(stack, fmt.Sprintf("%s (%s:%d)", frame.Function, frame.File, frame.Line))

		if !more {
			break
		}
	}

	return stack
}
