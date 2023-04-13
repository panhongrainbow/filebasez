package speedyArray

import (
	"fmt"
	"github.com/panhongrainbow/filebasez/shm"
)

type Error string

func (e Error) Error() string {
	return string(e)
}

const (
	ErrOverwriteBeyondSize = Error("overwrite failed because beyond shm size")
	ErrNotAlignWithMemory  = Error("not align with the memory size boundary")
)

type SpeedyArrayInt32 struct {
	shiftMap map[int32][]int64
	opts     Opts
}

type Opts struct {
	//
	ShmKey int64

	//
	width  uint64
	length uint64
}

func NewSpeedyArrayInt32(opts Opts) (array SpeedyArrayInt32, err error) {
	//
	estimateSize := shm.DefualtMinShmSize + (opts.width*opts.length)*4

	//
	shmOts := shm.Vopts{
		Key:  opts.ShmKey,
		Size: int64(estimateSize),
	}
	err = shm.NewShm(shmOts)
	if err != nil {
		return
	}

	//
	array = SpeedyArrayInt32{
		shiftMap: make(map[int32][]int64, opts.length),
		opts:     opts,
	}

	//
	return
}

func DeleteSpeedyArrayInt32(shmKey int64) (err error) {
	err = shm.DeleteShm(shmKey)
	return
}

func (array SpeedyArrayInt32) AppendArrayInt32(elements ...int32) (err error) {
	//
	if len(elements) <= 0 {
		return
	}

	//
	if array.shiftMap[elements[0]] == nil {
		//
		array.shiftMap[elements[0]] = make([]int64, 0, 5)
	}

	//
	var newElements = make([]int32, array.opts.width, array.opts.width)
	copy(newElements, elements)

	//
	err = shm.WriteInt32s(array.opts.ShmKey, newElements...)
	if err != nil {
		return
	}

	//
	var shmOffset int64
	shmOffset, err = shm.ReadOffset(array.opts.ShmKey)
	if err != nil {
		return
	}

	//
	offset := shmOffset - int64(array.opts.width*4)

	//
	array.shiftMap[elements[0]] = append(array.shiftMap[elements[0]], offset)

	//
	fmt.Println(">>>", shmOffset, offset, array.shiftMap[elements[0]])

	//
	return
}

func (array SpeedyArrayInt32) OverwriteArrayInt32ByShift(shmShift int64, elements ...int32) (err error) {
	//
	var shmSize int64
	shmSize, err = shm.ReadSize(array.opts.ShmKey)
	if shmShift > shmSize {
		err = ErrOverwriteBeyondSize
		return
	}

	//
	shouldAlign := (shmShift - shm.DefualtMinShmSize) % int64(4*array.opts.width)
	if shouldAlign != 0 {
		err = ErrNotAlignWithMemory
		return
	}

	//
	err = shm.OverwriteInt32sByShift(array.opts.ShmKey, shmShift, elements...)
	if err != nil {
		return
	}

	//
	return
}

func (array SpeedyArrayInt32) ReadArrayInt32ByShift(shmShift int64) (elements []int32, err error) {
	elements = make([]int32, array.opts.width)
	err = shm.ReadInt32s(array.opts.ShmKey, shmShift, elements)
	return
}

func (array SpeedyArrayInt32) ReadArrayInt32ByElement(firstElement int32) (elements [][]int32, err error) {

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
