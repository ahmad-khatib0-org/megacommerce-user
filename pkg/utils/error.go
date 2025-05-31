package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	pb "github.com/ahmad-khatib0-org/megacommerce-proto/gen/go/common/v1"
	"github.com/ahmad-khatib0-org/megacommerce-user/pkg/models"
)

const maxErrorLength = 1024

type AppError struct {
	Id string `json:"id"`
	// Message to be display to the end user without debugging information
	Message string `json:"message"`
	// Internal error string to help the developer
	DetailedError string `json:"detailed_error"`
	// The RequestId that's also set in the header
	RequestId string `json:"request_id,omitempty"`
	// The grpc status code
	StatusCode int            `json:"status_code,omitempty"`
	Params     map[string]any `json:"params"`
	// The function where it happened in the form of Struct.Func
	Where string `json:"-"`
	// Whether translation for the error should be skipped.
	SkipTranslation bool  `json:"-"`
	Wrapped         error `json:"-"`
}

func (er *AppError) Error() string {
	if er == nil {
		return ""
	}

	var sb strings.Builder

	// render the error information
	if er.Where != "" {
		sb.WriteString(er.Where)
		sb.WriteString(": ")
	}

	if er.Message != models.NoTranslation {
		sb.WriteString(er.Message)
	}

	// only render the detailed error when it's present
	if er.DetailedError != "" {
		if er.Message != models.NoTranslation {
			sb.WriteString(", ")
		}
		sb.WriteString(er.DetailedError)
	}

	// render the wrapped error
	err := er.Wrapped
	if err != nil {
		sb.WriteString(", ")
		sb.WriteString(err.Error())
	}

	res := sb.String()
	if len(res) > maxErrorLength {
		res = res[:maxErrorLength] + "..."
	}
	return res
}

func (er *AppError) Translate(tf models.TranslateFunc) {
	if er.SkipTranslation {
		return
	}

	if tf == nil {
		er.Message = er.Id
		return
	} else {
		er.Message = tf(er.Id)
	}
}

func (er *AppError) ToJSON() string {
	// turn the wrapped error into a detailed message
	detailed := er.DetailedError
	defer func() {
		er.DetailedError = detailed
	}()

	er.wrappedToDetailed()

	b, _ := json.Marshal(er)
	return string(b)
}

func (er *AppError) wrappedToDetailed() {
	if er.Wrapped == nil {
		return
	}

	if er.DetailedError != "" {
		er.DetailedError += ", "
	}

	er.DetailedError += er.Wrapped.Error()
}

func (er *AppError) Unwrap() error {
	return er.Wrapped
}

func (er *AppError) Wrap(err error) *AppError {
	er.Wrapped = err
	return er
}

func (er *AppError) WipeDetailed() {
	er.Wrapped = nil
	er.DetailedError = ""
}

// AppErrorFromJSON will try to decode the input into an AppError.
func AppErrorFromJSON(r io.Reader) error {
	data, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	var er AppError
	err = json.NewDecoder(bytes.NewReader(data)).Decode(&er)
	if err != nil {
		return fmt.Errorf("failed to decode JSON payload into AppError. Body: %s, err: %v", string(data), err)
	}

	return &er
}

func (ae *AppError) Default() *AppError {
	return &AppError{
		Id:              "",
		Message:         "",
		DetailedError:   "",
		RequestId:       "",
		StatusCode:      0,
		Where:           "",
		SkipTranslation: false,
		Params:          make(map[string]any),
		Wrapped:         nil,
	}
}

func NewAppError(where string, id string, params map[string]any, details string, status int) *AppError {
	ap := &AppError{
		Id:            id,
		Params:        params,
		Message:       id,
		Where:         where,
		DetailedError: details,
		StatusCode:    status,
	}
	// ap.Translate() // TODO: add translate
	return ap
}

func AppErrorFromProto(protoErr *pb.AppError) *AppError {
	if protoErr == nil {
		ae := &AppError{}
		return ae.Default()
	}

	return &AppError{
		Id:              protoErr.Id,
		Message:         protoErr.Message,
		DetailedError:   protoErr.DetailedError,
		RequestId:       protoErr.RequestId,
		StatusCode:      int(protoErr.StatusCode),
		Where:           protoErr.Where,
		SkipTranslation: protoErr.SkipTranslation,
		Params:          AppErrorConvertProtoParams(protoErr),
	}
}

func AppErrorConvertProtoParams(ae *pb.AppError) map[string]any {
	if ae.Params == nil {
		return nil
	}

	switch {
	case ae.Params.MessageIs(&pb.StringMap{}):
		var sm pb.StringMap
		if err := ae.Params.UnmarshalTo(&sm); err == nil {
			params := make(map[string]any, len(sm.Data))
			for k, v := range sm.Data {
				params[k] = v
			}
			return params
		}

	case ae.Params.MessageIs(&pb.NestedStringMap{}):
		var nsm pb.NestedStringMap
		if err := ae.Params.UnmarshalTo(&nsm); err == nil {
			params := make(map[string]any, len(nsm.Data))
			for k, v := range nsm.Data {
				params[k] = v
			}
			return params
		}

	default:
		return nil // TODO: fatal here
	}

	return nil
}
