package main

import (
	"net/http"
	"fmt"
)

func handleRootPath(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "hei1")
}

func handleRootPath2(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "hei2")
}

func handleRootPath3(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "hei3")
}
