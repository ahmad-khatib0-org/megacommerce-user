package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"

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
		// If the request exceeded FileSettings.MaxFileSize a plain error gets returned.
		// Convert it into an AppError.
		if string(data) == "http: request body too large\n" {
			return errors.New("The request was too large. Consider asking your System Admin to raise the FileSettings.MaxFileSize setting.")
		}

		return fmt.Errorf("failed to decode JSON payload into AppError. Body: %s, err: %v", string(data), err)
	}

	return &er
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
