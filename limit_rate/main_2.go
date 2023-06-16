package main

//func main() {
//	// TODO: Turn back to 3.
//	ipLimiter := tollbooth.NewLimiter(100, &limiter.ExpirableOptions{DefaultExpirationTTL: time.Hour})
//
//	// TODO: Turn back to 10.
//	globalLimiter := NewConcurrentLimiter(3)
//
//	http.Handle("/", globalLimiter.LimitConcurrentRequests(ipLimiter, HelloHandler))
//	log.Fatal(http.ListenAndServe(":8080", nil))
//}
//
//func HelloHandler(w http.ResponseWriter, req *http.Request) {
//	time.Sleep(10 * time.Second)
//	w.Write([]byte("Hello, World!"))
//}
//
//type ConcurrentLimiter struct {
//	max     int
//	current int
//	mut     sync.Mutex
//}
//
//func NewConcurrentLimiter(limit int) *ConcurrentLimiter {
//	return &ConcurrentLimiter{
//		max: limit,
//	}
//}
//
//func (limiter *ConcurrentLimiter) LimitConcurrentRequests(lmt *limiter.Limiter,
//	handler func(http.ResponseWriter, *http.Request)) http.Handler {
//
//	middle := func(w http.ResponseWriter, r *http.Request) {
//
//		limiter.mut.Lock()
//		maxHit := limiter.current == limiter.max
//
//		if maxHit {
//			limiter.mut.Unlock()
//			http.Error(w, http.StatusText(429), http.StatusTooManyRequests)
//			return
//		}
//
//		limiter.current += 1
//		limiter.mut.Unlock()
//
//		defer func() {
//			limiter.mut.Lock()
//			limiter.current -= 1
//			limiter.mut.Unlock()
//		}()
//
//		// There's no rate-limit error, serve the next handler.
//		handler(w, r)
//	}
//
//	return tollbooth.LimitHandler(lmt, http.HandlerFunc(middle))
//}
