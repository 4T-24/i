package ctf

import (
	"instancer/internal/env"

	ctfd "github.com/ctfer-io/go-ctfd/api"
	"github.com/sirupsen/logrus"
)

type Client struct {
	*ctfd.Client

	Queue chan *ctfd.Challenge
}

func New() *Client {
	c := env.Get()

	nonce, session, err := ctfd.GetNonceAndSession(c.CTFd.URL)
	if err != nil {
		logrus.Fatalf("Failed getting nonce and session: %s", err)
	}

	client := ctfd.NewClient(c.CTFd.URL, nonce, session, c.CTFd.Token)
	return &Client{
		Client: client,
		Queue:  make(chan *ctfd.Challenge),
	}
}
