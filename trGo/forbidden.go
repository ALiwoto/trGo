package trGo

const (
	forbiddenR01 = '\\'
)

func isForbiddenR(r rune) bool {
	return r == forbiddenR01
}
