package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/didip/tollbooth/v7"
	"github.com/didip/tollbooth/v7/limiter"
	"github.com/gin-gonic/gin"
)

func main() {
	txtSplit := strings.Split("PIESEAT 航天宏图", " ")
	for _, s := range txtSplit {
		fmt.Println("s: ", s)
	}
	return

	newLimiter := tollbooth.NewLimiter(1, &limiter.ExpirableOptions{
		DefaultExpirationTTL: time.Hour,
	})
	newLimiter = newLimiter.
		SetBasicAuthUsers([]string{"user1"}).
		SetIPLookups([]string{"RemoteAddr", "X-Forwarded-For", "X-Real-IP"}).
		SetHeaders(nil).
		SetBurst(10)

	engine := gin.Default()
	//engine.Use(LimitDefaultGin(newLimiter))
	engine.Use(LimitByKeys(newLimiter))

	engine.GET("/ping", func(c *gin.Context) {
		//c.String(200, "pong"+time.Now().Format(time.TimeOnly))
	})

	go func() {
		for i := 0; i < 15; i++ {
			go TestRequest("token_1")
			//time.Sleep(1 * time.Second)
		}
	}()

	go func() {
		fmt.Println("=========================================== send request ============================================")
		time.Sleep(time.Second * 2)
		for i := 0; i < 15; i++ {
			go TestRequest("token_1")
		}
	}()

	log.Fatal(engine.Run(":8089"))
}

func LimitDefaultGin(limiter *limiter.Limiter) gin.HandlerFunc {
	return func(context *gin.Context) {
		if limiterErr := tollbooth.LimitByRequest(limiter, context.Writer, context.Request); limiterErr != nil {
			context.Writer.WriteHeader(http.StatusTooManyRequests)
			context.Writer.Write([]byte("Too many requests, please try again later."))
			context.Abort()
			return
		}
		context.Next()
	}
}

func LimitByKeys(limiter *limiter.Limiter) gin.HandlerFunc {
	return func(context *gin.Context) {
		t := context.Query("token")
		if limiterErr := tollbooth.LimitByKeys(limiter, []string{t}); limiterErr != nil {
			context.Writer.WriteHeader(http.StatusTooManyRequests)
			context.Writer.Write([]byte("Too many requests, please try again later."))
			context.Abort()
			return
		}
		context.Next()
	}
}

func LimitByIp(limiter *limiter.Limiter) gin.HandlerFunc {
	return func(context *gin.Context) {
		if limiterErr := tollbooth.LimitFuncHandler(limiter, func(writer http.ResponseWriter, request *http.Request) {

		}); limiterErr != nil {
			context.Writer.WriteHeader(http.StatusTooManyRequests)
			context.Writer.Write([]byte("Too many requests, please try again later."))
			context.Abort()
			return
		}
		context.Next()
	}
}

func TestRequest(key string) {
	request, err := http.NewRequest(http.MethodGet, "http://user1@localhost:8089/ping?token="+key, nil)
	if err != nil {
		panic(err)
	}

	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		panic(err)
	}
	_, err = io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	//fmt.Printf("token:%s resp: %+v\n", key, string(bytes))
}
