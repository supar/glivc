package main

import (
	"github.com/go-martini/martini"
	"net/http"
	"time"
	
	"crypto/sha1"
)

// Logger returns a middleware handler that logs the request as it goes in and the response as it goes out.
func LogMiddle() martini.Handler {
	return func(res http.ResponseWriter, req *http.Request, c martini.Context) {
		var (
			hasher = sha1.New()
			start = time.Now()
			
			key = []byte(start.Format(time.RFC3339Nano))
			
			key_hash	[]byte
		)

		// Write key data
		hasher.Write(key)
		key_hash = hasher.Sum(nil)

		addr := req.Header.Get("X-Real-IP")
		if addr == "" {
			addr = req.Header.Get("X-Forwarded-For")
			if addr == "" {
				addr = req.RemoteAddr
			}
		}
		
		log.Info("%x Started %s %s for %s", key_hash[:5], req.Method, req.URL.Path, addr)

		rw := res.(martini.ResponseWriter)
		c.Next()

		log.Info("%x Completed %v %s in %v",key_hash[:5],  rw.Status(), http.StatusText(rw.Status()), time.Since(start))
	}
}