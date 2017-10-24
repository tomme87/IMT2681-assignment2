package main

import "net/http"

const (
	BasePath = "/exchange"
)

func main() {
	http.HandleFunc(BasePath, handleRootPath)
}
