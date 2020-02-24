package store

import "context"

type DataStore interface {
	Connect(ctx context.Context) error
	Disconnect(ctx context.Context) error
}
