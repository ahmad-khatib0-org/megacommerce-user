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
	"github.com/ahmad-khatib0-org/megacommerce-shared-go/pkg/models"
	"github.com/ahmad-khatib0-org/megacommerce-shared-go/pkg/utils"
	intModels "github.com/ahmad-khatib0-org/megacommerce-user/pkg/models"
	"google.golang.org/grpc/codes"
)

func (c *Controller) Login(context ctxPkg.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	startTime := time.Now()
	path := "users.controller.Login"
	errBuilder := func(e *models.AppError) (*pb.LoginResponse, error) {
		return &pb.LoginResponse{Response: &pb.LoginResponse_Error{Error: models.AppErrorToProto(e)}}, nil
	}

	ctx, ctxErr := models.ContextGet(context)
	if ctxErr != nil {
		duration := time.Since(startTime).Seconds()
		c.metricsCollector.RecordLoginRequest(false, duration)
		return errBuilder(ctxErr)
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

	ar := models.AuditRecordNew(ctx, intModels.EventNameLogin, models.EventStatusFail)
	defer c.ProcessAudit(ar)

	if err := intModels.LoginRequestIsValid(ctx, req); err != nil {
		duration := time.Since(startTime).Seconds()
		c.metricsCollector.RecordLoginRequest(false, duration)
		return errBuilder(err)
	}

	user, err := c.store.UsersGetByEmail(ctx, req.GetEmail())
	if err != nil {
		duration := time.Since(startTime).Seconds()
		c.metricsCollector.RecordLoginRequest(false, duration)
		if err.ErrType == models.DBErrorTypeNoRows {
			errors := &models.AppErrorErrorsArgs{Err: err, ErrorsInternal: map[string]*models.AppErrorError{"email": {ID: "email.not_found"}}}
			return errBuilder(models.NewAppError(ctx, path, "email.not_found", nil, err.Details, int(codes.NotFound), errors))
		} else {
			return internalErr(err, err.Details)
		}
	}

	if user.GetAuthService() != "" {
		duration := time.Since(startTime).Seconds()
		c.metricsCollector.RecordLoginRequest(false, duration)
		return errBuilder(models.NewAppError(ctx, path, "user.login.use_auth_service.error", map[string]any{"AuthService": user.GetAuthService()}, "", int(codes.InvalidArgument), nil))
	}
	if err := utils.PasswordCheck(user.GetPassword(), req.GetPassword()); err != nil {
		duration := time.Since(startTime).Seconds()
		c.metricsCollector.RecordLoginRequest(false, duration)
		errors := &models.AppErrorErrorsArgs{Err: err, ErrorsInternal: map[string]*models.AppErrorError{"password": {ID: "user.login.password.error"}}}
		return errBuilder(models.NewAppError(ctx, path, "user.login.password.error", nil, "", int(codes.InvalidArgument), errors))
	}

	// TODO: handle if this user is using mobile or not
	expiry := c.config().Security.GetAccessTokenExpiryWebInHours()
	body := map[string]any{
		"subject":      user.GetId(),
		"remember":     true,
		"remember_for": expiry * 60 * 60,
		"context": map[string]any{
			"lang":       ctx.AcceptLanguage,
			"email":      user.GetEmail(),
			"first_name": user.GetFirstName(),
		},
	}

	oauthPayload, marErr := json.Marshal(body)
	if marErr != nil {
		return internalErr(marErr, "failed to marshal json payload")
	}

	reqURL := fmt.Sprintf("%s/oauth2/auth/requests/login/accept?login_challenge=%s", c.config().Oauth.GetOauthAdminUrl(), req.GetLoginChallenge())
	oauthReq, reqErr := http.NewRequestWithContext(ctx.Context, http.MethodPut, reqURL, bytes.NewReader(oauthPayload))
	if reqErr != nil {
		return internalErr(reqErr, "failed to build login/accept HTTP request to send to OAuth service")
	}
	oauthReq.Header.Set("Content-Type", "application/json")

	// TODO: inclue the response time as metrics
	start := time.Now()
	resp, respErr := utils.HTTPRequestWithRetry(c.httpClient, oauthReq, 3)
	duration := time.Since(start)
	if respErr != nil {
		c.log.Errorf("HTTP %s %s failed: %v (took %s)", oauthReq.Method, oauthReq.URL, err, duration)
		totalDuration := time.Since(startTime).Seconds()
		c.metricsCollector.RecordLoginRequest(false, totalDuration)
		return internalErr(respErr, "failed to request OAuth server to accept login")
	}
	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusInternalServerError {
		totalDuration := time.Since(startTime).Seconds()
		c.metricsCollector.RecordLoginRequest(false, totalDuration)
		return internalErr(respErr, "failed to request OAuth server to accept login")
	}

	if resp.StatusCode != http.StatusOK {
		totalDuration := time.Since(startTime).Seconds()
		c.metricsCollector.RecordLoginRequest(false, totalDuration)
		var resErr intModels.OAuthErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&resErr); err != nil {
			return internalErr(err, "failed to unmarshal login/accept response from Oauth service error")
		}
		errors := &models.AppErrorErrorsArgs{
			Err: respErr,
			ErrorsInternal: map[string]*models.AppErrorError{
				"error":             {ID: "login.error"},
				"error_description": {ID: intModels.GetOAuthRequestErrMsgID(ctx.AcceptLanguage, resErr.Error, resErr.ErrorDescription)},
			},
		}
		return errBuilder(models.NewAppError(ctx, path, "login.error", nil, "", int(codes.InvalidArgument), errors))
	}

	var result struct {
		RedirectTo string `json:"redirect_to"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		totalDuration := time.Since(startTime).Seconds()
		c.metricsCollector.RecordLoginRequest(false, totalDuration)
		return internalErr(err, "failed to unmarshal response from login/accept Oauth service")
	}
	if result.RedirectTo == "" {
		totalDuration := time.Since(startTime).Seconds()
		c.metricsCollector.RecordLoginRequest(false, totalDuration)
		return internalErr(err, "received an empty redirect_url from OAuth service login/accept")
	}

	ar.Success()

	totalDuration := time.Since(startTime).Seconds()
	c.metricsCollector.RecordLoginRequest(true, totalDuration)

	meta := map[string]string{"redirect_to": result.RedirectTo}
	return sucBuilder(&pbSh.SuccessResponseData{Metadata: meta})
}
