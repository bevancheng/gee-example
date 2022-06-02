package gee

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strings"
)

func Recovery() HandlerFunc {
	return func(c *Context) {
		defer func() {
			if err := recover(); err != nil {
				message := fmt.Sprintf("%s", err)
				log.Printf("%s\n\n", trace(message))
				c.Fail(http.StatusInternalServerError, "Internal Server Error")
			}
		}()
		c.Next()
	}

}

func trace(message string) string {
	var pcs [32]uintptr
	//Callers用来返回调用栈的程序计数器，第0个Caller是Callers本身，第一个是上层trace，第二个是再上层defer func
	//为了日志简介，跳过前三个Caller
	n := runtime.Callers(3, pcs[:])

	var str strings.Builder
	str.WriteString(message + "\nTraceback:")
	for _, pc := range pcs[:n] {
		//FuncForPC获取对应函数
		fn := runtime.FuncForPC(pc)
		//FileLine获取调用该函数的文件名和行号
		file, line := fn.FileLine(pc)
		str.WriteString(fmt.Sprintf("\n\t%s:%d", file, line))
	}
	return str.String()
}
