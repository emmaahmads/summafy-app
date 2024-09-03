package util

import (
	"fmt"
	"runtime"

	"github.com/gin-gonic/gin"
)

func MyGinLogger(msg ...string) {
	pc := make([]uintptr, 1)
	runtime.Callers(2, pc)
	fmt.Fprintln(gin.DefaultWriter, "MyGinLogger: ", runtime.FuncForPC(pc[0]).Name(), msg)
}
