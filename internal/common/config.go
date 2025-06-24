package common

import (
	"context"
	"fmt"
	"time"

	com "github.com/ahmad-khatib0-org/megacommerce-proto/gen/go/common/v1"
	"github.com/ahmad-khatib0-org/megacommerce-user/pkg/models"
)

func (cc *CommonClient) ConfigGet() (*com.Config, *models.InternalError) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	res, err := cc.client.ConfigGet(ctx, &com.ConfigGetRequest{})
	if err != nil {
		return nil, &models.InternalError{
			Temp: false,
			Err:  err,
			Msg:  "failed to get configurations from common service",
			Path: "user.common.ConfigGet",
		}
	}

	switch res := res.Response.(type) {
	case *com.ConfigGetResponse_Data:
		return res.Data, nil
	case *com.ConfigGetResponse_Error:
		err := models.AppErrorFromProto(res.Error)
		return nil, &models.InternalError{
			Temp: false,
			Err:  err,
			Msg:  "failed to get configurations from common service",
			Path: "user.common.ConfigGet",
		}
	}

	return nil, nil
}

// TODO: complete this listener
func (cc *CommonClient) ConfigListener(clientID string) *models.InternalError {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	stream, err := cc.client.ConfigListener(ctx, &com.ConfigListenerRequest{ClientId: clientID})
	if err != nil {
		return &models.InternalError{
			Temp: false,
			Err:  err,
			Msg:  "failed to register a listener call",
			Path: "user.common.ConfigListener",
		}
	}

	for {
		res, err := stream.Recv()
		if err != nil {
			fmt.Println("config listener stream closed")
			break
		}

		switch x := res.Response.(type) {
		case *com.ConfigListenerResponse_Data:
			fmt.Println("Config changed: ", x.Data)
		case *com.ConfigListenerResponse_Error:
			fmt.Println("Error received: ", x.Error.Message)
		}
	}

	return nil
}
