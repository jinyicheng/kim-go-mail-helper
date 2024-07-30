package mailHelper

import (
	"crypto/tls"
	"encoding/json"
	"strconv"
	"strings"
	"time"
)

type Mail struct {
	Configs []Config
}

type Config struct {
	Name        string              `json:"name"`
	Host        string              `json:"host"`
	Port        int                 `json:"port"`
	Username    string              `json:"username"`
	Password    string              `json:"password"`
	Priority    *int                `json:"priority"`
	Subject     *string             `json:"subject"`
	Sender      *ConfigSender       `json:"sender"`
	From        ConfigFrom          `json:"from"`
	To          *[]ConfigTo         `json:"to"`
	ReplyTo     *[]ConfigReplyTo    `json:"reply_to"`
	Cc          *[]ConfigCc         `json:"cc"`
	Bcc         *[]ConfigBcc        `json:"bcc"`
	Body        *ConfigBody         `json:"body"`
	Attachments *[]ConfigAttachment `json:"attachments"`
}
type ConfigSender struct {
	Address string  `json:"address"`
	Name    *string `json:"name"`
}
type ConfigFrom struct {
	Address string  `json:"address"`
	Name    *string `json:"name"`
}
type ConfigTo struct {
	Address string  `json:"address"`
	Name    *string `json:"name"`
}
type ConfigReplyTo struct {
	Address string  `json:"address"`
	Name    *string `json:"name"`
}
type ConfigCc struct {
	Address string  `json:"address"`
	Name    *string `json:"name"`
}
type ConfigBcc struct {
	Address string  `json:"address"`
	Name    *string `json:"name"`
}
type ConfigBody struct {
	IsHtml  bool   `json:"is_html"`
	Content string `json:"content"`
}
type ConfigAttachment struct {
	Path    string  `json:"path"`
	Name    *string `json:"name"`
	IsEmbed bool    `json:"is_embed"`
}

func (m Mail) Get() map[string]Config {
	configMap := make(map[string]Config, len(m.Configs))
	for _, config := range m.Configs {
		configMap[config.Name] = config
	}
	return configMap
}

func (c Config) Send(message *Message) error {
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
	if c.To != nil {
		var to []string
		for _, item := range *c.To {
			if item.Name != nil {
				to = append(to, m.FormatAddress(item.Address, *item.Name))
			} else {
				to = append(to, item.Address)
			}
		}
		m.SetHeader("To", strings.Join(to, ","))
	}
	//section Reply To
	if c.ReplyTo != nil {
		var replyTo []string
		for _, item := range *c.ReplyTo {
			if item.Name != nil {
				replyTo = append(replyTo, m.FormatAddress(item.Address, *item.Name))
			} else {
				replyTo = append(replyTo, item.Address)
			}
		}
		m.SetHeader("Reply-To", strings.Join(replyTo, ","))
	}
	//section Cc
	if c.Cc != nil {
		var cc []string
		for _, item := range *c.Cc {
			if item.Name != nil {
				cc = append(cc, m.FormatAddress(item.Address, *item.Name))
			} else {
				cc = append(cc, item.Address)
			}
		}
		m.SetHeader("Cc", strings.Join(cc, ","))
	}
	//section Bcc
	if c.Bcc != nil {
		var bcc []string
		for _, item := range *c.Bcc {
			if item.Name != nil {
				bcc = append(bcc, m.FormatAddress(item.Address, *item.Name))
			} else {
				bcc = append(bcc, item.Address)
			}
		}
		m.SetHeader("Bcc", strings.Join(bcc, ","))
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
	if c.Attachments != nil {
		for _, attachment := range *c.Attachments {
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
	v, _ := json.Marshal(c)
	print(string(v))
	return d.DialAndSend(m)
}
