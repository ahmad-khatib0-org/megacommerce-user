package controller

import (
	"bytes"
	ctxPkg "context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	pbSh "github.com/ahmad-khatib0-org/megacommerce-proto/gen/go/shared/v1"
	pb "github.com/ahmad-khatib0-org/megacommerce-proto/gen/go/users/v1"
	"github.com/ahmad-khatib0-org/megacommerce-user/internal/store"
	"github.com/ahmad-khatib0-org/megacommerce-user/pkg/models"
	"github.com/ahmad-khatib0-org/megacommerce-user/pkg/utils"
	"google.golang.org/grpc/codes"
)

func (c *Controller) Login(context ctxPkg.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	path := "users.controller.Login"
	errBuilder := func(e *models.AppError) (*pb.LoginResponse, error) {
		return &pb.LoginResponse{Response: &pb.LoginResponse_Error{Error: models.AppErrorToProto(e)}}, nil
	}
	ctx, err := models.ContextGet(context)
	if err != nil {
		return errBuilder(err)
	}
	internalErr := func(err error, details string) (*pb.LoginResponse, error) {
		return errBuilder(models.NewAppError(ctx, path, models.ErrMsgInternal, nil, details, int(codes.Internal), &models.AppErrorErrorsArgs{Err: err}))
	}
	sucBuilder := func(data *pbSh.SuccessResponseData) (*pb.LoginResponse, error) {
		return &pb.LoginResponse{Response: &pb.LoginResponse_Data{Data: data}}, nil
	}

	rctx, cancel := ctxPkg.WithTimeout(ctxPkg.Background(), time.Second*12)
	defer cancel()
	ctx.Context = rctx

	ar := models.AuditRecordNew(ctx, models.EventNameLogin, models.EventStatusFail)
	defer c.ProcessAudit(ar)

	if err = models.LoginRequestIsValid(ctx, req); err != nil {
		return errBuilder(err)
	}

	user, errDB := c.store.UsersGetByEmail(ctx, req.GetEmail())
	if errDB != nil {
		if errDB.ErrType == store.DBErrorTypeNoRows {
			return errBuilder(models.NewAppError(ctx, path, "email.not_found", nil, errDB.Details, int(codes.NotFound), nil))
		} else {
			return internalErr(errDB, errDB.Details)
		}
	}

	if user.GetAuthService() != "" {
		return errBuilder(models.NewAppError(ctx, path, "user.login.use_auth_service.error", map[string]any{"AuthService": user.GetAuthService()}, "", int(codes.InvalidArgument), nil))
	}
	if err := utils.PasswordCheck(user.GetPassword(), req.GetPassword()); err != nil {
		return errBuilder(models.NewAppError(ctx, path, "user.login.password.error", nil, "", int(codes.InvalidArgument), nil))
	}

	// TODO: handle if this user is using mobile or not
	expiry := c.config().Security.GetAccessTokenExpiryWebInHours()
	body := map[string]any{
		"subject":      user.GetEmail(),
		"remember":     true,
		"remember_for": expiry * 60 * 60,
	}

	oauthPayload, marErr := json.Marshal(body)
	if marErr != nil {
		return internalErr(marErr, "failed to marshal json payload")
	}

	reqURL := fmt.Sprintf("%s/oauth2/auth/requests/login/accept?login_challenge=%s", c.config().Oauth.GetOauthAdminUrl(), req.GetLoginChallenge())
	oauthReq, reqErr := http.NewRequestWithContext(ctx.Context, http.MethodPut, reqURL, bytes.NewReader(oauthPayload))
	if reqErr != nil {
		return internalErr(reqErr, "failed to build an HTTP request to send to OAuth service")
	}
	oauthReq.Header.Set("Content-Type", "application/json")

	// TODO: inclue the response time as metrics
	start := time.Now()
	resp, respErr := utils.HTTPRequestWithRetry(c.httpClient, oauthReq, 3)
	duration := time.Since(start)
	if respErr != nil {
		c.log.Errorf("HTTP %s %s failed: %v (took %s)", oauthReq.Method, oauthReq.URL, err, duration)
		return internalErr(respErr, "failed to request OAuth server to accept login")
	}
	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusInternalServerError {
		return internalErr(respErr, "failed to request OAuth server to accept login")
	}

	if resp.StatusCode != http.StatusOK {
		var resErr models.OAuthErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&resErr); err != nil {
			return internalErr(err, "failed to unmarshal json payload from Oauth service error")
		}
		errorsInt := map[string]*models.AppErrorError{
			"error":       {ID: "login.error"},
			"description": {ID: models.GetOAuthRequestErrMsgID(ctx.AcceptLanguage, resErr.Error, resErr.ErrorDescription)},
		}
		errors := &models.AppErrorErrorsArgs{Err: respErr, ErrorsInternal: errorsInt}
		return errBuilder(models.NewAppError(ctx, path, "login.error", nil, "", int(codes.InvalidArgument), errors))
	}

	var result struct {
		RedirectTo string `json:"redirect_to"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return internalErr(err, "failed to unmarshal json payload from Oauth service response")
	}

	ar.Success()
	return sucBuilder(&pbSh.SuccessResponseData{})
}
