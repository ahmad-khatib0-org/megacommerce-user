package controller

import (
	"context"
	"strings"
	"testing"

	userPb "github.com/ahmad-khatib0-org/megacommerce-proto/gen/go/user/v1"
	"github.com/ahmad-khatib0-org/megacommerce-user/pkg/models"
	"github.com/ahmad-khatib0-org/megacommerce-user/pkg/utils"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func (th *TestHelper) getValidSignupSupplierRequest(t *testing.T) (*userPb.User, *userPb.SupplierCreateRequest) {
	t.Helper()
	validReq := &userPb.SupplierCreateRequest{
		Username:   "username",
		FirstName:  "firstname",
		LastName:   "lastname",
		Email:      "test@gmail.com",
		Password:   "pass1a$A",
		Membership: "free",
	}

	validSupplier := &userPb.User{
		Username:   utils.NewPointer(validReq.GetUsername()),
		FirstName:  utils.NewPointer(validReq.GetFirstName()),
		LastName:   utils.NewPointer(validReq.GetLastName()),
		Email:      utils.NewPointer(validReq.GetEmail()),
		Password:   utils.NewPointer(validReq.GetPassword()),
		Membership: utils.NewPointer(validReq.GetMembership()),
	}

	validUser, err := models.SignupSupplierRequestPreSave(th.Supplier1.Ctx, validSupplier)
	require.Nil(t, err)
	return validUser, validReq
}

func TestSignupSupplier(t *testing.T) {
	th, err := NewTestHelper(t)
	require.Nil(t, err)
	defer th.TearDown()

	mockDependencies := func(s *userPb.SupplierCreateRequest) {
		th.store.On("SignupSupplier", mock.AnythingOfType("*models.Context"), mock.MatchedBy(func(u *userPb.User) bool {
			return u.GetEmail() == s.GetEmail() &&
				u.GetUsername() == s.GetUsername()
		})).Return(nil)

		th.tasker.On("SendVerifyEmail", mock.Anything, utils.WithMockDebug("models.TaskSendVerifyEmailPayload", func(p *models.TaskSendVerifyEmailPayload) bool {
			return p.Email == s.GetEmail()
		}), mock.Anything).Return(nil)
	}

	t.Run("signup supplier successfully!", func(t *testing.T) {
		ctx := th.withContext(context.Background())

		_, supplier := th.getValidSignupSupplierRequest(t)
		mockDependencies(supplier)

		_, err := th.controller.CreateSupplier(ctx, supplier)
		require.Nil(t, err)
	})

	tests := map[string]struct {
		input   func() *userPb.SupplierCreateRequest
		expects string
	}{
		"invalid username": {
			input: func() *userPb.SupplierCreateRequest {
				_, sup := th.getValidSignupSupplierRequest(t)
				sup.Username = "##.n"
				return sup
			},
			expects: "user.create.username.valid.error",
		},
		"short username": {
			input: func() *userPb.SupplierCreateRequest {
				_, sup := th.getValidSignupSupplierRequest(t)
				sup.Username = "n"
				return sup
			},
			expects: "user.create.username.error",
		},
		"long username": {
			input: func() *userPb.SupplierCreateRequest {
				_, sup := th.getValidSignupSupplierRequest(t)
				sup.Username = strings.Repeat("s", models.UserNameMaxLength+1)
				return sup
			},
			expects: "user.create.username.error",
		},
		"missing first name": {
			input: func() *userPb.SupplierCreateRequest {
				_, sup := th.getValidSignupSupplierRequest(t)
				sup.FirstName = ""
				return sup
			},
			expects: "user.create.first_name.error",
		},
		"missing last name": {
			input: func() *userPb.SupplierCreateRequest {
				_, sup := th.getValidSignupSupplierRequest(t)
				sup.LastName = ""
				return sup
			},
			expects: "user.create.last_name.error",
		},
		"short password": {
			input: func() *userPb.SupplierCreateRequest {
				_, sup := th.getValidSignupSupplierRequest(t)
				sup.Password = "pass"
				return sup
			},
			expects: "password.min_length",
		},
		"password missing a lowercase character": {
			input: func() *userPb.SupplierCreateRequest {
				_, sup := th.getValidSignupSupplierRequest(t)
				sup.Password = "PASS$44E"
				return sup
			},
			expects: "password.lowercase",
		},
		"password missing an uppercase character": {
			input: func() *userPb.SupplierCreateRequest {
				_, sup := th.getValidSignupSupplierRequest(t)
				sup.Password = "pass$44e"
				return sup
			},
			expects: "password.uppercase",
		},
		"password missing a number character": {
			input: func() *userPb.SupplierCreateRequest {
				_, sup := th.getValidSignupSupplierRequest(t)
				sup.Password = "Pass$word"
				return sup
			},
			expects: "password.numbers",
		},
		"password missing a symbol character": {
			input: func() *userPb.SupplierCreateRequest {
				_, sup := th.getValidSignupSupplierRequest(t)
				sup.Password = "Pass4444"
				return sup
			},
			expects: "password.symbols",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ctx := th.withContext(context.Background())

			sup := tc.input()
			res, err := th.controller.CreateSupplier(ctx, sup)
			require.Nil(t, err)
			require.NotNil(t, res.Response)

			errRes, ok := res.Response.(*userPb.SupplierCreateResponse_Error)
			require.True(t, ok)
			require.Equal(t, tc.expects, errRes.Error.Id)
		})
	}
}
