package mailHelper

import (
	"crypto/tls"
	"encoding/json"
	"log"
	"strconv"
	"time"
)

type MailBox struct {
	Mails []Mail
}

type Mail struct {
	Name        string            `json:"name"`
	Host        string            `json:"host"`
	Port        int               `json:"port"`
	Username    string            `json:"username"`
	Password    string            `json:"password"`
	Priority    *int              `json:"priority"`
	Subject     *string           `json:"subject"`
	Sender      *MailContact      `json:"sender"`
	From        MailContact       `json:"from"`
	To          []*MailContact    `json:"to"`
	ReplyTo     []*MailContact    `json:"reply_to"`
	Cc          []*MailContact    `json:"cc"`
	Bcc         []*MailContact    `json:"bcc"`
	Body        *MailBody         `json:"body"`
	Attachments []*MailAttachment `json:"attachments"`
	Env         string            `json:"env"`
}
type MailContact struct {
	Address string  `json:"address"`
	Name    *string `json:"name"`
}
type MailBody struct {
	IsHtml  bool   `json:"is_html"`
	Content string `json:"content"`
}
type MailAttachment struct {
	Path    string  `json:"path"`
	Name    *string `json:"name"`
	IsEmbed bool    `json:"is_embed"`
}

func (m MailBox) Get() map[string]Mail {
	mailMap := make(map[string]Mail, len(m.Mails))
	for _, mail := range m.Mails {
		mailMap[mail.Name] = mail
	}
	return mailMap
}

func (c Mail) Send(message *Message) error {
	d := NewDialer(c.Host, c.Port, c.Username, c.Password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	var m *Message
	if message != nil {
		m = message
	} else {
		m = NewMessage()
	}
	//section sender
	if c.Sender != nil {
		if c.Sender.Name != nil {
			m.SetHeader("Sender", m.FormatAddress(c.Sender.Address, *c.Sender.Name))
		} else {
			m.SetHeader("Sender", c.Sender.Address)
		}
	}
	//section From
	if c.From.Name != nil {
		m.SetHeader("From", m.FormatAddress(c.From.Address, *c.From.Name))
	} else {
		m.SetHeader("From", c.From.Address)
	}

	//section To
	if len(c.To) > 0 {
		var to []string
		for _, item := range c.To {
			if item.Name != nil {
				to = append(to, m.FormatAddress(item.Address, *item.Name))
			} else {
				to = append(to, item.Address)
			}
		}
		m.SetHeader("To", to...)
	}
	//section Reply To
	if len(c.ReplyTo) > 0 {
		var replyTo []string
		for _, item := range c.ReplyTo {
			if item.Name != nil {
				replyTo = append(replyTo, m.FormatAddress(item.Address, *item.Name))
			} else {
				replyTo = append(replyTo, item.Address)
			}
		}
		m.SetHeader("Reply-To", replyTo...)
	}
	//section Cc
	if len(c.Cc) > 0 {
		var cc []string
		for _, item := range c.Cc {
			if item.Name != nil {
				cc = append(cc, m.FormatAddress(item.Address, *item.Name))
			} else {
				cc = append(cc, item.Address)
			}
		}
		m.SetHeader("Cc", cc...)
	}
	//section Bcc
	if len(c.Bcc) > 0 {
		var bcc []string
		for _, item := range c.Bcc {
			if item.Name != nil {
				bcc = append(bcc, m.FormatAddress(item.Address, *item.Name))
			} else {
				bcc = append(bcc, item.Address)
			}
		}
		m.SetHeader("Bcc", bcc...)
	}
	// section X-Date
	now := time.Now()
	m.SetDateHeader("X-Date", now)
	m.SetHeader("X-Date-2", m.FormatDate(now))
	// section Priority
	if c.Priority != nil {
		m.SetHeader("X-Priority", strconv.Itoa(*c.Priority))
	}
	// section Subject
	if c.Subject != nil {
		m.SetHeader("Subject", *c.Subject)
	}
	// section Body
	if c.Body != nil {
		if c.Body.IsHtml {
			m.SetBody("text/html", c.Body.Content)
		} else {
			m.SetBody("text/plain", c.Body.Content)
		}

	}
	// section Attachment
	if len(c.Attachments) > 0 {
		for _, attachment := range c.Attachments {
			if attachment.Name != nil {
				if attachment.IsEmbed {
					m.Embed(attachment.Path, Rename(*attachment.Name))
				} else {
					m.Attach(attachment.Path, Rename(*attachment.Name))
				}
			} else {
				if attachment.IsEmbed {
					m.Embed(attachment.Path)

				} else {
					m.Attach(attachment.Path)
				}
			}
		}
	}
	if c.Env == "debug" {
		v, _ := json.Marshal(c)
		log.Println(string(v))
	}
	return d.DialAndSend(m)
}
