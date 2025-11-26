package models

import "github.com/ahmad-khatib0-org/megacommerce-shared-go/pkg/models"

type TaskName string

type TaskFunc func()

const (
	TaskNameEmailBatching          TaskName = "email_batching"
	TaskNameSendVerifyEmail        TaskName = "send_verify_email"
	TaskNameSendPasswordResetEmail TaskName = "send_password_reset_email"
)

type TaskSendVerifyEmailPayload struct {
	Ctx     *models.Context `json:"ctx"`
	Email   string          `json:"email"`
	Token   string          `json:"token"`
	TokenID string          `json:"token_id"`
	Hours   int             `json:"hours"`
}

type TaskSendPasswordResetEmailPayload struct {
	Ctx     *models.Context `json:"ctx"`
	Email   string          `json:"email"`
	Token   string          `json:"token"`
	TokenID string          `json:"token_id"`
	Hours   int             `json:"hours"`
}
