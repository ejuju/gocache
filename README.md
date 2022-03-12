# Simple cache for Golang

---

Package `gocache` provides a simple cache for web services

- It uses in-memory storage by default (in a `map[string]*Item`) & temporary disk storage for items larger than a certain treshold
- It is concurrency-friendly (using sync.RWMutex)
- Key-Value
- It encodes data as JSON under the hood
- It tries to be as developer-friendly as possible

---

### Features

#### Storage

- [x] Set lifespan on data
- [x] In memory cache
- [x] Caching large data to temporary files
- [ ] Something similar to SQL transactions (rollback / commit)

#### Web

- [x] Handle HTTP request/response caching (by using some of the request data to generate unique ID and storing prepared HTTP response)

#### Logging

- [ ] Request log
- [ ] Error log

#### Querying

- [ ] Filters for read requests
