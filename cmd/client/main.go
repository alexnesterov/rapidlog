package main

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

func main() {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	req, err := http.NewRequest(
		http.MethodGet,
		"http://localhost:1508/",
		nil,
	)
	if err != nil {
		fmt.Println("Ошибка формирования запроса:", err)
		return
	}

	res, err := client.Do(req)
	if err != nil {
		fmt.Println("Ошибка отправки запроса:", err)
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
