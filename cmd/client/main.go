package main

import (
	"fmt"
	"io"
	"net/http"
)

func main() {
	res, err := http.Get("http://localhost:1508/")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() { _ = res.Body.Close() }()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Ошибка чтения:", err)
		return
	}

	fmt.Println(string(body))
}
