// Package client provides the Firestore client used by Firego.
package client

import "cloud.google.com/go/firestore"

// Client wraps the Google Cloud Firestore client used by Firego.
type Client struct {
	firestore *firestore.Client
}

// New creates a Firego client backed by client.
func New(client *firestore.Client) *Client {
	return &Client{
		firestore: client,
	}
}
