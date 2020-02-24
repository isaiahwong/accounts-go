package user

import (
	"context"

	"github.com/isaiahwong/auth-go/internal/models"
)

// Repo defines user repository operations
type Repo interface {
	Save(context.Context, *models.User) string
	Find(context.Context, string) []*models.User
	FindOne(context.Context) *models.User
}
