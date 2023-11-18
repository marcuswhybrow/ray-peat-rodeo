package cache

import (
	"fmt"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

type HTTPCache struct {
	// Known values derrived from previous HTTP responses.
	cache map[string]map[string]string

	// Cached HTTP requests mapped to response handlers.
	//
	// This avoids duplicating in-flight HTTP requests.
	responders map[*http.Request][]Responder

	// Cached HTTP responses made in lieu of a cache entry.
	//
	// Uncached keys may use cached HTTP responses to generate their value.
	responses map[*http.Request]*http.Response

	// Tracking url/key requests to purge unused cache values.
	requestsMade map[string][]string

	// Tracks url/ley request that weren't found in cache
	requestsMissed map[string][]string
}

func NewHTTPCache(cache map[string]map[string]string) *HTTPCache {
	return &HTTPCache{
		cache:        cache,
		responders:   map[*http.Request][]Responder{},
		responses:    map[*http.Request]*http.Response{},
		requestsMade: map[string][]string{},
	}
}

// Returns cache purged of unrequest entries
func (h *HTTPCache) GetRequestsMade() map[string]map[string]string {
	requestsMade := map[string]map[string]string{}
	for url, keys := range h.requestsMade {
		requestsMade[url] = map[string]string{}
		for _, key := range keys {
			requestsMade[url][key] = h.cache[url][key]
		}
	}
	return requestsMade
}

func (h *HTTPCache) GetRequestsMissed() map[string][]string {
	return h.requestsMissed
}

func (h *HTTPCache) insert(url string, key string, val string) {
	keys, urlExists := h.cache[url]
	if !urlExists {
		keys = map[string]string{}
	}

	_, keyExists := keys[key]
	if keyExists {
		panic(fmt.Sprintf("HTTP cache key '%v' already exists for URL '%v'", key, url))
	}

	keys[key] = val
	h.cache[url] = keys
}

func (h *HTTPCache) GetJSON(url string, key string, handler ResponseHandler) chan string {
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("Accept", "application/json; charset=UTF-8")
	if err != nil {
		panic(fmt.Sprintf("Failed to instantiate JSON request for URL '%v': %v", url, err))
	}
	return h.request(req, key, handler)
}

func (h *HTTPCache) Get(url string, key string, handler ResponseHandler) chan string {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(fmt.Sprintf("Failed to instantiate GET request for URL '%v': %v", url, err))
	}
	return h.request(req, key, handler)
}

func (h *HTTPCache) request(req *http.Request, key string, handler ResponseHandler) chan string {
	valCh := make(chan string, 1)

	url := req.URL.String()

	h.requestsMade[url] = append(h.requestsMade[url], key)

	keysFromCache, ok := h.cache[url]
	if !ok {
		keysFromCache = map[string]string{}
	}

	val, exists := keysFromCache[key]
	if exists {
		valCh <- val
		return valCh
	}

	h.requestsMissed[url] = append(h.requestsMissed[url], key)

	res := h.responses[req]
	if res != nil {
		val := handler(res)
		h.insert(url, key, val)
		valCh <- val
		return valCh
	}

	// HTTP reponse pending
	responders, existingResponders := h.responders[req]
	responders = append(responders, Responder{Handler: handler, Channel: valCh})
	h.responders[req] = responders

	if !existingResponders {
		go func() {
			res, err := http.DefaultClient.Do(req)
			if err != nil {
				panic(fmt.Sprintf("Failed to GET HTTP response for URL: %v\n%v", url, err))
			}

			h.responses[req] = res
			for _, deferredHandler := range h.responders[req] {
				val := deferredHandler.Handler(res)
				h.insert(url, key, val)
				deferredHandler.Channel <- val
			}
		}()
	}

	return valCh
}

func GetH1OrTitle(res *http.Response) string {
	defer res.Body.Close()
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		panic(fmt.Sprintf("Failed to parse HTML response for url '%v': %v", res.Request.URL.String(), err))
	}
	selection := doc.Find("h1")
	if selection.Length() > 0 {
		return selection.First().Text()
	}

	return doc.Find("title").Text()
}

type Responder struct {
	Handler ResponseHandler
	Channel chan string
}

type ResponseHandler = func(*http.Response) string
