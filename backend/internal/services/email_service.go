package services

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log/slog"

	"github.com/buildboard/backend/config"
	"github.com/wneessen/go-mail"
)

// EmailService handles sending emails via SMTP.
type EmailService struct {
	config *config.Config
	client *mail.Client
}

// EmailData contains the data needed to send an email.
type EmailData struct {
	To      string
	Subject string
	HTML    string
	Text    string // Optional plain text version
}

// NewEmailService creates a new email service instance.
// If SMTP credentials are not configured, it returns a service with nil client
// that will skip email sending.
func NewEmailService(cfg *config.Config) (*EmailService, error) {
	if cfg.SMTPUsername == "" || cfg.SMTPPassword == "" {
		slog.Info("SMTP credentials not configured, email service disabled")
		return &EmailService{config: cfg}, nil
	}

	client, err := mail.NewClient(cfg.SMTPHost,
		mail.WithPort(cfg.SMTPPort),
		mail.WithSMTPAuth(mail.SMTPAuthPlain),
		mail.WithUsername(cfg.SMTPUsername),
		mail.WithPassword(cfg.SMTPPassword),
		mail.WithTLSPolicy(mail.TLSMandatory),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create email client: %w", err)
	}

	return &EmailService{
		config: cfg,
		client: client,
	}, nil
}

// SendEmail sends an email with the provided data.
func (es *EmailService) SendEmail(data EmailData) error {
	if es.client == nil {
		slog.Warn("Skipping email send (not configured)", "subject", data.Subject, "to", data.To)
		return fmt.Errorf("email service not configured")
	}

	m := mail.NewMsg()

	// Set sender
	if err := m.From(fmt.Sprintf("%s <%s>", es.config.FromName, es.config.FromEmail)); err != nil {
		return fmt.Errorf("failed to set sender: %w", err)
	}

	// Set recipient
	if err := m.To(data.To); err != nil {
		return fmt.Errorf("failed to set recipient: %w", err)
	}

	// Set subject
	m.Subject(data.Subject)

	// Set plain text as the primary body
	if data.Text != "" {
		m.SetBodyString(mail.TypeTextPlain, data.Text)
	}

	// Set HTML as alternative (preferred by most email clients)
	m.AddAlternativeString(mail.TypeTextHTML, data.HTML)

	// Send the email using DialAndSend
	if err := es.client.DialAndSend(m); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	slog.Info("Email sent", "subject", data.Subject, "to", data.To)
	return nil
}

// GenerateEngagementToken creates a secure random token for email engagement tracking.
// Returns a URL-safe base64 encoded string of 32 random bytes (43 characters).
func GenerateEngagementToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("failed to generate random token: %w", err)
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// SendOTPEmail sends a one-time password to the user for verification.
// The trackingToken parameter is used to track email opens via a 1x1 pixel.
func (es *EmailService) SendOTPEmail(email, otp, trackingToken string) error {
	html := fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
  <meta http-equiv="Content-Type" content="text/html; charset=utf-8">
</head>
<body>
<div class="wrapper" style='font-family: "system-ui",-apple-system,BlinkMacSystemFont,Segoe UI,Roboto,sans-serif; max-width: 600px; margin: 0 auto; padding: 1rem;'>
  <div class="container" style="max-width: 100%%; width: 100%%; margin: 0; padding: 0;">
    <table>
      <tbody>
        <tr>
          <td>
            <div class="section" style="padding: .5rem 1rem;">
              <p>Hi there,</p>
              <p>Thanks for signing up for early access to BuildBoard! Please use the verification code below to complete your signup:</p>
              <pre style="text-align: center; background-color: #ebebeb; font-size: 1.5em; border-radius: 4px; padding: 8px 0;"><b>%s</b></pre>
              <p>This code will expire in 15 minutes.</p>
              <p>Tip: you can triple-click the box to copy-paste the whole thing.</p>
              <p>- The BuildBoard Team</p>
            </div>
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</div>
<div style="height:1px;background:url('%s/api/e/%s.png')"></div>
</body>
</html>`, otp, es.config.BackendURL, trackingToken)

	plainText := fmt.Sprintf("Hi there,\n\nThanks for signing up for early access to BuildBoard! Please use the verification code below to complete your signup:\n\n%s\n\nThis code will expire in 15 minutes.\n\nTip: you can triple-click the box to copy-paste the whole thing.\n\n- The BuildBoard Team", otp)

	return es.SendEmail(EmailData{
		To:      email,
		Subject: "Your BuildBoard Verification Code",
		HTML:    html,
		Text:    plainText,
	})
}
