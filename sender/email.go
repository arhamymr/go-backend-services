package sender

import (
	"context"
	"os"
	"time"

	"github.com/mailersend/mailersend-go"
)

func Email(from mailersend.From, recipients []mailersend.Recipient, subject, text, html string, tags []string) (*mailersend.Response, error) {

	ms := mailersend.NewMailersend(os.Getenv("MAILER_SENDER_TOKEN"))

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	message := ms.Email.NewMessage()

	message.SetFrom(from)
	message.SetRecipients(recipients)
	message.SetSubject(subject)
	message.SetHTML(html)
	message.SetText(text)
	message.SetTags(tags)

	// TODO: Implement retry mecanism logic
	return ms.Email.Send(ctx, message)
}
