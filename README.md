# kim-go-mail-helper
邮件助手，简化邮件发送接收的部分实践
``` golang
// 以下是一段示例，这段示例结合了kim-go-xxx-helper系列的kim-go-config-helper，在发送邮件的时候选择合适的邮箱
func notifyByMail(subject string, to []mailHelper.MailContact, cc *[]mailHelper.MailContact, bcc *[]mailHelper.MailContact, attachmentPath string) {
	mail := app.Mail["common"]
	newSubject := subject
	mail.Subject = &newSubject
	mail.Body.Content = subject
	mail.To = &to
	if cc != nil {
		mail.Cc = cc
	}
	if bcc != nil {
		mail.Bcc = bcc
	}
	if attachmentPath != "" {
		mail.Attachments = &[]mailHelper.MailAttachment{
			{Path: attachmentPath},
		}
	}
	err := mail.Send(mailHelper.NewMessage())
	if err != nil {
		log.Println(err)
	}
}
```