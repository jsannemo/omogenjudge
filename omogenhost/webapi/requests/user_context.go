package requests

import (
	"context"
	"crypto"
	"crypto/hmac"
	"encoding/base64"
	_ "golang.org/x/crypto/sha3"
	"strconv"
	"strings"
)

const userKey = "user"

var contextHmac = hmac.New(crypto.SHA3_256.New, []byte("DEFAULT DO NOT USE"))

func Configure(key string) {
	contextHmac = hmac.New(crypto.SHA3_256.New, []byte(key))
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
		panic(err)
	}
	return &User{userId}
}

func WithUser(ctx context.Context, user *User) context.Context {
	return context.WithValue(ctx, userKey, user)
}

func GetUser(ctx context.Context) *User {
	return ctx.Value(userKey).(*User)
}

func IsAnonymous(ctx context.Context) bool {
	return GetUser(ctx).UserId == 0
}
