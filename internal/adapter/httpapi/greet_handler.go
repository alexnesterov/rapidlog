package httpapi

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

func Greet(w http.ResponseWriter, r *http.Request) {
	result := &strings.Builder{}

	today := time.Now().Format("02.01 Mon 06")
	fmt.Fprintf(result, "%s\n\n", today)

	agent := r.Header.Get("User-Agent")
	fmt.Fprintf(result, "User-Agent: %s\n\n", agent)

	values := r.URL.Query()
	fmt.Fprintf(result, "%s\n", values)

	for k, v := range values {
		fmt.Fprintf(result, "%s: %s\n", k, v)
	}

	name := values.Get("name")
	if name == "" {
		name = "world"
	}
	fmt.Fprintf(result, "\nHello, %s!\n\n", name)

	if _, err := w.Write([]byte(result.String())); err != nil {
		log.Println(err)
	}
}
