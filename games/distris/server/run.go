package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ServiceWeaver/weaver"
	"github.com/ubntc/go/games/distris/api/command"
)

func Run(address string) error {
	root := weaver.Init(context.Background())
	opts := weaver.ListenerOptions{LocalAddress: "localhost:12345"}
	lis, err := root.Listener("hello", opts)
	if err != nil {
		return err
	}
	fmt.Printf("game listener available on %v\n", lis)

	// Get routers as our main entrance to the weaver network
	// This will connect the HTTP world to the weaver insides.
	router, err := weaver.Get[Router](root)
	if err != nil {
		return err
	}

	// Serve the /game endpoint to receive commands via HTTP.
	http.HandleFunc("/game", func(w http.ResponseWriter, r *http.Request) {
		value := r.FormValue("command")
		result, err := router.Send(r.Context(), command.Command(value))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		fmt.Fprintf(w, "Send: %s, result: %v\n", value, result)
	})
	return http.Serve(lis, nil)
}
