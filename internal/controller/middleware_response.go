package controller

import (
	"context"

	shPb "github.com/ahmad-khatib0-org/megacommerce-proto/gen/go/shared/v1"
	"github.com/ahmad-khatib0-org/megacommerce-user/pkg/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// processAppError sanitize the returned app error
func processAppError(err *shPb.AppError) *shPb.AppError {
	if err == nil {
		return nil
	}

	return &shPb.AppError{
		Id:              err.Id,
		Message:         err.Message,
		RequestId:       err.RequestId,
		StatusCode:      err.StatusCode,
		SkipTranslation: err.SkipTranslation,
		Errors:          err.Errors,
		ErrorsNested:    err.ErrorsNested,
	}
}

func (c *Controller) responseInterceptor(ctx context.Context, req any, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	res, err := handler(ctx, req)
	if err != nil {
		if st, ok := status.FromError(err); ok {
			if md, ok := metadata.FromIncomingContext(ctx); ok {
				lang := utils.GetAcceptedLanguageFromGrpcCtx(ctx, md, c.cfg.Localization.GetAvailableLocales(), c.cfg.Localization.GetDefaultClientLocale())
				return resp, status.Errorf(st.Code(), getGrpcErrMsg(lang, st.Code()))
			}
		}

		return res, err
	}

	if m, ok := res.(proto.Message); ok {
		msgRef := m.ProtoReflect()

		// look for the oneof named "response"
		oneof := msgRef.Descriptor().Oneofs().ByName("response")
		if oneof != nil {
			field := msgRef.WhichOneof(oneof)
			if field != nil && field.Name() == "error" {
				v := msgRef.Get(field)
				if v.Message().IsValid() {
					if appErr, ok := v.Message().Interface().(*shPb.AppError); ok {
						newErr := processAppError(appErr)
						msgRef.Set(field, protoreflect.ValueOf(newErr.ProtoReflect()))
					}
				}
			}
		}
	}

	return res, nil
}
