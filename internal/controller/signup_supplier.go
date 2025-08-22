package controller

import (
	"bytes"
	"context"
	"fmt"
	"time"

	shPb "github.com/ahmad-khatib0-org/megacommerce-proto/gen/go/shared/v1"
	pb "github.com/ahmad-khatib0-org/megacommerce-proto/gen/go/users/v1"
	"github.com/ahmad-khatib0-org/megacommerce-user/internal/files"
	"github.com/ahmad-khatib0-org/megacommerce-user/internal/store"
	"github.com/ahmad-khatib0-org/megacommerce-user/internal/worker"
	"github.com/ahmad-khatib0-org/megacommerce-user/pkg/models"
	"github.com/ahmad-khatib0-org/megacommerce-user/pkg/utils"
	"github.com/hibiken/asynq"
	"github.com/minio/minio-go/v7"
	"github.com/oklog/ulid/v2"
	"google.golang.org/grpc/codes"
)

func (c *Controller) CreateSupplier(context context.Context, req *pb.SupplierCreateRequest) (*pb.SupplierCreateResponse, error) {
	path := "user.controller.SignupSupplier"
	errBuilder := func(e *models.AppError) (*pb.SupplierCreateResponse, error) {
		return &pb.SupplierCreateResponse{Response: &pb.SupplierCreateResponse_Error{Error: models.AppErrorToProto(e)}}, nil
	}
	ctx, err := models.ContextGet(context)
	if err != nil {
		return errBuilder(err)
	}

	internalErr := func(err error) *models.AppError {
		return models.NewAppError(ctx, path, models.ErrMsgInternal, nil, "", int(codes.Internal), &models.AppErrorErrorsArgs{Err: err})
	}

	ar := models.AuditRecordNew(ctx, models.EventNameSupplierCreate, models.EventStatusFail)
	models.AuditEventDataParameter(ar, "supplier", models.SignupSupplierRequestAuditable(req))
	defer c.ProcessAudit(ar)

	sanitized := models.SignupSupplierRequestSanitize(req)
	if err = models.SignupSupplierRequestIsValid(ctx, sanitized, c.cfg.Password); err != nil {
		return errBuilder(err)
	}

	if sanitized.Image != nil {
		if imgErr := files.AttachmentsValidateSizeAndTypes(&files.AttachmentValidationConfig{
			Files:        []*shPb.Attachment{sanitized.Image},
			MaxSize:      models.UserImageMaxSizeBytes,
			AllowedTypes: models.UserImageAllowedTypes,
			Unit:         files.FileSizeUnitMB,
		}); imgErr != nil {
			errors := &models.AppErrorErrorsArgs{ErrorsInternal: map[string]*models.AppErrorError{"image": imgErr.Err}}
			err = models.NewAppError(ctx, path, imgErr.Err.ID, imgErr.Err.Params, "", int(codes.InvalidArgument), errors)
			return errBuilder(err)
		}
	}

	dbPay, err := models.SignupSupplierRequestPreSave(
		ctx,
		&pb.User{
			Username:   utils.NewPointer(sanitized.GetUsername()),
			FirstName:  utils.NewPointer(sanitized.GetFirstName()),
			LastName:   utils.NewPointer(sanitized.GetLastName()),
			Email:      utils.NewPointer(sanitized.GetEmail()),
			Membership: utils.NewPointer("free"),
			Password:   utils.NewPointer(req.GetPassword()),
			Roles:      []string{string(models.RoleIDSupplierAdmin)},
		},
	)
	if err != nil {
		return errBuilder(err)
	}

	token := &utils.Token{}
	tokenData, errTok := token.GenerateToken(time.Duration(time.Hour * time.Duration(c.cfg.Security.GetTokenConfirmationExpiryInHours())))
	if errTok != nil {
		return errBuilder(internalErr(errTok))
	}

	if sanitized.GetImage() != nil {
		imgContent := sanitized.GetImage().GetData()
		imgName := ulid.Make().String()
		_, imgErr := c.objStorage.PutObject(ctx.Context, c.cfg.File.GetAmazonS3Bucket(), imgName, bytes.NewReader(imgContent), int64(len(imgContent)), minio.PutObjectOptions{
			ContentType: sanitized.GetImage().GetMime(),
		})
		if imgErr != nil {
			return errBuilder(internalErr(err))
		}

		dbPay.Image = utils.NewPointer(fmt.Sprintf("%s/%s", c.cfg.File.GetAmazonS3Bucket(), imgName))
		im := &pb.UserImageMetadata{
			Mime:      sanitized.Image.Mime,
			Height:    int32(sanitized.Image.Crop.Height),
			Widht:     int32(sanitized.Image.Crop.Width),
			SizeBytes: int64(sanitized.Image.FileSize),
		}
		dbPay.ImageMetadata = im
	}

	if err := c.store.SignupSupplier(ctx, dbPay, tokenData); err != nil {
		if err.ErrType == store.DBErrorTypeUniqueViolation {
			id := "user.create.email.not_unique"
			details := fmt.Sprintf("the email %s is already in use", dbPay.GetEmail())
			errors := &models.AppErrorErrorsArgs{Err: err, ErrorsInternal: map[string]*models.AppErrorError{"email": {ID: id}}}
			return errBuilder(models.NewAppError(ctx, path, id, nil, details, int(codes.AlreadyExists), errors))
		} else {
			return errBuilder(internalErr(err))
		}
	}

	optoins := []asynq.Option{asynq.MaxRetry(10), asynq.ProcessIn(time.Second * 10), asynq.Queue(worker.QueuePriorityCritical)}
	taskPayload := &models.TaskSendVerifyEmailPayload{
		Ctx:     ctx,
		Email:   dbPay.GetEmail(),
		Token:   tokenData.Token,
		TokenID: tokenData.ID,
		Hours:   int(c.cfg.Security.GetTokenConfirmationExpiryInHours()),
	}

	// TODO: handle error
	if err := c.tasker.SendVerifyEmail(context, taskPayload, optoins...); err != nil {
		return errBuilder(internalErr(err))
	}

	ar.AuditEventDataResultState(models.SignupSupplierRequestResultState(dbPay))
	ar.Success()

	return &pb.SupplierCreateResponse{Response: &pb.SupplierCreateResponse_Data{}}, nil
}
