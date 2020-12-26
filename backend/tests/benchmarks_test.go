package tests

import (
	"fmt"
	"testing"
	"unsafe"
)

func TestEntityIDMemory(t *testing.T) {
	s := "-NDveu-9Q"
	fmt.Println("Size of id:", unsafe.Sizeof(s))
}

func TestWorldSize(t *testing.T) {
	s := "-NDveu-9Q"
	fmt.Println("Size of id:", unsafe.Sizeof(s))
}

//-NDveu-9Q
