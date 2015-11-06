package main

import (
	"log"
	"net/http"
	"golang.org/x/net/websocket"
	"fmt"
	"os"
)

const server = ":6969"

var active = make(map[string]*websocket.Conn)

type Peer struct  {
	client *websocket.Conn
	name string
}

type Message struct {
	Text string `json:"message"`
	Name string `json:"name"`
}

type Online struct {
	Users []string `json:"users"`
}

func Echo(ws *websocket.Conn) {
	var message Message
	defer ws.Close()

	for {
		if err := websocket.JSON.Receive(ws, &message); err != nil {
			Log(err)
			return
		}

		active[message.Name] = ws

		if message.Text == "Register" {
			if err := websocket.JSON.Send(ws, online()); err != nil {
				Log(err)
				delete(active, message.Name)
				return
			}
			continue
		}

		go sendAll(message)

	}
}

func sendAll(message *Message) {
	resp := &Message{
				Text: message.Text,
				Name: message.Name,
			}
			for n, v := range(active) {
				if err := websocket.JSON.Send(v, resp); err != nil {
					Log(err)
					delete(active, n)
				}
			}
}

func online() *Online {
	var names []string
	for k, _ := range active {
		names = append(names, k)
	}

	return &Online{
		Users:names,
	}
}

func Log(err error) {
	fmt.Fprint(os.Stdout, err.Error() + "\n")
}

func main() {
	log.Println("Server started on", server)
	http.Handle("/", websocket.Handler(Echo))
	if err := http.ListenAndServe(server, nil); err != nil {
		panic(err)
	}
}
