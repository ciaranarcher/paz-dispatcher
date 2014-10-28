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

// PaxMessage ...
type PaxMessage struct {
	Callsign    string `json:"callsign"`
	Subject     string `json:"subject"`
	Description string `json:"description"`
}

// Notification ...
type Notification struct {
	Callsign     string `json:"callsign"`
	Notification string `json:"notification"`
}

func main() {
	done := make(chan bool)

	// From radio queue to API
	go func() {
		conn, err := redis.Dial("tcp", "10.16.2.74:6379")
		if err != nil {
			panic(err)
		}
		defer conn.Close()

		fmt.Println("started ticket creation worker")
		for {
			msg, err := readNext(conn)
			if err != nil {
				time.Sleep(1 * time.Second) // Wait a moment before trying again
			} else {
				code, err := sendMessage(msg)
				if err != nil {
				} else {
					fmt.Println("Send ticket:", code)
				}
			}
		}
	}()

	// To notifications queue to radio queue
	go func() {
		conn, err := redis.Dial("tcp", "10.16.2.74:6379")
		if err != nil {
			panic(err)
		}
		defer conn.Close()

		classicConn, err := redis.Dial("tcp", ":6379")
		if err != nil {
			panic(err)
		}
		defer classicConn.Close()

		fmt.Println("started notifications worker")
		for {
			notification, err := readNotification(classicConn)
			if err != nil {
				time.Sleep(1 * time.Second) // Wait a moment before trying again
			} else {
				err := enqueueNotification(conn, notification)
				if err != nil {
					fmt.Println("error:", err)
				} else {
					fmt.Println("Enqueued notification:", notification)
				}
			}
		}
	}()

	<-done
}

func enqueueNotification(conn redis.Conn, notification Notification) error {

	json, err := json.Marshal(notification)
	if err != nil {
		return err
	}

	_, err = conn.Do("RPUSH", "paz:outq", json)
	if err != nil {
		return err
	}

	return nil
}

func sendMessage(msg PaxMessage) (string, error) {
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

func readNotification(conn redis.Conn) (Notification, error) {
	var notification Notification

	data, err := conn.Do("RPOP", "store:paz:notifications")
	if err != nil {
		return notification, err
	}

	if data == nil {
		return notification, errors.New("No notification data found")
	}

	err = json.Unmarshal(data.([]byte), &notification)
	if err != nil {
		return notification, err
	}

	return notification, nil
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
