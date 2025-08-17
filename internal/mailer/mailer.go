package mailer

import (
	"strconv"
	"time"

	com "github.com/ahmad-khatib0-org/megacommerce-proto/gen/go/common/v1"
	"github.com/ahmad-khatib0-org/megacommerce-user/internal/store"
	"github.com/ahmad-khatib0-org/megacommerce-user/pkg/utils"
	"github.com/k3a/html2text"
	"github.com/throttled/throttled/v2"
	"github.com/vanng822/go-premailer/premailer"
	mail "github.com/xhit/go-simple-mail/v2"
)

type Mailer struct {
	config                  func() *com.Config
	store                   store.UsersStore
	templateContainer       *TemplateContainer
	perHourEmailRateLimiter *throttled.GCRARateLimiterCtx
	perDayEmailRateLimiter  *throttled.GCRARateLimiterCtx
	EmailBatching           *EmailBatchingJob
}

type MailerArgs struct {
	ConfigFn          func() *com.Config
	Store             store.UsersStore
	TemplateContainer *TemplateContainer
}

type mailData struct {
	from    string
	to      string
	subject string
	replyTo string
	body    string
	files   []*mail.File
}

func NewMailer(ma *MailerArgs) MailerService {
	return &Mailer{config: ma.ConfigFn, store: ma.Store, templateContainer: ma.TemplateContainer}
}

func (m *Mailer) send(md *mailData) error {
	server := mail.NewSMTPClient()

	port, _ := strconv.Atoi(m.config().Email.GetSmtpPort())

	server.Host = m.config().Email.GetSmtpServer()
	server.Port = port
	server.ConnectTimeout = time.Duration(m.config().Email.GetSmtpServerTimeout()) * time.Second
	server.SendTimeout = time.Duration(m.config().Email.GetSmtpServerTimeout()) * time.Second
	server.KeepAlive = false
	server.Password = m.config().Email.GetSmtpPassword()
	server.Username = m.config().Email.GetSmtpUsername()
	server.Encryption = getEncryptionType(m.config().GetEmail().GetConnectionSecurity())

	client, err := server.Connect()
	if err != nil {
		return err
	}

	email := mail.NewMSG()
	if md.replyTo == "" {
		email.SetReplyTo(md.replyTo)
	}

	for _, file := range md.files {
		email.Attach(file)
	}

	html, err := m.inlineCSS(md.body)
	if err != nil {
		return err
	}

	plain := html2text.HTML2Text(html)
	if md.from == "" {
		md.from = m.config().Email.GetFeedbackEmail()
	}

	email.SetFrom(md.from)
	email.AddTo(md.to)
	email.SetBody(mail.TextPlain, plain)
	email.AddAlternative(mail.TextHTML, html)
	email.SetSubject(md.subject)
	email.SetDate(utils.EmailDateHeader(time.Now()))

	return email.Send(client)
}

// inlineCSS takes an email string (html) and returns the same string but with
// injecting the email styles inline to be compatible with most email sender providers
func (m *Mailer) inlineCSS(tmp string) (string, error) {
	opt := premailer.Options{
		RemoveClasses:     false,
		CssToAttributes:   false,
		KeepBangImportant: true,
	}

	prem, err := premailer.NewPremailerFromString(tmp, &opt)
	if err != nil {
		return "", err
	}

	htm, err := prem.Transform()
	if err != nil {
		return "", err
	}

	return htm, nil
}

func getEncryptionType(t string) mail.Encryption {
	switch t {
	case "tls":
		return mail.EncryptionSTARTTLS
	case "ssl":
		return mail.EncryptionSSLTLS
	case "none":
		return mail.EncryptionNone
	default:
		return mail.EncryptionNone
	}
}

func (m *Mailer) Store() store.UsersStore {
	return m.store
}

func (m *Mailer) SetStore(st store.UsersStore) {
	m.store = st
}

func (m *Mailer) GetPerDayEmailRateLimiter() *throttled.GCRARateLimiterCtx {
	return m.perDayEmailRateLimiter
}

func (m *Mailer) GetPerHourEmailRateLimiter() *throttled.GCRARateLimiterCtx {
	return m.perHourEmailRateLimiter
}

type MailerService interface {
	GetPerDayEmailRateLimiter() *throttled.GCRARateLimiterCtx
	GetPerHourEmailRateLimiter() *throttled.GCRARateLimiterCtx
	SendVerifyEmail(lang, email, token, tokenID string, hours int) error
	InitEmailBatching()
}
