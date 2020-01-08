package xgin

import (
	"github.com/jinares/gopkg/xtools"
)

type (
	gray struct {
		index     int ////表示上一次选择的服务器
		cw        int //表示当前调度的权值
		gcd       int //当前所有权重的最大公约数 比如 2，4，8 的最大公约数为：2
		maxWeight int //最大权重
		len       int //Upstream列表的长度
		weight    int
	}
)

func newGray(weight int) *gray {
	g := &gray{}
	g.SetW(weight)
	return g
}

//GetGrayWeigth GetGrayWeigth
func (h *gray) SetW(weigth int) {
	h.weight = weigth
	if weigth < 0 || weigth > 100 {
		return
	}
	h.weight = weigth
	h.gcd = xtools.GCD(100-h.weight, h.weight)
	h.maxWeight = 100
	h.len = 2
}

//IsProxy ==0 执行代理
func (h *gray) IsProxy() int {
	if h.weight == 0 {
		//不执行代理
		return 1
	}
	for {
		h.index = (h.index + 1) % h.len
		if h.index == 0 {
			h.cw = h.cw - h.gcd
			if h.cw <= 0 {
				h.cw = h.maxWeight
				if h.cw == 0 {
					return 0
				}
			}
		}
		if h.index >= h.len {
			h.index = -1
			return 0
		}
		w := h.weight
		if h.index != 1 {
			w = 100 - h.weight
			if w >= h.cw {
				return 1
			}
		} else {
			if w >= h.cw {
				return 0
			}
		}
	}

}
