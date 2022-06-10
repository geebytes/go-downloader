package pb

import (
	"fmt"
	"net/url"
)

// type RequestProxy struct {
// 	ProxyUrl string
// }
// type RequestInstance struct {
// 	Url            string                 // Set request URL
// 	Header         map[string]string      // Set request header
// 	Method         string                 // Set request Method
// 	Body           []byte                 // Set request body
// 	Params         map[string]string      // Set request query params
// 	Proxy          *Proxy                 // Set request proxy addr
// 	Cookies        map[string]string      // Set request cookie
// 	Meta           map[string]interface{} // Set other data
// 	AllowRedirects bool                   // Set if allow redirects. default is true
// 	MaxRedirects   int                    // Set max allow redirects number
// 	BodyReader     io.Reader              // Set request body reader

// 	ID string
// }
type Option func(r *Request)

func RequestWithRequestProxy(proxy *Proxy) Option {
	return func(r *Request) {
		r.Proxy = proxy
	}
}

// func RequestWithRequestHeader(header map[string]string) Option {
// 	return func(r *Request) {

// 		r.Header = header
// 	}
// }
// func RequestWithRequestCookies(cookies map[string]string) Option {
// 	return func(r *Request) {
// 		r.Cookies = cookies
// 	}
// }

// func RequestWithRequestMeta(meta map[string]interface{}) Option {
// 	return func(r *Request) {
// 		r.Meta = meta
// 	}
// }
// func RequestWithAllowRedirects(allowRedirects bool) Option {
// 	return func(r *Request) {
// 		r.AllowRedirects = allowRedirects
// 		if !allowRedirects {
// 			r.MaxRedirects = 0
// 		}
// 	}
// }
// func RequestWithMaxRedirects(maxRedirects int) Option {
// 	return func(r *Request) {
// 		r.MaxRedirects = maxRedirects
// 		if maxRedirects <= 0 {
// 			r.AllowRedirects = false
// 		}
// 	}
// }

// updateQueryParams update url query  params
func (r *Request) updateQueryParams() {
	defer func() {
		if p := recover(); p != nil {
			fmt.Printf("panic recover! p: %v", p)
		}
	}()
	if len(r.Params) != 0 {
		u, err := url.Parse(r.Url)
		if err != nil {
			panic(fmt.Sprintf("set request query params err %s", err.Error()))
		}
		q := u.Query()
		for _, param := range r.Params {
			q.Set(param.Key, param.Value)
		}
		u.RawQuery = q.Encode()
		r.Url = u.String()
	}
}
func NewRequest(url string, method string, opts ...Option) *Request {
	r := &Request{
		Url:    url,
		Method: method,
	}

	for _, o := range opts {
		o(r)
	}
	r.updateQueryParams()
	return r

}
