package requests

import (
	"context"
	"crypto"
	"crypto/hmac"
	"encoding/base64"
	_ "golang.org/x/crypto/sha3"
	"net/http"
	"strconv"
	"strings"
)

var contextHmac = hmac.New(crypto.SHA3_256.New, []byte("DEFAULT DO NOT USE"))

// TODO: configure secret key
func Configure(key string) {
	contextHmac = hmac.New(crypto.SHA3_256.New, []byte(key))
}

type User struct {
	UserId int64
}

func SerializeUser(user *User) string {
	if user.UserId == 0 {
		return ""
	}
	data := strconv.FormatInt(user.UserId, 10)
	return withMac([]byte(data))
}

func DeserializeUser(dataWithMac string) *User {
	data, ok := authenticated(dataWithMac)
	if !ok {
		return nil
	}
	userId, err := strconv.ParseInt(string(data), 10, 64)
	if err != nil {
		// Authenticated the data but it was incorrect
		panic(err)
	}
	return &User{userId}
}

const userContextKey = "user"
const userCookieKey = "user"

func ClearOutgoingUser(ctx context.Context) {
	cookie := &http.Cookie{
		Name:   userCookieKey,
		MaxAge: -1,
	}
	AddOutgoingCookie(ctx, cookie)
}

func AddOutgoingUser(ctx context.Context, user *User) {
	cookie := &http.Cookie{
		Name:  userCookieKey,
		Value: SerializeUser(user),
		Path: "/",
		// TODO: flip to secure in production!
		Secure:   false,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}
	AddOutgoingCookie(ctx, cookie)
}

func WithIncomingUser(ctx context.Context) context.Context {
	user := &User{}
	cookie := GetCookie(ctx, userCookieKey)
	if cookie != nil {
		user = DeserializeUser(cookie.Value)
	}
	return context.WithValue(ctx, userContextKey, user)
}

func GetUser(ctx context.Context) *User {
	return ctx.Value(userContextKey).(*User)
}

func IsAnonymous(ctx context.Context) bool {
	return GetUser(ctx).UserId == 0
}

func computeMac(data []byte) string {
	mac := contextHmac.Sum(data)
	return base64.StdEncoding.EncodeToString(mac)
}

func withMac(data []byte) string {
	return base64.StdEncoding.EncodeToString(data) + ":" + computeMac(data)
}

func authenticated(dataAndMac string) ([]byte, bool) {
	split := strings.IndexRune(dataAndMac, ':')
	if split == -1 {
		return nil, false
	}
	encodedData := dataAndMac[:split]
	data, err := base64.StdEncoding.DecodeString(encodedData)
	if err != nil {
		return nil, false
	}
	mac := dataAndMac[split+1:]
	if mac != computeMac(data) {
		return nil, false
	}
	return data, true
}
