package middleware

import (
	"log"
	"net/http"
	"runtime/debug"
	"time"
)

type customWriter struct {
	http.ResponseWriter
	statusCode int
}

func (cw *customWriter) WriteHeader(code int) {
	cw.statusCode = code
	cw.ResponseWriter.WriteHeader(code)
}

type logProvider struct {
	log *log.Logger
}

func NewLogger (l *log.Logger) *logProvider {
	return &logProvider{
		l,
	}
}

func LoggingHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		cw := &customWriter{w, http.StatusOK}
		t1 := time.Now()
		next.ServeHTTP(cw, r)
		t2 := time.Now()

		log.Printf("[%d] [%s:] [%q] [%v]\n",  cw.statusCode, r.Method, r.RequestURI, t2.Sub(t1))
	}
	return http.HandlerFunc(fn)
}

// Prevent abnormal shutdown while panic
func RecoverHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("panic: %+v", err)
				log.Println(string(debug.Stack()))
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

