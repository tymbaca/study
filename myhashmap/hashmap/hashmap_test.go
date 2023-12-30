package hashmap

import (
	"fmt"
	"math/rand"
	"strconv"
	"testing"
)

func Test_Set(t *testing.T) {
	m := New()
	count := 100000
	for i := 0; i < count; i++ {
		key := fmt.Sprintf("%d", i)
		val := fmt.Sprintf("%d", rand.Intn(60))
		m.Set(key, val)
	}
	for i := 0; i < 100; i++ {
		key := fmt.Sprintf("%d", rand.Intn(count))
		val := fmt.Sprintf("%d", rand.Intn(60))
		m.Set(key, val)
	}
	fmt.Printf("items count: %d\n", len(m.Items()))
	fmt.Printf("buckets count: %d\n", len(m.buckets))
	fmt.Printf("first 100 buckets: %v ...\n", m.buckets[:100])
	if len(m.buckets) > count*2 {
		t.Logf("Too much buckets: %d", len(m.buckets))
	}
}

func Test_All(t *testing.T) {
	m := New()
	testKey := "Hello"
	testVal := "World"

	val, ok := m.Get(testKey)
	if ok {
		t.Errorf("unexpectedly found")
	}

	val, ok = m.Remove(testKey)
	if ok {
		t.Errorf("unexpectedly found and deleted")
	}

	m.Set(testKey, testVal)
	val, ok = m.Get(testKey)
	if !ok {
		t.Errorf("unexpectedly not found")
	}
	if val != testVal {
		t.Errorf("unmatch vals, expected: %s, found: %s", testVal, val)
	}

	newVal := "Mom"
	m.Set(testKey, newVal)
	val, ok = m.Get(testKey)
	if !ok {
		t.Errorf("unexpectedly not found")
	}
	if val != newVal {
		t.Errorf("unmatch vals, expected: %s, found: %s", newVal, val)
	}

	val, ok = m.Remove(testKey)
	if !ok {
		t.Errorf("unexpectedly not found and not deleted")
	}
}

func Test_Get(t *testing.T) {
	m := New()
	vals := fill(m, 10)
	for key, expectedVal := range vals {
		actualVal, ok := m.Get(strconv.Itoa(key))
		if !ok {
			t.Fatalf("cannot find %s key", strconv.Itoa(key))
		}
		if actualVal != expectedVal {
			t.Errorf(
				"value unmatch: expected: %s, got: %s",
				expectedVal,
				actualVal,
			)
		}
	}
}

func Test_Remove(t *testing.T) {
	m := New()
	count := 10
	expectedVals := fill(m, count)

	for key := 0; key < count; key++ {
		val, ok := m.Remove(strconv.Itoa(key))
		if !ok {
			t.Errorf("Must found %d but didn't", key)
		} else if val != expectedVals[key] {
			t.Errorf(
				"value unmatch: expected: %s, got: %s",
				expectedVals[key],
				val,
			)
		}
		t.Logf("removed key: %v got value: %v", key, val)
	}

	// Get nonexistent key
	key := strconv.Itoa(count + 10)
	val, ok := m.Get(key)
	if ok {
		t.Errorf("error: found by nonexistent key: %s, got val: %s", key, val)
	}
}

// Will return inserted values. Values will be inserted in "0", "1", "2", ... "n" keys.
func fill(m *HashMap, n int) []string {
	vals := []string{}
	for i := 0; i < n; i++ {
		key := fmt.Sprintf("%d", i)
		val := fmt.Sprintf("%d", rand.Intn(60))
		m.Set(key, val)
		vals = append(vals, val)
	}
	return vals
}
