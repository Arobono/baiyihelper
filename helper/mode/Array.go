package mode

import (
	"errors"
	"sort"

	"github.com/Arobono/baiyihelper/helper/cal"
	"github.com/gogf/gf/v2/util/gconv"
)

type Array struct {
	data []interface{} // 泛型数组
	size int           // 元素数量
}

type ArrayInterface interface {
	AddIndex(int, interface{}) // 插入元素
	Add(interface{})           // 添加到末尾
	AddFirst(interface{})      // 添加到首位

	Remove(index int) (value interface{}, err error)            // 删除
	RemoveFirst() (interface{}, error)                          // 删除
	RemoveLast() (interface{}, error)                           // 删除
	Clear()                                                     //清空数组
	RemoveElement(value interface{}) (e interface{}, err error) //删除指定元素

	Find(interface{}) int      // 查找元素返回第一个索引
	FindAll(interface{}) []int // 查找元素返回所有索引
	Contains(interface{}) bool // 查找是否存在元素

	Get(int) interface{}  //索引查询
	Set(int, interface{}) //索引修改

	// 基本方法
	Len() int               // 获得数组容量
	Count() int             // 获得元素个数
	IsNull() bool           // 查看数组是否为空
	GetData() []interface{} //切片取数据，移除nil项

	Where(cfunc func(interface{}) bool) *Array               //筛选
	FirstOrDefault(cfunc func(interface{}) bool) interface{} //筛选出一项

	Min() (min float64)                    //废弃 只能用在泛型数组是[]float64的情况
	Max() (max float64)                    //废弃 只能用在泛型数组是[]float64的情况
	OrderBy() (newArr []float64)           //废弃 只能用在泛型数组是[]float64的情况
	OrderByDescending() (newArr []float64) //废弃 只能用在泛型数组是[]float64的情况

	SumFloat64(cfunc func(interface{}) float64) float64 //统计Float64求和
	AvgFloat64(cfunc func(interface{}) float64) float64 //统计Float64求平均
	Take(num int) (res *Array)                          //只取几个数据
	Skip(num int) (res *Array)                          //跳过几个数据
}

// 获得自定义数组，参数为数组的初始长度
func GetArray(capacity int) *Array {
	arr := &Array{}
	arr.data = make([]interface{}, capacity)
	arr.size = 0

	// grades := []int{78, 92, 99, 37, 81}
	// min := From(grades).Min()
	// fmt.Println(min)

	return arr
}

// 获得数组容量
func (a *Array) Len() int {
	return len(a.data)
}

// 获得数组元素个数
func (a *Array) Count() int {
	return a.size
}

// 判断数组是否为空
func (a *Array) IsNull() bool {
	return a.size == 0
}

// newCapacity 新数组容量逻辑：声明新的数组，将原数组的值 copy 到新数组中
func (a *Array) resize(newCapacity int) {
	newArr := make([]interface{}, newCapacity)
	for i := 0; i < a.size; i++ {
		newArr[i] = a.data[i]
	}
	a.data = newArr
}

//判断索引是否越界
func (array *Array) checkIndex(index int) bool {
	if index < 0 || index >= array.size {
		return true
	}

	return false
}

func (a *Array) GetData() []interface{} {
	newArr := make([]interface{}, a.size)
	for i := 0; i < a.size; i++ {
		newArr[i] = a.data[i]
	}
	a.data = newArr
	return a.data
}

// 获得元素的首个索引，不存在则返回 -1
func (a *Array) Find(element interface{}) int {
	for i := 0; i < a.size; i++ {
		if element == a.data[i] {
			return i
		}
	}
	return -1
}

// 获得元素的所有索引，返回索引组成的切片
func (a *Array) FindAll(element interface{}) (indexes []int) {
	for i := 0; i < a.size; i++ {
		if element == a.data[i] {
			indexes = append(indexes, i)
		}
	}
	return
}

// 查看数组是否存在元素，返回 bool
func (a *Array) Contains(element interface{}) bool {
	if a.Find(element) == -1 {
		return false
	}
	return true
}

// 获得索引对应元素，需要判断索引有效范围
func (a *Array) Get(index int) interface{} {
	if index < 0 || index > a.size-1 {
		panic("Get failed, index is illegal.")
	}
	return a.data[index]
}

//修改索引对应元素值
func (a *Array) Set(index int, element interface{}) {
	if index < 0 || index > a.size-1 {
		panic("Set failed, index is illegal.")
	}
	a.data[index] = element
}

//添加（根据index插入
func (a *Array) AddIndex(index int, element interface{}) {
	if index < 0 || index > a.Len() {
		panic("Add failed, require index >= 0 and index <= capacity")
	}
	// 数组已满则扩容
	if a.size == len(a.data) {
		a.resize(2 * a.size)
	}
	// 将插入的索引位置之后的元素后移，腾出插入位置
	for i := a.size - 1; i > index; i-- {
		a.data[i+1] = a.data[i]
	}
	a.data[index] = element
	// 维护数组元素的数量
	a.size++
}

//添加（最后一位
func (a *Array) Add(element interface{}) {
	a.AddIndex(a.size, element)
}

//添加（第一位
func (a *Array) AddFirst(element interface{}) {
	a.AddIndex(0, element)
}

// 删除 index 位置的元素，并返回
func (array *Array) RemoveIndex(index int) (value interface{}, err error) {
	if array.checkIndex(index) {
		err = errors.New("Remove failed. Index is illegal.")
		return
	}

	value = array.data[index]
	for i := index + 1; i < array.size; i++ {
		//数据全部往前挪动一位,覆盖需要删除的元素
		array.data[i-1] = array.data[i]
	}

	array.size--
	array.data[array.size] = nil //loitering objects != memory leak

	cap := array.Len()
	if array.size == cap/4 && cap/2 != 0 {
		array.resize(cap / 2)
	}
	return
}

//删除数组首个元素
func (array *Array) RemoveFirst() (interface{}, error) {
	return array.RemoveIndex(0)
}

//删除末尾元素
func (array *Array) RemoveLast() (interface{}, error) {
	return array.RemoveIndex(int(array.size - 1))
}

//从数组中删除指定元素
func (array *Array) Remove(value interface{}) (e interface{}, err error) {
	index := array.Find(value)
	if index != -1 {
		e, err = array.RemoveIndex(index)
	}
	return
}

//清空数组
func (array *Array) Clear() {
	array.data = make([]interface{}, array.size)
	array.size = 0
}

func (a *Array) Where(cfunc func(interface{}) bool) *Array {
	newArray := GetArray(10)
	for _, item := range a.data {
		if item != nil && cfunc(item) {
			newArray.Add(item)
		}
	}
	return newArray
}

func (a *Array) FirstOrDefault(cfunc func(interface{}) bool) interface{} {
	var item interface{}
	for _, v := range a.data {
		if v != nil && cfunc(v) {
			return v
		}
	}
	return item
}

// 取最小值 float64
func (a *Array) Min() (min float64) {
	min = 0
	for _, v := range a.data {
		if v != nil && v.(float64) < min {
			min = v.(float64)
		}
	}
	return
}

// 取最大值 float64
func (a *Array) Max() (max float64) {
	max = 0
	for _, v := range a.data {
		if v != nil && v.(float64) > max {
			max = v.(float64)
		}
	}
	return
}

//排序 float64
func (a *Array) OrderBy() (newArr []float64) {
	newArr = make([]float64, a.size)
	for i := 0; i < a.size; i++ {
		newArr[i] = gconv.Float64(a.data[i])
	}
	sort.Float64s(newArr)
	return
}

//降序 float64
func (a *Array) OrderByDescending() (newArr []float64) {
	newArr = make([]float64, a.size)
	for i := 0; i < a.size; i++ {
		newArr[i] = gconv.Float64(a.data[i])
	}
	sort.Sort(sort.Reverse(sort.Float64Slice(newArr)))
	return
}

func (a *Array) SumFloat64(cfunc func(interface{}) float64) float64 {
	var sum float64 = 0
	for _, item := range a.data {
		if item != nil {
			sum += cfunc(item)
		}
	}
	return sum
}

func (a *Array) AvgFloat64(cfunc func(interface{}) float64) float64 {
	var sum = a.SumFloat64(cfunc)
	var avg = cal.Divide(sum, float64(a.size))
	return avg
}

func (a *Array) Take(num int) (res *Array) {
	if a.size < num {
		num = a.size
	}
	newArr := make([]interface{}, num)
	for i := 0; i < num; i++ {
		newArr[i] = a.data[i]
	}
	res.data = newArr
	res.size = num
	return
}

func (a *Array) Skip(num int) (res *Array) {

	var size = 0
	if num < a.Count() {
		size = a.Count() - num
	}
	newArr := make([]interface{}, size)
	var p = 0
	for i := 0; i < a.size; i++ {
		if i >= num {
			newArr[p] = a.data[i]
			p++
		}

	}
	res.data = newArr
	res.size = size
	return
}
