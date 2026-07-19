// Package client provides the Firestore client used by Firego.
package client

import "cloud.google.com/go/firestore"

type Client struct {
	firestore *firestore.Client
}

func New(client *firestore.Client) *Client {
	return &Client{
		firestore: client,
	}
}
