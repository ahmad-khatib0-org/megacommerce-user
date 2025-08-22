package controller

import (
	"context"

	shPb "github.com/ahmad-khatib0-org/megacommerce-proto/gen/go/shared/v1"
	"google.golang.org/grpc"
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
		Params:          err.Params,
		NestedParams:    err.NestedParams,
	}
}

func (c *Controller) responseInterceptor(ctx context.Context, req any, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	res, err := handler(ctx, req)
	if err != nil {
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
