//go:build wasip2

package syscall

import "github.com/ydnar/wasm-tools-go/wasi/cli/environment"

func Getwd() (string, error) {
	result := environment.InitialCWD()
	if wd, ok := result.Some(); ok {
		return wd, nil
	}
	return "", EINVAL
}
