package mailgun

import (
	"context"
	"crop_connect/constant"
	"crop_connect/util"
	"errors"
	"fmt"
	"time"

	"github.com/mailgun/mailgun-go/v3"
)

type Function interface {
	SendOneMailUsingTemplate(subject string, template string, receipentEmail string, plainText string, variable map[string]string) (string, string, error)
}

type Mailgun struct {
	Mailgun     *mailgun.MailgunImpl
	EmailDomain string
	EmailSender string
}

func Init(emailDomain string, emailSender string, mailgunKey string) Function {
	mg := mailgun.NewMailgun(emailDomain, mailgunKey)

	return &Mailgun{
		Mailgun:     mg,
		EmailSender: util.ReplaceUnderScoreWithSpace(emailSender),
		EmailDomain: emailDomain,
	}
}

func (mg *Mailgun) SendOneMailUsingTemplate(subject string, template string, receipentEmail string, plainText string, variable map[string]string) (string, string, error) {
	if !util.CheckStringOnArray([]string{constant.MailgunForgotPasswordTemplate}, template) {
		return "", "", errors.New("template tidak tersedia")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	mail := mg.Mailgun.NewMessage(
		mg.EmailSender, // From
		subject,        // Subject
		plainText,      // Plain-text
		receipentEmail, // Recipients
	)

	mail.SetTemplate(template)

	for key, value := range variable {
		if err := mail.AddTemplateVariable(key, value); err != nil {
			return "", "", err
		}

	}

	response, id, err := mg.Mailgun.Send(ctx, mail)
	fmt.Println(response, id, err)

	return response, id, err
}
