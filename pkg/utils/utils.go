package utils

import (
	"fmt"

	compb "github.com/ahmad-khatib0-org/megacommerce-proto/gen/go/common/v1"
)

type Utils struct {
	allTrans map[string]map[string]string
}

type UtilsArgs struct {
	AllTrans map[string]*compb.TranslationElements
}

func NewUtils(ua *UtilsArgs) (*Utils, *AppError) {
	u := &Utils{}
	if err := u.initTrans(ua); err != nil {
		return nil, err
	}

	return u, nil
}

func (u *Utils) Tr(lang string, key string) string {
	t, ok := u.allTrans[lang]
	if !ok {
		panic(fmt.Errorf("the specified lang: %s is not existed, please see the supported languages", lang))
	}

	value, ok := t[key]
	if !ok {
		panic(fmt.Errorf("the specified key: %s is not existed", key))
	}

	return value
}
