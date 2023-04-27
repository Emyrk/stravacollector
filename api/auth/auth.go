package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/Emyrk/strava/database"
	"github.com/google/uuid"
)

type Authentication struct {
	DB database.Store
}

func New(db database.Store) *Authentication {
	return &Authentication{
		DB: db,
	}
}

func (a *Authentication) ValidateSession(ctx context.Context, now time.Time, id uuid.UUID, secret string) (*database.ApiToken, error) {
	secretBytes, err := hex.DecodeString(secret)
	if err != nil {
		return nil, fmt.Errorf("decode secret from hex: %w", err)
	}

	token, err := a.DB.GetToken(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get api token: %w", err)
	}

	hashed := HashSecret(id, secretBytes)

	return nil
}

func (a *Authentication) CreateSession(ctx context.Context, tokenName string, athlete *database.AthleteLogin) (string, database.ApiToken, error) {
	id := uuid.New()
	secret, token, err := GenerateToken(id)
	if err != nil {
		return "", database.ApiToken{}, fmt.Errorf("create token: %w", err)
	}
	lifetime := time.Hour * 24 * 7
	time.Now().Add(lifetime)
	apiToken, err := a.DB.InsertAPIToken(ctx,
		database.InsertAPITokenParams{
			ID:              id,
			Name:            tokenName,
			AthleteID:       athlete.AthleteID,
			HashedToken:     token,
			ExpiresAt:       time.Now().Add(lifetime),
			LifetimeSeconds: int64(lifetime.Seconds()),
		})
	if err != nil {
		return "", apiToken, err
	}

	return secret, apiToken, nil
}

// GenerateToken uses id as a salt.
func GenerateToken(id uuid.UUID) (string, string, error) {
	secret := make([]byte, 32)
	_, err := rand.Read(secret)
	if err != nil {
		return "", "", err
	}
	hashed := HashSecret(id, secret)
	return hex.EncodeToString(secret), hex.EncodeToString(hashed[:]), nil
}

func HashSecret(id uuid.UUID, secret []byte) []byte {
	hashed := sha256.Sum256(append([]byte(id.String()), secret...))
	return hashed[:]
}

func TokenString(id uuid.UUID, hashed string) string {
	return fmt.Sprintf("%s:%s", id.String(), hashed)
}
