package main

import (
	"log"
	"net/http"
	"golang.org/x/net/websocket"
	"encoding/json"
)

const server = "localhost:6969"

type JSONRequest struct {
	Msg string `json:"message"`
	Name string `json:"name"`
}

func init() {
	log.Println("Initialyzing")
}

func Echo(ws *websocket.Conn) {
	var reqJSON JSONRequest
	defer ws.Close()

	for {
		if err := websocket.JSON.Receive(ws, &reqJSON); err != nil {
			panic(err)
			resp := &JSONRequest{
				Msg: "Cannot parse request",
				Name: "Message error",
			}
			websocket.JSON.Send(ws, resp)
			return
		}
		out, _ := json.Marshal(reqJSON)
		log.Println(string(out))

		resp := &JSONRequest{
				Msg: "Message accepted",
				Name: "Success",
			}
		if err := websocket.JSON.Send(ws, resp); err != nil {
			log.Println(err.Error())
			return
		}
	}
}

func main() {
	log.Println("Server started on", server)
	http.Handle("/", websocket.Handler(Echo))

	if err := http.ListenAndServe(server, nil); err != nil {
		panic(err)
	}
}