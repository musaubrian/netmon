package main

import (
	"bytes"
	"html/template"
)

func serverLocTempl(loc ServerLocation) (bytes.Buffer, error) {
	var templContents bytes.Buffer

	templ, err := template.ParseFiles("./template/server_loc.html")
	if err != nil {
		return templContents, err
	}
	err = templ.Execute(&templContents, loc)
	if err != nil {
		return templContents, err
	}

	return templContents, nil
}
