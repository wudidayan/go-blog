package main

import (
	"fmt"
	"net/http"

	"github.com/wudidayan/go-blog/models"
	"github.com/wudidayan/go-blog/pkg/gredis"
	"github.com/wudidayan/go-blog/pkg/logging"
	"github.com/wudidayan/go-blog/pkg/setting"
	"github.com/wudidayan/go-blog/routers"
)

func main() {
	setting.Setup("conf/app.ini")
	models.Setup()
	logging.Setup()
	gredis.Setup()

	router := routers.InitRouter()
	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", setting.ServerSetting.HttpPort),
		Handler:        router,
		ReadTimeout:    setting.ServerSetting.ReadTimeout,
		WriteTimeout:   setting.ServerSetting.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}

	s.ListenAndServe()
}
