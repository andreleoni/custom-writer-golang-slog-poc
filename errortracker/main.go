package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"log/slog"
	"net/http"
	"os"
)

var errors = make(chan []byte)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	file, err := os.OpenFile("errors.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Only records logs on mongodb
	go func() {
		for {
			errbyte := <-errors

			if _, err := file.Write(errbyte); err != nil {
				log.Fatal(err)
			}
		}
	}()

	http.HandleFunc("/log", func(w http.ResponseWriter, r *http.Request) {
		var bodyBytes []byte
		var err error

		if r.Body != nil {
			bodyBytes, err = ioutil.ReadAll(r.Body)

			if err != nil {
				fmt.Printf("Body reading error: %v", err)
				return
			}

			errors <- bodyBytes

			defer r.Body.Close()
		}
	})

	fmt.Println("Listening on localhost:9090...")

	log.Fatal(http.ListenAndServe(":9090", nil))
}
