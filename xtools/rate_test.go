package xtools

import (
	"fmt"
	"testing"
	"time"
)

func TestAcRate_Rate(t *testing.T) {
	r := NewAcRate(60)
	r.SetKeyFunc(func() interface{} {
		return time.Now().Format("0601021504")
	})
	for i := 1; i <= 80000; i++ {
		r.Succ()
	}
	for i := 1; i <= 20; i++ {
		r.Fail()
	}
	fmt.Println(r.Rate())
	for i := 1; i <= 2000; i++ {
		r.Fail()
	}
	fmt.Println(r.Rate())
}
func BenchmarkNewAcRate(b *testing.B) {
	b.StopTimer() //调用该函数停止压力测试的时间计数
	r := NewAcRate(60)
	r.SetKeyFunc(func() interface{} {
		return time.Now().Format("0601021504")
	})
	//做一些初始化的工作,例如读取文件数据,数据库连接之类的,
	//这样这些时间不影响我们测试函数本身的性能
	for i := 1; i <= 80000; i++ {
		r.Succ()
	}
	for i := 1; i <= 20; i++ {
		r.Fail()
	}
	//fmt.Println(r.Rate())
	b.StartTimer() //重新开始时间
	for i := 0; i < b.N; i++ {
		r.Rate()

	}
}
func BenchmarkAcRate_Fail(b *testing.B) {
	b.StopTimer() //调用该函数停止压力测试的时间计数
	r := NewAcRate(60)
	r.SetKeyFunc(func() interface{} {
		return time.Now().Format("0601021504")
	})
	//做一些初始化的工作,例如读取文件数据,数据库连接之类的,
	//这样这些时间不影响我们测试函数本身的性能
	for i := 1; i <= 80000; i++ {
		r.Succ()
	}
	for i := 1; i <= 20; i++ {
		r.Fail()
	}
	b.StartTimer() //重新开始时间
	for i := 0; i < b.N; i++ {
		r.Fail()

	}
}
