package common

import (
	"context"
	"time"

	pb "github.com/ahmad-khatib0-org/megacommerce-proto/gen/go/common/v1"
	"github.com/ahmad-khatib0-org/megacommerce-user/pkg/models"
)

func (cc *CommonClient) TranslationsGet() (map[string]*pb.TranslationElements, *models.InternalError) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	res, err := cc.client.TranslationsGet(ctx, &pb.TranslationsGetRequest{})
	if err != nil {
		return nil, &models.InternalError{
			Temp: false,
			Err:  err,
			Msg:  "failed to get translations from the common service",
			Path: "user.common.TranslationsGet",
		}
	}

	if res.Error != nil {
		return nil, &models.InternalError{
			Temp: false,
			Err:  models.AppErrorFromProto(res.Error),
			Msg:  "failed to get translations from the common service",
			Path: "user.common.TranslationsGet",
		}
	}

	return res.Data, nil
}
