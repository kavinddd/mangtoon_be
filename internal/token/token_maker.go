package token

import "time"

// Maker is an interface for managing tokens
type Maker interface {
	CreateToken(username string, roles []string, duration time.Duration) (string, *Payload, error)
	// VerifyToken checks if the token is valid or not, if so it returns Payload, otherwise, error
	VerifyToken(token string) (*Payload, error)
}
