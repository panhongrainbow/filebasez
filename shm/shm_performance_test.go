package shm

import (
	"testing"
)

// Benchmark_shm_NewShm benchmarks the performance of writing offset to the shared memory segment.
func Benchmark_shm_WriteOffset(b *testing.B) {
	// The testShmKey is the shared memory key for testing
	var testShmKey int64 = 7

	// Create a new shared memory segment with Key=testShmKey and Size=1024
	opts := Vopts{
		Key:  testShmKey,
		Size: 50,
	}

	// Create a new shared memory segment
	_ = NewShm(opts)

	// Delete the shared memory segment with Key=testShmKey
	defer func() {
		_ = DeleteShm(testShmKey)
	}()

	// Reset the timer
	b.ResetTimer()

	// Write the offset 9223372036854775807 to the shared memory segment with Key=1 for b.N times
	for i := 0; i < b.N; i++ {
		_ = WriteOffset(testShmKey, 9223372036854775807)
	}
}

// Benchmark_shm_ReadOffset benchmarks the performance of reading offset from the shared memory segment.
func Benchmark_shm_ReadOffset(b *testing.B) {
	// The testShmKey is the shared memory key for testing
	var testShmKey int64 = 8

	// Create a new shared memory segment with Key=testShmKey and Size=1024
	opts := Vopts{
		Key:  testShmKey,
		Size: 50,
	}

	// Create a new shared memory segment
	_ = NewShm(opts)

	// Delete the shared memory segment with Key=testShmKey
	defer func() {
		_ = DeleteShm(testShmKey)
	}()

	// Write the offset 9223372036854775807 to the shared memory segment with Key=1
	_ = WriteOffset(testShmKey, 9223372036854775807)

	// Reset the timer
	b.ResetTimer()
	// Read the offset from the shared memory segment with Key=1 for b.N times
	for i := 0; i < b.N; i++ {
		_, _ = ReadOffset(testShmKey)
	}
}
