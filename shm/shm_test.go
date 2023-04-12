package shm

import (
	"github.com/stretchr/testify/require"
	"testing"
)

/*
Test_Check_Shm_Basic_Function checks the basic functionalities of a data segment implementation,
including valid and invalid cases, and checks for specific error messages under different conditions,
such as negative or zero size or negative key value.
*/
func Test_Check_Shm_Basic_Function(t *testing.T) {
	/*
		Test for a data segment implementation that tests basic functions such as write, read, and delete.
		It uses the Go testing package to check for errors and expected values.
	*/
	t.Run("test for basic functions in valid cases", func(t *testing.T) {
		// Set the options for the Vopts struct with key=1 and size=1024
		opts := Vopts{
			Key:  1,
			Size: 1024,
		}

		// Create a new segment with the specified options and return its id and any errors
		sg, err := newWithReturnId(opts)
		require.NoError(t, err)
		// Check that the segment's key matches the specified key (1) and its id is greater than 0
		require.Equal(t, int64(1), sg.key)
		require.Greater(t, sg.id, int64(0))

		// Write the bytes "abcde" to the segment with the specified id and return the length and any errors
		length, err := sg.writeWithId([]byte("abcde"))
		require.NoError(t, err)
		// Check that the length of the written bytes is 5
		require.Equal(t, int64(5), length)

		// Set the offset of the segment to 0 and create a new byte slice of length 5
		sg.offset = 0
		var data = make([]byte, 5)

		// Read the bytes from the segment with the specified id into the byte slice and return the length and any errors
		length, err = sg.readWithId(data)
		require.NoError(t, err)
		// Check that the length of the read bytes is 5
		require.Equal(t, int64(5), length)

		// Check that the bytes in the slice match the expected values
		require.Equal(t, byte('a'), data[0])
		require.Equal(t, byte('b'), data[1])
		require.Equal(t, byte('c'), data[2])
		require.Equal(t, byte('d'), data[3])
		require.Equal(t, byte('e'), data[4])

		// Delete the segment with the specified id and check for any errors
		err = sg.deleteWithId()
		require.NoError(t, err)
	})

	// Test basic functionality in invalid cases
	t.Run("Test basic functions in invalid cases", func(t *testing.T) {
		// Test negative shm key
		t.Run("negative shm key", func(t *testing.T) {
			// Set up Vopts struct with negative shm key
			opts := Vopts{
				Key: -1,
			}
			// Create new shared memory segment with specified options
			sg, err := newWithReturnId(opts)
			require.Equal(t, ErrNegativeOrZeroShmKey, err)
			// Attempt to delete the shared memory segment and ensure that the expected error is returned
			err = sg.deleteWithId()
			require.Equal(t, ErrShmEmptyPoint, err)
		})

		// Test zero shm key
		t.Run("zero shm key", func(t *testing.T) {
			// Set up Vopts struct with zero shm key
			opts := Vopts{
				Key: 0,
			}
			// Create new shared memory segment with specified options
			sg, err := newWithReturnId(opts)
			require.Equal(t, ErrNegativeOrZeroShmKey, err)
			// Attempt to delete the shared memory segment and ensure that the expected error is returned
			err = sg.deleteWithId()
			require.Equal(t, ErrShmEmptyPoint, err)
		})

		// Test negative shm size
		t.Run("negative shm size", func(t *testing.T) {
			// Set up Vopts struct with negative shm size
			opts := Vopts{
				Key:  1,
				Size: -1,
			}
			// Create new shared memory segment with specified options
			sg, err := newWithReturnId(opts)
			require.Equal(t, ErrNegativeOrZeroSize, err)
			// Attempt to delete the shared memory segment and ensure that the expected error is returned
			err = sg.deleteWithId()
			require.Equal(t, ErrShmEmptyPoint, err)
		})

		// Test zero shm size
		t.Run("zero shm size", func(t *testing.T) {
			// Set up Vopts struct with zero shm size
			opts := Vopts{
				Key:  1,
				Size: 0,
			}
			// Create new shared memory segment with specified options
			sg, err := newWithReturnId(opts)
			require.Equal(t, ErrNegativeOrZeroSize, err)
			// Attempt to delete the shared memory segment and ensure that the expected error is returned
			err = sg.deleteWithId()
			require.Equal(t, ErrShmEmptyPoint, err)
		})

		// Test negative shm flag
		t.Run("negative shm flag", func(t *testing.T) {
			// Non-Public configuration parameters will be overridden by the program during execution
			opts := Vopts{
				Key:  1,
				Size: 1024,
				flag: -1, // be overridden by StatusIpcCreate | StatusIpcExclusive
			}
			// Create a new shared memory segment with the specified options
			sg, err := newWithReturnId(opts)
			require.NoError(t, err)

			// Delete the shared memory segment
			err = sg.deleteWithId()
			require.NoError(t, err)
		})

		// Test negative shm parameter
		t.Run("negative shm parameter", func(t *testing.T) {
			// Non-Public configuration parameters will be overridden by the program during execution
			opts := Vopts{
				Key:       1,
				Size:      1024,
				parameter: 99999, // be overridden by 0600
			}
			// Create a new shared memory segment with the specified options
			sg, err := newWithReturnId(opts)
			require.NoError(t, err)

			// Delete the shared memory segment
			err = sg.deleteWithId()
			require.NoError(t, err)
		})
	})
}

/*
Test_Check_Shm_Extension_Function checks the extension functionalities of a data segment implementation,
including testing of many functions: NewShm,InfoShm,WriteOffset,ReadOffset,ReadSize and DeleteShm.
*/
func Test_Check_Shm_Extension_Function(t *testing.T) {
	// Creates shared memory segment, verifying returned information and deleting segment
	t.Run("Detailed inspection of every aspect in InfoShm function", func(t *testing.T) {
		// Create a new shared memory segment with Key=1 and Size=1024
		opts := Vopts{
			Key:  1,
			Size: 1024,
		}
		err := NewShm(opts)
		require.NoError(t, err)

		// Get information about the shared memory segment with Key=1
		var info Vinfo
		info, err = InfoShm(1) // <<<<< <<<<< <<<<< <<<<< <<<<< main test sample
		require.NoError(t, err)

		// Verify the information returned by InfoShm()
		require.Equal(t, uint16(1), info.Major)
		require.Equal(t, uint16(2), info.Minor)
		require.Equal(t, uint16(3), info.Patch)
		require.Equal(t, int64(1), info.Key)
		require.NotEqual(t, int64(0), info.Id)
		require.Equal(t, int64(1024), info.Size)

		// StatusIpcCreate 512 + StatusIpcExclusive 1024 = 1536
		require.Equal(t, int32(1536), info.Flag)

		require.Equal(t, int8(0), info.Parameter[0])
		require.Equal(t, int8(6), info.Parameter[1])
		require.Equal(t, int8(0), info.Parameter[2])
		require.Equal(t, int8(0), info.Parameter[3])
		require.Equal(t, int64(50), info.Offset)
		require.Equal(t, int32(0), info.Type)

		// Delete the shared memory segment with Key=1
		err = DeleteShm(1) // <<<<< <<<<< <<<<< assistant test sample
		require.NoError(t, err)
	})

	// Test WriteOffset function by creating shared memory segment, writing offset value and verifying information
	t.Run("Detailed inspection of every aspect in WriteOffset function", func(t *testing.T) {
		// Create a new shared memory segment with Key=1 and Size=1024
		opts := Vopts{
			Key:  1,
			Size: 1024,
		}
		err := NewShm(opts) // <<<<< <<<<< <<<<< assistant test sample
		require.NoError(t, err)

		// Call the WriteOffset function with a specific offset value (9223372036854775807)
		err = WriteOffset(1, 9223372036854775807) // <<<<< <<<<< <<<<< <<<<< <<<<< main test sample

		// Get information about the shared memory segment with Key=1
		var info Vinfo
		info, err = InfoShm(1) // <<<<< <<<<< <<<<< assistant test sample
		require.NoError(t, err)

		// Verify the information returned by InfoShm()
		require.Equal(t, int8(0), info.Parameter[3])
		info, _ = InfoShm(1)
		require.Equal(t, int64(9223372036854775807), info.Offset)

		// Delete the shared memory segment with Key=1
		err = DeleteShm(1) // <<<<< <<<<< <<<<< assistant test sample
		require.NoError(t, err)
	})

	// Test ReadOffset function by creating shared memory segment, writing offset value, reading and verifying offset value
	t.Run("Detailed inspection of every aspect in ReadOffset function", func(t *testing.T) {
		// Create a new shared memory segment with Key=1 and Size=1024
		opts := Vopts{
			Key:  1,
			Size: 1024,
		}
		err := NewShm(opts) // <<<<< <<<<< <<<<< assistant test sample
		require.NoError(t, err)

		// Call the WriteOffset function with a specific offset value (9223372036854775807)
		err = WriteOffset(1, 9223372036854775807) // <<<<< <<<<< <<<<< assistant test sample

		// Read the offset value from the shared memory segment
		var offset int64
		offset, err = ReadOffset(1) // <<<<< <<<<< <<<<< <<<<< <<<<< main test sample
		require.NoError(t, err)
		require.Equal(t, int64(9223372036854775807), offset)

		// Delete the shared memory segment with Key=1
		err = DeleteShm(1) // <<<<< <<<<< <<<<< assistant test sample
		require.NoError(t, err)
	})

	// Test WriteInt32s function by creating shared memory segment, writing and reading int32 values
	t.Run("Detailed inspection of every aspect in WriteInt32s function", func(t *testing.T) {
		// Create shared memory segment with key=1 and size=1024 bytes
		opts := Vopts{
			Key:  1,
			Size: 1024,
		}
		err := NewShm(opts) // <<<<< <<<<< <<<<< assistant test sample
		require.NoError(t, err)

		// Defer deleting the shared memory segment with key=1 until the end of the function
		defer func() {
			// Delete the shared memory segment with Key=1
			err = DeleteShm(1) // <<<<< <<<<< <<<<< assistant test sample
			require.NoError(t, err)
		}()

		// Write five int32 values to the shared memory segment starting at offset (defaultMinShmSize + 0)
		err = WriteInt32s(1, 1, 2, 3, 4, 5)
		require.NoError(t, err)

		// Read the current offset value from the shared memory segment and ensure it is (defaultMinShmSize+20)
		var offset int64
		offset, err = ReadOffset(1)
		require.NoError(t, err)
		require.Equal(t, int64(DefualtMinShmSize+20), offset)

		// Read five int32 values from the shared memory segment starting at offset 0 and ensure they are the values we wrote
		values := make([]int32, 5)
		err = ReadInt32s(1, 0, values)
		require.NoError(t, err)
		require.Equal(t, []int32{1, 2, 3, 4, 5}, values)

		// Write another five int32 values to the shared memory segment starting at offset (defaultMinShmSize + 20)
		err = WriteInt32s(1, 6, 7, 8, 9, 10)
		require.NoError(t, err)

		// Read the current offset value from the shared memory segment and ensure it is (defaultMinShmSize+40)
		offset, err = ReadOffset(1)
		require.NoError(t, err)
		require.Equal(t, int64(DefualtMinShmSize+40), offset)

		// Read five int32 values from the shared memory segment starting at offset 20 and ensure they are the values we wrote
		values = make([]int32, 5)
		err = ReadInt32s(1, 20, values)
		require.NoError(t, err)
		require.Equal(t, []int32{6, 7, 8, 9, 10}, values)
	})
}
