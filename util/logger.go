package util

import (
	"fmt"
	"runtime"

	"github.com/gin-gonic/gin"
)

func MyLogger(msg ...string) {
	pc := make([]uintptr, 1)
	runtime.Callers(2, pc)
	fmt.Fprintln(gin.DefaultWriter, "MyLogger: ", runtime.FuncForPC(pc[0]).Name(), msg)
}
