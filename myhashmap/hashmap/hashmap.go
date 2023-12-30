package hashmap

import (
	"errors"
	"fmt"

	"github.com/cespare/xxhash"
)

const (
	maxBucketLen        = 8
	defaultBucketsCount = 8
)

var (
	ErrBucketIsFull = errors.New("bucket is full")
)

type T = string

type HashMap struct {
	buckets []Bucket
}

func New() *HashMap {
	bucks := make([]Bucket, defaultBucketsCount)

	for i := range bucks {
		bucks[i] = make(Bucket, 0, maxBucketLen)
	}

	return &HashMap{
		buckets: bucks,
	}
}

// n is bucket count
func newHashMapN(n int) *HashMap {
	bucks := make([]Bucket, n)

	for i := range bucks {
		bucks[i] = make(Bucket, 0, maxBucketLen)
	}

	return &HashMap{
		buckets: bucks,
	}
}

// TODO: any type
func (h *HashMap) Set(key, val T) {

	// find bucket id and create item
	// WARN: it may panic (if len() will return 0)

	countUntilPanic := 0
	err := h.set(key, val)
	// if bucket is full -- increase the map size (more buckets) and try again
	for errors.Is(err, ErrBucketIsFull) {
		h.resizeAndEvac()
		err = h.set(key, val)

		if countUntilPanic > 4 {
			panic("you are stuck at evac :)")
		}
		countUntilPanic++
	}
}

func (h *HashMap) set(key, val T) error {
	buckIdx := h.getBucketIndex(key)

	buck := h.buckets[buckIdx]
	item := hashItem{key: key, val: val}

	i, ok := buck.find(key)
	if ok {
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

func (h *HashMap) Get(key T) (T, bool) {
	buckIdx := h.getBucketIndex(key)
	bucket := h.buckets[buckIdx]
	i, ok := bucket.find(key)
	if !ok {
		return "", false
	}

	item := bucket[i]
	return item.val, true
}

// Remove removes item and returns val and true, Returns false if key
// is not in HashMap
func (h *HashMap) Remove(key T) (T, bool) {
	buckIdx := h.getBucketIndex(key)
	bucket := h.buckets[buckIdx]
	i, ok := bucket.find(key)
	if !ok {
		return "", false
	}

	oldVal := bucket[i].val
	// Extra check
	if bucket[i].key != key {
		panic("unmatch keys while removing")
	}

	// Remove item at index i by re-slicing
	// h.buckets[buckIdx] = append(
	// 	h.buckets[buckIdx][:i],
	// 	h.buckets[buckIdx][i+1:]...,
	// )
	bucket = append(
		bucket[:i],
		bucket[i+1:]...,
	)

	return oldVal, true
}

func (h *HashMap) BucketCount() int {
	return len(h.buckets)
}

func (h *HashMap) resizeAndEvac() {
	evac(h)
}

func (h *HashMap) Items() []hashItem {
	items := []hashItem{}
	for _, buck := range h.buckets {
		for _, item := range buck {
			items = append(items, item)
		}
	}
	return items
}

func (h *HashMap) getBucketIndex(key T) uint64 {
	hash := getHash(key)
	buckIdx := hash % uint64(len(h.buckets))
	return buckIdx
}

type Bucket []hashItem

// find finds target in bucket by it's key. Returns (index, true) if found.
// Returns (0, false) if not found.
func (b Bucket) find(target T) (int, bool) {
	index := 0
	found := false

	for i, item := range b {
		if item.key == target {
			index = i
			found = true
			break
		}
	}
	return index, found
}

type hashItem struct {
	key T
	val T
}

func getHash(v any) uint64 {
	s := str(v)
	h := xxhash.New()
	h.Write([]byte(s))
	extra := 0
	// Hash itself extra times
	for i := 0; i < extra; i++ {
		h.Write(h.Sum(nil))
	}
	return h.Sum64()
}

func str(v any) string {
	switch v.(type) {
	case string:
		return v.(string)
	default:
		return fmt.Sprintf("%v", v)
	}
}

// evac creates x4 more buckets and evacuates all items to that buckets (via Set() method)
// Reassigns HashMap.buckets to new buckets slice
func evac(old *HashMap) {
	// if 0 buckets - create 8 bucket
	prevCount := len(old.buckets)
	// default factor
	factor := 8
	if prevCount > 1000 {
		factor = 4
	}
	if prevCount > 20000 {
		factor = 2
	}
	newHashMap := newHashMapN(prevCount * factor)

	items := old.Items()
	for _, item := range items {
		err := newHashMap.set(item.key, item.val)
		if err != nil {
			panic(err)
		}
	}
	old.buckets = newHashMap.buckets
	return
}
