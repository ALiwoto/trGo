// Bot.go Project
// Copyright (C) 2021 Sayan Biswas, ALiwoto
// This file is subject to the terms and conditions defined in
// file 'LICENSE', which is part of the source code.

package trGo

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	ws "github.com/ALiwoto/StrongStringGo/strongStringGo"
)

// DetectLanguage will detect the language of a text.
func DetectLanguage(text string) *Lang {
	m := map[string]string{
		userAgentKey:      userAgentValue,
		acceptKey:         acceptValue,
		acceptLanguageKey: acceptLanguageValue,
		refererKey:        refererValue,
		contentTypeKey:    contentTypeValue,
		originKey:         originValue,
		connectionKey:     connectionValue,
		teKey:             teValue,
		qKey:              text,
	}

	data, errJ := json.Marshal(m)
	if errJ != nil {
		log.Println(errJ)
		return nil
	}

	reader := bytes.NewReader(data)
	resp, errH := http.Post(dHostUrl, contentTypeValue, reader)

	if errH != nil {
		log.Println(errH)
	}

	defer resp.Body.Close()

	b, errB := ioutil.ReadAll(resp.Body)
	if errB != nil {
		log.Println(errB)
	}

	//log.Println(string(b))
	str := ws.Qsb(b)
	strs := str.SplitStr(preLeft, preRight)
	if len(strs) <= ws.BaseOneIndex {
		return nil
	}

	b = []byte(strs[ws.BaseOneIndex].GetValue())

	var l Lang
	errJ = json.Unmarshal(b, &l)
	if errJ != nil {
		log.Println(errJ)
	}

	return &l
}
