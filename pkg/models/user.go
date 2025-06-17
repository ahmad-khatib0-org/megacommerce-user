package models

type UserType string

const (
	UserTypeSupplier UserType = "supplier"
	UserTypeBuyer    UserType = "buyer"
)

const (
	UserEmailMaxLength    = 256
	UserFirstNameMaxRunes = 64
	UserFirstNameMinRunes = 2
	UserLastNameMaxRunes  = 64
	UserLastNameMinRunes  = 2
	UserAuthDataMaxLength = 128
	UserNameMaxLength     = 64
	UserNameMinLength     = 2
	UserPasswordMaxLength = 72
	UserLocaleMaxLength   = 5
	UserTimezoneMaxRunes  = 256
	UserRolesMaxLength    = 256
)
