package mailer

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/armanjr/termustat/api/models"
	"github.com/mailgun/mailgun-go/v4"
	"go.uber.org/zap"
)

// EmailTemplate holds the rendered subject and body.
type EmailTemplate struct {
	Subject string
	Body    string
}

// Mailer defines the interface for sending emails.
type Mailer interface {
	// SendEmail sends an HTML email.
	SendEmail(to, subject, body string) error
	// RenderTemplate loads and renders the specified email template.
	RenderTemplate(tplName string, data interface{}) (*EmailTemplate, error)
	// SendVerificationEmail sends an email using the verification template.
	SendVerificationEmail(user *models.User, token string) error
	// SendPasswordResetEmail sends a password reset email.
	SendPasswordResetEmail(user *models.User, resetToken string) error
}

// MailerConfig holds configuration values for the mailer.
type MailerConfig struct {
	Domain  string
	APIKey  string
	Sender  string
	TplPath string
}

type mailerImpl struct {
	mg            *mailgun.MailgunImpl
	sender        string
	tplPath       string
	logger        *zap.Logger
	templateCache map[string]*template.Template
	cacheMutex    sync.RWMutex
}

// NewMailer creates a new Mailer implementation using Mailgun.
func NewMailer(cfg MailerConfig, logger *zap.Logger) Mailer {
	mg := mailgun.NewMailgun(cfg.Domain, cfg.APIKey)
	return &mailerImpl{
		mg:            mg,
		sender:        cfg.Sender,
		tplPath:       cfg.TplPath,
		logger:        logger,
		templateCache: make(map[string]*template.Template),
	}
}

// SendEmail sends an HTML email using Mailgun.
func (m *mailerImpl) SendEmail(to, subject, body string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	message := m.mg.NewMessage(m.sender, subject, "", to)
	message.SetHtml(body)

	_, id, err := m.mg.Send(ctx, message)
	if err != nil {
		m.logger.Error("failed to send email",
			zap.String("to", to),
			zap.String("subject", subject),
			zap.Error(err))
		return fmt.Errorf("failed to send email to %s: %w", to, err)
	}

	m.logger.Info("email sent successfully", zap.String("id", id))
	return nil
}

// RenderTemplate loads and renders an email template from the templates directory.
func (m *mailerImpl) RenderTemplate(tplName string, data interface{}) (*EmailTemplate, error) {
	if tplName == "" {
		return nil, fmt.Errorf("template name cannot be empty")
	}

	m.cacheMutex.RLock()
	tpl, found := m.templateCache[tplName]
	m.cacheMutex.RUnlock()

	if !found {
		tplPath := filepath.Join(m.tplPath, tplName)
		var err error
		tpl, err = template.New(filepath.Base(tplName)).
			Funcs(template.FuncMap{
				"safeHTML": func(s string) template.HTML { return template.HTML(s) },
			}).
			ParseFiles(tplPath)
		if err != nil {
			m.logger.Error("failed to parse template",
				zap.String("template", tplName),
				zap.Error(err))
			return nil, fmt.Errorf("failed to parse template %s: %w", tplName, err)
		}
		m.cacheMutex.Lock()
		m.templateCache[tplName] = tpl
		m.cacheMutex.Unlock()
	}

	var subjectBuf, bodyBuf bytes.Buffer
	if err := tpl.ExecuteTemplate(&subjectBuf, "subject", data); err != nil {
		m.logger.Error("failed to render subject",
			zap.String("template", tplName),
			zap.Error(err))
		return nil, fmt.Errorf("failed to render subject template %s: %w", tplName, err)
	}
	if err := tpl.ExecuteTemplate(&bodyBuf, "body", data); err != nil {
		m.logger.Error("failed to render body",
			zap.String("template", tplName),
			zap.Error(err))
		return nil, fmt.Errorf("failed to render body template %s: %w", tplName, err)
	}

	return &EmailTemplate{
		Subject: strings.TrimSpace(subjectBuf.String()),
		Body:    strings.TrimSpace(bodyBuf.String()),
	}, nil
}

// SendVerificationEmail sends a verification email to the user.
// the auth service should pass the token or the mailer should be provided with a frontend URL to generate it.
func (m *mailerImpl) SendVerificationEmail(user *models.User, token string) error {
	verificationURL := "http://example.com/verify?token=" + token

	tpl, err := m.RenderTemplate("email/verification_email.html", struct {
		Name            string
		VerificationURL string
	}{
		Name:            user.FirstName + " " + user.LastName,
		VerificationURL: verificationURL,
	})
	if err != nil {
		return err
	}

	return m.SendEmail(user.Email, tpl.Subject, tpl.Body)
}

// SendPasswordResetEmail sends a password reset email to the user.
func (m *mailerImpl) SendPasswordResetEmail(user *models.User, resetToken string) error {
	resetURL := "http://example.com/reset?token=" + resetToken

	tpl, err := m.RenderTemplate("email/password_reset_email.html", struct {
		Name     string
		ResetURL string
	}{
		Name:     user.FirstName + " " + user.LastName,
		ResetURL: resetURL,
	})
	if err != nil {
		return err
	}

	return m.SendEmail(user.Email, tpl.Subject, tpl.Body)
}
