package client

import (
	"net/http"
	"time"
)

type NUClient struct {
	Client *http.Client
}

func NewNUClient() *NUClient {
	return &NUClient{
		Client: &http.Client{
			Timeout: 25 * time.Second,
		},
	}
}
