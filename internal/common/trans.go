package common

import (
	"context"
	"fmt"
	"time"

	pb "github.com/ahmad-khatib0-org/megacommerce-proto/gen/go/common/v1"
	"github.com/ahmad-khatib0-org/megacommerce-user/pkg/utils"
	"google.golang.org/grpc/status"
)

func (cc *CommonClient) TranslationsGet() (map[string]*pb.TranslationElements, *utils.AppError) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	res, err := cc.client.TranslationsGet(ctx, &pb.TranslationsGetRequest{})
	if err != nil {
		sc, _ := status.FromError(err)
		return nil, utils.NewAppError("user.common.TranslationsGet", "failed_get_common_trans", nil, fmt.Sprintf("failed to get translations %v", err), int(sc.Code()))
	}

	if res.Error != nil {
		err := utils.AppErrorFromProto(res.Error)
		return nil, err
	} else {
		return res.Data, nil
	}
}
