package xtools

import (
	"fmt"
	"testing"
	"time"
)

//
//func TestNewLazy(t *testing.T) {
//
//	tt := time.NewTimer(2 * time.Second)
//	for {
//		select {
//		case <-tt.C:
//			fmt.Println("11111...")
//		default:
//			fmt.Println("2...")
//			time.Sleep(1 * time.Second)
//		}
//	}
//}

var (
	topic = NewTopics(200, 3, func(keys ...string) (m map[string]interface{}, err error) {
		data := map[string]interface{}{}
		for _, item := range keys {
			data[item] = Guid()
		}
		return data, nil
	})
)

func TestLazyTopics_Get(t *testing.T) {
	for u := 0; u < 100; u++ {
		go func() {
			for i := 0; i < 300; i++ {
				topic.Get(RandID(10), 1*time.Second)
			}
		}()
	}
	for i := 0; i < 300; i++ {
		topic.Get(RandID(10), 1*time.Second)
	}
	val, err := topic.Get(RandID(10), 1*time.Second)
	fmt.Println(val, err)
}
func BenchmarkLazyTopics_Get(b *testing.B) {

	b.StopTimer()
	b.SetBytes(100)
	b.SetParallelism(5000)

	b.StartTimer()
	for idx := 0; idx < b.N; idx++ {
		topic.Get(RandID(10), 1*time.Second)
	}
	b.ReportAllocs()

}
