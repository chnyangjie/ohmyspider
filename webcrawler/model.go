package webcrawler

type CrawlRequest struct {
}
type HTTPRequest struct {
	URL       string
	Method    HTTPMethod
	URLParams map[string]string
	Callback  CrawlCallback
	Proxy     *Proxy
	Extra     map[string]interface{}
}
type Proxy struct {
	Host string
	Port int
}

type HTTPMethod string

type CrawlCallback func(request HTTPRequest, response []byte)

var (
	HTTP_METHOD_GET  = HTTPMethod("GET")
	HTTP_METHOD_POST = HTTPMethod("POST")
)
