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
)

// https://telegra.ph/Lang-Codes-03-19-3

// Translate will translate the specified text value
// tp english.
func Translate(lang *Lang, to, text string) (*WotoTr, error) {
	if ws.IsEmpty(&text) {
		return nil, errors.New("text cannot be empty")
	}

	if lang == nil {
		lang = DetectLanguage(text)
	}

	if lang.Data == nil || len(lang.Data.Detections) == ws.BaseIndex {
		TranslateD(ws.AutoStr, to, text)
	}

	l1 := lang.Data.Detections[ws.BaseIndex]

	return TranslateD(l1.TheLang, to, text)
}

func TranslateD(fr, to, text string) (*WotoTr, error) {
	uText := strings.TrimSpace(text)

	var err error
	text, err = trGoogle(fr, to, text)
	if err != nil {
		return nil, err
	}

	w := WotoTr{
		UserText:     uText,
		OriginalText: text,
		From:         fr,
		To:           to,
	}

	return parseGData(&w)
}

func TrGnuTxt(fr, to, text string) (string, error) {
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
	text := wTr.OriginalText
	test := ws.Split(text, ws.BracketOpen, ws.Bracketclose)
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

	for _, s := range test {
		if accepted(s) {
			original = append(original, s)
		}
	}

	parseGparams(original, wTr)

	if wTr.WrongFrom {
		var err error
		text, err = trGoogle(wTr.From, wTr.To, wTr.UserText)
		if err != nil {
			return nil, err
		}

		w := WotoTr{
			OriginalText: text,
			From:         wTr.From,
			To:           wTr.To,
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
func AparseGparamsOLD(value []string, wTr *WotoTr) []string {
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
	if wTr.Road == nil {
		wTr.Road = make(map[int]bool)
	}

	index := ws.BaseIndex

	for _, c := range wTr.OriginalText {
		if string(c) == ws.LineEscape {
			wTr.Road[index] = false
		}
		if string(c) == ws.Point {
			wTr.Road[index] = true
		}
		index++
	}
	tmp := strings.Join(value, DY_WOTO_TEXT)

	tmp = strings.ReplaceAll(tmp, NullN, ws.EMPTY)
	tmp = strings.ReplaceAll(tmp, NullCValueR, ws.EMPTY)
	tmp = strings.ReplaceAll(tmp, GenericStr, ws.EMPTY)
	tmp = strings.ReplaceAll(tmp, NullCValue, ws.EMPTY)
	tmp = strings.ReplaceAll(tmp, NeQ, ws.EMPTY)
	strs := strings.Split(tmp, DY_WOTO_TEXT)
	final := make([]string, ws.BaseIndex)
	strMap := make(map[string]bool)
	lastStr := ws.EMPTY
	for _, current := range strs {
		tmp = current
		if current == lastStr {
			continue
		}

		if strings.HasPrefix(current, DoubleQS) {
			current = strings.TrimPrefix(current, DoubleQS)
		} else {
			lastStr = ws.EMPTY
			continue
		}

		if strings.Contains(current, MiddleWave) {
			current = strings.Split(current, MiddleWave)[ws.BaseOneIndex]
		}

		if strings.HasSuffix(current, DoubleQSP) {
			// optional
			current = strings.TrimSuffix(current, DoubleQSP)

			if strMap[current] {
				continue
			} else {
				strMap[current] = true
			}

			final = append(final, current)
			lastStr = tmp
			continue
		}

		if strings.HasSuffix(current, DoubleQS) {
			current = strings.TrimSuffix(current, DoubleQS)
		} else {
			lastStr = ws.EMPTY
			continue
		}

		if strMap[current] {
			continue
		} else {
			strMap[current] = true
		}

		lastStr = tmp
		// log.Println(current)
		// log.Println(strMap)
		final = append(final, current)
	}

	log.Println(value[1])
	return final
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
	if wTr.Road == nil {
		wTr.Road = make(map[int]bool)
	}

	index := ws.BaseIndex

	for _, c := range wTr.OriginalText {
		if string(c) == ws.LineEscape {
			wTr.Road[index] = false
		}
		if string(c) == ws.Point {
			wTr.Road[index] = true
		}
		index++
	}
	//tmp := strings.Join(value, DY_WOTO_TEXT)

	//tmp = strings.ReplaceAll(tmp, NullN, ws.EMPTY)
	//tmp = strings.ReplaceAll(tmp, NullCValueR, ws.EMPTY)
	//tmp = strings.ReplaceAll(tmp, GenericStr, ws.EMPTY)
	//tmp = strings.ReplaceAll(tmp, NullCValue, ws.EMPTY)
	//tmp = strings.ReplaceAll(tmp, NeQ, ws.EMPTY)
	//strs := strings.Split(tmp, DY_WOTO_TEXT)
	strMap := make(map[string]bool)
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
			log.Println(value[2])

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
			log.Println(value[2])

			continue
		}

		if !tSet {
			if isSeparator(current, wTr) {
				tSet = true
				continue
			}
			if strings.HasSuffix(current, NullAndCama) {
				continue
			} else if strings.HasSuffix(current, StrAndCama) {

				tmpStr, find := extractTextStr(current)
				if !find {
					continue
				}
				tmpStr = strings.TrimSpace(tmpStr)
				if ws.IsEmpty(&tmpStr) {
					continue
				}
				if wTr.From != wTr.To {
					if strings.EqualFold(tmpStr, wTr.UserText) {
						continue
					}
				}
				wTr.TranslatedText = append(wTr.TranslatedText, tmpStr)
			}
		} else {
			break
		}

		if strings.HasPrefix(current, DoubleQS) {
			current = strings.TrimPrefix(current, DoubleQS)
		} else {
			lastStr = ws.EMPTY
			continue
		}

		if strings.Contains(current, MiddleWave) {
			current = strings.Split(current, MiddleWave)[ws.BaseOneIndex]
		}

		if strings.HasSuffix(current, DoubleQSP) {
			// optional
			current = strings.TrimSuffix(current, DoubleQSP)

			if strMap[current] {
				continue
			} else {
				strMap[current] = true
			}

			lastStr = tmp
		}

		if strings.HasSuffix(current, DoubleQS) {
			current = strings.TrimSuffix(current, DoubleQS)
		} else {
			lastStr = ws.EMPTY
			continue
		}

		if strMap[current] {
			continue
		} else {
			strMap[current] = true
		}

		lastStr = tmp

	}

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
	wTr.HasWrongNess = true
}

func isWrongness(value string) bool {
	return strings.Contains(value, WrongNessOpen) &&
		strings.Contains(value, WrongNessClose)
}

func canBePronunciation(value string) bool {
	return strings.HasSuffix(value, CamaNullCama)
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

/* OLD_VERSION

func arrangeParams(values []string, wTr *WotoTr) {
	index := ws.BaseIndex
	for i, current := range values {
		if i == ws.BaseIndex {
			if trLang.IsLang(current) {
				if current != wTr.From {
					wTr.WrongFrom = true
					wTr.From = current
					return
				}
			}
		}
		if strings.Contains(current, WrongNessOpen) {
			wTr.HasWrongNess = true
			current = strings.ReplaceAll(current, WrongNessOpen, ws.EMPTY)
			current = strings.ReplaceAll(current, WrongNessClose, ws.EMPTY)
			current = strings.ReplaceAll(current, WrongNessClose, ws.EMPTY)
			current = strings.ReplaceAll(current, QuetUnicode, ws.SingleQ)
			current = strings.TrimPrefix(current, ws.BackSlash)
			current = strings.ReplaceAll(current, ws.BackSlash, ws.EMPTY)
			//wTr.CorrectedValue = current
		} else {
			if wTr.Road != nil {
				if !wTr.Road[index] {
					current += ws.LineEscape
				} else {
					current = strings.TrimPrefix(current, ws.LineEscape)
					current = strings.TrimSuffix(current, ws.LineEscape)
					current += ws.Point
				}
			}

			current = strings.ReplaceAll(current, ThreeE, ws.EMPTY)
			current = strings.ReplaceAll(current, QuetUnicode, ws.SingleQ)
			current = strings.ReplaceAll(current, CeeE, ws.EMPTY)
			current = strings.ReplaceAll(current, ws.DoubleBackSlash, ws.EMPTY)
			current = strings.ReplaceAll(current, ws.BackSlash, ws.EMPTY)
			wTr.TranslatedText = append(wTr.TranslatedText, current)
		}
	}
}

*/

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
	if strings.Contains(text, ws.BracketOpen) {
		text = strings.ReplaceAll(text, ws.BracketOpen, ws.ParaOpen)
	}
	if strings.Contains(text, ws.Bracketclose) {
		text = strings.ReplaceAll(text, ws.Bracketclose, ws.ParaClose)
	}
	if strings.Contains(text, ws.Star) {
		text = strings.ReplaceAll(text, ws.Star, ws.EMPTY)
	}
	if strings.Contains(text, ws.LineEscape) {
		text = strings.ReplaceAll(text, ws.LineEscape, ws.SPACE_VALUE)
	}
	if strings.Contains(text, ws.R_ESCAPE) {
		text = strings.ReplaceAll(text, ws.R_ESCAPE, ws.SPACE_VALUE)
	}
	if strings.Contains(text, ws.DoubleQ) {
		text = strings.ReplaceAll(text, ws.DoubleQ, ws.DoubleQJ)
	}
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
	l := len(value) - ws.BaseOneIndex
	if l == ws.BaseIndex {
		return
	}

	for i, s := range value {
		if find {
			if s == '\\' {
				if i == l {
					return ws.EMPTY, false // not found
				}
				if value[i+ws.BaseOneIndex] == '"' {
					return // found
				}
			} else if s == '"' {
				// it means before it, there wasn't any
				// back slash character, so we are safe!
				return // found
			}
			str += string(s)
			continue
		}
		if s == ws.CHAR_STR {
			find = true
		}
	}
	return
}
