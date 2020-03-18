package xlog

import (
	"github.com/jinares/gopkg/xtools"
	"strings"
)

type Msg struct {
	data []string
}

func NewMsg(str string) *Msg {
	return &Msg{data: []string{str}}
}
func (h *Msg) Add(mstr interface{}) *Msg {
	str := ""
	switch vv := mstr.(type) {
	case Msg:
		str = vv.String()
	case string:
		str = vv
	default:
		str = xtools.JSONToStr(vv)

	}
	h.data = append(h.data, str)
	return h
}
func (h *Msg) BR() *Msg {
	return h.Add("\n")

}
func (h *Msg) String() string {
	return strings.Join(h.data, " ")
}

func (h *Msg) MarshalJSON() ([]byte, error) {
	return []byte("\"" + h.String() + "\""), nil
}
func (h *Msg) UnmarshalJSON(data []byte) error {
	str := string(data)
	h.data = strings.Split(str, " ")
	return nil
}
