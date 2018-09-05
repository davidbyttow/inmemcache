package inmemcache_test

import (
	"testing"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/davidbyttow/inmemcache"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestAdd(t *testing.T) {
	mc := inmemcache.New()
	item := randomItem()
	assert.NoError(t, mc.Add(item))
	assert.Equal(t, memcache.ErrNotStored, mc.Add(item))
}

func TestCompareAndSwap(t *testing.T) {
	// TODO: Implement me
}

func TestDecrement(t *testing.T) {
	// TODO: Implement me
}

func TestDelete(t *testing.T) {
	mc := inmemcache.New()
	item := randomItem()
	assert.NoError(t, mc.Set(item))
	assert.NoError(t, mc.Delete(item.Key))
	_, err := mc.Get(item.Key)
	assert.Equal(t, memcache.ErrCacheMiss, err)
}

func TestDeleteAll(t *testing.T) {
	mc := inmemcache.New()
	items := make([]*memcache.Item, 10)
	for i := 0; i < 10; i++ {
		items[i] = randomItem()
		assert.NoError(t, mc.Add(items[i]))
	}
	assert.NoError(t, mc.DeleteAll())
	for _, item := range items {
		_, err := mc.Get(item.Key)
		assert.Equal(t, memcache.ErrCacheMiss, err)
	}
}

func TestFlushAll(t *testing.T) {
	mc := inmemcache.New()
	assert.NoError(t, mc.FlushAll())
}

func TestGet(t *testing.T) {
	mc := inmemcache.New()
	copy, err := mc.Get("missing")
	assert.Nil(t, copy)
	assert.Equal(t, err, memcache.ErrCacheMiss)
	item := randomItem()
	item.Expiration = 1
	assert.NoError(t, mc.Set(item))
	copy, err = mc.Get(item.Key)
	assert.NoError(t, err)
	assert.Equal(t, *copy, *item)
	// Poor man's expiration test :(
	time.Sleep(time.Millisecond * 1200)
	copy, err = mc.Get(item.Key)
	assert.Equal(t, err, memcache.ErrCacheMiss)
	assert.Nil(t, copy)
}

func TestGetMulti(t *testing.T) {
	mc := inmemcache.New()
	items := make([]*memcache.Item, 10)
	keys := make([]string, 0)
	for i := 0; i < 10; i++ {
		items[i] = randomItem()
		keys = append(keys, items[i].Key)
		assert.NoError(t, mc.Add(items[i]))
	}
	keys = append(keys, "missing")
	itemMap, err := mc.GetMulti(keys)
	assert.NoError(t, err)
	assert.Equal(t, 10, len(itemMap))
	for _, item := range items {
		assert.Equal(t, *item, *itemMap[item.Key])
	}
}

func TestIncrement(t *testing.T) {
	// TODO: Implement me
}

func TestReplace(t *testing.T) {
	mc := inmemcache.New()
	item := randomItem()
	assert.Equal(t, memcache.ErrNotStored, mc.Replace(item))
	assert.NoError(t, mc.Set(item))
	item.Value = []byte("replaced")
	assert.NoError(t, mc.Replace(item))
	copy, err := mc.Get(item.Key)
	assert.NoError(t, err)
	assert.Equal(t, *copy, *item)
}

func TestSet(t *testing.T) {
	mc := inmemcache.New()
	item := randomItem()
	assert.NoError(t, mc.Set(item))
	copy, err := mc.Get(item.Key)
	assert.NoError(t, err)
	assert.Equal(t, *copy, *item)
	copy.Expiration = 1
	assert.NoError(t, mc.Set(copy))
	copy, err = mc.Get(item.Key)
	assert.NoError(t, err)
	assert.Equal(t, int32(1), copy.Expiration)
}

func randomItem() *memcache.Item {
	return &memcache.Item{
		Key:   randUUID(),
		Value: []byte(randUUID()),
	}
}

func randUUID() string {
	return uuid.Must(uuid.NewRandom()).String()
}
