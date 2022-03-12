# Experimental Golang cache for HTTP response cache

`gdb` is a JSON-based data storage engine

Planned features:

## v1.0

- [ ] In memory cache
- [ ] Caching large data to temporary files
- [ ] Storing mutations in a history to be able to reconstruct database from a history

## v2.0

- [ ] Make gdb a http.Handler to run database on standalone server/routine if needed
- [ ] Allow for easily subscribing to updates via websocket
