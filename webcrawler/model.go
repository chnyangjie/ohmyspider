package webcrawler

type CrawlRequest struct {
}
type HTTPRequest struct {
	URL       string
	Method    HTTPMethod
	URLParams map[string]string
	Headers   map[string]string
	Callback  CrawlCallback
	Proxy     *Proxy
	Extra     map[string]interface{}
}
type Proxy struct {
	Host string
	Port int
}

type CrawlCallback func(request HTTPRequest, response []byte)

type HTTPMethod string

var (
	HTTP_METHOD_GET  = HTTPMethod("GET")
	HTTP_METHOD_POST = HTTPMethod("POST")
)
