package token

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

const minSecretKeyLength = 32

//types of error returned by VerifyToken function

var (
	ErrInvalidToken = errors.New("token is invalid")
)

type JwtMaker struct {
	secretKey string
}

func (maker *JwtMaker) CreateToken(username string, roles []string, duration time.Duration) (string, *Payload, error) {
	payload, err := NewPayload(username, roles, duration)
	if err != nil {
		return "", payload, err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	signedJwt, err := token.SignedString([]byte(maker.secretKey))

	if err != nil {
		return "", payload, fmt.Errorf("cannot sign jwt: %w", err)
	}

	return signedJwt, payload, nil
}

func (maker *JwtMaker) VerifyToken(signedJwt string) (*Payload, error) {
	// is a function received by ParseWithClaims that will get them secret key to parse the token
	// we can put some validation here before handing them the secret key
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		//_, ok := token.Method.(*jwt.SigningMethodHMAC)
		//if !ok {
		//	return nil, ErrInvalidToken
		//}
		return []byte(maker.secretKey), nil
	}
	token, err := jwt.ParseWithClaims(signedJwt, &Payload{}, keyFunc,
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}),
		jwt.WithExpirationRequired(),
		jwt.WithIssuedAt(),
	)

	if err != nil {
		return nil, err
	}

	payload, ok := token.Claims.(*Payload)
	if !ok {
		return nil, ErrInvalidToken
	}
	return payload, nil
}

func NewJwtMaker(secretKey string) (Maker, error) {
	if len(secretKey) < minSecretKeyLength {
		return nil, fmt.Errorf("invalid key size, must be at least %d characters", minSecretKeyLength)
	}

	return &JwtMaker{secretKey: secretKey}, nil
}
