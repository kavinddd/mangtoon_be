package token

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"time"
)

// Payload contains payload of the token
type Payload struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	Roles     []string  `json:"roles"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

func (p Payload) GetExpirationTime() (*jwt.NumericDate, error) {
	return jwt.NewNumericDate(p.ExpiredAt), nil
}

func (p Payload) GetIssuedAt() (*jwt.NumericDate, error) {
	return jwt.NewNumericDate(p.IssuedAt), nil
}

func (p Payload) GetNotBefore() (*jwt.NumericDate, error) {
	return nil, nil
}

func (p Payload) GetIssuer() (string, error) {
	return "", fmt.Errorf("get issuer")
}

func (p Payload) GetSubject() (string, error) {
	return "", fmt.Errorf("get subject")
}

func (p Payload) GetAudience() (jwt.ClaimStrings, error) {

	return nil, fmt.Errorf("get audience")
}

// NewPayload creates a new token payload
func NewPayload(username string, roles []string, duration time.Duration) (*Payload, error) {
	tokenId, err := uuid.NewRandom()
	if err != nil {
		return nil, fmt.Errorf("cannot generate uuid")
	}

	payload := Payload{
		ID:        tokenId,
		Username:  username,
		Roles:     roles,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}

	return &payload, nil

}
