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

	title := models.Tr(lang, "templates.verify.title", map[string]any{"SiteName": m.config().GetMain().GetSiteName()})

	welcome := models.Tr(lang, "templates.verify.welcome", map[string]any{"SiteName": m.config().GetMain().GetSiteName()})
	welcome2 := models.Tr(lang, "templates.verify.part1", map[string]any{"SiteName": m.config().GetMain().GetSiteName()})
	click := models.Tr(lang, "templates.verify.part2", nil)
	redirect := models.Tr(lang, "templates.verify.part3", map[string]any{"SiteName": m.config().GetMain().GetSiteName()})
	note := models.Tr(lang, "templates.verify.part4", map[string]any{"Hours": hours})

	td.Props["Title"] = title
	td.Props["Welcome"] = welcome
	td.Props["Welcome2"] = welcome2
	td.Props["Click"] = click
	td.Props["Redirect"] = redirect
	td.Props["Note"] = note
	td.Props["Url"] = fmt.Sprintf("%s?token=%s&token_id=%s&email=%s", m.config().Security.GetEmailConfirmationUrl(), token, tokenID, email)

	body, err := m.templateContainer.RenderToString("verify_email", td)
	if err != nil {
		return err
	}

	return m.send(&mailData{to: email, subject: title, body: body})
}

func (m *Mailer) SendPasswordResetEmail(lang, email, token, tokenID string, hours int) error {
	td, err := m.NewTemplateData(lang)
	if err != nil {
		return err
	}

	title := models.Tr(lang, "templates.reset_password.title", map[string]any{"SiteName": m.config().GetMain().GetSiteName()})
	welcome := models.Tr(lang, "templates.welcome", map[string]any{"SiteName": m.config().GetMain().GetSiteName()})
	received := models.Tr(lang, "templates.reset_password.part1", nil)
	click := models.Tr(lang, "templates.click_on_link", nil)
	redirect := models.Tr(lang, "templates.reset_password.part2", map[string]any{"SiteName": m.config().GetMain().GetSiteName()})
	note := models.Tr(lang, "templates.reset_password.part3", map[string]any{"Hours": hours})

	td.Props["Title"] = title
	td.Props["Welcome"] = welcome
	td.Props["Received"] = received
	td.Props["Click"] = click
	td.Props["Redirect"] = redirect
	td.Props["Note"] = note
	td.Props["Url"] = fmt.Sprintf("%s?token=%s&token_id=%s&email=%s", m.config().Security.GetPasswordResetUrl(), token, tokenID, email)

	body, err := m.templateContainer.RenderToString("password_reset_email", td)
	if err != nil {
		return err
	}

	return m.send(&mailData{to: email, subject: title, body: body})
}
