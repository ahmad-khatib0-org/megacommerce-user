package store

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/ahmad-khatib0-org/megacommerce-user/pkg/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type DBErrorType string

const (
	DBErrorTypeNoRows              DBErrorType = "no_rows"
	DBErrorTypeUniqueViolation     DBErrorType = "unique_violation"
	DBErrorTypeForeignKeyViolation DBErrorType = "foreign_key_violation"
	DBErrorTypeNotNullViolation    DBErrorType = "not_null_violation"
	DBErrorTypeJsonMarshal         DBErrorType = "json_marshal"
	DBErrorTypeJsonUnmarshal       DBErrorType = "json_unmarshal"
	DBErrorTypeConnection          DBErrorType = "connection_exception"
	DBErrorTypePrivileges          DBErrorType = "insufficient_privilege"
	DBErrorTypeStartTransaction    DBErrorType = "start_transaction"
	DBErrorTypeInternal            DBErrorType = "insufficient_privilege"
)

type DBError struct {
	ErrType DBErrorType
	Err     error
	Msg     string
	Path    string
	Details string
}

func (de *DBError) Error() string {
	if de == nil {
		return "DBError <nil>"
	}

	var sb strings.Builder
	if de.Path != "" {
		sb.WriteString(fmt.Sprintf("path: %s", de.Path))
		sb.WriteString(", ")
	}

	if de.ErrType != "" {
		sb.WriteString(fmt.Sprintf("err_type: %s ", de.ErrType))
		sb.WriteString(", ")
	}

	if de.Msg != "" {
		sb.WriteString(fmt.Sprintf("msg: %s ", de.Msg))
		sb.WriteString(", ")
	}

	if de.Details != "" {
		sb.WriteString(fmt.Sprintf("details: %s ", de.Details))
		sb.WriteString(", ")
	}

	if de.Err != nil {
		sb.WriteString(fmt.Sprintf("err: %v ", de.Err))
	}

	return sb.String()
}

func StartTransactionError(err error, path string) *DBError {
	return &DBError{
		ErrType: DBErrorTypeStartTransaction,
		Err:     err,
		Msg:     "failed to start a db transaction",
		Path:    path,
	}
}

func JsonMarshalError(err error, path, msg string) *DBError {
	return &DBError{
		ErrType: DBErrorTypeJsonMarshal,
		Err:     err,
		Msg:     msg,
		Path:    path,
	}
}

// HandleDBError builds a DBError, rollback the transaction if
func HandleDBError(ctx *models.Context, err error, path string, tr pgx.Tx) *DBError {
	// TODO: track rolling back error
	if tr != nil {
		errRb := tr.Rollback(ctx.Context)
		fmt.Printf("error rolling back the transaction %v", errRb)
	}

	if err == nil {
		return nil
	}

	intErr := func() *DBError {
		return &DBError{ErrType: DBErrorTypeInternal, Path: path, Err: err, Msg: "database error"}
	}

	switch e := err.(type) {
	case *pgconn.PgError:
		// PostgreSQL-specific errors
		switch e.Code {
		// Constraint violations
		case "23505": // unique_violation
			return &DBError{
				ErrType: DBErrorTypeUniqueViolation,
				Err:     e,
				Path:    path,
				Msg:     parseDuplicateFieldDBError(e),
			}

		case "23503": // foreign_key_violation
			return &DBError{
				ErrType: DBErrorTypeForeignKeyViolation,
				Err:     e,
				Path:    path,
				Msg:     "referenced record is not found",
			}

		case "23502": // not_null_violation
			return &DBError{
				ErrType: DBErrorTypeNotNullViolation,
				Err:     e,
				Path:    path,
				Msg:     fmt.Sprintf("%s cannot be null ", parseDBFieldName(e)),
			}

			// Connection/availability errors
		case "08000", "08003", "08006": // connection exceptions
			return &DBError{
				ErrType: DBErrorTypeConnection,
				Err:     e,
				Path:    path,
				Msg:     "database connection exception",
			}

		// Permission errors
		case "42501": // insufficient_privilege
			return &DBError{
				ErrType: DBErrorTypePrivileges,
				Err:     e,
				Path:    path,
				Msg:     "insufficient permissions to preform an action",
			}
		}
	default:
		if errors.Is(err, pgx.ErrNoRows) {
			return &DBError{
				ErrType: DBErrorTypeNoRows,
				Path:    path,
				Err:     e,
				Msg:     "the requested resource is not found",
			}
		}
		return intErr()
	}

	return intErr()
}

// Extract the duplicate field from error detail
// Example: "Key (email)=(test@example.com) already exists.
func parseDuplicateFieldDBError(err *pgconn.PgError) string {
	parts := strings.Split(err.Detail, ")=(")
	if len(parts) > 0 {
		field := strings.TrimPrefix(parts[0], "Key (")
		return fmt.Sprintf("%s already exits ", field)
	}

	return err.Detail
}

// Extract field name from error message
// Example: "null value in column \"email\" violates not-null constraint
func parseDBFieldName(err *pgconn.PgError) string {
	re := regexp.MustCompile(`column "(.+?)"`)
	matches := re.FindStringSubmatch(err.Message)
	if len(matches) > 1 {
		return matches[0]
	}
	return "field"
}
