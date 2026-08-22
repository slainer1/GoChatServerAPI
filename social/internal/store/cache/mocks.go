package cache

import (
	"context"
	"main/internal/store"

	"github.com/stretchr/testify/mock"
)

// **************************************************
// File for creating a mock cache storage for testing
// ***************************************************
func NewMockStore() Storage {
	return Storage{
		Users: &MockUserStore{},
	}
}

type MockUserStore struct {
	mock.Mock
}

func (m *MockUserStore) Get(ctx context.Context, userID int64) (*store.User, error) {
	args := m.Called(ctx, userID)
	return nil, args.Error(1)
}
func (m *MockUserStore) Set(ctx context.Context, user *store.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}
