package gmail

import (
	"fmt"
	"net/smtp"
	"testing"

	"github.com/jordan-wright/email"
)

func TestGmail_sendEmail(t *testing.T) {

	e := email.NewEmail()
	e.From = "jia <feijiajidangao@gmail.com>"
	e.To = []string{"jfe@nube-io.com"}
	e.Subject = "test"
	e.Text = []byte("Text Body is, of course, supported!")
	e.HTML = []byte("<h1>Fancy HTML is supported, too!</h1>")
	err := e.Send("smtp.gmail.com:587", smtp.PlainAuth("", "feijiajidangao@gmail.com", "wakzbowhbpbcdenu", "smtp.gmail.com"))
	fmt.Println(err)
}
