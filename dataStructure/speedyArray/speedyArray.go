package speedyArray

import (
	"github.com/panhongrainbow/filebasez/shm"
)

// Define the default capacity for the shiftMap
const (
	defaultOverlapsFirstElement = 5
)

// Define the error messages
const (
	ErrOverwriteBeyondSize = Error("overwrite failed because beyond shm size")
	ErrNotAlignWithMemory  = Error("not align with the memory size boundary")
	ErrTruncateData        = Error("data is truncated")
	ErrWasteMemorySpace    = Error("waste memory space")
)

// Error Defines a new Error type as a string
type Error string

// Error is made for the Error type to return the error message as a string.
func (e Error) Error() string {
	return string(e)
}

// SpdArrayInt32 is a struct that contains a shiftMap and options.
type SpdArrayInt32 struct {
	shiftMap map[int32][]int64
	opts     Opts
}

// Opts contains options for SpeedyArrayInt32.
type Opts struct {
	// ShmKey represents the shared memory key
	ShmKey int64

	// Width and Length represent the Width and Length of the array respectively
	Width  uint64
	Length uint64
}

// NewSpeedyArrayInt32 creates a new instance of SpdArrayInt32 with the given options.
func NewSpeedyArrayInt32(opts Opts) (array SpdArrayInt32, err error) {
	// estimateSize calculates the estimated size of the shared memory based on the Width and Length of the array
	estimateSize := shm.DefualtMinShmSize + (opts.Width*opts.Length)*4

	// shmOts is an instance of Vopts with the given shared memory key and estimated size
	shmOts := shm.Vopts{
		Key:  opts.ShmKey,
		Size: int64(estimateSize),
	}
	// create a new shared memory with the given options
	err = shm.NewShm(shmOts)
	if err != nil {
		return
	}

	// create a new instance of SpdArrayInt32 with the given options and an empty shiftMap
	array = SpdArrayInt32{
		shiftMap: make(map[int32][]int64, opts.Length),
		opts:     opts,
	}

	// return the new instance of SpdArrayInt32 and the error value
	return
}

// DeleteSpeedyArrayInt32 deletes a shared memory segment with the given key
func DeleteSpeedyArrayInt32(shmKey int64) (err error) {
	// delete the shared memory segment with the given key
	err = shm.DeleteShm(shmKey)
	// return any error that occurred
	return
}

// AppendArrayInt32 appends int32 elements to a shared memory segment associated with a SpeedyArrayInt32 instance
func (array SpdArrayInt32) AppendArrayInt32(elements ...int32) (err error) {
	// check if there are any elements to append
	if len(elements) <= 0 {
		return
	}

	// check if the shiftMap for the first element is nil
	if array.shiftMap[elements[0]] == nil {
		// create a new int64 slice with capacity 5 for the first element
		array.shiftMap[elements[0]] = make([]int64, 0, 5)
	}

	// create a new int32 slice with the Width specified in the options
	var newElements = make([]int32, array.opts.Width, array.opts.Width)
	// copy the given elements to the newElements slice
	copy(newElements, elements)

	// write the newElements to the shared memory segment with the given key
	err = shm.AppendInt32s(array.opts.ShmKey, newElements...)
	if err != nil {
		return
	}

	// read the offset value for the given key
	var shmOffset int64
	shmOffset, err = shm.ReadOffset(array.opts.ShmKey)
	if err != nil {
		return
	}

	// calculate the new offset value
	offset := shmOffset - int64(array.opts.Width*4)

	// append the new offset value to the shiftMap for the first element
	array.shiftMap[elements[0]] = append(array.shiftMap[elements[0]], offset-shm.DefualtMinShmSize)

	// Check if the elements are truncated
	if len(elements) > int(array.opts.Width) {
		err = ErrTruncateData
		return
	}

	// return any error that occurred
	return
}

/*
Unique overwrites a shared memory array with given elements
and maintains the uniqueness of the first element in the whole array.
*/
func (array SpdArrayInt32) Unique(elements ...int32) (err error) {
	// check if there are any elements to append
	if len(elements) <= 0 {
		return
	}

	/*
		Check if the shiftMap for the first element is nil.
		If so, create a new int64 slice with capacity 5 for the first element.
		It will alter the shm offset value.
	*/
	if array.shiftMap[elements[0]] == nil {
		// Create a new int64 slice with capacity 5 for the first element
		array.shiftMap[elements[0]] = make([]int64, 0, 5)

		// Append the new elements to the shared memory
		err = array.AppendArrayInt32(elements...)
	}

	/*
		Check if the shiftMap for the first element is nil.
		If not, overwrite the shared memory with the given elements.
		It won't alter the shm offset value.
	*/
	if len(array.shiftMap[elements[0]]) >= 1 {
		/*
			Find the offset values for the first element.
			There should be only one offset value.
		*/
		shmShift := array.shiftMap[elements[0]][0]
		// Overwrite the shared memory with the given elements
		err = shm.OverwriteOrAppendInt32sByShift(array.opts.ShmKey, shm.DefualtMinShmSize+shmShift, false, elements...)
		if err != nil {
			// If first element is not unique, return ErrWasteMemorySpace
			if len(array.shiftMap[elements[0]]) > 1 {
				err = ErrWasteMemorySpace
			}
		}
	}

	// Return any error that occurred
	return
}

// ReadRowInInt32ByShift obtains an int32 array from the shared memory space by an offset.
func (array SpdArrayInt32) ReadRowInInt32ByShift(shmShift int64) (elements []int32, err error) {
	elements = make([]int32, array.opts.Width)
	err = shm.ReadRowInInt32s(array.opts.ShmKey, shmShift, elements)
	return
}

// ReadRowInInt32ByFirstElement obtains an int32 array from the shared memory space by the first element.
func (array SpdArrayInt32) ReadRowInInt32ByFirstElement(firstElement int32) (elements [][]int32, err error) {
	// Check if shiftMap has an entry for firstElement. If not, creates an empty slice for it.
	if array.shiftMap[firstElement] == nil {
		// Create an empty slice if no entry exists
		array.shiftMap[firstElement] = make([]int64, 0, 5)
	}

	// Check if the shiftMap for the firstElement is empty. If so, returns an empty slice
	if len(array.shiftMap[firstElement]) <= 0 {
		return
	}

	// Create a slice of slices to hold the retrieved int32 arrays
	elements = make([][]int32, 0, len(array.shiftMap[firstElement]))

	// Read the int32 arrays for all offsets in shiftMap[firstElement] and appends them to elements
	for i := 0; i < len(array.shiftMap[firstElement]); i++ {
		// Read int32 array at offset and handle any error
		var raw []int32
		// Read the row in int32s array by shift
		raw, err = array.ReadRowInInt32ByShift(array.shiftMap[firstElement][i])
		if err != nil {
			return
		}
		elements = append(elements, raw)

	}

	// Return the elements and any error
	return
}
