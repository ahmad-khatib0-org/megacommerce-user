package models

type EventName string

const (
	EventNameSupplierCreate    = "supplier_create"
	EventNameCustomerCreate    = "customer_create"
	EventNameEmailConfirmation = "email_confirmation"
	EventNamePasswordForgot    = "password_forgot"
	EventNameLogin             = "login"
)

type TokenType string

const (
	TokenTypePasswordReset     TokenType = "password_reset"
	TokenTypeEmailConfirmation TokenType = "email_confirmation"
)
