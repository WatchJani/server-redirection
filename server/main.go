package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	f := os.Args

	port := f[1]

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.RemoteAddr)
	})

	http.ListenAndServe(":"+port, nil)
}
