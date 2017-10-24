package main

import "net/http"

const (
	BasePath = "/exchange"
)

func main() {
	http.HandleFunc(BasePath, handleRootPath)
	http.HandleFunc(BasePath + "/", handleRootPath2)
	http.HandleFunc(BasePath + "/lol", handleRootPath3)
	http.ListenAndServe(":8080", nil)
}
