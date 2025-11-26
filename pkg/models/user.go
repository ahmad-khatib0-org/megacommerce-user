// Package models contains models for user, config, validation....
package models

type UserType string

const (
	UserTypeSupplier UserType = "supplier"
	UserTypeCustomer UserType = "customer"
)

const (
	UserEmailMaxLength    = 256
	UserNameMaxLength     = 64
	UserNameMinLength     = 2
	UserFirstNameMaxRunes = 64
	UserFirstNameMinRunes = 2
	UserLastNameMaxRunes  = 64
	UserLastNameMinRunes  = 2
	UserAuthDataMaxLength = 128
	UserPasswordMaxLength = 72
	UserPasswordMinLength = 8
	UserLocaleMaxLength   = 5
	UserTimezoneMaxRunes  = 256
	UserRolesMaxLength    = 256
	UserImageMaxSizeBytes = 1024 * 1024 * 2
)

var UserImageAllowedTypes = []string{"image/png", "image/webp", "image/jpeg", "image/jpg"}
