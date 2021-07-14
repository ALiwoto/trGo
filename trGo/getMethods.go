package trGo

import ws "github.com/ALiwoto/StrongStringGo/strongStringGo"

func (w *WotoTr) isTrEmpty() bool {
	for _, current := range w.TranslatedText {
		if !ws.IsEmpty(&current) {
			return false
		}
	}
	return true
}
