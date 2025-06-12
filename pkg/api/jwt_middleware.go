package api

import (
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/valyala/fasthttp"
)

var JWTSecret string

// JWTMiddleware is a FastHTTP middleware that checks for a valid JWT in the Authorization header.
func JWTMiddleware(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		header := string(ctx.Request.Header.Peek("Authorization"))
		if !strings.HasPrefix(header, "Bearer ") {
			ctx.SetStatusCode(fasthttp.StatusUnauthorized)
			ctx.SetBodyString(`{"error":"missing or invalid Authorization header"}`)
			return
		}
		tokenString := strings.TrimPrefix(header, "Bearer ")
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(JWTSecret), nil
		})
		if err != nil || !token.Valid {
			ctx.SetStatusCode(fasthttp.StatusUnauthorized)
			ctx.SetBodyString(`{"error":"invalid or expired token"}`)
			return
		}
		next(ctx)
	}
}
