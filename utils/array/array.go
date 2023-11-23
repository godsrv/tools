package array

//@author: lipper
//@function: IsContain
//@description: 检查某个值是否在数组中
//@param: items []T, item T
//@return: bool
func IsContain[T comparable](items []T, item T) bool {
	for _, v := range items {
		if v == item {
			return true
		}
	}
	return false
}

//@author: lipper
//@function: RemoveRep
//@description: 数组去重
//@param: items []T
//@return: []T
func RemoveRep[T comparable](items []T) []T {
	result := []T{}
	tempMap := map[T]byte{} // 存放不重复主键
	for _, v := range items {
		l := len(tempMap)
		tempMap[v] = 0
		if len(tempMap) != l {
			result = append(result, v)
		}
	}
	return result
}

//@author: lipper
//@function: Reverse
//@description: 切片反转
//@param: s []T
//@return: []T
func ReverseSlice[T comparable](s []T) []T {
	var r []T
	for i := len(s) - 1; i >= 0; i-- {
		r = append(r, s[i])
	}
	return r
}
