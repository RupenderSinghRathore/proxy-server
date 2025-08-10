package response

import "net/http"

type Resp struct {
	Header http.Header
	Body   []byte
}
