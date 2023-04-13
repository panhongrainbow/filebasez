package speedyArray

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_Check_SpeedyArrayInt32(t *testing.T) {
	//
	var testShmKey int64 = 2

	//
	opts := Opts{
		//
		ShmKey: testShmKey,

		//
		width:  5,
		length: 30,
	}

	//
	array, err := NewSpeedyArrayInt32(opts)
	require.NoError(t, err)

	//
	defer func() {
		array = SpeedyArrayInt32{}
		err = DeleteSpeedyArrayInt32(testShmKey)
		require.NoError(t, err)
	}()

	//
	err = array.AppendArrayInt32(3, 2, 3, 4, 5, 6, 7, 8, 9)
	require.NoError(t, err)

	//
	result, err := array.ReadArrayInt32ByShift(0)
	fmt.Println(result, err)

	//
	err = array.AppendArrayInt32(3, 3, 5, 7, 9, 11, 13, 15, 17)
	require.NoError(t, err)

	//
	result, err = array.ReadArrayInt32ByShift(20)
	fmt.Println(result, err)

	//
	results, err := array.ReadArrayInt32ByElement(3)
	fmt.Println(results, err)

	/*err := shm.WriteInt32s(2, 1, 2, 3, 4, 5)
	require.NoError(t, err)

	values := make([]int32, 5)
	err = shm.ReadInt32s(2, 0, values)
	fmt.Println(values)
	status, _ := shm.ReadOffset(2)
	fmt.Println(status)

	err = shm.WriteInt32s(2, 6, 7, 8, 9, 10)
	values = make([]int32, 5)
	err = shm.ReadInt32s(2, 20, values)
	fmt.Println(values)
	status, _ = shm.ReadOffset(2)
	fmt.Println(status)*/

	/*err = array.AppendArrayInt32(1, 2, 3, 4, 5)
	fmt.Println(err)*/

}
