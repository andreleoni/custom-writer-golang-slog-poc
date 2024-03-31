package main

import (
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

type CustomWriter struct{}

func (CustomWriter) Write(b []byte) (n int, err error) {
	go http.Post("http://errortracker:9090/log", "application/json", strings.NewReader(string(b)))

	return 1, nil
}

type User struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (u User) LogValue() slog.Value {
	u.Email = "[REDACTED]"

	jsonUser, err := json.Marshal(u)
	if err != nil {
		log.Fatal(err)
	}

	return slog.StringValue(string(jsonUser))
}

func main() {
	for {
		logger := slog.New(slog.NewJSONHandler(&CustomWriter{}, nil))

		userName := fmt.Sprint("Name", randomStringRunes(4))
		userEmail := fmt.Sprint("Email", randomStringRunes(4))

		user := User{Name: userName, Email: userEmail}

		logger.Info(
			fmt.Sprint("Random error: ", randomStringRunes(10)),
			"user", user,
		)

		time.Sleep(1 * time.Second)
	}
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randomStringRunes(n int) string {
	b := make([]rune, n)

	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}

	return string(b)
}
