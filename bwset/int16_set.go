// Code generated by "bwsetter -type=int16"; DO NOT EDIT; bwsetter: go get -type=int16 -set=Int16 -test%!(EXTRA string=github.com/baza-winner/bwcore/bwsetter)

package bwset

import (
	"encoding/json"
	"sort"
	"strconv"
)

// Int16 - множество значений типа int16 с поддержкой интерфейсов Stringer и MarshalJSON
type Int16 map[int16]struct{}

// Int16From - конструктор Int16
func Int16From(kk ...int16) Int16 {
	result := Int16{}
	for _, k := range kk {
		result[k] = struct{}{}
	}
	return result
}

// Int16FromSlice - конструктор Int16
func Int16FromSlice(kk []int16) Int16 {
	result := Int16{}
	for _, k := range kk {
		result[k] = struct{}{}
	}
	return result
}

// Int16FromSet - конструктор Int16
func Int16FromSet(s Int16) Int16 {
	result := Int16{}
	for k, _ := range s {
		result[k] = struct{}{}
	}
	return result
}

// Copy - создает независимую копию
func (v Int16) Copy() Int16 {
	return Int16FromSet(v)
}

// ToSlice - возвращает в виде []int16
func (v Int16) ToSlice() []int16 {
	result := _int16Slice{}
	for k, _ := range v {
		result = append(result, k)
	}
	sort.Sort(result)
	return result
}

func _Int16ToSliceTestHelper(kk []int16) []int16 {
	return Int16FromSlice(kk).ToSlice()
}

// String - поддержка интерфейса Stringer
func (v Int16) String() string {
	result, _ := json.Marshal(v)
	return string(result)
}

// MarshalJSON - поддержка интерфейса MarshalJSON
func (v Int16) MarshalJSON() ([]byte, error) {
	result := []interface{}{}
	for _, k := range v.ToSlice() {
		result = append(result, k)
	}
	return json.Marshal(result)
}

// ToSliceOfStrings - возвращает []string строковых представлений элементов множества
func (v Int16) ToSliceOfStrings() []string {
	result := []string{}
	for k, _ := range v {
		result = append(result, strconv.FormatInt(int64(k), 10))
	}
	sort.Strings(result)
	return result
}

// Has - возвращает true, если множество содержит заданный элемент, в противном случае - false
func (v Int16) Has(k int16) bool {
	_, ok := v[k]
	return ok
}

/*
HasAny - возвращает true, если множество содержит хотя бы один из заданныx элементов, в противном случае - false.
HasAny(<пустой набор/множесто>) возвращает false
*/
func (v Int16) HasAny(kk ...int16) bool {
	for _, k := range kk {
		if _, ok := v[k]; ok {
			return true
		}
	}
	return false
}

/*
HasAnyOfSlice - возвращает true, если множество содержит хотя бы один из заданныx элементов, в противном случае - false.
HasAny(<пустой набор/множесто>) возвращает false
*/
func (v Int16) HasAnyOfSlice(kk []int16) bool {
	for _, k := range kk {
		if _, ok := v[k]; ok {
			return true
		}
	}
	return false
}

/*
HasAnyOfSet - возвращает true, если множество содержит хотя бы один из заданныx элементов, в противном случае - false.
HasAny(<пустой набор/множесто>) возвращает false
*/
func (v Int16) HasAnyOfSet(s Int16) bool {
	for k, _ := range s {
		if _, ok := v[k]; ok {
			return true
		}
	}
	return false
}

/*
HasEach - возвращает true, если множество содержит все заданные элементы, в противном случае - false.
HasEach(<пустой набор/множесто>) возвращает true
*/
func (v Int16) HasEach(kk ...int16) bool {
	for _, k := range kk {
		if _, ok := v[k]; !ok {
			return false
		}
	}
	return true
}

/*
HasEachOfSlice - возвращает true, если множество содержит все заданные элементы, в противном случае - false.
HasEach(<пустой набор/множесто>) возвращает true
*/
func (v Int16) HasEachOfSlice(kk []int16) bool {
	for _, k := range kk {
		if _, ok := v[k]; !ok {
			return false
		}
	}
	return true
}

/*
HasEachOfSet - возвращает true, если множество содержит все заданные элементы, в противном случае - false.
HasEach(<пустой набор/множесто>) возвращает true
*/
func (v Int16) HasEachOfSet(s Int16) bool {
	for k, _ := range s {
		if _, ok := v[k]; !ok {
			return false
		}
	}
	return true
}

// Add - добавляет элементы в множество v
func (v Int16) Add(kk ...int16) {
	for _, k := range kk {
		v[k] = struct{}{}
	}
}

func (v Int16) _AddTestHelper(kk ...int16) Int16 {
	result := v.Copy()
	result.Add(kk...)
	return result
}

// AddSlice - добавляет элементы в множество v
func (v Int16) AddSlice(kk []int16) {
	for _, k := range kk {
		v[k] = struct{}{}
	}
}

func (v Int16) _AddSliceTestHelper(kk []int16) Int16 {
	result := v.Copy()
	result.AddSlice(kk)
	return result
}

// AddSet - добавляет элементы в множество v
func (v Int16) AddSet(s Int16) {
	for k, _ := range s {
		v[k] = struct{}{}
	}
}

func (v Int16) _AddSetTestHelper(s Int16) Int16 {
	result := v.Copy()
	result.AddSet(s)
	return result
}

// Del - удаляет элементы из множествa v
func (v Int16) Del(kk ...int16) {
	for _, k := range kk {
		delete(v, k)
	}
}

func (v Int16) _DelTestHelper(kk ...int16) Int16 {
	result := v.Copy()
	result.Del(kk...)
	return result
}

// DelSlice - удаляет элементы из множествa v
func (v Int16) DelSlice(kk []int16) {
	for _, k := range kk {
		delete(v, k)
	}
}

func (v Int16) _DelSliceTestHelper(kk []int16) Int16 {
	result := v.Copy()
	result.DelSlice(kk)
	return result
}

// DelSet - удаляет элементы из множествa v
func (v Int16) DelSet(s Int16) {
	for k, _ := range s {
		delete(v, k)
	}
}

func (v Int16) _DelSetTestHelper(s Int16) Int16 {
	result := v.Copy()
	result.DelSet(s)
	return result
}

// Union - возвращает результат объединения двух множеств. Исходные множества остаются без изменений
func (v Int16) Union(s Int16) Int16 {
	result := v.Copy()
	result.AddSet(s)
	return result
}

// Intersect - возвращает результат пересечения двух множеств. Исходные множества остаются без изменений
func (v Int16) Intersect(s Int16) Int16 {
	result := Int16{}
	for k, _ := range v {
		if _, ok := s[k]; ok {
			result[k] = struct{}{}
		}
	}
	return result
}

// Subtract - возвращает результат вычитания двух множеств. Исходные множества остаются без изменений
func (v Int16) Subtract(s Int16) Int16 {
	result := Int16{}
	for k, _ := range v {
		if _, ok := s[k]; !ok {
			result[k] = struct{}{}
		}
	}
	return result
}

type _int16Slice []int16

func (v _int16Slice) Len() int {
	return len(v)
}

func (v _int16Slice) Swap(i int, j int) {
	v[i], v[j] = v[j], v[i]
}

func (v _int16Slice) Less(i int, j int) bool {
	return v[i] < v[j]
}
