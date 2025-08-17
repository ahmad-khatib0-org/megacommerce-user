// Package mailer sends emails
package mailer

import (
	"fmt"

	"github.com/ahmad-khatib0-org/megacommerce-user/pkg/models"
)

func (m *Mailer) SendVerifyEmail(lang, email, token, tokenID string, hours int) error {
	td, err := m.NewTemplateData(lang)
	if err != nil {
		return err
	}

	title, err := models.Tr(lang, "templates.verify.title", map[string]any{"SiteName": m.config().GetMain().GetSiteName()})
	if err != nil {
		return err
	}
	welcome, err := models.Tr(lang, "templates.verify.welcome", map[string]any{"SiteName": m.config().GetMain().GetSiteName()})
	if err != nil {
		return err
	}
	welcome2, err := models.Tr(lang, "templates.verify.part1", map[string]any{"SiteName": m.config().GetMain().GetSiteName()})
	if err != nil {
		return err
	}
	click, err := models.Tr(lang, "templates.verify.part2", nil)
	if err != nil {
		return err
	}
	redirect, err := models.Tr(lang, "templates.verify.part3", map[string]any{"SiteName": m.config().GetMain().GetSiteName()})
	if err != nil {
		return err
	}
	note, err := models.Tr(lang, "templates.verify.part4", map[string]any{"Hours": hours})
	if err != nil {
		return err
	}

	td.Props["Title"] = title
	td.Props["Welcome"] = welcome
	td.Props["Welcome2"] = welcome2
	td.Props["Click"] = click
	td.Props["Redirect"] = redirect
	td.Props["Note"] = note
	td.Props["Url"] = fmt.Sprintf("%s?token=%s&token_id=%s&email=%s", m.config().Security.GetEmailConfirmationUrl(), token, tokenID, email)

	body, err := m.templateContainer.RenderToString("verify_email", *td)
	if err != nil {
		return err
	}

	return m.send(&mailData{to: email, subject: title, body: body})
}
