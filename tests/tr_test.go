package tests

import (
	"testing"

	"github.com/ALiwoto/trGo/trGo"
	"github.com/ALiwoto/trGo/trGo/trLang"
)

//---------------------------------------------------------

//======================CORRECT-TESTS======================
//---------------------------------------------------------
//-----------------Tests which have correct from-----------
//---------------------------------------------------------

func TestCorrectTr(t *testing.T) {
	// after separator, always search for kind {?}
	//tr, err := trGo.TranslateD("en", "ja", "what are you doing here??\n I hat you!!")
	//tr, err := trGo.TranslateD("en", "ja", "what")
	//tr, err := trGo.TranslateD("ja", "en", "yasashii")
	//tr, err := trGo.TranslateD("en", "ru", "Hello, what are you doing here, mother?")
	//tr, err := trGo.TranslateD("en", "ru", "wha should I do rigt now?")
	//tr, err := trGo.TranslateD("en", "ru", "wha shold I do rigt now?")
	tr, err := trGo.TranslateD("en", "ru", "wha shold I do rigt nuw?")
	//tr, err := trGo.TranslateD("en", "ru", "wha shold I do rigt nuw? \n I really love you!")
	//tr, err := trGo.TranslateD("en", "ru", "organaization")
	//tr, err := trGo.TranslateD("en", "ru", "\n\r\n     \n\n     \n")
	//tr, err := trGo.TranslateD("ja", "ru", "こんにちわ")

	if err != nil {
		t.Fatal(err)
	}

	for i, tl := range tr.Translations {
		t.Log(i, ": "+tl)
	}

	t.Log("TranslatedPronunciation: " + tr.TranslatedPronunciation)
}

//---------------------------------------------------------

//======================WRONG-TESTS========================
//---------------------------------------------------------
//-----------------Tests which have wrong from-------------
//---------------------------------------------------------

func TestWrongTr(t *testing.T) {
	//tr, err := trGo.TranslateD("en", "ja", "what are you doing here??\n I hat you!!")
	//tr, err := trGo.TranslateD("en", "ja", "what")
	//tr, err := trGo.TranslateD("auto", "en", "yasashii")
	//tr, err := trGo.TranslateD("en", "ru", "Hello, what are you doing here, mother?")
	//tr, err := trGo.TranslateD("en", "ru", "wha should I do rigt now?")
	//tr, err := trGo.TranslateD("en", "ru", "wha shold I do rigt now?")
	//tr, err := trGo.TranslateD("en", "ru", "wha shold I do rigt nuw?")
	//tr, err := trGo.TranslateD("en", "ru", "wha shold I do rigt nuw? \n I really love you!")
	//tr, err := trGo.TranslateD("en", "ru", "what should I do right now? \n I really love you!")
	//tr, err := trGo.TranslateD("ja", "ru", "organaization")
	//tr, err := trGo.TranslateD("en", "ru", "\n\r\n     \n\n     \n")
	//tr, err := trGo.TranslateD("en", "ru", "こんにちわ")
	//tr, err := trGo.TranslateD("ru", "ja", "Организация")
	//tr, err := trGo.TranslateD("ru", "en", "組織")
	//tr, err := trGo.TranslateD("ru", "en", "組織")
	//tr, err := trGo.TranslateD("ja", "en", "nami")
	//tr, err := trGo.TranslateD("ru", "en", "波")
	//tr, err := trGo.TranslateD("ja", "ru", "波") // колыхаться
	//tr, err := trGo.TranslateD("ja", "ru", "колыхаться")
	//tr, err := trGo.TranslateD("fr", "ru", "колыхаться")
	tr, err := trGo.TranslateD("fr", "ru", "Kolykhat'sya")

	if err != nil {
		t.Fatal(err)
	}

	for i, tl := range tr.Translations {
		t.Log(i, ": "+tl)
	}

	t.Log("TranslatedPronunciation: " + tr.TranslatedPronunciation)
}

//---------------------------------------------------------

//======================TRLANG-TESTS========================
//---------------------------------------------------------
//-----------------Tests related to trLang package---------
//---------------------------------------------------------

func TestTrLang(t *testing.T) {
	short := trLang.ExtractShortLang("eNglish")
	if short == nil {
		t.Fatal("short was nil!")
	}

	t.Log(short)

	short = trLang.ExtractShortLang("JaPANESE")
	if short == nil {
		t.Fatal("short was nil!")
	}

	t.Log(short)

	short = trLang.ExtractShortLang("JA")
	if short == nil {
		t.Fatal("short was nil!")
	}

	t.Log(short)

	short = trLang.ExtractShortLang("123456789784524541254874512548764")
	if short != nil {
		t.Fatal("short wasn't nil!")
	}

	t.Log(short)
}

//---------------------------------------------------------
