package mail

import (
	"backend/configs"
	"bytes"
	"fmt"
	"html/template"
	"mime"
	"net/smtp"
	"net/url"
	"path/filepath"
	"strings"
)

type Service struct {
	smtpHost    string
	smtpPort    string
	smtpUser    string
	smtpPass    string
	templateDir string
	verifyURL   string
}

// NewService khởi tạo Service với cấu hình từ configs.Config
func NewService(cfg *configs.Config) *Service {
	templateDir := "templates"
	if strings.TrimSpace(templateDir) == "" {
		templateDir = "templates"
	}

	return &Service{
		smtpHost:    strings.TrimSpace(cfg.MailHost),
		smtpPort:    strings.TrimSpace(cfg.MailPort),
		smtpUser:    strings.TrimSpace(cfg.MailUser),
		smtpPass:    cfg.MailPass,
		templateDir: templateDir,
		verifyURL:   strings.TrimSpace(cfg.VerifyURL),
	}
}

func (s *Service) BuildVerificationURL(email, code string) string {
	if s.verifyURL == "" {
		return fmt.Sprintf("/api/v1/auth/verify?email=%s&code=%s", url.QueryEscape(email), url.QueryEscape(code))
	}

	sep := "?"
	if strings.Contains(s.verifyURL, "?") {
		sep = "&"
	}

	return s.verifyURL + sep + "email=" + url.QueryEscape(email) + "&code=" + url.QueryEscape(code)
}

// SendVerificationEmail gửi email xác thực tài khoản với link chứa code xác thực
func (s *Service) SendVerificationEmail(to, firstName, verifyURL string) error {
	body, err := s.renderTemplate("verification.html", VerificationEmailData{
		FirstName: firstName,
		VerifyURL: verifyURL,
	})
	if err != nil {
		return err
	}

	return s.sendHTML(to, "Xac thuc tai khoan Smart Parking", body)
}

// SendPasswordResetEmail gửi email đặt lại mật khẩu với link chứa code reset
func (s *Service) SendPasswordResetEmail(to, firstName, code string) error {
	body, err := s.renderTemplate("reset-password.html", ResetPasswordEmailData{
		FirstName: firstName,
		CodeID:    code,
	})
	if err != nil {
		return err
	}

	return s.sendHTML(to, "Dat lai mat khau Smart Parking", body)
}

// RenderVerifiedPage trả về HTML của trang đã xác thực, có thể dùng để hiển thị sau khi user click link xác thực email
func (s *Service) RenderVerifiedPage(year int) (string, error) {
	return s.renderTemplate("verified.html", VerifiedPageData{Year: year})
}

// RenderTemplate đọc file template, thực thi với data và trả về kết quả HTML
func (s *Service) renderTemplate(firstName string, data any) (string, error) {
	tplPath := filepath.Join(s.templateDir, firstName)

	tpl, err := template.ParseFiles(tplPath)
	if err != nil {
		return "", fmt.Errorf("parse template %s: %w", firstName, err)
	}

	var buf bytes.Buffer
	if err := tpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("execute template %s: %w", firstName, err)
	}

	return buf.String(), nil
}

// SendHTML xây dựng email với header và body HTML, sau đó gửi qua SMTP
func (s *Service) sendHTML(to, subject, htmlBody string) error {
	if strings.TrimSpace(to) == "" {
		return fmt.Errorf("recipient email is required")
	}
	if s.smtpHost == "" || s.smtpPort == "" || s.smtpUser == "" || s.smtpPass == "" {
		return fmt.Errorf("smtp credentials are not configured")
	}

	from := s.smtpUser
	headers := map[string]string{
		"From":         from,
		"To":           to,
		"Subject":      mime.BEncoding.Encode("UTF-8", subject),
		"MIME-Version": "1.0",
		"Content-Type": "text/html; charset=UTF-8",
	}

	var msg strings.Builder
	for k, v := range headers {
		msg.WriteString(k)
		msg.WriteString(": ")
		msg.WriteString(v)
		msg.WriteString("\r\n")
	}
	msg.WriteString("\r\n")
	msg.WriteString(htmlBody)

	addr := s.smtpHost + ":" + s.smtpPort
	auth := smtp.PlainAuth("", s.smtpUser, s.smtpPass, s.smtpHost)

	if err := smtp.SendMail(addr, auth, from, []string{to}, []byte(msg.String())); err != nil {
		return fmt.Errorf("send email failed: %w", err)
	}

	return nil
}
