package server

import (
	commonPb "github.com/ahmad-khatib0-org/megacommerce-proto/gen/go/common/v1"
)

func (a *App) initTrans() map[string]*commonPb.TranslationElements {
	trans, err := a.commonClient.TranslationsGet()
	if err != nil {
		a.done <- err
	}

	return trans
}
