package accounts

import (
	"context"

	"github.com/isaiahwong/accounts-go/internal/models"
)

// Repo defines accounts repository operations
type Repo interface {
	Update(c context.Context, filter interface{}, update interface{}) (int, error)
	Save(c context.Context, u *models.Account) (string, error)
	Find(c context.Context, f interface{}, opts ...interface{}) ([]*models.Account, error)
	FindOne(c context.Context, f interface{}, opts ...interface{}) (*models.Account, error)
}
