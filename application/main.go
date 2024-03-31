package main

import (
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/andreleoni/random"
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

		userName := fmt.Sprint("Name", random.Letters(4))
		userEmail := fmt.Sprint("Email", random.Letters(4))

		user := User{Name: userName, Email: userEmail}

		logger.Info(
			fmt.Sprint("Random error: ", random.Letters(10)),
			"user", user,
		)

		time.Sleep(1 * time.Second)
	}
}
