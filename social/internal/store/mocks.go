package store

import (
	"context"
	"database/sql"
	"time"
)

//**************************************************
//File for creating a mock storage for testing
//***************************************************
// To Use, add a storage for testing posts, users, etc., something that needs to be in local storage
// make a struct interface containing all methods for the feature impl.
// return nil and/or customize the methods for specific edge cases or scenarios
//***************************************************

//impl
//		Create(context.Context, *sql.Tx, *User) error
//		GetByID(context.Context, int64) (*User, error)
//		GetByEmail(context.Context, string) (*User, error)
//		CreateAndInvite(ctx context.Context, user *User, token string, exp time.Duration) error
//		Activate(context.Context, string) error
//		Delete(context.Context, int64) error//

func NewMockStore() Storage {
	return Storage{
		Users: &MockUserStore{},
	}
}

type MockUserStore struct{}

func (m *MockUserStore) Create(ctx context.Context, tx *sql.Tx, u *User) error {
	return nil
}
func (m *MockUserStore) GetByID(context.Context, int64) (*User, error) {
	//add user in brackets to get by id EX-> ID: 1
	return &User{}, nil
}
func (m *MockUserStore) GetByEmail(context.Context, string) (*User, error) {
	//add user in brackets to get by email EX-> email: email@email.com
	return &User{}, nil
}
func (m *MockUserStore) CreateAndInvite(ctx context.Context, user *User, token string, exp time.Duration) error {
	return nil
}
func (m *MockUserStore) Activate(ctx context.Context, t string) error {
	return nil
}
func (m *MockUserStore) Delete(context.Context, int64) error {
	return nil
}
