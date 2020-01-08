package xgin

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"

	"github.com/jinares/gopkg/xlog"
	"github.com/jinares/gopkg/xtools"
)

type (
	Rule struct {
		IsRuning     int               `json:"IsRuning,omitempty" yaml:"IsRuning"` //0:有效 1:无效
		Location     string            `json:"Location",omitempty yaml:"Location"`
		Rewrite      string            `json:"Rewrite,omitempty"yaml:"Rewrite"`
		ProxyWeigth  int               `json:"ProxyWeigth,omitempty" yaml:"ProxyWeigth"`
		ProxyPass    string            `json:"ProxyPass,omitempty" yaml:"ProxyPass"`
		ProxyRewrite string            `json:"ProxyRewrite,omitempty"yaml:"ProxyRewrite"`
		ProxyHost    string            `json:"ProxyHost,omitempty" yaml:"ProxyHost"`
		Params       map[string]string `json:"Params,omitempty" yaml:"Params"` // 和 Location 一起进行匹配g
		gray         *gray
	}
	Locations []*Rule
)

func (h *Rule) GetLocationVal() string {
	if strings.HasPrefix(h.Location, "$") {
		return h.Location[1:]
	}
	return h.Location
}
func (h *Rule) IsProxy() bool {
	if h.ProxyPass == "" {
		return false
	}
	if h.ProxyWeigth <= 0 {
		return false
	}
	if h.gray == nil {
		h.gray = newGray(h.ProxyWeigth)
	}

	//IsGray ==0  执行代理
	if h.gray.IsProxy() == 0 {
		return true
	}
	return false

}
func (h *Rule) matching(url string, param map[string]string) bool {
	ispre := strings.HasPrefix(h.Location, "$")
	if (ispre && strings.HasPrefix(url, h.Location[1:])) || (ispre == false && url == h.Location) {
		if len(h.Params) == 0 {
			return true
		}
		tmpRet := true
		for key, val := range h.Params {
			tmpVal, isok := param[key]
			tmpPrefix := strings.HasPrefix(val, "$")
			if isok == false {
				tmpRet = false
			} else if tmpPrefix == true && strings.HasPrefix(tmpVal, val[1:]) == false {
				tmpRet = false
			} else if tmpPrefix == false && val != "" && val != tmpVal {
				tmpRet = false
			}
		}
		return tmpRet
	}
	return false
}
func (h *Locations) Matching(r *http.Request) *Rule {
	params := r.URL.RawQuery
	if strings.ToLower(r.Method) == "post" {
		contenttype := r.Header.Get("Content-Type")
		xlog.CtxLog(r.Context()).Tracef("url:%s content-type:%s", r.URL.Path, contenttype)
		//application/x-www-form-urlencoded
		urlencode := "application/x-www-form-urlencoded"
		if contenttype == urlencode {
			bodyBytes, _ := ioutil.ReadAll(r.Body)
			r.Body.Close() //  must close
			r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
			params = string(bodyBytes)
		}
	}
	xlog.CtxLog(r.Context()).Tracef("url:%s method:%s, params:%s", r.URL.Path, r.Method, params)
	pars := xtools.UrlDecode(params)
	url := r.URL.Path
	sort.Stable(h)
	for _, item := range *h {
		if item.IsRuning != 0 {
			continue
		}
		if item.matching(url, pars) {
			return item
		}
	}
	return nil
}

// Len 长度计算
func (a Locations) Len() int { // 重写 Len() 方法
	return len(a)
}

// Swap 交换
func (a Locations) Swap(i, j int) { // 重写 Swap() 方法
	a[i], a[j] = a[j], a[i]
}

// Less 比较
func (a Locations) Less(i, j int) bool { // 重写 Less() 方法， 从大到小排序
	itemi := a[i]
	itemj := a[j]
	ihasPrefix := strings.HasPrefix(itemi.Location, "$")
	jhasPrefix := strings.HasPrefix(itemj.Location, "$")
	if (ihasPrefix == false && jhasPrefix == false) || (ihasPrefix == true && jhasPrefix == true) {
		if len(itemi.Location) > len(itemj.Location) {
			return true
		} else if len(itemi.Location) == len(itemj.Location) {
			if len(itemi.Params) > len(itemj.Params) {
				return true
			}
		}
		return false
	}
	if ihasPrefix == false && jhasPrefix == true {
		return true
	} else if ihasPrefix == true && jhasPrefix == false {
		return false
	}
	return false

}
