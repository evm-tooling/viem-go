package utils_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/ChefBingbong/viem-go/utils"
)

var _ = Describe("LruMap", func() {
	Describe("Basic operations", func() {
		It("should set and get values", func() {
			cache := utils.NewLruMap[int](10)
			cache.Set("key1", 42)

			value, ok := cache.Get("key1")
			Expect(ok).To(BeTrue())
			Expect(value).To(Equal(42))
		})

		It("should return false for missing keys", func() {
			cache := utils.NewLruMap[int](10)

			_, ok := cache.Get("missing")
			Expect(ok).To(BeFalse())
		})

		It("should update existing keys", func() {
			cache := utils.NewLruMap[int](10)
			cache.Set("key1", 1)
			cache.Set("key1", 2)

			value, ok := cache.Get("key1")
			Expect(ok).To(BeTrue())
			Expect(value).To(Equal(2))
			Expect(cache.Size()).To(Equal(1))
		})

		It("should delete keys", func() {
			cache := utils.NewLruMap[int](10)
			cache.Set("key1", 42)

			deleted := cache.Delete("key1")
			Expect(deleted).To(BeTrue())

			_, ok := cache.Get("key1")
			Expect(ok).To(BeFalse())
		})

		It("should check if key exists", func() {
			cache := utils.NewLruMap[int](10)
			cache.Set("key1", 42)

			Expect(cache.Has("key1")).To(BeTrue())
			Expect(cache.Has("key2")).To(BeFalse())
		})
	})

	Describe("LRU eviction", func() {
		It("should evict oldest item when over capacity", func() {
			cache := utils.NewLruMap[int](3)
			cache.Set("a", 1)
			cache.Set("b", 2)
			cache.Set("c", 3)

			Expect(cache.Size()).To(Equal(3))

			// Adding a new item should evict "a" (oldest)
			cache.Set("d", 4)

			Expect(cache.Size()).To(Equal(3))
			Expect(cache.Has("a")).To(BeFalse())
			Expect(cache.Has("b")).To(BeTrue())
			Expect(cache.Has("c")).To(BeTrue())
			Expect(cache.Has("d")).To(BeTrue())
		})

		It("should refresh item on get", func() {
			cache := utils.NewLruMap[int](3)
			cache.Set("a", 1)
			cache.Set("b", 2)
			cache.Set("c", 3)

			// Access "a" to make it recently used
			cache.Get("a")

			// Add new item - should evict "b" now (oldest after "a" was accessed)
			cache.Set("d", 4)

			Expect(cache.Has("a")).To(BeTrue()) // "a" was refreshed
			Expect(cache.Has("b")).To(BeFalse()) // "b" was evicted
			Expect(cache.Has("c")).To(BeTrue())
			Expect(cache.Has("d")).To(BeTrue())
		})

		It("should refresh item on set update", func() {
			cache := utils.NewLruMap[int](3)
			cache.Set("a", 1)
			cache.Set("b", 2)
			cache.Set("c", 3)

			// Update "a" to make it recently used
			cache.Set("a", 10)

			// Add new item - should evict "b"
			cache.Set("d", 4)

			Expect(cache.Has("a")).To(BeTrue())
			Expect(cache.Has("b")).To(BeFalse())
		})
	})

	Describe("Keys", func() {
		It("should return keys in MRU order", func() {
			cache := utils.NewLruMap[int](5)
			cache.Set("a", 1)
			cache.Set("b", 2)
			cache.Set("c", 3)

			keys := cache.Keys()
			Expect(keys).To(Equal([]string{"c", "b", "a"}))
		})

		It("should reflect access order", func() {
			cache := utils.NewLruMap[int](5)
			cache.Set("a", 1)
			cache.Set("b", 2)
			cache.Set("c", 3)

			cache.Get("a") // Access "a"

			keys := cache.Keys()
			Expect(keys).To(Equal([]string{"a", "c", "b"}))
		})
	})

	Describe("Clear", func() {
		It("should remove all items", func() {
			cache := utils.NewLruMap[int](5)
			cache.Set("a", 1)
			cache.Set("b", 2)
			cache.Set("c", 3)

			cache.Clear()

			Expect(cache.Size()).To(Equal(0))
			Expect(cache.Has("a")).To(BeFalse())
		})
	})

	Describe("String values", func() {
		It("should work with string values", func() {
			cache := utils.NewLruMap[string](5)
			cache.Set("key1", "value1")
			cache.Set("key2", "value2")

			value, ok := cache.Get("key1")
			Expect(ok).To(BeTrue())
			Expect(value).To(Equal("value1"))
		})
	})
})
