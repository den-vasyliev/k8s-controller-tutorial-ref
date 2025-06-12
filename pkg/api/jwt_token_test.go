package api

import (
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/require"
	"github.com/valyala/fasthttp"
)

func TestTokenHandler_IssuesValidJWT(t *testing.T) {
	JWTSecret = "test-secret"
	ctx := &fasthttp.RequestCtx{}
	TokenHandler(ctx)
	resp := string(ctx.Response.Body())
	require.Contains(t, resp, "token")

	// Extract token
	tokenStr := resp[strings.Index(resp, ":")+2 : len(resp)-2]
	parsed, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(JWTSecret), nil
	})
	require.NoError(t, err)
	require.True(t, parsed.Valid)
}

func TestJWTMiddleware_ValidAndInvalidToken(t *testing.T) {
	JWTSecret = "test-secret"
	claims := jwt.MapClaims{
		"sub": "testuser",
		"exp": time.Now().Add(time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(JWTSecret))
	require.NoError(t, err)

	// Valid token
	ctx := &fasthttp.RequestCtx{}
	ctx.Request.Header.Set("Authorization", "Bearer "+tokenStr)
	called := false
	mw := JWTMiddleware(func(ctx *fasthttp.RequestCtx) { called = true })
	mw(ctx)
	require.True(t, called, "middleware should call next handler for valid token")

	// Invalid token
	ctx2 := &fasthttp.RequestCtx{}
	ctx2.Request.Header.Set("Authorization", "Bearer invalidtoken")
	called = false
	mw(ctx2)
	require.False(t, called, "middleware should not call next handler for invalid token")
	require.Equal(t, fasthttp.StatusUnauthorized, ctx2.Response.StatusCode())

	// Expired token
	expiredClaims := jwt.MapClaims{
		"sub": "testuser",
		"exp": time.Now().Add(-time.Hour).Unix(),
	}
	expiredToken := jwt.NewWithClaims(jwt.SigningMethodHS256, expiredClaims)
	expiredStr, err := expiredToken.SignedString([]byte(JWTSecret))
	require.NoError(t, err)
	ctx3 := &fasthttp.RequestCtx{}
	ctx3.Request.Header.Set("Authorization", "Bearer "+expiredStr)
	called = false
	mw(ctx3)
	require.False(t, called, "middleware should not call next handler for expired token")
	require.Equal(t, fasthttp.StatusUnauthorized, ctx3.Response.StatusCode())
}
