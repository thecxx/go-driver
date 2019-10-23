package unexpected

import (
	"bytes"
	"fmt"
	"os"
	"runtime"
)

// Check error
func CheckError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "err: %s", err.Error())
		os.Exit(-1)
	}
}

// Watch panic
type PanicReceiver func(pv interface{}, fun, gid, stack string)

func WatchPanic(fun string, handler PanicReceiver) {
	if pv := recover(); pv != nil {
		buf := make([]byte, 4096)
		num := runtime.Stack(buf, false)
		gid := "0"
		sts := 0
		// No head found
		if bytes.HasPrefix(buf, []byte{'g', 'o', 'r', 'o', 'u', 't', 'i', 'n', 'e', ' '}) {
			i := 10
			for ; ; i++ {
				if buf[i] < '0' || '9' < buf[i] {
					break
				}
			}
			gid = string(buf[10:i])
			sts = bytes.IndexByte(buf, '\n') + 1
		}
		// Call handler
		handler(pv, fun, gid, string(buf[sts:num]))
	}
}
