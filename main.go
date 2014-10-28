package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/garyburd/redigo/redis"
)

// PaxMessage representation
type PaxMessage struct {
	Callsign    string `json:"callsign"`
	Subject     string `json:"subject"`
	Description string `json:"description"`
}

func main() {

	conn, err := redis.Dial("tcp", ":6379")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	msg, err := readNext(conn)
	if err != nil {
		fmt.Println("error:", err)
	} else {
		code, err := sendMessage(msg)
		if err != nil {
			fmt.Println("error:", err)
		} else {
			fmt.Println("Response from API", code)
		}
	}
}

func sendMessage(msg PaxMessage) (string, error) {
	fmt.Println("send", msg)

	json, err := json.Marshal(msg)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", "https://support.zendesk.dev/requests/emergency/create.json", bytes.NewReader(json))
	req.Header["Content-Type"] = []string{"application/json"}
	req.Header["Accept"] = []string{"application/json"}

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	return resp.Status, nil
}

func readNext(conn redis.Conn) (PaxMessage, error) {
	var message PaxMessage

	data, err := conn.Do("RPOP", "paz:inq")
	if err != nil {
		return message, err
	}

	if data == nil {
		return message, errors.New("No data found")
	}

	err = json.Unmarshal(data.([]byte), &message)
	if err != nil {
		return message, err
	}

	return message, nil
}
