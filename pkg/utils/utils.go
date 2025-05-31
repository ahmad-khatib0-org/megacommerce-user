package utils

import v1 "github.com/ahmad-khatib0-org/megacommerce-proto/gen/go/common/v1"

type Utils struct {
	Tr []v1.TranslationElement
}

type UtilsArgs struct {
	Tr []v1.TranslationElement
}

func NewUtils(ua *UtilsArgs) *Utils {
	return &Utils{
		Tr: ua.Tr,
	}
}
