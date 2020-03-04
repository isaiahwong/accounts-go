package mail

import (
	"github.com/isaiahwong/accounts-go/api/client"
)

// NewMailClient returns a new MailServiceClient
func NewMailClient(opts ...client.Option) (MailServiceClient, error) {
	conn, err := client.CreateClient(opts...)
	if err != nil {
		return nil, err
	}
	client := NewMailServiceClient(conn)
	return client, nil
}
