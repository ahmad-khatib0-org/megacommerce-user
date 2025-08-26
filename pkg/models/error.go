package models

import (
	"fmt"
	"maps"
	"strings"

	shared "github.com/ahmad-khatib0-org/megacommerce-proto/gen/go/shared/v1"
)

const (
	ErrMsgInternal = "server.internal.error"
)

const maxErrorLength = 1024

type InternalError struct {
	Temp bool   `json:"temp"`
	Err  error  `json:"err"`
	Msg  string `json:"msg"`
	Path string `json:"path"`
}

func (ie InternalError) Error() string {
	var sb strings.Builder

	if ie.Path != "" {
		sb.WriteString(ie.Path)
		sb.WriteString(": ")
	}

	if ie.Msg != "" {
		sb.WriteString(ie.Msg)
		sb.WriteString(", ")
	}

	if ie.Err != nil {
		sb.WriteString(fmt.Sprintf("%v", ie.Err))
	}

	sb.WriteString(fmt.Sprintf(", temp: %t", ie.Temp))

	return sb.String()
}

func (ie *InternalError) Message() string {
	return ie.Msg
}

type AppErrorError struct {
	ID     string
	Params map[string]any
}

type AppErrorErrorsArgs struct {
	Err                  error
	ErrorsInternal       map[string]*AppErrorError            // a map of: field_name: error
	ErrorsNestedInternal map[string]map[string]*AppErrorError // same as Errors, but nested
}

type AppError struct {
	Ctx                  *Context
	ID                   string
	Message              string         // displayed to end user, without debugging information
	DetailedError        string         // Internal error string to help the developer
	StatusCode           int            // The grpc status code
	IDParams             map[string]any // params passed to templates trans
	Where                string         // err location in the form of Struct.Func
	SkipTranslation      bool
	Err                  error
	Errors               map[string]string
	ErrorsNested         map[string]map[string]string
	ErrorsInternal       map[string]*AppErrorError            // a map of: field_name: error
	ErrorsNestedInternal map[string]map[string]*AppErrorError // same as Errors, but nested
}

func (ae *AppError) Error() string {
	if ae == nil {
		return ""
	}

	var sb strings.Builder

	// render the error information
	if ae.Where != "" {
		sb.WriteString(fmt.Sprintf("%s :", ae.Where))
	}

	if ae.Message != "" {
		sb.WriteString(fmt.Sprintf("%s ,", ae.Message))
	}

	// only render the detailed error when it's present
	if ae.DetailedError != "" {
		sb.WriteString(fmt.Sprintf("%s ,", ae.DetailedError))
	}

	// render the wrapped error
	if ae.Err != nil {
		sb.WriteString(fmt.Sprintf("%s ", ae.Err.Error()))
	}

	res := sb.String()
	if len(res) > maxErrorLength {
		res = res[:maxErrorLength] + "..."
	}

	return res
}

func (ae *AppError) Translate(tf TranslateFunc) {
	if ae.SkipTranslation {
		return
	}

	if tf == nil {
		ae.Message = ae.ID
		return
	} else {
		tr := tf(ae.Ctx.AcceptLanguage, ae.ID, ae.IDParams)
		ae.Message = tr
	}

	if len(ae.ErrorsInternal) > 0 {
		errors := make(map[string]string, len(ae.ErrorsInternal))
		for k, v := range ae.ErrorsInternal {
			errors[k] = tf(ae.Ctx.AcceptLanguage, v.ID, v.Params)
		}
		ae.Errors = errors
	}
}

func (ae *AppError) Unwrap() error {
	return ae.Err
}

func (ae *AppError) Wrap(err error) *AppError {
	ae.Err = err
	return ae
}

func (ae *AppError) WipeDetailed() {
	ae.Err = nil
	ae.DetailedError = ""
}

func AppErrorDefault() *AppError {
	return &AppError{
		ID:                   "",
		Message:              "",
		DetailedError:        "",
		StatusCode:           0,
		Where:                "",
		SkipTranslation:      false,
		IDParams:             make(map[string]any),
		Ctx:                  &Context{},
		Errors:               make(map[string]string),
		ErrorsInternal:       make(map[string]*AppErrorError),
		ErrorsNested:         make(map[string]map[string]string),
		ErrorsNestedInternal: make(map[string]map[string]*AppErrorError),
		Err:                  nil,
	}
}

func NewAppError(
	ctx *Context,
	where string,
	id string,
	idParams map[string]any,
	details string,
	status int,
	errors *AppErrorErrorsArgs,
) *AppError {
	if errors == nil {
		errors = &AppErrorErrorsArgs{}
	}

	ap := &AppError{
		Ctx:                  ctx,
		ID:                   id,
		IDParams:             idParams,
		Message:              id,
		Where:                where,
		DetailedError:        details,
		StatusCode:           status,
		Err:                  errors.Err,
		ErrorsInternal:       errors.ErrorsInternal,
		ErrorsNestedInternal: errors.ErrorsNestedInternal,
	}

	ap.Translate(Tr)
	return ap
}

func AppErrorFromProto(ctx *Context, ae *shared.AppError) *AppError {
	if ae == nil {
		return AppErrorDefault()
	}

	errors, errorsNested := AppErrorConvertProtoParams(ae)
	return &AppError{
		ID:              ae.GetId(),
		Ctx:             ctx,
		Message:         ae.GetMessage(),
		DetailedError:   ae.GetDetailedError(),
		StatusCode:      int(ae.GetStatusCode()),
		Where:           ae.GetWhere(),
		SkipTranslation: ae.GetSkipTranslation(),
		Errors:          errors,
		ErrorsNested:    errorsNested,
	}
}

func AppErrorToProto(e *AppError) *shared.AppError {
	nested := make(map[string]*shared.StringMap, len(e.ErrorsNested))

	if len(e.ErrorsNested) > 0 {
		for k, v := range e.ErrorsNested {
			nested[k] = &shared.StringMap{Data: v}
		}
	}

	return &shared.AppError{
		Id:              e.ID,
		RequestId:       e.Ctx.RequestID,
		Message:         e.Message,
		DetailedError:   e.DetailedError,
		StatusCode:      int32(e.StatusCode),
		Where:           e.Where,
		SkipTranslation: e.SkipTranslation,
		Errors:          &shared.StringMap{Data: e.Errors},
		ErrorsNested:    &shared.NestedStringMap{Data: nested},
	}
}

func AppErrorConvertProtoParams(ae *shared.AppError) (map[string]string, map[string]map[string]string) {
	if ae.Errors == nil && ae.ErrorsNested == nil {
		return nil, nil
	}

	shallowCount := 0
	nestedCount := 0
	if ae.Errors != nil && len(ae.Errors.Data) > 0 {
		shallowCount = len(ae.Errors.Data)
	}
	if ae.ErrorsNested != nil && len(ae.ErrorsNested.Data) > 0 {
		nestedCount = len(ae.ErrorsNested.Data)
	}

	shallow := make(map[string]string, shallowCount)
	nested := make(map[string]map[string]string, nestedCount)

	if ae.Errors != nil && len(ae.Errors.Data) > 0 {
		maps.Copy(shallow, ae.Errors.Data)
	}

	if ae.ErrorsNested != nil && len(ae.ErrorsNested.Data) > 0 {
		for k, v := range ae.ErrorsNested.Data {
			nested[k] = v.Data
		}
	}

	return shallow, nested
}
