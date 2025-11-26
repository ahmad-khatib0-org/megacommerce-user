package models

import (
	"testing"
	"time"

	v1 "github.com/ahmad-khatib0-org/megacommerce-proto/gen/go/users/v1"
	"github.com/ahmad-khatib0-org/megacommerce-shared-go/pkg/models"
	"github.com/ahmad-khatib0-org/megacommerce-shared-go/pkg/utils"
	"github.com/stretchr/testify/require"
)

func getValidEmailConfirmationRequet(t *testing.T, hours uint) *v1.EmailConfirmationRequest {
	t.Helper()
	token := &utils.Token{}

	expiry := time.Duration(hours) * time.Hour
	tokenData, err := token.GenerateToken(expiry)
	require.Nil(t, err)
	require.NotNil(t, tokenData)

	expectedExpiry := time.Now().Add(expiry)

	// Compare with tolerance (1 second to account for execution time)
	require.WithinDuration(t, expectedExpiry, tokenData.Expiry, time.Second)

	return &v1.EmailConfirmationRequest{Email: "test@email.com", Token: tokenData.Token, TokenId: tokenData.ID}
}

func TestEmailConfirmationIsValid(t *testing.T) {
	hours := 24
	ctx := models.ContextForTesting()

	t.Run("invalid email ", func(t *testing.T) {
		req := getValidEmailConfirmationRequet(t, uint(hours))
		req.Email = "invalid@@"
		err := EmailConfirmationIsValid(ctx, req)
		require.NotNil(t, err)
		require.Equal(t, "email_confirm.email.error", err.ID)
	})

	t.Run("missing token", func(t *testing.T) {
		req := getValidEmailConfirmationRequet(t, uint(hours))
		req.Token = ""
		err := EmailConfirmationIsValid(ctx, req)
		require.NotNil(t, err)
		require.Equal(t, "email_confirm.token.error", err.ID)
	})

	t.Run("invalid token_id", func(t *testing.T) {
		req := getValidEmailConfirmationRequet(t, uint(hours))
		req.TokenId = "invalid@@"
		err := EmailConfirmationIsValid(ctx, req)
		require.NotNil(t, err)
		require.Equal(t, "email_confirm.token_id.error", err.ID)
	})
}
