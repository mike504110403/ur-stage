package email

import "testing"

func TestSendEmail(t *testing.T) {
	Init("http://localhost:3465/send-email")
	if err := SendEmail(EmailRequest{
		To:      "mike504110403@gmail.com",
		Subject: "Test",
		Body:    "Test",
		Image:   "",
	}); err != nil {
		t.Errorf("SendEmail() error = %v", err)
	}
}
