package xtools

import (
	"fmt"
	"testing"
	"time"
)

var (
	topic2 = NewTopics2(200, 30, func(keys ...string) (m map[string]interface{}, err error) {
		data := map[string]interface{}{}
		for _, item := range keys {
			data[item] = Guid()
		}
		return data, nil
	})
)

func TestLazyTopics2_Get(t *testing.T) {
	for u := 0; u < 100; u++ {
		go func() {
			for i := 0; i < 300; i++ {
				topic2.Get(RandID(10), 1*time.Second)
			}
		}()
	}
	for i := 0; i < 300; i++ {
		topic2.Get(RandID(10), 1*time.Second)
	}
	val, err := topic2.Get(RandID(10), 1*time.Second)
	fmt.Println(val, err)

}
func BenchmarkLazyTopics2_Get(b *testing.B) {

	b.StopTimer()
	b.SetBytes(100)
	b.SetParallelism(5000)

	b.StartTimer()
	for idx := 0; idx < b.N; idx++ {
		topic2.Get(RandID(10), 1*time.Second)
	}
	b.ReportAllocs()

}

/*
var pars struct {
			Telecom string `json:"telecom"`
		}
		err := utils.ToJson(param.PayTypeExt, &pars)
		if err != nil {
			return xtype.Result(
				errcode.ErrParams, nil, "paytypeext参数错误",
			)
		}
		if utils.Has(pars.Telecom, "0", "1", "2", "3") == false {
			param.ext3.MEvent = "0"
		} else {

			param.ext3.MEvent = pars.Telecom
		}

				"telecom":          jsonext3.MEvent,
*/
