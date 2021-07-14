package tests

import (
	"log"
	"testing"

	"github.com/ALiwoto/trGo/trGo"
)

func TestTr(t *testing.T) {
	// after separator, always search for kind {?}
	//tr, err := trGo.TranslateD("en", "ja", "what are you doing here??\n I hat you!!")
	//tr, err := trGo.TranslateD("en", "ja", "what")
	//tr, err := trGo.TranslateD("en", "ru", "what are you doing here, mother?")
	//tr, err := trGo.TranslateD("en", "ru", "wha should I do rigt now?")
	//tr, err := trGo.TranslateD("en", "ru", "wha shold I do rigt now?")
	tr, err := trGo.TranslateD("en", "ru", "wha shold I do rigt nuw?")
	//tr, err := trGo.TranslateD("en", "ru", "organaization")
	if err != nil {
		log.Fatal(err)
	}
	log.Println(tr.TranslatedText)
}
