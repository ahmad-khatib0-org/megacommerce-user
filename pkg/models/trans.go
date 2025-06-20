package models

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"regexp"

	pb "github.com/ahmad-khatib0-org/megacommerce-proto/gen/go/common/v1"
)

const (
	NoTranslation string = "<untranslated>"
)

type transStoreElement struct {
	Temp    *template.Template
	HasVars bool
}

var (
	transStore         map[string]map[string]*transStoreElement
	ErrTrMissingParams = errors.New("translation requires params to be passed, but none received")
)

type TranslateFunc = func(lang, transId string, params map[string]any) (string, error)

func TranslationsInit(translations map[string]*pb.TranslationElements) error {
	trans := make(map[string]map[string]*transStoreElement, len(translations))

	for lang, langTrans := range translations {
		t := make(map[string]*transStoreElement, len(langTrans.Trans))
		for _, v := range langTrans.Trans {
			trID := v.GetId()
			if trID == "" {
				return fmt.Errorf("encountered an empty translation key, key: %s value: %s", trID, v.Tr)
			}

			trVal := v.GetTr()
			if trVal == "" {
				return fmt.Errorf("encountered an empty translation value, value: %s key: %s", trVal, trID)
			}

			goTemp, hasVars := parseTranslateIdVars(trVal)
			tmp, err := template.New("msg").Parse(goTemp)
			if err != nil {
				return fmt.Errorf("an error occurred while trying to parse a translation template %v", err)
			}
			t[trID] = &transStoreElement{Temp: tmp, HasVars: hasVars}
		}
		trans[lang] = t
	}

	transStore = trans
	return nil
}

// Tr translate a given message id with the passed params (if passed)
func Tr(lang string, id string, params map[string]any) (string, error) {
	if transStore == nil {
		panic("trans keys are not initialized, call models.TranslationsInit on server init")
	}

	value, ok := transStore[lang][id]
	if !ok {
		panic(fmt.Errorf("the specified key: %s is not existed", id))
	}

	var buf bytes.Buffer
	if value.HasVars && len(params) == 0 {
		panic("this translation id: %s, must has params associated with it")
	}

	if err := value.Temp.Execute(&buf, params); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// parseTranslateIdVars converts `{{Var}}` → `{{.Var}}` and returns if vars were present.
func parseTranslateIdVars(id string) (string, bool) {
	re := regexp.MustCompile(`{{\s*([a-zA-Z0-9_]+)\s*}}`)
	has := re.MatchString(id)
	result := re.ReplaceAllString(id, "{{.$1}}")
	return result, has
}
