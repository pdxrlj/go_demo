package main

import (
	"fmt"
	"log"
	"net/http"

	limiter "github.com/ulule/limiter/v3"
	mhttp "github.com/ulule/limiter/v3/drivers/middleware/stdlib"
	"github.com/ulule/limiter/v3/drivers/store/memory"
)

func main() {

	// Define a limit rate to 4 requests per hour.
	rate, err := limiter.NewRateFromFormatted("1-S")
	if err != nil {
		log.Fatal(err)
		return
	}

	// Create a redis client.
	store := memory.NewStore()

	// Create a new middleware with the limiter instance.
	middleware := mhttp.NewMiddleware(limiter.New(store, rate, limiter.WithTrustForwardHeader(true)))
	//newMiddlewareGin := mgin.NewMiddleware(limiter.New(store, rate, limiter.WithTrustForwardHeader(true)))

	// Launch a simple server.
	http.Handle("/", middleware.Handler(http.HandlerFunc(index)))
	fmt.Println("Server is running on port 7777...")
	log.Fatal(http.ListenAndServe(":7777", nil))

}

func index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_, err := w.Write([]byte(`{"message": "ok"}`))
	if err != nil {
		log.Fatal(err)
	}
}
