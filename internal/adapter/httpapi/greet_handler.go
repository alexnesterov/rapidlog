// Package httpapi contains HTTP API handlers
package httpapi

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func Greet(w http.ResponseWriter, r *http.Request) {
	today := time.Now().Format("02.01 Mon 06")

	if _, err := fmt.Fprintf(w, "%s\nHello World!", today); err != nil {
		log.Println(err)
	}
}
