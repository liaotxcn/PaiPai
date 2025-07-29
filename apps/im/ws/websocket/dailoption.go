/**
 * @author: dn-jinmin/dn-jinmin
 * @doc:
 */

package websocket

import "net/http"

type DailOptions func(option *dailOption)

type dailOption struct {
	pattern string
	header  http.Header
}

func newDailOptions(opts ...DailOptions) dailOption {
	o := dailOption{
		pattern: "/ws",
		header:  nil,
	}

	for _, opt := range opts {
		opt(&o)
	}

	return o
}
func WithClientPatten(pattern string) DailOptions {
	return func(opt *dailOption) {
		opt.pattern = pattern
	}
}
func WithClientHeader(header http.Header) DailOptions {
	return func(opt *dailOption) {
		opt.header = header
	}
}


