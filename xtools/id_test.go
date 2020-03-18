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

	var ext map[string]string
	err := ToJSON(`{}`, &ext)
	if err != nil {
		fmt.Println(err)
		return
	}
	ext["mAid"] = ""
	fmt.Println(ext, JSONToStr(ext))
}
func BenchmarkHashID(b *testing.B) {
	b.StopTimer()
	gid := Guid()
	//fmt.Println(HashID(gid, 10))
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		HashID(gid, 100)

	}
}
