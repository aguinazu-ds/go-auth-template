package mailer

import (
	"bytes"
	"context"
	"embed"
	"fmt"
	"log/slog"
	"os"
	"text/template"

	brevo "github.com/getbrevo/brevo-go/lib"
)

//go:embed "templates"
var templateFS embed.FS

var br *brevo.APIClient

type WelcomeEmailData struct {
	ActivationUrl string
}

func SendEmailUsingTemplate(to string, subject string, nameSender string, emailSender string, templateNameFile string, data interface{}) error {
	template, err := template.New(templateNameFile).ParseFS(templateFS, "templates/"+templateNameFile)
	if err != nil {
		return err
	}

	htmlBody := new(bytes.Buffer)
	err = template.ExecuteTemplate(htmlBody, "htmlBody", data)
	if err != nil {
		return err
	}

	sender := &brevo.SendSmtpEmailSender{
		Name:  nameSender,
		Email: emailSender,
	}

	_, resp, err := br.TransactionalEmailsApi.SendTransacEmail(context.Background(), brevo.SendSmtpEmail{
		Sender: sender,
		To: []brevo.SendSmtpEmailTo{
			{
				Email: to,
			},
		},
		HtmlContent: htmlBody.String(),
		Subject:     subject,
	})
	if err != nil {
		return err
	}
	slog.Info(fmt.Sprintf("Response: %v", resp))
	return nil
}

func Init() error {
	brevoApiKey := os.Getenv("BREVO_API_KEY")

	var ctx context.Context

	config := brevo.NewConfiguration()
	config.AddDefaultHeader("api-key", brevoApiKey)

	br = brevo.NewAPIClient(config)
	result, _, err := br.AccountApi.GetAccount(ctx)
	if err != nil {
		return err
	}

	slog.Info(fmt.Sprintf("Brevo account connection success: %v", result))
	return nil
}
