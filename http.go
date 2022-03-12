package gocache

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"time"
)

//
type RequestIdentifier interface {
	Identify(r *http.Request) string
	ItemTTL() time.Duration
}

//
func NewRequestIdentifier(itemTTL time.Duration, identifyFunc func(*http.Request) string) *DefaultRequestIdentifier {
	return &DefaultRequestIdentifier{
		itemTTL:      itemTTL,
		identifyFunc: identifyFunc,
	}
}

//
type DefaultRequestIdentifier struct {
	itemTTL      time.Duration
	identifyFunc func(*http.Request) string
}

//
func (dri *DefaultRequestIdentifier) Identify(r *http.Request) string {
	return dri.identifyFunc(r)
}

//
func (dri *DefaultRequestIdentifier) ItemTTL() time.Duration {
	return dri.itemTTL
}

//
type CachedResponse struct {
	StatusCode int
	Header     http.Header
	Cookies    []*http.Cookie
	Body       []byte
}

// CacheHTTP can be used to filter requests so that if they have a cached response already they are answered directly
func CacheHTTP(in http.Handler, identifier RequestIdentifier, cacheConfig CacheConfig) (http.Handler, error) {
	// init response cache
	c, err := NewCache(cacheConfig)
	if err != nil {
		return nil, err
	}
	c.Start()
	defer c.Stop()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// identify request
		id := identifier.Identify(r)
		if id == "" {
			in.ServeHTTP(w, r)
			return
		}

		// get cached response
		item, err := c.ReadOne(ReadOneRequest{
			ItemID: id,
		})
		// if error serve http anyway but log err
		if err != nil && !errors.Is(err, ErrUnknownID) {
			c.logErr(err)
			in.ServeHTTP(w, r)
			return
		}

		// first request with this id
		if errors.Is(err, ErrUnknownID) {
			// record request result
			resrec := httptest.NewRecorder()
			// serve http and record response
			in.ServeHTTP(resrec, r)

			// write result to cache
			resp := resrec.Result()
			statusCode := resp.StatusCode
			header := resp.Header.Clone()
			cookies := resp.Cookies()
			body, err := ioutil.ReadAll(resp.Body)
			// err = no body?
			// set cookies
			for _, cookie := range cookies {
				http.SetCookie(w, cookie)
			}
			// set status codes
			w.WriteHeader(statusCode)
			// set header fields
			for key, values := range header {
				for _, val := range values {
					w.Header().Add(key, val)
				}
			}
			if err == nil {
				w.Write(body)
			}

			// cache response data
			c.WriteOne(WriteOneRequest{
				ItemID: id,
				Expiry: time.Now().Add(identifier.ItemTTL()),
				Value: CachedResponse{
					StatusCode: statusCode,
					Header:     header,
					Body:       body,
					Cookies:    cookies,
				},
			})
			return
		}

		// decode cached response
		cachedres := CachedResponse{}
		err = item.DecodeInto(&cachedres)
		// if error serve http anyway but log err
		if err != nil {
			c.logErr(err)
			in.ServeHTTP(w, r)
			return
		}

		// respond to request with cached data
		for _, cookie := range cachedres.Cookies {
			http.SetCookie(w, cookie)
		}
		w.WriteHeader(cachedres.StatusCode)
		for key, values := range cachedres.Header {
			for _, val := range values {
				w.Header().Add(key, val)
			}
		}
		w.Write(cachedres.Body)
	}), nil
}
