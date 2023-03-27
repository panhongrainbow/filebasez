package shm

import (
	"fmt"
	"testing"
)

func Test_Check_Shm(t *testing.T) {
	opts := Vopts{
		size: 1024,
	}

	segment, err := New(opts)
	fmt.Println(segment)
	fmt.Println(err)

	opts = Vopts{
		key:  1,
		size: 1024,
	}

	segment, err = New(opts)
	fmt.Println(segment)
	fmt.Println(err)
}
