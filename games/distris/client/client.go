package client

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/ubntc/go/games/distris/api/command"
)

type Client struct {
	address string
	gameUrl string
	api     *http.Client
	counter uint64
}

func New(address string) *Client {
	gameUrl, err := url.JoinPath("http://", address, "game")
	if err != nil {
		log.Fatalln(err)
	}

	return &Client{
		address: address,
		gameUrl: gameUrl,
		api:     http.DefaultClient,
	}
}

func (c *Client) Send(ctx context.Context, cmd command.Command) error {
	c.counter++
	values := url.Values{}
	values.Add("command", url.QueryEscape(string(cmd)))
	res, err := c.api.PostForm(c.gameUrl, values)
	if err != nil {
		return err
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	text := strings.TrimSpace(string(body))
	fmt.Printf("response: %v (#%d)\n\r", text, c.counter)
	return nil
}
