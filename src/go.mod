module github.com/alexkarpovich/lst-api/src

go 1.17

require golang.org/x/crypto v0.0.0-20211117183948-ae814b36b871

require (
	github.com/go-redis/redis/v8 v8.11.4
	github.com/golang-jwt/jwt/v4 v4.2.0
	github.com/gorilla/handlers v1.5.1
	github.com/gorilla/mux v1.8.0
	github.com/jmoiron/sqlx v1.3.4
	github.com/justinas/alice v1.2.0
	github.com/lib/pq v1.10.4
)

require gopkg.in/alexcesaro/quotedprintable.v3 v3.0.0-20150716171945-2caba252f4dc // indirect

require (
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/felixge/httpsnoop v1.0.1 // indirect
	gopkg.in/gomail.v2 v2.0.0-20160411212932-81ebce5c23df
)
