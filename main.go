package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

// Opens the websocket and shows received messages and when the server closes the socket
const homeHTML = `<!DOCTYPE html>
<html lang="en">
    <head>
        <title>WebSocket Example</title>
    </head>
    <body>
        <pre id="fileData"></pre>
        <script type="text/javascript">
            (function() {
                var data = document.getElementById("fileData");
                var conn = new WebSocket("ws://127.0.0.1:8080/ws");
                conn.onclose = function(evt) {
                    data.textContent = data.textContent + 'Connection closed';
                }
                conn.onmessage = function(evt) {
                    data.textContent = data.textContent + evt.data;
                }
            })();
        </script>
    </body>
</html>
`

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

// Upgrades the http request to a websocket, sends 25 json messages, then close the websocket
func wsEndpoint(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	data := struct {
		Name  string `json:"name,omitempty"`
		Index int    `json:"index,omitempty"`
	}{
		Name: "dan",
	}

	for data.Index = 1; data.Index <= 25; data.Index++ {
		if err := ws.WriteJSON(&data); err != nil {
			log.Println(err)
			return
		}
		time.Sleep(100 * time.Millisecond)
	}

	ws.Close()
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { fmt.Fprint(w, homeHTML) })
	http.HandleFunc("/ws", wsEndpoint)
	fmt.Println("Listening on 8080, open http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
