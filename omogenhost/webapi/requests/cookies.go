package requests

import (
	"context"
	"github.com/google/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"net/http"
)

const cookiesKey = "cookies"

func GetCookie(ctx context.Context, name string) *http.Cookie {
	cookies := ctx.Value(cookiesKey).([]*http.Cookie)
	for _, cookie := range cookies {
		if cookie.Name == name {
			return cookie
		}
	}
	return nil
}

func AddOutgoingCookie(ctx context.Context, cookie *http.Cookie) {
	logger.Infof("append cookie: %v", cookie.String())
	grpc.SetHeader(ctx, metadata.Pairs("set-cookie", cookie.String()))
}

func ParseIncomingCookies(ctx context.Context) context.Context {
	md, has := metadata.FromIncomingContext(ctx)
	if has {
		cookies := md.Get("cookie")
		h := http.Request{
			Header: http.Header{},
		}
		for _, cookie := range cookies {
			h.Header.Set("cookie", cookie)
		}
		return context.WithValue(ctx, cookiesKey, h.Cookies())
	}
	return ctx
}
