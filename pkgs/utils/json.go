package utils

import (
	"encoding/json"
	"log"
)

func JsonStatus(message string) []byte {
	m, _ := json.Marshal(struct {
		Message string `json:"message"`
	}{
		Message: message,
	})
	log.Printf("%s", string(m))
	return m
}
