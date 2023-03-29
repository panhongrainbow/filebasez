package shm

import (
	"fmt"
	"testing"
)

func Test_Check_Shm(t *testing.T) {

	/*opts := Vopts{
		Key:  1,
		Size: 1024,
	}

	segment, err := New(opts)
	fmt.Println(segment)
	fmt.Println(err)

	length, err := segment.WriteWithId([]byte("12345"))
	fmt.Println(length)
	fmt.Println(err)*/

	segment := Vsegment{}

	segment.offset = 0
	segment.id = 18546746
	segment.size = 1024

	var data = make([]byte, 5)
	length1, err := segment.ReadWithId(data)
	fmt.Println(length1)
	fmt.Println(err)
	fmt.Println(data)
}
