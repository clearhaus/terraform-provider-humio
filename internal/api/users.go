package api

import (
	"context"
)

// User represents a Humio user
type User struct {
	ID       string
	Username string
	FullName string
	Email    string
	IsRoot   bool
}

// Users provides operations for managing users
type Users struct {
	client *Client
}

const currentUserQuery = `
query CurrentUser {
  currentUser {
    id
    username
    fullName
    email
    isRoot
  }
}
`

// currentUserResponse represents the response from current user query
type currentUserResponse struct {
	CurrentUser struct {
		ID       string `json:"id"`
		Username string `json:"username"`
		FullName string `json:"fullName"`
		Email    string `json:"email"`
		IsRoot   bool   `json:"isRoot"`
	} `json:"currentUser"`
}

// GetCurrent returns the current authenticated user
func (u *Users) GetCurrent() (*User, error) {
	var resp currentUserResponse
	err := u.client.Query(context.Background(), currentUserQuery, nil, &resp)
	if err != nil {
		return nil, err
	}

	return &User{
		ID:       resp.CurrentUser.ID,
		Username: resp.CurrentUser.Username,
		FullName: resp.CurrentUser.FullName,
		Email:    resp.CurrentUser.Email,
		IsRoot:   resp.CurrentUser.IsRoot,
	}, nil
}
