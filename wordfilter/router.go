package wordfilter

import (
	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
)

func NewInterfaceHttp() *InterfaceHttp {
	return &InterfaceHttp{}
}

func (http *InterfaceHttp) Run() {
	router := gin.Default()
	router.HEAD("/info", func(context *gin.Context) {
		context.JSON(200, "ok")
	})
	router.GET("/api/reload", http.reload)
	router.POST("/api/find-words", http.search)
	addr := NewConfig().Addr
	_ = endless.ListenAndServe(addr, router)
}
