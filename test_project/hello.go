package main

import (
	"fmt"
	"net/http"
	"io"

	"golang.org/x/example/hello/reverse"

)

func main() {
	fmt.Println("Hello, World!")
	fmt.Println(reverse.String("Hello"))
	resp, err := http.Get("http://example.com/")
	if err != nil {
		fmt.Println("Could not get page")

	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
    fmt.Println("could not read body:", err)
    return
}
	fmt.Println(string(body))
}
