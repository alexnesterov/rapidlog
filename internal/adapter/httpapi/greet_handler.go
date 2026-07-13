// Package httpapi contains HTTP API handlers
package httpapi

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func Greet(w http.ResponseWriter, r *http.Request) {
	if _, err := fmt.Fprintf(w, "Hello World! %s", time.Now()); err != nil {
		log.Println(err)
	}
}
