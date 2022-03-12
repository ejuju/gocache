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
type DefaultRequestIdentifier struct {
	TTL time.Duration
}

//
func (dri *DefaultRequestIdentifier) Identify(r *http.Request) string {
	return r.UserAgent()
}

type CachedResponse struct {
	StatusCode int
	Header     http.Header
	Cookies    []*http.Cookie
	Body       []byte
}

// WrapHandler can be used to filter requests so that if they have a cached response already they are answered directly
func (c *Cache) WrapHandler(in http.Handler, identifier RequestIdentifier, cacheConfig CacheConfig) (http.Handler, error) {
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

		// get cached response
		item, err := c.ReadOne(ReadOneRequest{
			ItemID: id,
		})

		// handle cache read err
		if err != nil && !errors.Is(err, ErrUnknownID) {
			c.logErr(err)
			// serve http if error
			in.ServeHTTP(w, r)
			return
		}

		// first request with this id
		isFirstID := errors.Is(err, ErrUnknownID)
		if isFirstID {
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
			if err != nil {
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
				return
			}
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

		// handle decoding err
		if err != nil {
			c.logErr(err)
			// serve http if error
			in.ServeHTTP(w, r)
			return
		}

		// respond
		// set cookies
		for _, cookie := range cachedres.Cookies {
			http.SetCookie(w, cookie)
		}
		// set status code
		w.WriteHeader(cachedres.StatusCode)
		// set header fields
		for key, values := range cachedres.Header {
			for _, val := range values {
				w.Header().Add(key, val)
			}
		}
		// set body
		w.Write(cachedres.Body)

	}), nil
}
