// Bot.go Project
// Copyright (C) 2021 Sayan Biswas, ALiwoto
// This file is subject to the terms and conditions defined in
// file 'LICENSE', which is part of the source code.

package trGo

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"unicode"

	"net/http"
	"net/url"
	"strconv"
	"strings"

	ws "github.com/ALiwoto/StrongStringGo/strongStringGo"
	"github.com/ALiwoto/trGo/trGo/trLang"
)

// https://telegra.ph/Lang-Codes-03-19-3

// Translate will translate the specified text value
// tp english.
func TranslateIt(lang *Lang, to, text string) (*WotoTr, error) {
	if ws.IsEmpty(&text) {
		return nil, errors.New("text cannot be empty")
	}

	toptr := trLang.ExtractShortLang(to)
	if toptr == nil {
		return nil, errors.New("language " + to + " is unrecognized")
	} else {
		to = *toptr
	}

	if lang == nil {
		lang = DetectLanguage(text)
	}

	best := lang.GetBest()

	if best == nil || !best.Reliable {
		return translateD(ws.AutoStr, to, text)
	}

	return translateD(best.TheLang, to, text)
}

func translateD(fr, to, text string) (*WotoTr, error) {
	uText := strings.TrimSpace(text)
	if ws.IsEmpty(&uText) {
		return nil, errors.New("[function TranslateD] " +
			"in package [trGo]: " + " text cannot be empty")
	}

	var err error
	text, err = trGoogle(fr, to, text)
	if err != nil {
		return nil, err
	}

	w := WotoTr{
		UserText:     uText,
		originalText: text,
		From:         fr,
		To:           to,
	}

	return parseGData(&w)
}

func Translate(to, text string) (*WotoTr, error) {
	return TranslateIt(nil, to, text)
}

func TranslateText(fr, to, text string) (*WotoTr, error) {
	frptr := trLang.ExtractShortLang(fr)
	if frptr == nil {
		return nil, errors.New("language " + to + " is unrecognized")
	} else {
		fr = *frptr
	}

	wTr, err := Translate(to, text)
	if err != nil {
		return nil, err
	}

	if !wTr.HasWrongFrom && wTr.From != fr {
		wTr.HasWrongFrom = true
	} else if wTr.HasWrongFrom && wTr.From == fr {
		wTr.HasWrongFrom = false
	}

	return wTr, nil
}

func SimpleTranslate(fr, to, text string) (string, error) {
	urlG := fmt.Sprintf(gnuHostUrl, fr, to, url.QueryEscape(text))
	resp, err := http.Get(urlG)
	if err != nil {
		return ws.EMPTY, err
	}

	defer resp.Body.Close()

	b, errB := ioutil.ReadAll(resp.Body)
	if errB != nil {
		return ws.EMPTY, errB
	}

	var g gnuTranslate
	errJ := json.Unmarshal(b, &g)
	if errJ != nil {
		return ws.EMPTY, errJ
	}

	if !ws.IsEmpty(&g.Err) {
		return ws.EMPTY, errors.New("an error from server: \"" +
			g.Err + "\"")
	}

	return g.Result, nil
}

// The default IDL of gRPC, Google’s RPC framework,
// is Protobuf — a data serialization method also created by Google
// that makes data transmission quick, efficient,
// and accessible to any language.
// gRPC: gRPC is a modern open source high performance
// Remote Procedure Call (RPC) framework that can run in
// any environment.
// It can efficiently connect services in and across data centers
// with pluggable support for load balancing, tracing,
// health checking and authentication.
// It is also applicable in last mile of distributed computing to
// connect devices, mobile applications and browsers to
// backend services.
//
//  > see also: https://en.wikipedia.org/wiki/Interface_description_language
//  > see also: https://grpc.io/docs/
//  > see also: https://developers.google.com/protocol-buffers
func parseGData(wTr *WotoTr) (*WotoTr, error) {
	text := ws.Qss(wTr.originalText)
	text.LockSpecial()

	myStrs := text.SplitStr(ws.BracketOpen, ws.Bracketclose)
	original := make([]string, ws.BaseIndex)
	accepted := func(v string) bool {
		if v == NonEscapeN {
			return false
		}
		if v == NonEscapeNV {
			return false
		}
		if strings.Contains(v, HttpRm) {
			return false
		}
		if strings.Contains(v, E4Value) {
			return false
		}
		if strings.HasPrefix(v, NullCValue) {
			tmpN := strings.ReplaceAll(v, NullCValue, ws.EMPTY)
			_, errN := strconv.Atoi(tmpN)
			if errN == nil {
				return false
			}
		}
		if strings.HasPrefix(v, DiValue) {
			return false
		}
		tmp := strings.ReplaceAll(v, ws.SPACE_VALUE, ws.EMPTY)
		tmp = strings.ReplaceAll(tmp, ws.CAMA, ws.EMPTY)
		if len(tmp) == ws.BaseIndex {
			return false
		}
		if tmp == ws.ParaClose || tmp == AkCloseQ {
			return false
		}
		if !strings.Contains(tmp, ws.NullStr) &&
			!strings.Contains(tmp, ws.DoubleQ) {
			return false
		}
		if strings.Contains(v, WrbFr) {
			return false
		}
		tmp = strings.ReplaceAll(tmp, ws.LineEscape, ws.EMPTY)
		_, errI := strconv.Atoi(tmp)
		return errI != nil
	}

	tmpV := ws.EMPTY
	for _, s := range myStrs {
		s.UnlockSpecial()
		tmpV = s.GetValue()
		if accepted(tmpV) {
			original = append(original, tmpV)
		}
	}

	log.Println(original)

	parseGparams(original, wTr)

	if wTr.wrongFrom {
		var err error
		var textStr string
		textStr, err = trGoogle(wTr.From, wTr.To, wTr.UserText)
		if err != nil {
			return nil, err
		}

		w := WotoTr{
			UserText:     wTr.UserText,
			originalText: textStr,
			From:         wTr.From,
			To:           wTr.To,
			HasWrongFrom: true,
		}
		wTr = &w

		// call yourself another time
		wTr, err = parseGData(wTr)
		if err != nil {
			return nil, err
		}

		return wTr, nil
	}

	return wTr, nil
}

// There are no key-value pairs in raw protobuf, just values assigned to field numbers. With batchexecute, I think Google is mapping protobuf messages to JSON in a special way. There is documentation on this, but it doesn’t quite match up to what we see here. This is how I think the above message would be mapped to JSON in batchexecute
func parseGparams(value []string, wTr *WotoTr) {
	//null,
	//null,
	// \"ja\"
	// \n,null,
	// null,\"Konnichiwa. Ohayou Minna\",null,null,null,
	// \"konnichiwa。\",
	// \"konnichiwa。\",\"こんにちは。\"
	// \"Ohayou Minna\",
	// \"Ohayou Minna\",\"みんなおはよう\"
	// \n,\"ja\",1,\"en\",
	// \"konnichiwa. ohayou minna \",\"en\",\"ja\",true
	// \n",null,null,null,"generic"
	//log.Println(wTr.OriginalText)
	if isWrongFrom(value, wTr) {
		return
	}

	//tmp := strings.Join(value, DY_WOTO_TEXT)

	//tmp = strings.ReplaceAll(tmp, NullN, ws.EMPTY)
	//tmp = strings.ReplaceAll(tmp, NullCValueR, ws.EMPTY)
	//tmp = strings.ReplaceAll(tmp, GenericStr, ws.EMPTY)
	//tmp = strings.ReplaceAll(tmp, NullCValue, ws.EMPTY)
	//tmp = strings.ReplaceAll(tmp, NeQ, ws.EMPTY)
	//strs := strings.Split(tmp, DY_WOTO_TEXT)
	//strMap := make(map[string]bool)
	lastStr := ws.EMPTY
	p1Set := false  // is original Pronunciation already set??
	p2Set := false  // is translated Pronunciation already set??
	tSet := false   // is trasnlated text already set??
	wCheck := false // is wrongness checked??
	isW := false    // is the current value a wrongness??
	tmp := ws.EMPTY
	for i, current := range value {
		current = strings.TrimSpace(current)

		tmp = current
		if current == lastStr || current == NullAndCama {
			continue
		} else {
			lastStr = tmp
		}

		if !wCheck {
			isW = isWrongness(current)
		}
		// check if pSet is true or not, if not, please try to
		// extract it from current element, if you couldn't extract it
		// at the end, you have to go for next element.
		if i == ws.BaseIndex && !p1Set && !isW {
			wTr.OriginalPronunciation, p1Set = getPronunciation(current)
			// Pronunciation field is mandatory, if you don't
			// find it at the first, you have to iterate over all
			// of the array elements to at the very list find it.
			// Tho we find that if we use only one word for our
			// original text, the first element will be our
			// pronunciation, and in contrary, if we use only
			// more than one word, it will be our second one.
			// Tho we can't tell this for sure, because we don't
			// know if google will continue to send the data with
			// the same algorithm or not (but we are sure that
			// the order of the data WILL NOT change in the future,
			// in ProtoBuff, order of data matters after all.)
			continue
		} else if !p1Set {
			p1Set = true
		}

		if !wCheck && isW {
			wStr, find := extractTextStr(current)
			if find && !ws.IsEmpty(&wStr) {
				setWrongNess(wStr, wTr)
			}

			wCheck = true
			continue
		}

		// check if p2Set is true or not, if not, please try to
		// extract it from current element, if you couldn't extract it
		// at the end, you have to go and see if you can
		// extract it in next element or not.
		// translated pronunciation is mandatory, it SHOULD
		// exist in the data.
		if !p2Set && canBePronunciation(current) {
			wTr.TranslatedPronunciation, p2Set = getPronunciation(current)
			// TranslatedPronunciation field is mandatory, if you don't
			// find it at the first, you have to iterate over all
			// of the array elements to at the very list find it.
			// Tho we find that if we use only one word for our
			// original text, the first element will be our
			// pronunciation, and in contrary, if we use only
			// more than one word, it will be our second one.
			// Tho we can't tell this for sure, because we don't
			// know if google will continue to send the data with
			// the same algorithm or not (but we are sure that
			// the order of the data WILL NOT change in the future,
			// in ProtoBuff, order of data matters after all.)
			//log.Println(value[2])

			continue
		}

		if !tSet {
			if isSeparator(current, wTr) {
				tSet = true
				continue
			}

			if strings.HasSuffix(current, NullAndCama) {
				log.Println(current)
				for strings.HasSuffix(current, NullAndCama) {
					current = strings.TrimSuffix(current, NullAndCama)
				}
				log.Println(current)
				if ws.IsEmpty(&current) {
					continue
				}

				if !strings.HasPrefix(current, ws.STR_SIGN) ||
					!strings.HasSuffix(current, StringAndCama) {
					continue
				}

				log.Println("success: " + current)
			}

			current = strings.TrimSpace(current)
			if strings.HasPrefix(current, CAMA) &&
				strings.HasSuffix(current, CAMA) {
				//tmpCheck := strings.ReplaceAll(current, wTr.To, ws.EMPTY)
				//tmpCheck = strings.ReplaceAll(tmpCheck, wTr.From, ws.EMPTY)
				tmpCheck := trLang.RemoveShortsWithStrs(current)
				//tmpCheck = strings.ReplaceAll(tmpCheck, TwoStr, ws.EMPTY)
				tmpCheck = strings.ReplaceAll(tmpCheck, TwoCama, ws.EMPTY)
				tmpCheck = strings.ReplaceAll(tmpCheck, ws.SPACE_VALUE,
					ws.EMPTY)
				if ws.IsEmpty(&tmpCheck) {
					continue
				}

				_, err := strconv.Atoi(tmpCheck)
				if err == nil {
					continue
				}

			}

			if strings.HasSuffix(current, StringAndCama) {

				tmpStr, find := extractTextStr(current)
				if !find {
					continue
				}
				tmpStr = strings.TrimSpace(tmpStr)
				if ws.IsEmpty(&tmpStr) || wTr.alreadyExists(tmpStr) {
					continue
				}

				if wTr.From != wTr.To {
					if strings.EqualFold(tmpStr, wTr.UserText) ||
						strings.EqualFold(tmpStr, wTr.OriginalPronunciation) {
						continue
					}

					if strings.EqualFold(tmpStr, wTr.OriginalPronunciation) {
						wTr.TranslatedPronunciation = ws.EMPTY
					}
				}

				//logStr(tmpStr)
				wTr.Translations = append(wTr.Translations, tmpStr)
			}
		} else {
			break
		}

	}
}

func isWrongFrom(value []string, wTr *WotoTr) bool {
	if isSimpleWrongFrom(value, wTr) {
		return true
	}

	if len(value) <= baseTenIndex {
		return value == nil
	}

	return false
}

func isSimpleWrongFrom(value []string, wTr *WotoTr) bool {
	nullCount := ws.BaseIndex
	for _, current := range value {
		if nullCount >= baseTwoIndex {
			txt, find := extractTextStr(current)
			if !find {
				continue
			}

			if trLang.IsLang(txt) {
				short := trLang.ExtractShortLang(txt)
				if short != nil {
					if !strings.EqualFold(*short, wTr.From) {
						wTr.wrongFrom = true
						wTr.From = *short
						return true
					}
				}
			}
		}

		current = strings.TrimSpace(current)
		if strings.EqualFold(NullAndCama, current) {
			nullCount++
		}
	}

	return false
}

func setWrongNess(value string, wTr *WotoTr) {
	if wTr == nil {
		return
	}
	value = strings.TrimSpace(value)

	log.Println(value)
	part := ws.EMPTY
	whole := ws.EMPTY
	var parts []string
	myStr := ws.Qss(value)
	myStrs := myStr.SplitStr(WrongNessOpen, WrongNessClose)

	var j int
	var another string
	l := len(myStrs) - ws.BaseOneIndex
	index := strings.Index(value, WrongNessOpen)
	if index != ws.BaseIndex {
		index = ws.BaseOneIndex
		whole = strings.TrimSpace(myStrs[ws.BaseIndex].GetValue())
	}
	for i := index; i <= l; i += baseTwoIndex {
		j = i + ws.BaseOneIndex
		part = strings.TrimSpace(myStrs[i].GetValue())
		if j > l {
			another = ws.EMPTY
		} else {
			another = strings.TrimSpace(myStrs[j].GetValue())
		}
		parts = append(parts, part)
		if i != index {
			whole += SPACE_VALUE
		}

		if ws.IsEmpty(&another) {
			whole += part
		} else if ws.IsEmpty(&part) {
			whole += another
		} else {
			whole += part + ws.SPACE_VALUE + another
		}
	}

	wTr.Corrected = &Corrected{
		CorrectedParts: parts,
		CorrectedValue: whole,
	}
	wTr.HasWrongness = true
}

func IsKind(value string) bool {
	return value == "adjective"
}

func isWrongness(value string) bool {
	return strings.Contains(value, WrongNessOpen) &&
		strings.Contains(value, WrongNessClose)
}

func canBePronunciation(value string) bool {
	if strings.HasSuffix(value, CamaNullCama) {
		value = strings.TrimSuffix(value, CamaNullCama)
		return len(value) >= baseTwoIndex
	}

	return false
}

func isSeparator(value string, wTr *WotoTr) bool {
	if wTr == nil || wTr.isTrEmpty() {
		return false
	}

	left := CamaAndStr + wTr.To + StrAndCama
	right := CamaAndStr + wTr.From + StrAndCama
	b1 := strings.HasPrefix(value, left) ||
		strings.HasSuffix(value, left)
	b2 := strings.HasPrefix(value, right) ||
		strings.HasSuffix(value, right)
	return b1 && b2
}

func getPronunciation(value string) (str string, find bool) {
	str, find = extractTextStr(value)
	if !find {
		return
	}

	var final string
	lastBad := badIgnore
	for i, current := range str {
		if i == ws.BaseIndex {
			if isBadIgnore(current) ||
				isBad(current) || isForbiddenR(current) {
				continue
			}
		} else {
			if i == len(str)-ws.BaseOneIndex {
				// a forbidden character at the end of the
				// pronunciation string is not allowed
				if isForbiddenR(current) {
					break
				}
			}
			if lastBad == current && current != badIgnore {
				continue
			} else if isBad(current) {
				lastBad = current
			} else {
				if !unicode.IsSpace(current) {
					lastBad = badIgnore
				}
			}
		}

		final += string(current)
	}

	return strings.TrimSpace(final), find
}

func trGoogle(fr, to, text string) (str string, err error) {
	body := strings.NewReader(googleFQ(fr, to, purify(text)))
	req, err := http.NewRequest(requestType, gHostUrl, body)
	if err != nil {
		return
	}

	req.Header.Set(userAgentGKey, userAgentGValue)
	req.Header.Set(acceptGKey, acceptGValue)
	req.Header.Set(acceptLanguageGKey, acceptLanguageGValue)
	req.Header.Set(refererGKey, refererGValue)
	req.Header.Set(xSameDomainGKey, xSameDomainGValue)
	req.Header.Set(xGoogBatchExecuteBgrGKey, xGoogBatchExecuteBgrGValue)
	req.Header.Set(contentTypeGKey, contentTypeGValue)
	req.Header.Set(originGKey, originGValue)
	req.Header.Set(gDNTGKey, gDNTGValue)
	req.Header.Set(connectionGKey, connectionGValue)

	// please notice that we don't need to
	// set cookies header

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()

	var b []byte

	b, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	str = string(b)
	if ws.IsEmpty(&str) {
		err = errors.New("an unexpected error in trGo package" +
			": Respond body is empty")
		return
	}

	return string(b), nil
}

func purify(text string) string {
	text = strings.ReplaceAll(text, ws.BracketOpen, ws.ParaOpen)
	text = strings.ReplaceAll(text, ws.Bracketclose, ws.ParaClose)
	text = strings.ReplaceAll(text, ws.Star, ws.EMPTY)
	text = strings.ReplaceAll(text, ws.LineEscape, ws.SPACE_VALUE)
	text = strings.ReplaceAll(text, ws.R_ESCAPE, ws.SPACE_VALUE)
	text = strings.ReplaceAll(text, ws.DoubleQ, ws.DoubleQJ)

	return text
}

func googleFQ(fr, to, text string) string {
	//sUrl := url.PathEscape(text)
	//tUrl, _ := url.Parse(text)
	//sUrl := tUrl.String()
	//return "f.req=%5B%5B%5B%22MkEWBc%22%2C%22%5B%5B%5C%22How%20are%20you%5C%22%2C%5C%22auto%5C%22%2C%5C%22fa%5C%22%2Ctrue%5D%2C%5Bnull%5D%5D%22%2Cnull%2C%22generic%22%5D%5D%5D&"
	// [[["MkEWBc","[[\"Hello\",\"auto\",\"ja\",true],[null]]",null,"generic"]]]&
	//return "[[[\"MkEWBc\", \"[[\"Hello\",\"auto\",\"ja\",true],[null]]\",null,\"generic\"]]]&\""
	return fReqGValue1 + url.QueryEscape(text) +
		fReqGValue2 + fr +
		fReqGValue3 + to + fReqGValue4
}

func extractTextStr(value string) (str string, find bool) {
	log.Println(value)
	l := len(value) - ws.BaseOneIndex
	if l == ws.BaseIndex {
		return
	}

	pre := false
	for i, s := range value {
		if find {
			// check if s is a forbidden rune or not (such as '\\')
			if isForbiddenR(s) {
				if i == l {
					return ws.EMPTY, false // not found
				}
				if value[i+ws.BaseOneIndex] == ws.CHAR_STR {
					return // found
				}
			} else if s == ws.CHAR_STR {
				// it means before it, there wasn't any
				// back slash character, so we are safe!
				return // found
			}
			str += string(s)
			continue
		}
		if isForbiddenR(s) {
			if !pre {
				pre = true
			}
		} else if pre {
			if s == ws.CHAR_STR {
				find = true
			}
		} else if s == ws.CHAR_STR && !pre {
			find = true
		}
	}

	return // found
}

func isBadIgnore(r rune) bool {
	return r == badIgnore
}

func isBad(r rune) bool {
	return r == bad01 || r == bad02
}

func isForbiddenR(r rune) bool {
	return r == forbiddenR01
}
