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

type Message struct {
	Text string `json:"message"`
	Name string `json:"name"`
}

type Online struct {
	Users []string `json:"users"`
}

func Echo(ws *websocket.Conn) {
	var reqJSON Message
	defer ws.Close()

	for {
		if err := websocket.JSON.Receive(ws, &reqJSON); err != nil {
			Log(err)
			return
		}

		active[reqJSON.Name] = ws

		if reqJSON.Text == "Register" {
			if err := websocket.JSON.Send(ws, online()); err != nil {
				Log(err)
				return
			}
			continue
		}

		resp := &Message{
			Text: reqJSON.Text,
			Name: reqJSON.Name,
		}
		go func() {
			for n, v := range(active) {
			if err := websocket.JSON.Send(v, resp); err != nil {
				Log(err)
				delete(active, n)
			}
		}
		}()

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
