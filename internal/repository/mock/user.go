package mock

import (
	"context"

	"github.com/eugenedhz/auth_service_test/internal/domain/models"
)

// Mock repo, not for testing, just to use :)
type UserRepository struct{}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

// Returns user by ID.
//
// In theory there might be ErrUserNotFound error :)
func (r *UserRepository) Get(ctx context.Context, userID string) (*models.User, error) {
	return &models.User{
		ID:    userID,
		Email: "example@mail.com",
	}, nil
}
