package middleware

type ctxKey string

const (
	Logger ctxKey = "logger"
	ReqID  ctxKey = "requestID"
)
