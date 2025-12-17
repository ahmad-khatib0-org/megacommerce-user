package controller

import (
	"context"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	com "github.com/ahmad-khatib0-org/megacommerce-proto/gen/go/common/v1"
	user "github.com/ahmad-khatib0-org/megacommerce-proto/gen/go/users/v1"
	"github.com/ahmad-khatib0-org/megacommerce-shared-go/pkg/logger"
	"github.com/ahmad-khatib0-org/megacommerce-shared-go/pkg/models"
	"github.com/ahmad-khatib0-org/megacommerce-shared-go/pkg/utils"
	"github.com/ahmad-khatib0-org/megacommerce-user/internal/common"
	mailerMocks "github.com/ahmad-khatib0-org/megacommerce-user/internal/mailer/mocks"
	storeMocks "github.com/ahmad-khatib0-org/megacommerce-user/internal/store/mocks"
	workerMocks "github.com/ahmad-khatib0-org/megacommerce-user/internal/worker/mocks"
	intModels "github.com/ahmad-khatib0-org/megacommerce-user/pkg/models"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/spf13/viper"
)

type TestingUser struct {
	Ctx  *models.Context
	User *user.User
}

type TestHelper struct {
	config     func() *com.Config
	srvCfg     *intModels.Config
	log        *logger.Logger
	common     *common.CommonClient
	controller *Controller
	store      *storeMocks.MockUsersStore
	tasker     *workerMocks.MockTaskDistributor
	mailer     *mailerMocks.MockMailerService
	Customer1  *TestingUser
	Customer2  *TestingUser
	Supplier1  *TestingUser
	Supplier2  *TestingUser
}

func NewTestHelper(tb testing.TB) (*TestHelper, *models.InternalError) {
	th := &TestHelper{}
	log, err := logger.InitLogger("dev")
	if err != nil {
		return nil, &models.InternalError{Err: err, Msg: "failed to initialize logger"}
	}
	th.log = log

	if err := th.initServiceConfig(); err != nil {
		return nil, err
	}
	if err := th.initCommonService(); err != nil {
		return nil, err
	}
	if err := th.initSharedConfig(); err != nil {
		return nil, err
	}
	if err := th.initTrans(); err != nil {
		return nil, err
	}

	store := storeMocks.NewMockUsersStore(tb)
	mailer := mailerMocks.NewMockMailerService(tb)
	tasker := workerMocks.NewMockTaskDistributor(tb)

	th.mailer = mailer
	th.store = store
	th.tasker = tasker
	th.controller = &Controller{
		config: th.config,
		log:    th.log,
		store:  store,
		tasker: tasker,
	}

	th.initUsers()
	return th, nil
}

func (th *TestHelper) TearDown() {
	if err := th.common.Close(); err != nil {
		th.log.Warnf("failed to close the common client listener ", err)
	}
}

func (th *TestHelper) initUsers() {
	th.Customer1 = &TestingUser{
		Ctx:  th.getContext(),
		User: th.createUser(intModels.UserTypeCustomer, models.RoleIDCustomer),
	}
	th.Customer2 = &TestingUser{
		Ctx:  th.getContext(),
		User: th.createUser(intModels.UserTypeCustomer, models.RoleIDCustomer),
	}
	th.Supplier1 = &TestingUser{
		Ctx:  th.getContext(),
		User: th.createUser(intModels.UserTypeSupplier, models.RoleIDSupplierAdmin),
	}
	th.Supplier2 = &TestingUser{
		Ctx:  th.getContext(),
		User: th.createUser(intModels.UserTypeSupplier, models.RoleIDSupplierAdmin),
	}
}

func (th *TestHelper) createUser(userType intModels.UserType, roleID models.RoleID) *user.User {
	return &user.User{
		Id:                 utils.NewIDPointer(),
		Username:           utils.NewPointer(utils.RandomUserName(6, 12)),
		Email:              utils.NewPointer(gofakeit.Email()),
		FirstName:          utils.NewPointer(gofakeit.FirstName()),
		LastName:           utils.NewPointer(gofakeit.LastName()),
		UserType:           utils.NewPointer(string(userType)),
		Password:           utils.NewPointer(gofakeit.Password(true, true, true, true, false, 8)),
		Roles:              []string{string(roleID)},
		IsEmailVerified:    utils.NewPointer(true),
		AuthData:           utils.NewPointer(""),
		AuthService:        utils.NewPointer(""),
		FailedAttempts:     utils.NewPointer(int32(0)),
		MfaActive:          utils.NewPointer(false),
		MfaSecret:          utils.NewPointer(""),
		LastPasswordUpdate: nil,
		LastPictureUpdate:  nil,
		Locale:             th.config().GetLocalization().DefaultClientLocale,
		LastActivityAt:     utils.NewPointer(utils.TimeGetMillis()),
		LastLogin:          utils.NewPointer(utils.TimeGetMillis()),
		Membership:         utils.NewPointer("free"),
		UpdatedAt:          nil,
		DeletedAt:          nil,
		CreatedAt:          utils.NewPointer(utils.TimeGetMillis()),
	}
}

func (th *TestHelper) getContext() *models.Context {
	ctx := &models.Context{
		RequestID:      utils.NewID(),
		IPAddress:      gofakeit.IPv4Address(),
		XForwardedFor:  gofakeit.IPv4Address(),
		UserAgent:      gofakeit.UserAgent(),
		AcceptLanguage: th.config().Localization.GetDefaultClientLocale(),
		Session: &models.Session{
			ID:        utils.NewID(),
			Token:     utils.NewID(),
			CreatedAt: utils.TimeGetMillis(),
			ExpiresAt: utils.TimeGetMillis() + time.Duration(time.Hour).Milliseconds(),
			UserID:    utils.NewID(),
			DeviceID:  utils.NewID(),
			IsOAuth:   gofakeit.Bool(),
		},
	}

	return ctx
}

func (th *TestHelper) withContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, models.ContextKeyMetadata, th.getContext())
}

func (th *TestHelper) initServiceConfig() *models.InternalError {
	_, b, _, _ := runtime.Caller(0)
	basePath := filepath.Join(filepath.Dir(b), "../..")
	configFilePath := filepath.Join(basePath, "config.dev.yaml")

	viper.SetConfigFile(configFilePath)
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		return &models.InternalError{Err: err, Msg: "failed to initialize service config"}
	}

	var c intModels.Config
	if err := viper.Unmarshal(&c); err != nil {
		return &models.InternalError{Err: err, Msg: "failed to initialize service config"}
	}

	th.srvCfg = &c
	return nil
}

func (th *TestHelper) initSharedConfig() *models.InternalError {
	config, err := th.common.ConfigGet(com.Environment_LOCAL)
	if err != nil {
		return &models.InternalError{Err: err, Msg: "failed to initialize shared config", Path: "users.controller.initSharedConfig"}
	}
	th.config = func() *com.Config { return config }
	return nil
}

func (th *TestHelper) initCommonService() *models.InternalError {
	srv, err := common.NewCommonClient(&common.CommonArgs{Config: th.srvCfg})
	if err != nil {
		return &models.InternalError{Err: err, Msg: "failed to initialize common service client", Path: "users.controller.initCommonService"}
	}
	th.common = srv
	return nil
}

func (th *TestHelper) initTrans() *models.InternalError {
	trans, err := th.common.TranslationsGet()
	if err != nil {
		return err
	}
	if err := models.TranslationsInit(trans, th.config().GetLocalization().GetDefaultClientLocale()); err != nil {
		return &models.InternalError{Err: err, Msg: "failed to init translations", Path: "user.controller.initTrans"}
	}
	return nil
}
