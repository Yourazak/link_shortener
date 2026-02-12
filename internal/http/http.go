package http

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"ozon_link_shortener/internal/service"

	"github.com/gin-gonic/gin"
)

const formatResp = `{"url": %s}`

func SaveUrl(ctx *gin.Context, s *service.Service) {
	req := service.Request{}
	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		log.Println(err)
		return
	}

	err = json.Unmarshal(body, &req)
	if err != nil {
		log.Println(err)
	}
	shortUrl := service.Shortener(req.LongUrl)

	err = service.WriteToStorage(s, shortUrl, req)
	if err != nil {
		log.Println(err)
		return
	}

	resp := fmt.Sprintf(formatResp, shortUrl)
	ctx.Writer.Write([]byte(resp))
	log.Println("Generated short URL:", shortUrl)
}

func GetUrl(ctx *gin.Context, s *service.Service) {
	shortUrl := ctx.Param("shortUrl")

	longUrl, err := service.GetFromStorage(s, shortUrl)
	if err != nil {
		log.Println(err)
		ctx.Writer.WriteHeader(http.StatusNotFound)
		return
	}
	log.Printf("Short URL: %s, Long URL: %s", shortUrl, longUrl)
	http.Redirect(ctx.Writer, ctx.Request, longUrl, http.StatusMovedPermanently)
}
