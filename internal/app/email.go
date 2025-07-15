package app

import (
	"bytes"
	"os"
	"strings"

	"github.com/resend/resend-go/v2"
)

func (a *appLayer) sendEmail(from, to, subject, template string, obj any) {
	resendApiKey, resendApiKeyExist := os.LookupEnv("RESEND_API_KEY")
	emailTemplate, exists := a.emailTemplates[template]

	if !resendApiKeyExist && !exists {
		return
	}

	var buf bytes.Buffer
	err := emailTemplate.Execute(&buf, obj)
	if err != nil {
		return
	}

	client := resend.NewClient(resendApiKey)

	params := &resend.SendEmailRequest{
		From:    from,
		To:      strings.Split(to, ","),
		Subject: subject,
		Html:    buf.String(),
	}

	client.Emails.Send(params)
}
