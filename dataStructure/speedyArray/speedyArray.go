package speedyArray

import (
	"github.com/panhongrainbow/filebasez/shm"
)

const (
	defaultOverlapsFirstElement = 5
)

const (
	ErrOverwriteBeyondSize = Error("overwrite failed because beyond shm size")
	ErrNotAlignWithMemory  = Error("not align with the memory size boundary")
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

	// width and length represent the width and length of the array respectively
	width  uint64
	length uint64
}

// NewSpeedyArrayInt32 creates a new instance of SpdArrayInt32 with the given options.
func NewSpeedyArrayInt32(opts Opts) (array SpdArrayInt32, err error) {
	// estimateSize calculates the estimated size of the shared memory based on the width and length of the array
	estimateSize := shm.DefualtMinShmSize + (opts.width*opts.length)*4

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
		shiftMap: make(map[int32][]int64, opts.length),
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

	// create a new int32 slice with the width specified in the options
	var newElements = make([]int32, array.opts.width, array.opts.width)
	// copy the given elements to the newElements slice
	copy(newElements, elements)

	// write the newElements to the shared memory segment with the given key
	err = shm.WriteInt32s(array.opts.ShmKey, newElements...)
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
	offset := shmOffset - int64(array.opts.width*4)

	// append the new offset value to the shiftMap for the first element
	array.shiftMap[elements[0]] = append(array.shiftMap[elements[0]], offset)

	// return any error that occurred
	return
}

// OverwriteArrayInt32ByShift overwrites a shared memory array with given elements at a specified shift, checking for alignment and size errors.
func (array SpdArrayInt32) OverwriteArrayInt32ByShift(shmShift int64, elements ...int32) (err error) {
	// Read the shared memory size
	var shmSize int64
	shmSize, err = shm.ReadSize(array.opts.ShmKey)
	if shmShift > shmSize {
		err = ErrOverwriteBeyondSize
		return
	}

	// Check if the shift is aligned with memory
	shouldAlign := (shmShift - shm.DefualtMinShmSize) % int64(4*array.opts.width)
	if shouldAlign != 0 {
		err = ErrNotAlignWithMemory
		return
	}

	// Overwrite the shared memory with the given elements at the given shift
	err = shm.OverwriteInt32sByShift(array.opts.ShmKey, shmShift, elements...)
	if err != nil {
		return
	}

	// Return any errors
	return
}

func (array SpdArrayInt32) ReadArrayInt32ByShift(shmShift int64) (elements []int32, err error) {
	elements = make([]int32, array.opts.width)
	err = shm.ReadInt32s(array.opts.ShmKey, shmShift, elements)
	return
}

func (array SpdArrayInt32) ReadArrayInt32ByElement(firstElement int32) (elements [][]int32, err error) {

	//
	if array.shiftMap[firstElement] == nil {
		//
		array.shiftMap[firstElement] = make([]int64, 0, 5)
	}

	//
	if len(array.shiftMap[firstElement]) <= 0 {
		return
	}

	//
	elements = make([][]int32, len(array.shiftMap[firstElement]))

	//
	for i := 0; i < len(array.shiftMap[firstElement]); i++ {
		//
		var raw []int32
		raw, err = array.ReadArrayInt32ByShift(array.shiftMap[firstElement][i])
		elements = append(elements, raw)

	}

	//
	return
}
