package models

type TokenType string

const (
	TokenTypePasswordReset     TokenType = "password_reset"
	TokenTypeEmailConfirmation TokenType = "email_confirmation"
)
