package email

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"mime"
	"net"
	"net/smtp"
	"path/filepath"
	"strings"
	"time"

	"deepspace/internal/config"

	"github.com/redis/go-redis/v9"
)

const (
	defaultDialTimeout = 10 * time.Second
)

var (
	ErrEmailDisabled     = errors.New("email disabled")
	ErrQueueUnavailable  = errors.New("queue unavailable")
	ErrInvalidRecipients = errors.New("invalid recipients")
	ErrMissingSubject    = errors.New("missing subject")
	ErrInvalidEmailType  = errors.New("invalid email type")
	ErrTemplateOverride  = errors.New("template override not allowed")
)

const (
	EmailTypeWelcome       = "welcome"
	EmailTypeResetPassword = "reset_password"
)

type Service struct {
	cfg       *config.Config
	redis     *redis.Client
	templates *template.Template
}

type EmailInput struct {
	Type         string            `json:"type"`
	To           []string          `json:"to"`
	Subject      string            `json:"subject"`
	HTMLTemplate string            `json:"html_template"`
	TemplateData map[string]any    `json:"template_data"`
	HTML         string            `json:"html"`
	Text         string            `json:"text"`
	ReplyTo      string            `json:"reply_to"`
	Headers      map[string]string `json:"headers"`
}

type QueueItem struct {
	Email EmailInput `json:"email"`
}

func New(cfg *config.Config) (*Service, error) {
	if cfg == nil {
		return nil, errors.New("config missing")
	}

	svc := &Service{cfg: cfg}

	if strings.TrimSpace(cfg.RedisURL) != "" {
		opt, err := redis.ParseURL(cfg.RedisURL)
		if err != nil {
			return nil, err
		}
		svc.redis = redis.NewClient(opt)
	}

	if strings.TrimSpace(cfg.EmailTemplateDir) != "" {
		patterns := []string{
			filepath.Join(cfg.EmailTemplateDir, "*.html"),
			filepath.Join(cfg.EmailTemplateDir, "**", "*.html"),
		}
		parsed := template.New("email")
		for _, pattern := range patterns {
			matches, _ := filepath.Glob(pattern)
			if len(matches) == 0 {
				continue
			}
			if _, err := parsed.ParseFiles(matches...); err != nil {
				return nil, err
			}
		}
		if len(parsed.Templates()) > 0 {
			svc.templates = parsed
		}
	}

	return svc, nil
}

func (s *Service) Send(ctx context.Context, input EmailInput) error {
	if s == nil || s.cfg == nil || !s.cfg.EmailEnabled {
		return ErrEmailDisabled
	}
	if !isValidEmailType(input.Type) {
		return ErrInvalidEmailType
	}
	if len(input.To) == 0 {
		return ErrInvalidRecipients
	}
	if strings.TrimSpace(input.Subject) == "" {
		return ErrMissingSubject
	}
	if strings.TrimSpace(input.HTML) != "" || strings.TrimSpace(input.HTMLTemplate) != "" {
		return ErrTemplateOverride
	}

	payload := input
	templateName := templateNameForType(payload.Type)
	if templateName == "" {
		return ErrInvalidEmailType
	}
	if payload.HTML == "" && templateName != "" {
		payload.HTMLTemplate = templateName
		htmlContent, err := s.RenderTemplate(payload.HTMLTemplate, payload.TemplateData)
		if err != nil {
			return err
		}
		payload.HTML = htmlContent
	}

	message, err := s.buildMessage(payload)
	if err != nil {
		return err
	}

	return s.sendSMTP(ctx, payload.To, message)
}

func isValidEmailType(value string) bool {
	switch strings.TrimSpace(value) {
	case EmailTypeWelcome, EmailTypeResetPassword:
		return true
	default:
		return false
	}
}

func templateNameForType(value string) string {
	switch strings.TrimSpace(value) {
	case EmailTypeWelcome:
		return "welcome.html"
	case EmailTypeResetPassword:
		return "reset-password.html"
	default:
		return ""
	}
}

func (s *Service) EnqueueBatch(ctx context.Context, inputs []EmailInput) error {
	if s == nil || s.redis == nil {
		return ErrQueueUnavailable
	}
	if len(inputs) == 0 {
		return nil
	}
	items := make([]any, 0, len(inputs))
	for _, input := range inputs {
		encoded, err := json.Marshal(QueueItem{Email: input})
		if err != nil {
			return err
		}
		items = append(items, string(encoded))
	}

	_, err := s.redis.LPush(ctx, s.cfg.RedisQueueKey, items...).Result()
	return err
}

func (s *Service) RenderTemplate(name string, data map[string]any) (string, error) {
	if s == nil || s.templates == nil {
		return "", fmt.Errorf("template not loaded")
	}

	var buf bytes.Buffer
	if data == nil {
		data = map[string]any{}
	}
	if err := s.templates.ExecuteTemplate(&buf, name, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func (s *Service) buildMessage(input EmailInput) ([]byte, error) {
	from := s.formatFrom()
	if strings.TrimSpace(from) == "" {
		return nil, errors.New("from address required")
	}

	headers := map[string]string{
		"From":         from,
		"To":           strings.Join(input.To, ", "),
		"Subject":      mime.QEncoding.Encode("utf-8", input.Subject),
		"MIME-Version": "1.0",
	}
	if input.ReplyTo != "" {
		headers["Reply-To"] = input.ReplyTo
	}
	for key, value := range input.Headers {
		headers[key] = value
	}

	var body bytes.Buffer
	if input.HTML != "" && input.Text != "" {
		boundary := fmt.Sprintf("alt_%d", time.Now().UnixNano())
		headers["Content-Type"] = "multipart/alternative; boundary=" + boundary
		writeHeaders(&body, headers)
		body.WriteString("\r\n")
		body.WriteString("--" + boundary + "\r\n")
		body.WriteString("Content-Type: text/plain; charset=utf-8\r\n\r\n")
		body.WriteString(input.Text)
		body.WriteString("\r\n--" + boundary + "\r\n")
		body.WriteString("Content-Type: text/html; charset=utf-8\r\n\r\n")
		body.WriteString(input.HTML)
		body.WriteString("\r\n--" + boundary + "--\r\n")
		return body.Bytes(), nil
	}

	if input.HTML != "" {
		headers["Content-Type"] = "text/html; charset=utf-8"
		writeHeaders(&body, headers)
		body.WriteString("\r\n")
		body.WriteString(input.HTML)
		return body.Bytes(), nil
	}

	headers["Content-Type"] = "text/plain; charset=utf-8"
	writeHeaders(&body, headers)
	body.WriteString("\r\n")
	body.WriteString(input.Text)
	return body.Bytes(), nil
}

func (s *Service) sendSMTP(ctx context.Context, recipients []string, msg []byte) error {
	addr := net.JoinHostPort(s.cfg.SMTPHost, fmt.Sprintf("%d", s.cfg.SMTPPort))
	var client *smtp.Client
	var err error

	if s.cfg.SMTPUseTLS {
		dialer := &net.Dialer{Timeout: defaultDialTimeout}
		conn, dialErr := tls.DialWithDialer(dialer, "tcp", addr, &tls.Config{ServerName: s.cfg.SMTPHost})
		if dialErr != nil {
			return dialErr
		}
		client, err = smtp.NewClient(conn, s.cfg.SMTPHost)
	} else {
		client, err = smtp.Dial(addr)
	}
	if err != nil {
		return err
	}
	defer func() {
		_ = client.Quit()
	}()

	if s.cfg.SMTPUser != "" {
		auth := smtp.PlainAuth("", s.cfg.SMTPUser, s.cfg.SMTPPassword, s.cfg.SMTPHost)
		if err := client.Auth(auth); err != nil {
			return err
		}
	}

	if err := client.Mail(s.cfg.EmailFromAddress); err != nil {
		return err
	}
	for _, to := range recipients {
		if err := client.Rcpt(to); err != nil {
			return err
		}
	}

	writer, err := client.Data()
	if err != nil {
		return err
	}
	_, err = writer.Write(msg)
	if err != nil {
		_ = writer.Close()
		return err
	}
	if err := writer.Close(); err != nil {
		return err
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return nil
	}
}

func (s *Service) formatFrom() string {
	name := strings.TrimSpace(s.cfg.EmailFromName)
	addr := strings.TrimSpace(s.cfg.EmailFromAddress)
	if addr == "" {
		return ""
	}
	if name == "" {
		return addr
	}
	return fmt.Sprintf("%s <%s>", mime.QEncoding.Encode("utf-8", name), addr)
}

func writeHeaders(buf *bytes.Buffer, headers map[string]string) {
	keys := make([]string, 0, len(headers))
	for key := range headers {
		keys = append(keys, key)
	}
	for _, key := range keys {
		buf.WriteString(key)
		buf.WriteString(": ")
		buf.WriteString(headers[key])
		buf.WriteString("\r\n")
	}
}
