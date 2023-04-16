package speedyArray

import (
	"github.com/panhongrainbow/filebasez/shm"
	"github.com/stretchr/testify/require"
	"testing"
)

/*
Test_Check_SpeedyArrayInt32 tests the functions of SpdArrayInt32.
It appends some rows to the array, and check first four rows in memory in the shared memory.
*/
func Test_Check_SpeedyArrayInt32(t *testing.T) {
	// create a new instance of SpdArrayInt32 with the given options

	// The testShmKey is the shared memory key for testing
	var testShmKey int64 = 20

	// The opts is options for creating a new instance of SpdArrayInt32
	opts := Opts{
		// ShmKey represents the shared memory key
		ShmKey: testShmKey,
		// Width and Length represent the Width and Length of the array respectively
		Width:  5,
		Length: 30,
	}

	// Create a new instance of SpdArrayInt32 with the given options
	array, err := NewSpeedyArrayInt32(opts)
	require.NoError(t, err, "create new speedy array failed")

	// The firstElement is the first element of the array
	var firstElement int32 = 20

	// AnotherFirstElement is the another first element of the array
	var AnotherFirstElement int32 = 50

	// UniqueElement is the unique element of the array
	var UniqueElement int32 = 70

	// Delete the shared memory segment with the given key
	defer func() {
		err := DeleteSpeedyArrayInt32(testShmKey)
		require.NoError(t, err)
	}()

	// Subtest: Append the many rows of the Int32s array
	t.Run("Test AppendArrayInt32", func(t *testing.T) {
		// Append the first row of the array
		err = array.AppendArrayInt32(firstElement, 11, 12, 13, 14, 15, 16, 17, 18, 19)
		require.Equal(t, ErrTruncateData, err, "append array error is not equal to the expected value")

		// Check the shm offset after appending the first row of the array
		var shmOffset int64
		shmOffset, err = shm.ReadOffset(testShmKey)
		require.NoError(t, err, "read offset failed")
		require.Equal(t, int64(shm.DefualtMinShmSize+20), shmOffset, "shm offset is not equal to the expected value")

		// Append the second row of the array
		err = array.AppendArrayInt32(firstElement, 21, 22, 23, 24, 25, 26, 27, 28, 29)
		require.Equal(t, ErrTruncateData, err, "append array error is not equal to the expected value")

		// Check the shm offset after appending the second row of the array
		shmOffset, err = shm.ReadOffset(testShmKey)
		require.NoError(t, err, "read offset failed")
		require.Equal(t, int64(shm.DefualtMinShmSize+20*2), shmOffset, "shm offset is not equal to the expected value")

		// append the third row of the array
		err = array.AppendArrayInt32(firstElement, 31, 32, 33, 34, 35, 36, 37, 38, 39)
		require.Equal(t, ErrTruncateData, err, "append array error is not equal to the expected value")

		// Check the shm offset after appending the third row of the array
		shmOffset, err = shm.ReadOffset(testShmKey)
		require.NoError(t, err, "read offset failed")
		require.Equal(t, int64(shm.DefualtMinShmSize+20*3), shmOffset, "shm offset is not equal to the expected value")

		// append the fourth row of the array
		err = array.AppendArrayInt32(AnotherFirstElement, 51, 52, 53, 54, 55, 56, 57, 58, 59)
		require.Equal(t, ErrTruncateData, err, "append array error is not equal to the expected value")

		// Check the shm offset after appending the fourth row of the array
		shmOffset, err = shm.ReadOffset(testShmKey)
		require.NoError(t, err, "read offset failed")
		require.Equal(t, int64(shm.DefualtMinShmSize+20*4), shmOffset, "shm offset is not equal to the expected value")
	})
	// Subtest: Read the first four rows in the Int32s array
	t.Run("Test ReadRowInInt32ByShift", func(t *testing.T) {
		// check the first row of the array
		var raw []int32
		raw, err = array.ReadRowInInt32ByShift(0) // 0 is the shift of the first row
		require.Equal(t, []int32{firstElement, 11, 12, 13, 14}, raw, "raw is not equal to the expected value")
		require.NoError(t, err)

		// check the second row of the array
		raw, err = array.ReadRowInInt32ByShift(20) // 20 is the shift of the second row
		require.Equal(t, []int32{firstElement, 21, 22, 23, 24}, raw, "raw is not equal to the expected value")
		require.NoError(t, err)

		// check the third row of the array
		raw, err = array.ReadRowInInt32ByShift(40) // 40 is the shift of the third row
		require.Equal(t, []int32{firstElement, 31, 32, 33, 34}, raw, "raw is not equal to the expected value")
		require.NoError(t, err)

		// check the fourth row of the array
		raw, err = array.ReadRowInInt32ByShift(60) // 60 is the shift of the fourth row
		require.Equal(t, []int32{AnotherFirstElement, 51, 52, 53, 54}, raw, "raw is not equal to the expected value")
		require.NoError(t, err)
	})
	/*
		Subtest: Read the first three rows in the Int32s array by using the first element.
		Then check the last row of the array by using the first element.
	*/
	t.Run("Test ReadRowInInt32ByFirstElement", func(t *testing.T) {
		// >>>>> Check the first three rows of the array by using the first element

		// Define the twoDimensionalArray variable
		var twoDimensionalArray [][]int32
		var err error

		// Read the first three rows of the array by using the first element, which is the variable firstElement
		twoDimensionalArray, err = array.ReadRowInInt32ByFirstElement(firstElement)
		require.Equal(t, [][]int32{{
			firstElement, 11, 12, 13, 14},
			{firstElement, 21, 22, 23, 24},
			{firstElement, 31, 32, 33, 34}},
			twoDimensionalArray, "twoDimensionalArray is not equal to the expected value")

		// Check if the error is equal to nil
		require.NoError(t, err)

		// >>>>> Check the last row of the array by using the first element, which is the variable AnotherFirstElement
		twoDimensionalArray, err = array.ReadRowInInt32ByFirstElement(AnotherFirstElement)
		require.Equal(t, [][]int32{{
			AnotherFirstElement, 51, 52, 53, 54}},
			twoDimensionalArray, "twoDimensionalArray is not equal to the expected value")

		// Check if the error is equal to nil
		require.NoError(t, err)
	})
	/*
		Unique is a method to maintain the uniqueness of the first element in the array.
		When use the Unique method many times to append the same first element, the array will only keep the unique first element.
	*/
	t.Run("Test Unique", func(t *testing.T) {
		// Append the unique row of the int32s array
		err = array.Unique(UniqueElement, 61, 62, 63, 64)
		require.NoError(t, err)

		// Check the shm offset after appending the second row of the array
		var shmOffset int64
		shmOffset, err = shm.ReadOffset(testShmKey)
		require.NoError(t, err, "read offset failed")
		require.Equal(t, int64(shm.DefualtMinShmSize+20*5), shmOffset, "shm offset is not equal to the expected value")

		// Check the unique row of the int32s array
		var raw []int32
		raw, err = array.ReadRowInInt32ByShift(80)
		require.Equal(t, []int32{UniqueElement, 61, 62, 63, 64}, raw, "raw is not equal to the expected value")
		require.NoError(t, err)

		// Overwrite the unique row of the int32s array
		err = array.Unique(UniqueElement, 71, 72, 73, 74)
		require.NoError(t, err)

		// Check the shm offset after overwriting the unique row of the array
		shmOffset, err = shm.ReadOffset(testShmKey)
		require.NoError(t, err, "read offset failed")
		require.Equal(t, int64(shm.DefualtMinShmSize+20*5), shmOffset, "shm offset is not equal to the expected value")

		// Check the unique row of the int32s array
		var raw2 []int32
		raw2, err = array.ReadRowInInt32ByShift(80)
		require.Equal(t, []int32{UniqueElement, 71, 72, 73, 74}, raw2, "raw is not equal to the expected value")
		require.NoError(t, err)
	})
}

func Test_Check_SpeedyArrayInt32_Unique(t *testing.T) {
	// create a new instance of SpdArrayInt32 with the given options

	// The testShmKey is the shared memory key for testing
	var testShmKey int64 = 21

	// The opts is options for creating a new instance of SpdArrayInt32
	opts := Opts{
		// ShmKey represents the shared memory key
		ShmKey: testShmKey,
		// Width and Length represent the Width and Length of the array respectively
		Width:  9,
		Length: 30,
	}

	// Create a new instance of SpdArrayInt32 with the given options
	array, err := NewSpeedyArrayInt32(opts)
	require.NoError(t, err, "create new speedy array failed")

	// UniqueElement is the unique element of the array
	var uniqueElement int32 = 70
	var secondUniqueElement int32 = 71

	// Delete the shared memory segment with the given key
	defer func() {
		err := DeleteSpeedyArrayInt32(testShmKey)
		require.NoError(t, err)
	}()

	// Append the unique row of the int32s array
	err = array.Unique(uniqueElement, 11, 12, 13, 14, 15, 16, 17, 18)
	require.NoError(t, err)

	// Overwrite the unique row of the int32s array
	err = array.Unique(uniqueElement, 12, 12, 14, 14, 16, 16, 18, 18)
	require.NoError(t, err)

	// Check the shm offset after appending the second row of the array
	var shmOffset int64
	shmOffset, err = shm.ReadOffset(testShmKey)
	require.NoError(t, err, "read offset failed")
	require.Equal(t, int64(shm.DefualtMinShmSize+36), shmOffset, "shm offset is not equal to the expected value")

	// Check the unique row of the int32s array
	var raw []int32
	raw, err = array.ReadRowInInt32ByShift(0)
	require.NotEqual(t, []int32{uniqueElement, 11, 12, 13, 14, 15, 16, 17, 18}, raw, "raw is not equal to the expected value")
	require.Equal(t, []int32{uniqueElement, 12, 12, 14, 14, 16, 16, 18, 18}, raw, "raw is not equal to the expected value")
	require.NoError(t, err)

	// Check the shm offset after appending the second row of the array
	shmOffset, err = shm.ReadOffset(testShmKey)
	require.NoError(t, err, "read offset failed")
	require.Equal(t, int64(shm.DefualtMinShmSize+36), shmOffset, "shm offset is not equal to the expected value")

	// Append the unique row of the int32s array
	err = array.Unique(secondUniqueElement, 21, 22, 23, 24, 25, 26, 27, 28)
	require.NoError(t, err)

	// Check the shm offset after appending the second row of the array
	shmOffset, err = shm.ReadOffset(testShmKey)
	require.NoError(t, err, "read offset failed")
	require.Equal(t, int64(shm.DefualtMinShmSize+36*2), shmOffset, "shm offset is not equal to the expected value")

	// Check the unique row of the int32s array
	raw, err = array.ReadRowInInt32ByShift(36)
	require.Equal(t, []int32{secondUniqueElement, 21, 22, 23, 24, 25, 26, 27, 28}, raw, "raw is not equal to the expected value")
	require.NoError(t, err)
}
