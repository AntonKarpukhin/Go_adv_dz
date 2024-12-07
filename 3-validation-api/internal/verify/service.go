package verify

import (
	"fmt"
	"github.com/jordan-wright/email"
	"net/http"
	"net/smtp"
	"os"
	"validation/pkg/response"
)

func CreateMail(hash string, w http.ResponseWriter) {
	verifyURL := fmt.Sprintf("http://localhost:8083/auth/verify/" + hash)
	emailSMTP := os.Getenv("SMTP_EMAIL")
	passwordSMTP := os.Getenv("SMTP_PASSWORD")

	e := email.NewEmail()
	e.From = "LembrarB <a.karpukhin21@gmail.com>"
	e.To = []string{"a.karpukhin21@mail.ru"}
	e.Subject = "Test mail!"
	e.Text = []byte(fmt.Sprintf("Click fucking anchor =) %s", verifyURL))
	e.HTML = []byte(fmt.Sprintf(`
		<!DOCTYPE html>
		<html>
		<head>
			<title>Click</title>
		</head>
		<body>
			<p>Click fucking anchor =):</p>
			<a href="%s">Verify Account</a>
		</body>
		</html>
	`, verifyURL))

	err := e.Send("smtp.gmail.com:587", smtp.PlainAuth("", emailSMTP, passwordSMTP, "smtp.gmail.com"))
	if err != nil {
		fmt.Println("Error sending email:", err)
	} else {
		response.JsonResponse(w, "Письмо отправлено", http.StatusOK)
		fmt.Println("Письмо отправлено")
	}
}
