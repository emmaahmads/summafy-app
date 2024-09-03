package util

import (
	"fmt"
	"runtime"

	"github.com/gin-gonic/gin"
)

func MyLogger(f string, msg ...string) {

	fmt.Fprintln(gin.DefaultWriter, "MyLogger: ", f, msg)
}

func MyFunc() string {
	pc := make([]uintptr, 1)
	runtime.Callers(2, pc)

	return runtime.FuncForPC(pc[0]).Name()
}
