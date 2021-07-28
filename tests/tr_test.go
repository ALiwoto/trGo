package tests

import (
	"io/ioutil"
	"log"
	"net/http"
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
	tr, err := trGo.TranslateText("en", "ru", "Hello, what are you doing here, mother?")
	//tr, err := trGo.TranslateD("en", "ru", "wha should I do rigt now?")
	//tr, err := trGo.TranslateD("en", "ru", "wha shold I do rigt now?")
	//tr, err := trGo.TranslateD("en", "ru", "wha shold I do rigt nuw?")
	//tr, err := trGo.TranslateD("en", "ru", "wha shold I do rigt nuw? \n I really love you!")
	//tr, err := trGo.TranslateD("en", "ru", "organaization")
	//tr, err := trGo.TranslateD("en", "ru", "\n\r\n     \n\n     \n")
	//tr, err := trGo.TranslateD("ja", "ru", "こんにちわ") // брух
	//tr, err := trGo.TranslateD("ru", "en", "брух")
	//tr, err := trGo.TranslateD("sr", "en", "брух")

	if err != nil {
		t.Fatal(err)
	}

	for i, tl := range tr.Translations {
		t.Log(i, ": "+tl)
		log.Println(i, ": "+tl)
	}

	t.Log("TranslatedPronunciation: " + tr.TranslatedPronunciation)
	log.Println("TranslatedPronunciation: " + tr.TranslatedPronunciation)
}

//---------------------------------------------------------

//======================WRONG-TESTS========================
//---------------------------------------------------------
//-----------------Tests which have wrong from-------------
//---------------------------------------------------------

func TestWrongTr(t *testing.T) {
	//tr, err := trGo.TranslateText("en", "ja", "what are you doing here??\n I hat you!!")
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
	//tr, err := trGo.TranslateText("fr", "ru", "Kolykhat'sya")
	//tr, err := trGo.TranslateText("ja", "en", "kimi ha nanimono da?")
	tr, err := trGo.TranslateText("ja", "en", "yasashii")

	if err != nil {
		t.Fatal(err)
		log.Fatal(err)
	}

	for i, tl := range tr.Translations {
		t.Log(i, ": "+tl)
		log.Println(i, ": "+tl)
	}

	t.Log("TranslatedPronunciation: " + tr.TranslatedPronunciation)
	log.Println("TranslatedPronunciation: " + tr.TranslatedPronunciation)
}

//---------------------------------------------------------

//======================TRLANG-TESTS========================
//---------------------------------------------------------
//-----------------Tests related to trLang package---------
//---------------------------------------------------------

func TestTrLang(t *testing.T) {
	//xhr.open('GET', 'https://translate.google.com/translate_a/single?client=t&sl=auto&tl=' + toLang + '&hl=en&dt=bd&dt=ex&dt=ld&dt=md&dt=qc&dt=rw&dt=rm&dt=ss&dt=t&dt=at&ie=UTF-8&oe=UTF-8&source=btn&ssel=0&tsel=0&kc=0&tk=5&q=' + encodeURIComponent(text));

	resp, err := http.Get("https://translate.google.com/translate_a/single?client=t&sl=auto&tl=ja&hl=en&dt=bd&dt=ex&dt=ld&dt=md&dt=qc&dt=rw&dt=rm&dt=ss&dt=t&dt=at&ie=UTF-8&oe=UTF-8&source=btn&ssel=0&tsel=0&kc=0&tk=5&q=hello")
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatal(err)
	}
	log.Println(string(body))

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
