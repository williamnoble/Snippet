package models

import (
	"errors"
	"time"
)

var ErrNoRecord = errors.New("models: no matching record found")
var ErrDuplicateEmail = errors.New("models: duplicate email")
var ErrInvalidCredentials = errors.New("invalid crendentials suppleid")

// Snippet is a Model representing a single Snippet
type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

// User represents a Single User
type User struct {
	ID             int
	Name           string
	Email          string
	HashedPassword []byte
	Created        time.Time
	Active         bool
}
