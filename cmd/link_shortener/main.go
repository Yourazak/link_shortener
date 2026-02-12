package main

import (
	"log"
	"net/http"
	"ozon_link_shortener/internal/app"
	h "ozon_link_shortener/internal/http"
	mw "ozon_link_shortener/internal/middleware"
	"ozon_link_shortener/internal/service"
	"ozon_link_shortener/internal/storage"

	"github.com/gin-gonic/gin"
)

func main() {
	s := service.Service{}
	s.Conf = app.InitEnv()
	s.Cache = storage.NewCache()

	if s.Conf.Storage == "postgres" {
		db, err := storage.NewDB(s.Conf.DBAddress)
		if err != nil {
			log.Fatal(err)
		}
		s.DB = db
	}

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(mw.LoggingMiddleware)
	r.Use(gin.Recovery())

	r.POST("/", mw.ValidatePostRequest, func(c *gin.Context) {
		h.SaveUrl(c, &s)
	})
	r.GET("/:shortUrl", mw.ValidateGetRequest, func(c *gin.Context) {
		h.GetUrl(c, &s)
	})

	server := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	app.GracefulShutdown(&s, server)

	server.ListenAndServe()
}
