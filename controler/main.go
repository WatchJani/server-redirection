package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"root/consistent_hashing"
)

type ReqRedirection struct {
	*consistent_hashing.ConsistentHashing
}

func NewReqRedirection() *ReqRedirection {
	c := consistent_hashing.NewConsistentHashing()
	if err := c.Load("./config.json"); err != nil {
		log.Println(err)
	}

	return &ReqRedirection{
		ConsistentHashing: c,
	}
}

func (d *ReqRedirection) Redirect(w http.ResponseWriter, r *http.Request) {
	var f struct {
		Key string `json:"key"`
	}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&f); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	newIpAddress := d.FindServer(f.Key).Addr
	fmt.Println(newIpAddress)
	http.Redirect(w, r, newIpAddress, http.StatusMovedPermanently)
}

func main() {
	redirection := NewReqRedirection()
	http.HandleFunc("POST /query", redirection.Redirect)

	http.ListenAndServe(":5000", nil)
}
