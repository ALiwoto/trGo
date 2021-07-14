package trGo

import (
	"strings"

	ws "github.com/ALiwoto/StrongStringGo/strongStringGo"
)

func (w *WotoTr) isTrEmpty() bool {
	for _, current := range w.Translations {
		if !ws.IsEmpty(&current) {
			return false
		}
	}
	return true
}

func (w *WotoTr) alreadyExists(value string) bool {
	for _, current := range w.Translations {
		if strings.EqualFold(current, value) {
			return true
		}
	}
	return false
}
