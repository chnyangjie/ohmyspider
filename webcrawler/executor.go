package webcrawler

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"sync"
)

var (
	CRAWLER_CHANNEL = make(chan HTTPRequest, 10000)
	wg              = sync.WaitGroup{}
)

func Start() {
	go func() {
		for {
			select {
			case req := <-CRAWLER_CHANNEL:
				{
					func() {
						defer wg.Done()
						log.Printf("Crawling %+v", req)
						resp, err := execute(req)
						if err != nil || resp == nil {
							log.Printf("Crawling %+v failed", req)
							return
						}
						if req.Callback != nil {
							req.Callback(req, resp)
						}
					}()
				}
			}
		}
	}()
}

func Crawl(req HTTPRequest) {
	wg.Add(1)
	CRAWLER_CHANNEL <- req
}

func Stop() {
	wg.Wait()
}
func execute(request HTTPRequest) ([]byte, error) {
	log.Printf("Executing %+v", request)
	client := &http.Client{}
	if request.Proxy != nil {
		proxyUrl, err := url.Parse(fmt.Sprintf("http://%v:%v", request.Proxy.Host, request.Proxy.Port))
		if err == nil {
			client.Transport = &http.Transport{
				Proxy: http.ProxyURL(proxyUrl),
			}
		}
	}
	values := url.Values{}
	_url := request.URL
	if len(request.URLParams) != 0 {
		// request.URLParams = map[string]string{}
		for k, v := range request.URLParams {
			values.Add(k, v)
		}
		paramsString := values.Encode()
		_url = fmt.Sprintf("%s?%s", request.URL, paramsString)
	}
	req, _ := http.NewRequest(string(request.Method), _url, nil)
	req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.4 Safari/605.1.15")
	parseFormErr := req.ParseForm()
	if parseFormErr != nil {
		log.Println(parseFormErr)
		return nil, parseFormErr
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Failure : ", err)
		return nil, err
	}
	respBody, _ := io.ReadAll(resp.Body)
	return respBody, nil
}
