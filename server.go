package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"sync"
)

type Game struct {
	clients []*Client
	joins   chan net.Conn
	Answers []string
	Stats   map[string]int
	Mutex   sync.RWMutex
}

func (game *Game) JoinConnection(connection net.Conn) {
	client := NewClient(connection, game)
	game.clients = append(game.clients, client)
}

func (game *Game) Listen() {
	go func() {
		for {
			conn := <-game.joins
			if conn != nil {
				game.JoinConnection(conn)
			}
		}
	}()
}

func (game *Game) IsGameMove(move string) bool {
	isGameMove := false
	for _, answer := range game.Answers {
		if answer == move {
			isGameMove = true
		}
	}
	return isGameMove
}
func (game *Game) SaveStats(answer string, move string) string {
	switch {
	case answer == move:
		game.Stats["D"] = game.Stats["D"] + 1
		return "DRAW"
	case answer == "Rock" && move == "SCISSORS":
		game.Stats["L"] = game.Stats["L"] + 1
		return "LOSE"
	case answer == "ROCK" && move == "PAPER":
		game.Stats["W"] = game.Stats["W"] + 1
		return "WIN"
	case answer == "SCISSORS" && move == "ROCK":
		game.Stats["W"] = game.Stats["W"] + 1
		return "WIN"
	case answer == "SCISSORS" && move == "PAPER":
		game.Stats["L"] = game.Stats["L"] + 1
		return "LOSE"
	case answer == "PAPER" && move == "SCISSORS":
		game.Stats["W"] = game.Stats["W"] + 1
		return "WIN"
	case answer == "PAPER" && move == "ROCK":
		game.Stats["L"] = game.Stats["L"] + 1
		return "LOSE"
	}
	return ""
}

func (game *Game) printStats() string {
	stats := game.Stats
	wins := strconv.Itoa(stats["W"])
	draws := strconv.Itoa(stats["D"])
	losses := strconv.Itoa(stats["L"])
	result := "W" + wins + " D" + draws + " L" + losses
	return result
}

func (game *Game) generateMoveAnswer(move string) string {
	count := len(game.Answers)
	randomNumber := rand.Intn(count)
	answer := game.Answers[randomNumber]
	result := game.SaveStats(answer, move)
	return result + " " + answer
}

func (game *Game) GenerateResponse(move string) string {
	switch {
	case move == "STATS":
		return game.printStats()
	case game.IsGameMove(move):
		game.Mutex.Lock()
		defer game.Mutex.Unlock()
		return game.generateMoveAnswer(move)
	case move == "QUIT":
		return "Until then stranger"
	}
	return "Not sure what you mean"
}

func newGame() *Game {
	game := &Game{
		clients: make([]*Client, 0),
		joins:   make(chan net.Conn),
		Answers: []string{"ROCK", "PAPER", "SCISSORS"},
		Stats:   map[string]int{"W": 0, "D": 0, "L": 0},
		Mutex:   sync.RWMutex{},
	}
	game.Listen()
	return game
}

type Client struct {
	reader     *bufio.Reader
	connection net.Conn
	game       *Game
}

func (client *Client) Read() {
	for {
		message, _ := client.reader.ReadString('\n')
		if message != "" {
			message = strings.TrimSpace(message)
			fmt.Println("Just read:", message)
			response := client.game.GenerateResponse(message)
			fmt.Println(response, "response")
			client.connection.Write([]byte(response + "\n"))
		}
	}
}

func NewClient(connection net.Conn, game *Game) *Client {
	reader := bufio.NewReader(connection)

	client := &Client{
		reader:     reader,
		connection: connection,
		game:       game,
	}

	go client.Read()
	return client
}

func main() {

	game := newGame()

	listener, err := net.Listen("tcp", ":1983")
	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		connection, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		game.joins <- connection
	}

}
