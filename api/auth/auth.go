package auth

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/go-jose/go-jose/v3"
	"github.com/go-jose/go-jose/v3/jwt"
	"github.com/google/uuid"

	"github.com/Emyrk/strava/api/auth/authkeys"
)

type Options struct {
	Lifetime  time.Duration
	SecretPEM []byte
	Issuer    string
}

type Authentication struct {
	Lifetime time.Duration
	Signer   jose.Signer
	Issuer   string
}

func New(opts Options) (*Authentication, error) {
	secretKey, err := authkeys.ParsePrivateKey(opts.SecretPEM)
	if err != nil {
		return nil, fmt.Errorf("parse private key: %w", err)
	}

	// Instantiate a signer using RSASSA-PSS (SHA512) with the given private key.
	signer, err := jose.NewSigner(jose.SigningKey{Algorithm: jose.PS512, Key: secretKey}, nil)
	if err != nil {
		return nil, fmt.Errorf("create signer: %w", err)
	}

	return &Authentication{
		Lifetime: opts.Lifetime,
		Signer:   signer,
		Issuer:   opts.Issuer,
	}, nil
}

// ValidateSession returns the athlete ID if the session is valid
func (a *Authentication) ValidateSession(payload string) (int64, error) {
	token, err := jwt.ParseSigned(payload)
	if err != nil {
		return -1, fmt.Errorf("parse token: %w", err)
	}

	claims := jwt.Claims{}
	err = token.Claims(a.Signer, &claims)
	if err != nil {
		return -1, fmt.Errorf("parse claims: %w", err)
	}

	err = claims.Validate(jwt.Expected{
		Issuer: a.Issuer,
		Time:   time.Now(),
	})
	if err != nil {
		return -1, fmt.Errorf("validate claims: %w", err)
	}

	id, err := strconv.ParseInt(claims.Subject, 10, 64)
	if err != nil {
		return -1, fmt.Errorf("parse subject: %w", err)
	}
	return id, nil
}

func (a *Authentication) CreateSession(ctx context.Context, athleteID int64) (string, error) {
	c := &jwt.Claims{
		Issuer:    a.Issuer,
		Subject:   fmt.Sprintf("%d", athleteID),
		Audience:  []string{a.Issuer},
		Expiry:    jwt.NewNumericDate(time.Now().Add(a.Lifetime)),
		NotBefore: jwt.NewNumericDate(time.Now().Add(time.Minute * -1)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ID:        uuid.NewString(),
	}
	payload, err := jwt.Signed(a.Signer).Claims(c).CompactSerialize()
	if err != nil {
		return "", fmt.Errorf("sign session: %w", err)
	}

	return payload, nil
}
