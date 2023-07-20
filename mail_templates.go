package main

import (
	"bytes"
	"html/template"
)

func serviceLocTempl(loc *ServiceLocation) (bytes.Buffer, error) {
	var templContents bytes.Buffer

	templ, err := template.ParseFiles("./templates/server_loc.html")
	if err != nil {
		return templContents, err
	}
	err = templ.Execute(&templContents, loc)

	return templContents, nil
}

func alertMailTempl(alert *Alert) (bytes.Buffer, error) {
	var templContents bytes.Buffer

	templ, err := template.ParseFiles("./templates/alert.html")
	if err != nil {
		return templContents, err
	}
	err = templ.Execute(&templContents, alert)

	return templContents, err
}

func backOnlineNotif(lg *LastLog) (bytes.Buffer, error) {

	var templContents bytes.Buffer

	templ, err := template.ParseFiles("./templates/back_online.html")
	if err != nil {
		return templContents, err
	}
	err = templ.Execute(&templContents, lg)

	return templContents, err
}
