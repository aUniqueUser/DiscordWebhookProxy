package main

import (
	"bytes"
	json2 "encoding/json"
	"github.com/joho/godotenv"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

// SendMessage We only pass a text - no other discord defined json template used in this implementation!
func SendMessage(webhookID string, webHookToken string, message string) {
	uri := "https://discord.com/api/webhooks/" + webhookID + "/" + webHookToken
	data := map[string]string{"content": message}
	json, _ := json2.Marshal(data)
	_, err := http.Post(uri, "application/json", bytes.NewBuffer(json))
	if err != nil {
		log.Println("[ERROR] Failed to send message to Discord")
	}
}

type Request struct {
	Message string `json:"message"`
}

func HandleRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		log.Println("[INFO] Received an invalid HTTP method")
		w.WriteHeader(405) // Forbidden method
		return
	}

	url := r.URL.String()
	if strings.HasPrefix(url, "/api/webhooks/") && len(url) > 13 {
		sub := url[14:len(url)]
		res := strings.Split(sub, "/")
		if len(res) == 2 {
			body, _ := io.ReadAll(r.Body)

			var request Request
			if err := json2.Unmarshal(body, &request); err != nil {
				log.Println("[ERROR] Failed to unmarshal JSON")
				request.Message = "```[PROXY USAGE ERROR] No message post/empty form was included```"
			}

			log.Printf("[INFO] Sending a message to discord (ID: %s), (MSG: %s)", res[0], request.Message)
			SendMessage(res[0], res[1], request.Message)
			w.WriteHeader(200)
			_, err := w.Write([]byte("Success!"))
			if err != nil {
				return
			}
			return
		}
	}

	w.WriteHeader(404)
	return
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("[FATAL] Failed to load .env file, pls check your configuration!")
	}
	err = http.ListenAndServe(":"+os.Getenv("WEBHOOK_PROXY_PORT"), http.HandlerFunc(HandleRequest))
	if err != nil {
		log.Fatal("[FATAL] Failed to start server")
	}
}
