package headers

import (
	"net/http"
	"strings"
)

type HeaderName = string

const (
	AIM                         HeaderName = "A-IM"
	Accept                      HeaderName = "Accept"
	AcceptCharset               HeaderName = "Accept-Charset"
	AcceptEncoding              HeaderName = "Accept-Encoding"
	AcceptLanguage              HeaderName = "Accept-Language"
	AcceptDatetime              HeaderName = "Accept-Datetime"
	AccessControlRequestMethod  HeaderName = "Access-Control-Request-Method"
	AccessControlRequestHeaders HeaderName = "Access-Control-Request-Headers"
	Authorization               HeaderName = "Authorization"
	CacheControl                HeaderName = "Cache-Control"
	Connection                  HeaderName = "Connection"
	ContentLength               HeaderName = "Content-Length"
	ContentMD5                  HeaderName = "Content-MD5"
	ContentType                 HeaderName = "Content-Type"
	Cookie                      HeaderName = "Cookie"
	Date                        HeaderName = "Date"
	Expect                      HeaderName = "Expect"
	Forwarded                   HeaderName = "Forwarded"
	From                        HeaderName = "From"
	Host                        HeaderName = "Host"
	IfMatch                     HeaderName = "If-Match"
	IfModifiedSince             HeaderName = "If-Modified-Since"
	IfNoneMatch                 HeaderName = "If-None-Match"
	IfRange                     HeaderName = "If-Range"
	IfUnmodifiedSince           HeaderName = "If-Unmodified-Since"
	MaxForwards                 HeaderName = "Max-Forwards"
	Origin                      HeaderName = "Origin"
	Pragma                      HeaderName = "Pragma"
	ProxyAuthorization          HeaderName = "Proxy-Authorization"
	Range                       HeaderName = "Range"
	Referer                     HeaderName = "Referer"
	TE                          HeaderName = "TE"
	UserAgent                   HeaderName = "User-Agent"
	Upgrade                     HeaderName = "Upgrade"
	Via                         HeaderName = "Via"
	Warning                     HeaderName = "Warning"
)

func Concat(h1, h2 http.Header) http.Header {
	for h, val := range h2 {
		h1.Add(h, strings.Join(val, ","))
	}
	return h1
}
