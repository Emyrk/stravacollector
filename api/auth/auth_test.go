package auth_test

import (
	"context"
	"testing"

	"github.com/Emyrk/strava/api/auth/authkeys"

	"github.com/Emyrk/strava/api/auth"
	"github.com/stretchr/testify/require"
)

func TestSignValidate(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	key, err := authkeys.GenerateKey()
	require.NoError(t, err, "generate key")

	a, err := auth.New(auth.Options{
		SecretPEM: authkeys.MarshalPrivateKey(key),
		Issuer:    "Test",
	})
	require.NoError(t, err, "create auth")

	const athID = int64(1)
	token, err := a.CreateSession(ctx, athID)
	require.NoError(t, err, "create session")

	ath, err := a.ValidateSession(token)
	require.NoError(t, err, "validate session")

	require.Equal(t, athID, ath, "athlete ID")
}
