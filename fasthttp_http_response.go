package main

import "github.com/valyala/fasthttp"

type fasthttpHTTPResponse struct {
	response *fasthttp.Response
}

func newFasthttpHTTPResponse(r *fasthttp.Response) httpResponse {
	return fasthttpHTTPResponse{r}
}

func (r fasthttpHTTPResponse) StatusCode() int {
	return r.response.StatusCode()
}

func (r fasthttpHTTPResponse) Header(key string) string {
	return string(r.response.Header.Peek(key))
}

func (r fasthttpHTTPResponse) Body() []byte {
	return r.response.Body()
}
