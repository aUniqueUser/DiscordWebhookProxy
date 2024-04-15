package main

import (
	"bytes"
	json2 "encoding/json"
	"fmt"
	"net/http"
)

func main() {
	data := map[string]string{"message": ":3"}
	json, _ := json2.Marshal(data)
	_, err := http.Post("http://localhost:3333/api/webhooks/[REDACTED]/[REDACTED]", "application/json", bytes.NewBuffer(json))
	if err != nil {
		fmt.Println(err)
		return
	}
}
