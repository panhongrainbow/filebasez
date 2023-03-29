package shm

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_Check_Shm(t *testing.T) {

	t.Run("test for basic functions", func(t *testing.T) {
		//
		opts := Vopts{
			Key:  1,
			Size: 1024,
		}
		sg, err := newWithReturnId(opts)
		require.NoError(t, err)
		require.Equal(t, int64(1), sg.key)
		require.Greater(t, sg.id, int64(0))

		//
		length, err := sg.writeWithId([]byte("abcde"))
		require.NoError(t, err)
		require.Equal(t, int64(5), length)

		//
		sg.offset = 0
		var data = make([]byte, 5)
		length, err = sg.readWithId(data)
		require.NoError(t, err)
		require.Equal(t, int64(5), length)
		fmt.Println(data)

		//
		err = sg.deleteWithId(sg.id)
		require.NoError(t, err)
	})

	t.Run("test for extension functions", func(t *testing.T) {
		opts := Vopts{
			Key:  1,
			Size: 1024,
		}
		err := NewShm(opts)
		require.NoError(t, err)

		//
		/*var sg = new(Vsegment)
		sg.key = 1
		sg.id = VsegmentMap[1]
		sg.size = 1024
		var data = make([]byte, 20)
		_, err = sg.readWithId(data)
		require.NoError(t, err)
		fmt.Println(data)*/
		ShmInfo(1)

	})
}
