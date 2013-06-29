package main

import (
	"log"
	"net/http"
)

func main() {
	TempList = new(GlobalTempList)
	go BackgroundTempUpdate()

	http.HandleFunc("/", IndexHandler)
	http.HandleFunc("/data", DataHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
