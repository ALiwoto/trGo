// Bot.go Project
// Copyright (C) 2021 Sayan Biswas, ALiwoto
// This file is subject to the terms and conditions defined in
// file 'LICENSE', which is part of the source code.

package trGo

// Lang struct contains data of the language.
type Lang struct {
	Data *LangData `json:"data"`
}

// LangData contains the data of the language.
type LangData struct {
	Detections []LangDetect `json:"detections"`
}

type Kind string

// LangDetect contains the detected languages.
type LangDetect struct {
	TheLang    string  `json:"language"`
	Reliable   bool    `json:"isReliable"`
	Confidence float32 `json:"confidence"`
}

// gnuTranslate contains necessary fields for
// using in gnu translation.
type gnuTranslate struct {
	Result string `json:"result"`
	Err    string `json:"error"`
}

// WotoTr contains necessary fields for results of a
// transtion operation.
type WotoTr struct {
	// Pronunciation of the original text
	OriginalPronunciation string

	// Pronunciation of the translated text
	TranslatedPronunciation string
	UserText                string
	OriginalText            string
	TranslatedText          []string
	From                    string
	To                      string
	Corrected               *Corrected
	HasWrongNess            bool
	WrongFrom               bool
	Road                    map[int]bool
}

type Corrected struct {
	// an array of the corrected parts of the
	// original input text
	CorrectedParts []string

	// whole string
	CorrectedValue string
}
