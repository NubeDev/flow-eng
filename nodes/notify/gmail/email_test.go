package gmail

import (
	"fmt"
	"net/smtp"
	"testing"

	"github.com/jordan-wright/email"
)

func TestGmail_sendEmail(t *testing.T) {

	e := email.NewEmail()
	e.From = "test@gmail.com"
	e.To = []string{"test@nube-io.com"}
	e.Subject = "test"
	e.Text = []byte("Text Body is, of course, supported!")
	e.HTML = []byte("<h1>Fancy HTML is supported, too!</h1>")
	err := e.Send("smtp.gmail.com:587", smtp.PlainAuth("", "test@gmail.com", "1111222233334444", "smtp.gmail.com"))
	fmt.Println(err)
}
