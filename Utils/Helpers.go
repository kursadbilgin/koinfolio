package Utils

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func RecoverPanic(c *gin.Context) {
	if rec := recover(); rec != nil {
		fmt.Printf("[RECOVERED PANIC] %+v", rec)
	}
	fmt.Printf("Received %s %s [Rsp: {%+v}]", c.Request, c.Request.URL, c.Request.Response)
}
