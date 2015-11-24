package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

func ClinetSays(message string) {
	fmt.Printf("[CLIENT]:", message)
}

func ServerSays(message string) {
	fmt.Println("[SERVER]:", message)
}
func receiveMessage(channel chan string) bool {
	for {
		message := <-channel
		message = strings.TrimSpace(message)
		if message == "Until then stranger" {
			ServerSays(message)
			return true
		} else if message != "" {
			ServerSays(message)
			return false
		}
	}
	return false
}

func client() {
	fmt.Println("[SERVER]:", "<wait for connection on TCP port 1983>")
	c, err := net.Dial("tcp", "127.0.0.1:1983")
	channel := make(chan string)
	reader := bufio.NewReader(c)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("[CLIENT]:", "<open connection>")
	ServerSays("Welcome stranger")

	go func() {
		for {
			response, _ := reader.ReadString('\n')
			if response != "" {
				channel <- response
			}
		}
	}()

	for {
		var msg string
		fmt.Printf("[CLIENT]: ")
		fmt.Scanf("%s\n", &msg)
		c.Write([]byte(msg + "\n"))
		if receiveMessage(channel) {
			break
		}
	}
	fmt.Println("[CLIENT]:", "<open connection>")
	ServerSays("<wait for next connection>")
	c.Close()
}

func main() {
	client()
}
