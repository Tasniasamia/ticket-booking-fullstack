package email
import (
	"crypto/tls"
	"fmt"
	"net"
	"net/smtp"
	"time"
	"ticketBooking/internal/config"
	"github.com/resend/resend-go/v3"
)
type EmailService struct {
	config config.Config
	env string
}
func NewEmailService(cfg config.Config) *EmailService {
	return &EmailService{
		config: cfg,
	}
}
func (s *EmailService) SendOTP(toEmail, otp, reason string) error {
	subject := "Your Verification OTP"
	body := fmt.Sprintf("Your OTP is: %s\n\nThis OTP is valid for 2 minutes.\nDo not share this code with anyone.\n\nReason: %s", otp, reason)
	return s.sendWithGmail(toEmail, subject, body)
}
func (s *EmailService) Send(toEmail, subject, body string) error {
	return s.sendWithGmail(toEmail, subject, body)
}
func (s *EmailService) sendWithGmail(to, subject, body string) error {
	from := s.config.SMTPEmail
	password := s.config.SMTPPassword
	host := s.config.SMTPHost
	port := s.config.SMTPPort
	if from == "" || password == "" {
		return fmt.Errorf("SMTP credentials not set")
	}
	if host == "" {
		host = "smtp.gmail.com"
	}
	if port == "" {
		port = "587"
	}
	msg := []byte(
		"From: " + from + "\r\n" +
			"To: " + to + "\r\n" +
			"Subject: " + subject + "\r\n" +
			"MIME-Version: 1.0\r\n" +
			"Content-Type: text/plain; charset=\"utf-8\"\r\n" +
			"\r\n" + body,
	)
	auth := smtp.PlainAuth("", from, password, host)
	addr := host + ":" + port
	// If port is 465, use direct SSL/TLS
	if port == "465" {
		return sendMailPort465(addr, auth, from, []string{to}, msg, host)
	}
	// Default to Port 587 (STARTTLS) with timeout
	err := sendMailPort587(addr, auth, from, []string{to}, msg, host)
	if err != nil {
		// If port 587 fails (e.g. blocked), try port 465 as fallback
		fallbackAddr := host + ":465"
		if fbErr := sendMailPort465(fallbackAddr, auth, from, []string{to}, msg, host); fbErr == nil {
			fmt.Println("OTP sent via Gmail SMTP (port 465 fallback) to:", to)
			return nil
		}
		return fmt.Errorf("gmail smtp error: %w", err)
	}
	fmt.Println("OTP sent via Gmail SMTP to:", to)
	return nil
}
func sendMailPort587(addr string, auth smtp.Auth, from string, to []string, msg []byte, host string) error {
	dialer := &net.Dialer{Timeout: 15 * time.Second}
	conn, err := dialer.Dial("tcp", addr)
	if err != nil {
		return fmt.Errorf("tcp dial error: %w", err)
	}
	defer conn.Close()
	client, err := smtp.NewClient(conn, host)
	if err != nil {
		return fmt.Errorf("smtp client error: %w", err)
	}
	defer client.Quit()
	tlsConfig := &tls.Config{
		InsecureSkipVerify: false,
		ServerName: host,
	}
	if err := client.StartTLS(tlsConfig); err != nil {
		return fmt.Errorf("starttls error: %w", err)
	}
	if err := client.Auth(auth); err != nil {
		return fmt.Errorf("smtp auth error: %w", err)
	}
	if err := client.Mail(from); err != nil {
		return fmt.Errorf("smtp mail from error: %w", err)
	}
	for _, recipient := range to {
		if err := client.Rcpt(recipient); err != nil {
			return fmt.Errorf("smtp rcpt error: %w", err)
		}
	}
	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("smtp data error: %w", err)
	}
	if _, err := w.Write(msg); err != nil {
		return fmt.Errorf("smtp write error: %w", err)
	}
	if err := w.Close(); err != nil {
		return fmt.Errorf("smtp close error: %w", err)
	}
	return nil
}
func sendMailPort465(addr string, auth smtp.Auth, from string, to []string, msg []byte, host string) error {
	dialer := &net.Dialer{Timeout: 15 * time.Second}
	tlsConfig := &tls.Config{
		InsecureSkipVerify: false,
		ServerName: host,
	}
	conn, err := tls.DialWithDialer(dialer, "tcp", addr, tlsConfig)
	if err != nil {
		return fmt.Errorf("tls dial error: %w", err)
	}
	defer conn.Close()
	client, err := smtp.NewClient(conn, host)
	if err != nil {
		return fmt.Errorf("smtp client error: %w", err)
	}
	defer client.Quit()
	if err := client.Auth(auth); err != nil {
		return fmt.Errorf("smtp auth error: %w", err)
	}
	if err := client.Mail(from); err != nil {
		return fmt.Errorf("smtp mail from error: %w", err)
	}
	for _, recipient := range to {
		if err := client.Rcpt(recipient); err != nil {
			return fmt.Errorf("smtp rcpt error: %w", err)
		}
	}
	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("smtp data error: %w", err)
	}
	if _, err := w.Write(msg); err != nil {
		return fmt.Errorf("smtp write error: %w", err)
	}
	if err := w.Close(); err != nil {
		return fmt.Errorf("smtp close error: %w", err)
	}
	return nil
}
func (s *EmailService) sendWithResend(to, subject, body string) error {
	apiKey := s.config.ResendAPIKey
	from := s.config.ResendFromEmail
	if apiKey == "" || from == "" {
		return fmt.Errorf("Resend API key or from email not set")
	}
	client := resend.NewClient(apiKey)
	params := &resend.SendEmailRequest{
		From: from,
		To: []string{to},
		Subject: subject,
		Text: body,
	}
	sent, err := client.Emails.Send(params)
	if err != nil {
		return fmt.Errorf("resend error: %w", err)
	}
	fmt.Println("OTP sent via Resend. Email ID:", sent.Id)
	return nil
}
