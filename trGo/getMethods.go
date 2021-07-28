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

//---------------------------------------------------------

func (l *Lang) GetBest() *LangDetect {
	if l.IsEmpty() {
		return nil
	}

	var best *LangDetect

	for _, d := range l.Data.Detections {
		if d.Reliable {
			if best != nil {
				if d.Confidence >= best.Confidence {
					best = &d
				}
			} else {
				best = &d
			}
		}
	}

	if best.Confidence < MinConfidence {
		return nil
	}

	return best
}

func (l *Lang) IsEmpty() bool {
	return l.Data == nil || len(l.Data.Detections) == 0
}
