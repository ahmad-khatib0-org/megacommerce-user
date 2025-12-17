package common

import (
	"context"
	"fmt"
	"time"

	com "github.com/ahmad-khatib0-org/megacommerce-proto/gen/go/common/v1"
	"github.com/ahmad-khatib0-org/megacommerce-shared-go/pkg/models"
)

func (cc *CommonClient) GetServiceEnv() com.Environment {
	env := cc.cfg.Service.Env
	switch env {
	case "local":
		return com.Environment_LOCAL
	case "dev":
		return com.Environment_DEV
	case "production":
		return com.Environment_PRODUCTION
	default:
		return com.Environment_DEV
	}
}

func (cc *CommonClient) ConfigGet(env com.Environment) (*com.Config, *models.InternalError) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	res, err := cc.client.ConfigGet(ctx, &com.ConfigGetRequest{Env: env})
	ie := func(err error, msg string) (*com.Config, *models.InternalError) {
		return nil, &models.InternalError{Err: err, Msg: msg, Path: "user.common.ConfigGet"}
	}

	if err != nil {
		return ie(err, "failed to get configurations from common service")
	}

	switch res := res.Response.(type) {
	case *com.ConfigGetResponse_Data:
		return res.Data, nil
	case *com.ConfigGetResponse_Error:
		err := models.AppErrorFromProto(nil, res.Error) // no need for ctx here
		return ie(err, "failed to get configurations from common service")
	}

	return nil, nil
}

func (cc *CommonClient) ConfigListener(clientID string) *models.InternalError {
	// TODO: complete this listener
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	stream, err := cc.client.ConfigListener(ctx, &com.ConfigListenerRequest{ClientId: clientID})
	path := "user.common.ConfigListener"
	if err != nil {
		return &models.InternalError{Err: err, Msg: "failed to register a listener call", Path: path}
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
