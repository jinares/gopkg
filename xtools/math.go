package xtools

//Decomposition 分解质因数
func Decomposition(data int) []int {
	var i int = 0
	n := data
	ret := []int{}
	for i = 2; i <= n; i++ {
		for n != i {
			if n%i == 0 {
				ret = append(ret, i)
				n = n / i
			} else {
				break
			}
		}
	}
	ret = append(ret, n)
	return ret

}

//GCD 最大公约数
func GCD(data ...int) int {
	tmp1 := []int{}
	tmp2 := []map[int]int{}
	for index, val := range data {
		tmp := Decomposition(val)
		if index == 0 {
			tmp1 = tmp
		}
		tmpItem := make(map[int]int)
		for index, val := range tmp {
			tmpItem[index] = val
		}
		tmp2 = append(tmp2, tmpItem)
	}
	tmpRet := 1
	for _, val := range tmp1 {
		tmpIsok := true
		for _, item := range tmp2 {
			tmpIsok2 := false
			for key, subVal := range item {
				if subVal == val {
					tmpIsok2 = true
					delete(item, key)
					break
				}
			}
			if tmpIsok2 == false {
				tmpIsok = false
			}

		}
		if tmpIsok {
			tmpRet = tmpRet * val
		}
	}

	return tmpRet
}
