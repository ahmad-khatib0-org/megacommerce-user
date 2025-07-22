package controller

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	pb "github.com/ahmad-khatib0-org/megacommerce-proto/gen/go/users/v1"
	"github.com/ahmad-khatib0-org/megacommerce-user/internal/store"
	"github.com/ahmad-khatib0-org/megacommerce-user/pkg/models"
	"github.com/ahmad-khatib0-org/megacommerce-user/pkg/utils"
)

// Builds a mock token with optional overrides
func getToken(t *testing.T, req *pb.EmailConfirmationRequest, opts ...func(token *pb.Token)) *pb.Token {
	t.Helper()

	pass, err := utils.PasswordHash(req.Token)
	require.Nil(t, err)
	token := &pb.Token{
		Id:        req.TokenId,
		Token:     pass,
		Type:      string(models.TokenTypeEmailConfirmation),
		Used:      false,
		CreatedAt: utils.TimeGetMillis(),
		ExpiresAt: time.Now().Add(1 * time.Hour).UnixMilli(),
	}

	for _, opt := range opts {
		opt(token)
	}

	return token
}

// Generates a valid EmailConfirmationRequest with a fresh token
func getValidEmailConfirmationRequest(t *testing.T, hours uint) *pb.EmailConfirmationRequest {
	t.Helper()

	token := &utils.Token{}
	expiry := time.Duration(hours) * time.Hour

	tokenData, err := token.GenerateToken(expiry)
	require.NoError(t, err)

	expectedExpiry := time.Now().Add(expiry)
	require.WithinDuration(t, expectedExpiry, tokenData.Expiry, time.Second)

	return &pb.EmailConfirmationRequest{
		Email:   "test@email.com",
		Token:   tokenData.Token,
		TokenId: tokenData.Id,
	}
}

func TestEmailConfirmation(t *testing.T) {
	th, err := NewTestHelper(t)
	fmt.Println(err)
	require.Nil(t, err)
	defer th.TearDown()

	hours := uint(th.config.Security.GetTokenConfirmationExpiryInHours())
	// req := getValidEmailConfirmationRequest(t, hours)

	t.Run("invalid email", func(t *testing.T) {
		req := getValidEmailConfirmationRequest(t, hours)
		req.Email = "invalid@@"

		res, err := th.controller.EmailConfirmation(th.withContext(context.Background()), req)
		require.NoError(t, err)

		errRes, ok := res.Response.(*pb.EmailConfirmationResponse_Error)
		require.True(t, ok)
		require.Equal(t, "email_confirm.email.error", errRes.Error.Id)
	})

	t.Run("expired token", func(t *testing.T) {
		req := getValidEmailConfirmationRequest(t, hours)
		expired := getToken(t, req, func(token *pb.Token) {
			token.ExpiresAt = time.Now().Add(-1 * time.Hour).UnixMilli()
		})

		th.store.On("TokensGet", mock.Anything, req.TokenId).Return(expired, nil)

		res, err := th.controller.EmailConfirmation(th.withContext(context.Background()), req)
		require.NoError(t, err)

		errRes, ok := res.Response.(*pb.EmailConfirmationResponse_Error)
		require.True(t, ok)
		require.Equal(t, "email_confirm.token.expired", errRes.Error.Id)
	})

	t.Run("token already used", func(t *testing.T) {
		req := getValidEmailConfirmationRequest(t, hours)
		used := getToken(t, req, func(token *pb.Token) {
			token.Used = true
		})

		th.store.On("TokensGet", mock.Anything, req.TokenId).Return(used, nil)

		res, err := th.controller.EmailConfirmation(th.withContext(context.Background()), req)
		require.NoError(t, err)

		successRes, ok := res.Response.(*pb.EmailConfirmationResponse_Data)
		require.True(t, ok)
		require.Contains(t, successRes.Data.GetMessage(), "already confirmed")
	})

	t.Run("wrong token", func(t *testing.T) {
		req := getValidEmailConfirmationRequest(t, hours)
		wrong := getToken(t, req, func(token *pb.Token) {
			pass, err := utils.PasswordHash("wrong-token")
			require.NoError(t, err)
			token.Token = pass
		})

		th.store.On("TokensGet", mock.Anything, req.TokenId).Return(wrong, nil)

		res, err := th.controller.EmailConfirmation(th.withContext(context.Background()), req)
		require.NoError(t, err)

		successRes, ok := res.Response.(*pb.EmailConfirmationResponse_Data)
		require.True(t, ok)
		require.Contains(t, successRes.Data.GetMessage(), "invalid")
	})

	t.Run("success", func(t *testing.T) {
		req := getValidEmailConfirmationRequest(t, hours)
		valid := getToken(t, req)

		th.store.On("TokensGet", mock.Anything, req.TokenId).Return(valid, nil)
		th.store.On("MarkEmailAsConfirmed", mock.Anything, req.TokenId).Return(nil)

		res, err := th.controller.EmailConfirmation(th.withContext(context.Background()), req)
		require.NoError(t, err)

		successRes, ok := res.Response.(*pb.EmailConfirmationResponse_Data)
		require.True(t, ok)
		require.Contains(t, successRes.Data.GetMessage(), "confirmed")
	})

	t.Run("db error on MarkEmailAsConfirmed", func(t *testing.T) {
		req := getValidEmailConfirmationRequest(t, hours)
		valid := getToken(t, req)

		th.store.On("TokensGet", mock.Anything, req.TokenId).Return(valid, nil)
		th.store.On("MarkEmailAsConfirmed", mock.Anything, req.TokenId).
			Return(&store.DBError{
				Path: "users.store.MarkEmailAsConfirmed",
				Msg:  "db failed",
			})

		res, err := th.controller.EmailConfirmation(th.withContext(context.Background()), req)
		require.NoError(t, err)

		errRes, ok := res.Response.(*pb.EmailConfirmationResponse_Error)
		require.True(t, ok)
		require.Equal(t, models.ErrMsgInternal, errRes.Error.Id)
	})
}

