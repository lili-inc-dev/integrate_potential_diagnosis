package constant

type contextKey uint64

const (
	CtxKeyAdmin contextKey = iota
	CtxKeyUser
)
