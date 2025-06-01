package utils

import (
	"fmt"

	"google.golang.org/grpc/codes"
)

func (u *Utils) initTrans(ua *UtilsArgs) *AppError {
	trans := make(map[string]map[string]string, len(ua.AllTrans))
	for lang, langTrans := range ua.AllTrans {
		t := make(map[string]string, len(langTrans.Trans))
		for _, v := range langTrans.Trans {
			trID := v.GetId()
			if trID == "" {
				return NewAppError("user.utils.NewUtils", "empty_translation_key", nil, fmt.Sprintf("encountered an empty translation key, key: %s value: %s", trID, v.Tr), int(codes.Internal))
			}
			t[trID] = v.Tr
		}
		trans[lang] = t
	}

	u.allTrans = trans
	return nil
}
