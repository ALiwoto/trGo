// Bot.go Project
// Copyright (C) 2021 Sayan Biswas, ALiwoto
// This file is subject to the terms and conditions defined in
// file 'LICENSE', which is part of the source code.

package trLang

import (
	"strings"
	"sync"

	ws "github.com/ALiwoto/StrongStringGo/strongStringGo"
)

var langList map[string]string  // lang list map
var langListR map[string]string // reverse lang list map
var listMutix *sync.RWMutex     // mutix of these maps

// initLang will initialize both reverse and normal
// language map for
func initLang() {
	initLangMap()
	initLangReseve()
}

// initLangMap will initialize the language map for
// wotoLang package.
func initLangMap() {
	if listMutix == nil {
		listMutix = &sync.RWMutex{}
	}

	if langList == nil {
		listMutix.Lock()
		langList = map[string]string{
			L_af:    "Afrikaans", // Afrikaans
			L_am:    "Amharic",   // Amharic
			L_ar:    "Arabic",
			L_az:    "Azerbaijani",
			L_be:    "Belarusian",
			L_bg:    "Bulgarian",
			L_bn:    "Bengali",
			L_bs:    "Bosnian",
			L_ca:    "Catalan",
			L_ceb:   "Chechen",
			L_co:    "Corsican",
			L_cs:    "Czech",
			L_cy:    "Welsh",
			L_da:    "Danish",
			L_de:    "German",
			L_el:    "Greek",
			L_en:    "English",
			L_eo:    "Esperanto",
			L_es:    "Spanish",
			L_et:    "Estonian",
			L_eu:    "Basque",
			L_fa:    "Persian",
			L_fi:    "Finnish",
			L_fr:    "French",
			L_fy:    "WesternFrisian",
			L_ga:    "Irish",
			L_gd:    "Gaelic",
			L_gl:    "Galician",
			L_gu:    "Gujarati",
			L_ha:    "Hausa",
			L_haw:   "haw", // ???
			L_hi:    "Hindi",
			L_hmn:   "hmn", // ???
			L_hr:    "Croatian",
			L_ht:    "Haitian",
			L_hu:    "Hungarian",
			L_hy:    "Armenian",
			L_id:    "Indonesian",
			L_ig:    "Igbo",
			L_is:    "Icelandic",
			L_it:    "Italian",
			L_iw:    "Hebrew",
			L_ja:    "Japanese",
			L_jw:    "jw", // ???
			L_ka:    "Georgian",
			L_kk:    "Kazakh",
			L_km:    "Central Khmer",
			L_kn:    "Kannada",
			L_ko:    "Korean",
			L_ku:    "Kurdish",
			L_ky:    "Kirghiz",
			L_la:    "Latin",
			L_lb:    "Luxembourgish",
			L_lo:    "Lao",
			L_lt:    "Lithuanian",
			L_lv:    "Latvian",
			L_mg:    "Malagasy",
			L_mi:    "Maori",
			L_mk:    "Macedonian",
			L_ml:    "Malayalam",
			L_mn:    "Mongolian",
			L_mr:    "Marathi",
			L_ms:    "Malay",
			L_mt:    "Maltese",
			L_my:    "Burmese",
			L_ne:    "Nepali",
			L_nl:    "Dutch",
			L_no:    "Norwegian",
			L_ny:    "Chichewa",
			L_pa:    "Punjabi",
			L_pl:    "Polish",
			L_ps:    "Pashto",
			L_pt:    "Portuguese",
			L_ro:    "Romanian",
			L_ru:    "Russian",
			L_sd:    "Sindhi",
			L_si:    "Sinhala",
			L_sk:    "Slovak",
			L_sl:    "Slovenian",
			L_sm:    "Samoan",
			L_sn:    "Shona",
			L_so:    "Somali",
			L_sq:    "Albanian",
			L_sr:    "Serbian",
			L_st:    "Southern Sotho",
			L_su:    "Sundanese",
			L_sv:    "Swedish",
			L_sw:    "Swahili",
			L_ta:    "Tamil",
			L_te:    "Telugu",
			L_tg:    "Tajik",
			L_th:    "Thai",
			L_tl:    "Tagalog",
			L_tr:    "Turkish",
			L_uk:    "Ukrainian",
			L_ur:    "Urdu",
			L_uz:    "Uzbek",
			L_vi:    "Vietnamese",
			L_xh:    "Xhosa",
			L_yi:    "Yiddish",
			L_yo:    "Yoruba",
			L_zh:    "Chinese",
			L_zh_CN: "Chinese_CN", // Chinese
			L_zh_TW: "Chinese_TW", // Chinese
			L_zu:    "zuZulu",
		}
		listMutix.Unlock()
	}
}

// initLangReseve will initialize the reverse map
// for wotoLang package.
func initLangReseve() {
	if langList == nil {
		return
	}
	listMutix.Lock()

	if langListR == nil {
		langListR = make(map[string]string)
	}

	for k, v := range langList {
		k = strings.ToLower(k)
		v = strings.ToLower(v)
		if langListR[v] != k {
			langListR[v] = k
		}
	}

	listMutix.Unlock()
}

// IsLang will check if a string value is a valid
// language or not.
func IsLang(value string) bool {
	l := len(value)
	if l <= ws.BaseIndex || l > len(L_zh_CN) {
		return false
	}

	initLang()
	s := langList[value]

	if ws.IsEmpty(&s) {
		s = langListR[value]
		return !ws.IsEmpty(&s)
	}

	return true
}

func ExtractShortLang(value string) *string {
	l := len(value)
	if l <= ws.BaseIndex || l >= MaxLenght {
		return nil
	}

	initLang()

	value = strings.ToLower(value)
	s := langList[value]
	if ws.IsEmpty(&s) {
		s = langListR[value]
		if ws.IsEmpty(&s) {
			return nil
		}

		return &s
	}

	return &value
}

func RemoveShortsWithStrs(value string) string {
	if ws.IsEmpty(&value) || len(value) <= ws.BaseOneIndex {
		return value
	}

	for c := range langList {
		value = strings.ReplaceAll(value,
			ws.STR_SIGN+c+ws.STR_SIGN, ws.EMPTY)
	}

	return value
}
