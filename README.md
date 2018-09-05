# inmemcache

inmemcache is an in-memory key:value store/cache similar to memcached that is
suitable for applications running on a single machine. This may typically be
used for testing applications that use memcache.

Under the hood, this library uses go get github.com/patrickmn/go-cache

### Installation

`go get github.com/davidbyttow/inmemcache`

### Usage

```go
import (
  "github.com/bradfitz/gomemcache/memcache"
	"github.com/davidbyttow/inmemcache"
)

func main() {
	mc := inmemcache.New()
  mc.Set(&memcache.Item{Key: "foo", Value: []byte("my value")})

  it, err := mc.Get("foo")
  ...
}
```

### Reference

`godoc` or https://godoc.org/github.com/bradfitz/gomemcache/memcache
