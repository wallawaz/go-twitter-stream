package main

import "net/http"

func runWebSocket(hub *Hub, wsAddr string) {
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})
	http.ListenAndServe(wsAddr, nil)
}
