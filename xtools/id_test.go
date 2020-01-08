package xtools

import (
	"fmt"
	"testing"
	"time"
)

func TestObjectIDCounter(t *testing.T) {
	id := ObjectIDCounter(20)
	fmt.Println("id:", id, len("18062212561352929260616643068753"))
	for index, _ := range [20]int{} {
		fmt.Println(time.Now().Format("060102150405")+ObjectIDCounter(20), index)
	}
}
