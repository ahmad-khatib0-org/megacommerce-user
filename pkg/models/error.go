package models

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"maps"
	"strings"

	shared "github.com/ahmad-khatib0-org/megacommerce-proto/gen/go/shared/v1"
)

const maxErrorLength = 1024

type InternalError struct {
	Temp bool   `json:"temp"`
	Err  error  `json:"err"`
	Msg  string `json:"msg"`
	Path string `json:"path"`
}

func (ie *InternalError) Error() string {
	var sb strings.Builder

	if ie.Path != "" {
		sb.WriteString(ie.Path)
		sb.WriteString(": ")
	}

	if ie.Msg != "" {
		sb.WriteString(ie.Msg)
		sb.WriteString(", ")
	}

	sb.WriteString(fmt.Sprintf(", temp: %t", ie.Temp))

	return sb.String()
}

func (ie *InternalError) Message() string {
	return ie.Msg
}

type AppError struct {
	Ctx *Context `json:"ctx"`
	Id  string   `json:"id"`
	// Message to be display to the end user without debugging information
	Message string `json:"message"`
	// Internal error string to help the developer
	DetailedError string `json:"detailed_error"`
	// The RequestId that's also set in the header
	RequestId string `json:"request_id,omitempty"`
	// The grpc status code
	StatusCode   int                          `json:"status_code,omitempty"`
	TrParams     map[string]any               `json:"tr_params"`
	Params       map[string]string            `json:"params,omitempty"`
	NestedParams map[string]map[string]string `json:"nested_params,omitempty"`
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

	if er.Message != NoTranslation {
		sb.WriteString(er.Message)
	}

	// only render the detailed error when it's present
	if er.DetailedError != "" {
		if er.Message != NoTranslation {
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

func (er *AppError) Translate(tf TranslateFunc) {
	if er.SkipTranslation {
		return
	}

	if tf == nil {
		er.Message = er.Id
		return
	} else {
		tr, err := tf(er.Ctx.AcceptLanguage, er.Id, er.TrParams)
		if err != nil {
			// Track error
		}
		er.Message = tr
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
		TrParams:        make(map[string]any),
		Params:          make(map[string]string),
		NestedParams:    make(map[string]map[string]string),
		Wrapped:         nil,
	}
}

func NewAppError(
	ctx *Context,
	where string,
	id string,
	trParams map[string]any,
	details string,
	status int,
	err error,
) *AppError {
	ap := &AppError{
		Ctx:           ctx,
		Id:            id,
		TrParams:      trParams,
		Message:       id,
		Where:         where,
		DetailedError: details,
		StatusCode:    status,
		Wrapped:       err,
	}

	ap.Translate(Tr)
	return ap
}

func AppErrorFromProto(ae *shared.AppError) *AppError {
	if ae == nil {
		ae := &AppError{}
		return ae.Default()
	}

	params, nested := AppErrorConvertProtoParams(ae)

	return &AppError{
		Id:              ae.Id,
		Message:         ae.Message,
		DetailedError:   ae.DetailedError,
		RequestId:       ae.RequestId,
		StatusCode:      int(ae.StatusCode),
		Where:           ae.Where,
		SkipTranslation: ae.SkipTranslation,
		Params:          params,
		NestedParams:    nested,
	}
}

func AppErrorToProto(e *AppError) *shared.AppError {
	nested := make(map[string]*shared.StringMap, len(e.NestedParams))

	if len(e.NestedParams) > 0 {
		for k, v := range e.NestedParams {
			nested[k] = &shared.StringMap{Data: v}
		}
	}

	return &shared.AppError{
		Id:              e.Id,
		Message:         e.Message,
		DetailedError:   e.DetailedError,
		StatusCode:      int32(e.StatusCode),
		Where:           e.Where,
		SkipTranslation: e.SkipTranslation,
		Params:          &shared.StringMap{Data: e.Params},
		NestedParams:    &shared.NestedStringMap{Data: nested},
	}
}

func AppErrorConvertProtoParams(ae *shared.AppError) (map[string]string, map[string]map[string]string) {
	if ae.Params == nil && ae.NestedParams == nil {
		return nil, nil
	}

	shallowCount := 0
	nestedCount := 0
	if ae.Params != nil && len(ae.Params.Data) > 0 {
		shallowCount = len(ae.Params.Data)
	}
	if ae.NestedParams != nil && len(ae.NestedParams.Data) > 0 {
		nestedCount = len(ae.NestedParams.Data)
	}

	shallow := make(map[string]string, shallowCount)
	nested := make(map[string]map[string]string, nestedCount)

	if ae.Params != nil && len(ae.Params.Data) > 0 {
		maps.Copy(shallow, ae.Params.Data)
	}

	if ae.NestedParams != nil && len(ae.NestedParams.Data) > 0 {
		for k, v := range ae.NestedParams.Data {
			nested[k] = v.Data
		}
	}

	return shallow, nested
}
