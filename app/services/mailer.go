package services

import (
	"bytes"
	"context"
	"github.com/armanjr/termustat/app/config"
	"github.com/mailgun/errors"
	"html/template"
	"path/filepath"
	"strings"
	"sync"

	"github.com/mailgun/mailgun-go/v4"
)

type EmailTemplate struct {
	Subject string
	Body    string
}

type Mailer struct {
	mg      *mailgun.MailgunImpl
	sender  string
	tplPath string
}

var (
	templateCache = make(map[string]*template.Template)
	templateMutex = &sync.RWMutex{}
)

func NewMailer(cfg *config.Config) *Mailer {
	mg := mailgun.NewMailgun(
		cfg.MailgunDomain,
		cfg.MailgunAPIKey,
	)
	return &Mailer{
		mg:      mg,
		sender:  "noreply@" + cfg.MailgunDomain,
		tplPath: "templates/email/",
	}
}

func (m *Mailer) SendEmail(to, subject, body string) error {
	ctx := context.Background()
	message := mailgun.NewMessage(m.sender, subject, "", to)
	message.SetHTML(body)
	_, _, err := m.mg.Send(ctx, message)
	return err
}

func (m *Mailer) RenderTemplate(tplName string, data interface{}) (*EmailTemplate, error) {
	if tplName == "" {
		return nil, errors.New("template name cannot be empty")
	}

	templateMutex.RLock()
	cachedTpl, exists := templateCache[tplName]
	templateMutex.RUnlock()

	var tpl *template.Template
	var err error

	if !exists {
		tplPath := filepath.Join(m.tplPath, tplName)

		tpl, err = template.New(filepath.Base(tplName)).
			Funcs(template.FuncMap{
				"safeHTML": func(s string) template.HTML {
					return template.HTML(s)
				},
			}).
			ParseFiles(tplPath)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to parse template %s", tplName)
		}

		templateMutex.Lock()
		templateCache[tplName] = tpl
		templateMutex.Unlock()
	} else {
		tpl = cachedTpl
	}

	var subjectBuf, bodyBuf bytes.Buffer

	if err := tpl.ExecuteTemplate(&subjectBuf, "subject", data); err != nil {
		return nil, errors.Wrapf(err, "failed to execute subject template %s", tplName)
	}

	if err := tpl.ExecuteTemplate(&bodyBuf, "body", data); err != nil {
		return nil, errors.Wrapf(err, "failed to execute body template %s", tplName)
	}

	subject := strings.TrimSpace(subjectBuf.String())
	body := strings.TrimSpace(bodyBuf.String())

	return &EmailTemplate{
		Subject: subject,
		Body:    body,
	}, nil
}
