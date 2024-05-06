# Prefix API server

## Launch the server

`docker compose up` 

In compose, the server will start with a PostgreSQL database backend.

Without docker using for instance `go run ./cmd/main.go` it start with a local in memory backend. 

## API

- `GET /words/` returns 200 and a list of words and their occurrence if any.

- `POST /words/{word}` returns 201 upon successful word insertion

- `GET /words/{prefix}` returns 200 if a matching prefix is found, 404 otherwise.

## Design

2 kinds of data store are implemented: `InMemory` and `PostGres`.

## In Memory 

The ideal data structure would be a **trie** (prefix tree), maybe in a next iteration. I choosed something simpler based on a combination of array and hashtable:

- The insertion complexity is `O(1)` (append to a slice, add a map entry or increment an existing one).

- The prefix search complexity is `O(n*log(n))` (sort the slice and binary search the first prefix).

If I had to make it production ready, I would:

- Think of an efficient way to persist the data in the file system.

- Implement a LRU policy admission controller to control the number of keys and avoid storage and memory exhaustion.

- Think about a solution to dispatch the words with a consistent hashing algorithm to be able to scale out with more replicas.

## PostGres

- No need to manage concurrency, the database connexion pool is concurrently safe.

- The heavy lifting is delegated to the database and the data is persistent.

- Indexing doesn't work with the wildcard matching operator `%`...
