package xtools

import (
	"errors"
	"fmt"
	"testing"
)

type (
	TestPerson struct {
		Name string `json:"name"`
		Data TestES `json:"data"`
	}
	TestES map[string]string
)

func (h *TestPerson) Check() {
	fmt.Println("check:", h.Name)
	val := h.Data["test"]
	fmt.Println("data init", val)
	if h.Data == nil {
		h.Data = TestES{}
	}
	h.Data["test"] = "1234567890"
	val = h.Data["test"]
	fmt.Println("data init", val)
}

func TestGetMapVal(t *testing.T) {
	data := map[string]string{
		"test1": JSONToStr(TestPerson{Name: "ares"}),
	}
	var item TestPerson
	item.Check()
	err := GetMapVal(func() error {
		val, isok := data["test1"]
		if isok == false {
			return errors.New("not found")
		}
		return ToJSON(val, &item)
	})
	fmt.Println(err, JSONToStr(item))
}
