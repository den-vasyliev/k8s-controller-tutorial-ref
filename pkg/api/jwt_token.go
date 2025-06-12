package api

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/valyala/fasthttp"
)

// TokenHandler issues a JWT for a test user (for local/dev use only).
func TokenHandler(ctx *fasthttp.RequestCtx) {
	claims := jwt.MapClaims{
		"sub": "testuser",
		"exp": time.Now().Add(time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(JWTSecret))
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.SetBodyString(`{"error":"could not generate token"}`)
		return
	}
	ctx.SetContentType("application/json")
	ctx.SetBodyString(`{"token":"` + tokenString + `"}`)
}
