package hashmap

import (
	"errors"
	"fmt"

	"github.com/cespare/xxhash"
)

const (
	maxBucketLen        = 4
	defaultBucketsCount = 8
)

var (
	ErrBucketIsFull = errors.New("bucket is full")
)

type HashMap struct {
	buckets []Bucket
}

func NewHashMap() *HashMap {
	bucks := make([]Bucket, defaultBucketsCount)

	for i := range bucks {
		bucks[i] = make(Bucket, 0, maxBucketLen)
	}

	return &HashMap{
		buckets: bucks,
	}
}

// TODO: any type
func (h *HashMap) Set(key, val string) error {

	// find bucket id and create item
	// WARN: it may panic (if len() will return 0)

	err := h.set(key, val)
	// if bucket is full -- increase the map size (more buckets) and try again
	for errors.Is(err, ErrBucketIsFull) {
		h.resize()
		err = h.set(key, val)
	}
	return nil
}

func (h *HashMap) set(key, val string) error {
	// hash of key
	hash := getHash(key)
	buckIdx := hash % uint64(len(h.buckets))

	buck := h.buckets[buckIdx]
	item := hashItem{key: key, val: val}

	if i, ok := buck.find(item); ok {
		h.buckets[buckIdx][i] = item
	} else if len(buck) == maxBucketLen {
		return ErrBucketIsFull
	} else if len(buck) > maxBucketLen {
		panic("bucket len is more then limit")
	} else {
		h.buckets[buckIdx] = append(h.buckets[buckIdx], item)
	}

	return nil
}

func (h *HashMap) Get(key string) {}

func (h *HashMap) BucketCount() int {
	return len(h.buckets)
}

func (h *HashMap) resize() {
	// if 0 buckets - create 8 bucket
	panic("unimplemented")
}

type Bucket []hashItem

// find finds target in bucket by it's key. Returns (index, true) if found.
// Returns (0, false) if not found.
func (b Bucket) find(target hashItem) (int, bool) {
	index := 0
	found := false

	for i, item := range b {
		if item.key == target.key {
			index = i
			found = true
			break
		}
	}
	return index, found
}

type hashItem struct {
	key string
	val string
}

func getHash(v any) uint64 {
	s := str(v)
	return xxhash.Sum64String(s)
}

func str(v any) string {
	switch v.(type) {
	case string:
		return v.(string)
	default:
		return fmt.Sprintf("%v", v)
	}
}
