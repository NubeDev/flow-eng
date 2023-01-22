package gmail

import (
	"fmt"
	"github.com/jordan-wright/email"
	"net/smtp"
	"testing"
)

func TestGmail_sendEmail(t *testing.T) {
	e := email.NewEmail()
	e.From = "nubeio <noreply@nube-io.com>"
	e.To = []string{"ap@nube-io.com"}
	e.Subject = "test"
	e.Text = []byte("Text Body is, of course, supported!")
	e.HTML = []byte("<h1>Fancy HTML is supported, too!</h1>")
	err := e.Send("smtp.gmail.com:587", smtp.PlainAuth("", "noreply@nube-io.com", "22222-11eb-111111111", "smtp.gmail.com"))
	fmt.Println(err)
}
