package services

import (
	"bytes"
	"context"
	"html/template"
	"path/filepath"
	"strings"
	"sync"

	"github.com/mailgun/errors"
	"github.com/mailgun/mailgun-go/v4"
	"go.uber.org/zap"
)

type MailerService interface {
	SendEmail(to, subject, body string) error
	RenderTemplate(tplName string, data interface{}) (*EmailTemplate, error)
}

type EmailTemplate struct {
	Subject string
	Body    string
}

type MailerConfig struct {
	Domain  string
	APIKey  string
	Sender  string
	TplPath string
}

type mailerService struct {
	mg            *mailgun.MailgunImpl
	sender        string
	tplPath       string
	logger        *zap.Logger
	templateCache map[string]*template.Template
	templateMutex *sync.RWMutex
}

func NewMailerService(config MailerConfig, logger *zap.Logger) MailerService {
	mg := mailgun.NewMailgun(config.Domain, config.APIKey)

	return &mailerService{
		mg:            mg,
		sender:        config.Sender,
		tplPath:       config.TplPath,
		logger:        logger,
		templateCache: make(map[string]*template.Template),
		templateMutex: &sync.RWMutex{},
	}
}

func (m *mailerService) SendEmail(to, subject, body string) error {
	ctx := context.Background()
	message := mailgun.NewMessage(m.sender, subject, "", to)
	message.SetHTML(body)

	_, _, err := m.mg.Send(ctx, message)
	if err != nil {
		m.logger.Error("Failed to send email",
			zap.String("to", to),
			zap.String("subject", subject),
			zap.Error(err))
		return errors.Wrapf(err, "failed to send email to %s", to)
	}

	m.logger.Info("Email sent successfully",
		zap.String("to", to),
		zap.String("subject", subject))
	return nil
}

func (m *mailerService) RenderTemplate(tplName string, data interface{}) (*EmailTemplate, error) {
	if tplName == "" {
		return nil, errors.New("template name cannot be empty")
	}

	m.templateMutex.RLock()
	cachedTpl, exists := m.templateCache[tplName]
	m.templateMutex.RUnlock()

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
			m.logger.Error("Failed to parse template",
				zap.String("template", tplName),
				zap.Error(err))
			return nil, errors.Wrapf(err, "failed to parse template %s", tplName)
		}

		m.templateMutex.Lock()
		m.templateCache[tplName] = tpl
		m.templateMutex.Unlock()
	} else {
		tpl = cachedTpl
	}

	var subjectBuf, bodyBuf bytes.Buffer

	if err := tpl.ExecuteTemplate(&subjectBuf, "subject", data); err != nil {
		m.logger.Error("Failed to execute subject template",
			zap.String("template", tplName),
			zap.Error(err))
		return nil, errors.Wrapf(err, "failed to execute subject template %s", tplName)
	}

	if err := tpl.ExecuteTemplate(&bodyBuf, "body", data); err != nil {
		m.logger.Error("Failed to execute body template",
			zap.String("template", tplName),
			zap.Error(err))
		return nil, errors.Wrapf(err, "failed to execute body template %s", tplName)
	}

	subject := strings.TrimSpace(subjectBuf.String())
	body := strings.TrimSpace(bodyBuf.String())

	return &EmailTemplate{
		Subject: subject,
		Body:    body,
	}, nil
}
