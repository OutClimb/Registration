package app

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/resend/resend-go/v2"
)

func (a *appLayer) sendEmail(from, to, subject, template string, obj any) {
	resendApiKey, resendApiKeyExist := os.LookupEnv("RESEND_API_KEY")
	emailTemplate, exists := a.emailTemplates[template]

	if !resendApiKeyExist && !exists {
		fmt.Println("Email - 1")
		return
	}

	var buf bytes.Buffer
	err := emailTemplate.Execute(&buf, obj)
	if err != nil {
		fmt.Println("Email - 2")
		fmt.Println(err)
		return
	}

	client := resend.NewClient(resendApiKey)

	params := &resend.SendEmailRequest{
		From:    from,
		To:      strings.Split(to, ","),
		Subject: subject,
		Html:    buf.String(),
	}

	_, err = client.Emails.Send(params)
	fmt.Println(err)
}
