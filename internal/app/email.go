package app

import (
	"bytes"
	"context"
	"log"
	"os"
	"strings"

	"github.com/mailerlite/mailerlite-go"
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
		ReplyTo: "info@outclimb.gay",
		Subject: subject,
		Html:    buf.String(),
	}

	client.Emails.Send(params)
}

func (a *appLayer) subscribeToNewsletter(firstName, lastName, email string) {
	mailerLiteApiKey, mailerLiteApiKeyExist := os.LookupEnv("MAILERLITE_API_KEY")
	if !mailerLiteApiKeyExist {
		return
	}

	client := mailerlite.NewClient(mailerLiteApiKey)

	ctx := context.Background()

	subscriber := &mailerlite.UpsertSubscriber{
		Email: email,
		Fields: map[string]interface{}{
			"name":      firstName,
			"last_name": lastName,
		},
	}

	_, _, err := client.Subscriber.Upsert(ctx, subscriber)
	if err != nil {
		log.Printf("Unable to subscribe individual to newsletter: %s\n", err)
	}
}
