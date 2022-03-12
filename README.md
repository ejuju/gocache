# Simple cache for Golang

---

Package `gocache` provides a simple cache

- It uses in-memory storage by default (in a `map[string]*Item`) & temporary disk storage for items larger than a certain treshold
- It is concurrency-friendly (using sync.RWMutex)
- Key-Value
- It encodes data as JSON under the hood
- It tries to be as developer-friendly as possible

---

### Roadmap

#### Storage

- [x] In memory cache
- [x] Caching large data to temporary files

#### Logging

- [ ] Request log

#### Querying

- [ ] Filters for read requests

#### Web

- [ ] Handle http request/response caching (by hashing request data for ID and storing response)
- [ ] Provide middleware / http.Handler to easily filter cached request
