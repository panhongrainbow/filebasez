package speedyArray

import (
	"fmt"
	"github.com/panhongrainbow/filebasez/shm"
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

func (array SpeedyArrayInt32) AppendArrayInt32(elements ...int32) (err error) {

	if array.shiftMap[elements[0]] == nil {
		//
		array.shiftMap[elements[0]] = make([]int64, 5)
	}

	err = shm.WriteInt32s(array.opts.ShmKey, elements...)
	if err != nil {
		return
	}

	var shmOffset int64
	shmOffset, err = shm.ReadOffset(array.opts.ShmKey)
	if err != nil {
		return
	}

	offset := shmOffset - int64(array.opts.width*4)

	array.shiftMap[elements[0]] = append(array.shiftMap[elements[0]], offset)
	fmt.Println(shmOffset, offset)

	return
}
