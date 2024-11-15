package qslice

import (
	"fmt"
	"github.com/jinzhu/copier"
	"reflect"
	"sort"
)

// 切片去重
func FilterDuplicates(slice []int) []int {
	seen := make(map[int]bool)
	result := []int{}
	for _, num := range slice {
		if !seen[num] {
			seen[num] = true
			result = append(result, num)
		}
	}
	return result
}

// 合并切片并去重
func MergeAndRemoveDuplicates(slice1, slice2 []int) []int {
	merged := append(slice1, slice2...)
	return FilterDuplicates(merged)
}

// 查找切片value的下标
func SearchValue(li []uint32, value uint32) (int, error) {
	for k, v := range li {
		if v == value {
			return k, nil
		}
	}
	return 0, fmt.Errorf("value not found")
}

// 在切片指定下标后插入新值
func Insert(slice []uint32, index int, value uint32) []uint32 {
	// 创建一个新的切片，将原始切片分为两部分
	first := append([]uint32{}, slice[:index+1]...)
	second := slice[index+1:]

	// 将值插入切片中间
	result := append(first, value)
	result = append(result, second...)
	return result
}

// 在切片头节点前插入值
func InsertAtHead(slice []uint32, value uint32) []uint32 {
	// 创建一个新的切片，长度比原切片多1
	newSlice := make([]uint32, len(slice)+1)
	// 将新值赋给新切片的第一个元素
	newSlice[0] = value
	// 使用切片操作符将原切片的所有元素追加到新切片的后面
	copy(newSlice[1:], slice)
	return newSlice
}

// 判断两个切片是否相等，忽略排序
func IsSlicesEqual(slice1, slice2 []int) bool {
	if len(slice1) != len(slice2) {
		return false
	}
	sortedSlice1 := make([]int, len(slice1))
	sortedSlice2 := make([]int, len(slice2))
	copy(sortedSlice1, slice1)
	copy(sortedSlice2, slice2)

	sort.Ints(sortedSlice1)
	sort.Ints(sortedSlice2)

	return reflect.DeepEqual(sortedSlice1, sortedSlice2)
}

// []uint32转[]int
func Uint32SliceToIntSlice(uintSlice []uint32) []int {
	intSlice := make([]int, len(uintSlice))
	_ = copier.Copy(&intSlice, uintSlice)
	return intSlice
}

// 是否在切片中
func IsInSliceInt(value int, slice []int) bool {
	for _, item := range slice {
		if item == value {
			return true
		}
	}
	return false
}

// 是否在切片中
func IsInSliceInt64(value int64, slice []int64) bool {
	for _, item := range slice {
		if item == value {
			return true
		}
	}
	return false
}

// 是否在切片中
func IsInSliceString(value string, slice []string) bool {
	for _, item := range slice {
		if item == value {
			return true
		}
	}
	return false
}

// 删除切片指定值
func RemoveValue(slice []int, value int) []int {
	result := []int{}
	for _, item := range slice {
		if item != value {
			result = append(result, item)
		}
	}
	return result
}

// 判断两个切片是否有交集
func CheckIntersection(slice1, slice2 []int) bool {
	seen := make(map[int]bool)
	// 将第一个切片的元素添加到map中
	for _, v := range slice1 {
		seen[v] = true
	}
	// 检查第二个切片的元素是否在map中存在
	for _, v := range slice2 {
		if seen[v] {
			return true
		}
	}
	return false
}

// 两个切片过滤交集
func DiffStrSlice(a, b []string) ([]string, []string) {
	diffA := make([]string, 0)
	diffB := make([]string, 0)
	mA := make(map[string]bool)
	mB := make(map[string]bool)

	for _, item := range a {
		mA[item] = true
	}
	for _, item := range b {
		mB[item] = true
	}

	for _, item := range a {
		if !mB[item] {
			diffB = append(diffB, item)
		}
	}
	for _, item := range b {
		if !mA[item] {
			diffA = append(diffA, item)
		}
	}
	return diffA, diffB
}
