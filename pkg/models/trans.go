package models

type TranslateFunc = func(transId string) string

const (
	NoTranslation string = "<untranslated>"
)
