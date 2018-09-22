package main

import (
	"context"
	"fmt"
	"gin_blog/models"
	"gin_blog/pkg/gredis"
	"gin_blog/pkg/logging"
	"gin_blog/pkg/setting"
	"gin_blog/routers"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

// @title 博客API
// @version v1
// @description 博客简单的api.

// @contact.name Blog Api
// @contact.url https://github.com/yin5th/gin_blog
// @contact.email 541304803@qq.com

// @securitydefinitions.oauth2.password OAuth2Password
// @tokenUrl /auth
// @in header
// @name token

// @host localost:8088
func main() {
	setting.Setup()
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

	//热重启
	go func() {
		if err := s.ListenAndServe(); err != nil {
			log.Printf("Listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}

	log.Println("server exiting")
}
