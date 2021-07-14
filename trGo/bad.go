package trGo

const (
	badIgnore = '.' // always ignore this bad
	bad01     = '?'
	bad02     = '!'
)

func isBadIgnore(r rune) bool {
	return r == badIgnore
}

func isBad(r rune) bool {
	return r == bad01 || r == bad02
}
