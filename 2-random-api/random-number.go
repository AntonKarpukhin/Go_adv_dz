package main

import (
	"fmt"
	"math/rand"
	"net/http"
)

type RandomAPI struct{}

func NewRandomAPI(router *http.ServeMux) {
	randData := &RandomAPI{}
	router.HandleFunc("/", randData.getRandomNumber)
}

func (randNum *RandomAPI) getRandomNumber(w http.ResponseWriter, r *http.Request) {
	num := rand.Intn(6)
	_, err := w.Write([]byte(fmt.Sprintf("%d", num)))
	if err != nil {
		fmt.Println(err.Error())
	}
	return
}
