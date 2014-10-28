package main

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/garyburd/redigo/redis"
)

// PaxMessage representation
type PaxMessage struct {
	Callsign, Subject, Description string
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
		sendMessage(msg)
	}

}

func sendMessage(msg PaxMessage) (string, error) {
	fmt.Println("send", msg)
	return "", nil
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
