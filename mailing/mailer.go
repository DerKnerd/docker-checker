package mailing

import (
	"bytes"
	"docker-checker/configuration"
	"github.com/hashicorp/go-version"
	"github.com/scorredoira/email"
	"html/template"
	"io/ioutil"
	"net/smtp"
	"strconv"
)

func SendMail(usedVersion version.Version, latestVersion version.Version, image configuration.Image, config configuration.EmailConfig) error {
	tmpl, err := template.New("email").ParseFiles("mailing/mail-body.gohtml")
	if err != nil {
		return err
	}

	type tmplData struct {
		UsedVersion   string
		LatestVersion string
		Image         string
	}

	buffer := bytes.NewBuffer([]byte{})
	err = tmpl.Execute(buffer, tmplData{
		UsedVersion:   usedVersion.String(),
		LatestVersion: latestVersion.String(),
		Image:         image.Name,
	})

	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(buffer)
	if err != nil {
		return err
	}

	message := email.NewHTMLMessage("New version for image "+image.Name, string(body))
	message.To = []string{config.To}
	message.From.Address = config.From
	message.BodyContentType = "text/html"

	var auth smtp.Auth
	if config.Username != "" && config.Password != "" {
		auth = smtp.PlainAuth("", config.Username, config.Password, config.Host)
	}

	return email.Send(config.Host+":"+strconv.Itoa(config.Port), auth, message)
}
