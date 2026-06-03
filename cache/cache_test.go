package cache

import (
	"sync"
	"testing"
	"time"
)

func TestMemoryDriver_PutGet(t *testing.T) {
	driver := NewMemoryDriver()

	if err := driver.Put("key1", "value1", 0); err != nil {
		t.Fatalf("Put failed: %v", err)
	}

	val, err := driver.Get("key1")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if val != "value1" {
		t.Errorf("expected 'value1', got %v", val)
	}
}

func TestMemoryDriver_GetMissing(t *testing.T) {
	driver := NewMemoryDriver()

	val, err := driver.Get("nonexistent")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != nil {
		t.Errorf("expected nil for missing key, got %v", val)
	}
}

func TestMemoryDriver_Forget(t *testing.T) {
	driver := NewMemoryDriver()
	driver.Put("key1", "value1", 0)

	if err := driver.Forget("key1"); err != nil {
		t.Fatalf("Forget failed: %v", err)
	}

	if driver.Has("key1") {
		t.Error("expected Has to return false after Forget")
	}
}

func TestMemoryDriver_Has(t *testing.T) {
	driver := NewMemoryDriver()

	if driver.Has("key1") {
		t.Error("expected Has to return false for missing key")
	}

	driver.Put("key1", "value1", 0)
	if !driver.Has("key1") {
		t.Error("expected Has to return true after Put")
	}
}

func TestMemoryDriver_Expiration(t *testing.T) {
	driver := NewMemoryDriver()

	driver.Put("key1", "value1", 1*time.Second)
	time.Sleep(1100 * time.Millisecond)

	val, err := driver.Get("key1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != nil {
		t.Errorf("expected nil for expired key, got %v", val)
	}
}

func TestMemoryDriver_Flush(t *testing.T) {
	driver := NewMemoryDriver()
	driver.Put("key1", "value1", 0)
	driver.Put("key2", "value2", 0)

	if err := driver.Flush(); err != nil {
		t.Fatalf("Flush failed: %v", err)
	}

	if driver.Has("key1") || driver.Has("key2") {
		t.Error("expected all keys to be removed after Flush")
	}
}

func TestMemoryDriver_IncrementDecrement(t *testing.T) {
	driver := NewMemoryDriver()
	driver.Put("counter", 10, 0)

	newVal, err := driver.Increment("counter", 5)
	if err != nil {
		t.Fatalf("Increment failed: %v", err)
	}
	if newVal != 15 {
		t.Errorf("expected 15, got %d", newVal)
	}

	newVal, err = driver.Decrement("counter", 3)
	if err != nil {
		t.Fatalf("Decrement failed: %v", err)
	}
	if newVal != 12 {
		t.Errorf("expected 12, got %d", newVal)
	}
}

func TestMemoryDriver_Forever(t *testing.T) {
	driver := NewMemoryDriver()

	if err := driver.Forever("key1", "value1"); err != nil {
		t.Fatalf("Forever failed: %v", err)
	}

	val, err := driver.Get("key1")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if val != "value1" {
		t.Errorf("expected 'value1', got %v", val)
	}
}

func TestMemoryDriver_Concurrent(t *testing.T) {
	driver := NewMemoryDriver()
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			key := "key"
			driver.Put(key, n, 0)
			driver.Get(key)
			driver.Increment(key, 1)
		}(i)
	}
	wg.Wait()
}

func TestLock_Get(t *testing.T) {
	driver := NewMemoryDriver()
	lock := NewLock(driver, "test-lock", 10*time.Second)

	acquired := lock.Get()
	if !acquired {
		t.Error("expected lock to be acquired")
	}

	// Try to acquire again - should return true (same owner)
	acquired2 := lock.Get()
	if !acquired2 {
		t.Error("expected lock to be acquired by same owner")
	}

	lock.Release()
}

func TestLock_Release(t *testing.T) {
	driver := NewMemoryDriver()
	lock := NewLock(driver, "test-lock", 10*time.Second)

	lock.Get()
	released := lock.Release()
	if !released {
		t.Error("expected Release to return true")
	}

	// Should be able to acquire again
	lock2 := NewLock(driver, "test-lock", 10*time.Second)
	if !lock2.Get() {
		t.Error("expected lock to be acquired after Release")
	}
}

func TestRateLimiter_TooManyAttempts(t *testing.T) {
	driver := NewMemoryDriver()
	limiter := NewRateLimiter(driver)

	for i := 0; i < 5; i++ {
		limiter.Hit("test-key", 60)
	}

	if !limiter.TooManyAttempts("test-key", 5) {
		t.Error("expected TooManyAttempts to return true after max attempts")
	}

	if limiter.TooManyAttempts("test-key", 10) {
		t.Error("expected TooManyAttempts to return false when under limit")
	}
}

func TestRateLimiter_Reset(t *testing.T) {
	driver := NewMemoryDriver()
	limiter := NewRateLimiter(driver)

	for i := 0; i < 5; i++ {
		limiter.Hit("test-key", 60)
	}

	limiter.Reset("test-key")

	if limiter.TooManyAttempts("test-key", 5) {
		t.Error("expected TooManyAttempts to return false after Reset")
	}
}

func TestSlidingWindow_Allow(t *testing.T) {
	driver := NewMemoryDriver()
	limiter := NewSlidingWindowRateLimiter(driver, time.Minute, 3)

	for i := 0; i < 3; i++ {
		if !limiter.Allow("test-key") {
			t.Errorf("expected Allow to return true for attempt %d", i+1)
		}
	}

	// 4th attempt should be blocked
	if limiter.Allow("test-key") {
		t.Error("expected Allow to return false after max attempts")
	}
}

func TestTokenBucket_Allow(t *testing.T) {
	limiter := NewTokenBucketRateLimiter(10, 5) // 10 tokens/sec, burst of 5

	for i := 0; i < 5; i++ {
		if !limiter.Allow("test-key") {
			t.Errorf("expected Allow to return true for attempt %d", i+1)
		}
	}

	// 6th attempt should be blocked (no tokens left)
	if limiter.Allow("test-key") {
		t.Error("expected Allow to return false after exhausting tokens")
	}
}

func TestTokenBucket_Refill(t *testing.T) {
	limiter := NewTokenBucketRateLimiter(100, 2) // Fast refill

	limiter.Allow("test-key")
	limiter.Allow("test-key")

	// Wait for refill
	time.Sleep(50 * time.Millisecond)

	if !limiter.Allow("test-key") {
		t.Error("expected Allow to return true after token refill")
	}
}
