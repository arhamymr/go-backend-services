package handlers

import (
	"go-backend-services/sender"
	"go-backend-services/types"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/mailersend/mailersend-go"
)

func TestMessaging(c echo.Context) error {
	subject := "Subject"
	text := "Greetings from the team, you got this message through MailerSend."
	html := "Greetings from the team, you got this message through MailerSend."

	from := mailersend.From{
		Name:  "Apsara Project",
		Email: "info@apsara.com",
	}

	recipients := []mailersend.Recipient{
		{
			Name:  "Recipient",
			Email: "arhamymr@gmail.com",
		},
	}

	tags := []string{"foo", "bar"}

	_, err := sender.Email(from, recipients, subject, text, html, tags)

	var response types.Response

	if err != nil {
		response = types.Response{
			Status:  http.StatusBadRequest,
			Data:    struct{}{},
			Message: err.Error(),
		}
		return c.JSON(http.StatusBadRequest, response)
	}

	response = types.Response{
		Status:  http.StatusOK,
		Data:    struct{}{},
		Message: "OK sended",
	}

	return c.JSON(http.StatusOK, response)
}
