package shm

import (
	"testing"
)

func Benchmark_shm_WriteOffset(b *testing.B) {
	// Create a new shared memory segment with Key=1 and Size=1024
	opts := Vopts{
		Key:  1,
		Size: 50,
	}
	_ = NewShm(opts)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = WriteOffset(1, 9223372036854775807)
	}

	// Delete the shared memory segment with Key=1
	_ = DeleteShm(1)
}

func Benchmark_shm_ReadOffset(b *testing.B) {
	// Create a new shared memory segment with Key=1 and Size=1024
	opts := Vopts{
		Key:  1,
		Size: 50,
	}
	_ = NewShm(opts)
	_ = WriteOffset(1, 9223372036854775807)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = ReadOffset(1)
	}

	// Delete the shared memory segment with Key=1
	_ = DeleteShm(1)
}
