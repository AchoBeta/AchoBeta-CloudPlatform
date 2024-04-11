## 关于json.marshal的源码分析以及反射的应用

在我们操作数据库的时候，我们需要将go的数据类型存储到redis中，而其中遇到的一个问题就是他们之间格式并不互通，所以需要一个转化的途径，好比怎么把大树变成家具，我们需要先对树进行统一处理成木材而后才可进行加工成指定家具

### json

json是一种轻量级的数据交换格式，是流行的最主要的数据交换之一，Marshal是将语言中的内存对象解析为json格式的字符串

例如：

```go
type Student struct {
	Name string `json:"name"`
	Id   string `json:"id"`
}

func TestHSet(t *testing.T) {
	stu := Student{
		Name: "hhah",
		Id: "1231231",
	}
	resByte ,err := json.Marshal(stu)
	if err != nil {
		tlog.Error(err)
		return
	}
	fmt.Println(resByte)
}
```

​	最后的执行结果为

```
=== RUN   TestHSet
[123 34 110 97 109 101 34 58 34 104 104 97 104 34 44 34 105 100 34 58 34 49 50 51 49 50 51 49 34 125]
--- PASS: TestHSet (0.00s)
PASS
```

json.Marshal的返回结果为byte数组，当然不只是结构体，还可以是别的类型，或者别的复合结构

其中，复合类型需要进行复杂的递归流程直至基本类型才能生成字节序列，简单类型可以直接生成字节序列。每种类型都对应一种编码函数，以下为源码部分：

```go
switch t.Kind() {
	case reflect.Bool:
		return boolEncoder
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return intEncoder
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return uintEncoder
	case reflect.Float32:
		return float32Encoder
	case reflect.Float64:
		return float64Encoder
	case reflect.String:
		return stringEncoder
	case reflect.Interface:
		return interfaceEncoder
	case reflect.Struct:
		return newStructEncoder(t)
	case reflect.Map:
		return newMapEncoder(t)
	case reflect.Slice:
		return newSliceEncoder(t)
	case reflect.Array:
		return newArrayEncoder(t)
	case reflect.Pointer:
		return newPtrEncoder(t)
	default:
		return unsupportedTypeEncoder
```

其中例如基础类型bool，int等基础类型只是简单的写入buffer里。除 String 内各种转义处理会略麻烦，其余都相对简单。需要注意：Channel， complex 以及函数不能被编码json字符串。当然，循环的数据结构也不行，它会导致marshal陷入死循环。

对于为什么函数不可被编码json字符串，因为反射类型没有函数类型

### 反射

go中结构体序列化json的过程，利用了语言内置的反射的特性来实现不同数据类型的通用逻辑

在编程语言里，每个对象都有自己的type以及自己的value，type表示类型，每个自定义的结构体都是一个不同的类型，可以用typeof来返回类型，而value则装载具体对象的值，可用valueof来返回。例如

```go
type Student struct {
	Name string `json:"name"`
	Id   string `json:"id"`
}

func TestHSet(t *testing.T) {
	stu := Student{
		Name: "1231",
		Id:   "1231231",
	}
	fmt.Println(reflect.TypeOf(stu))
  fmt.Println(reflect.ValueOf(stu))
}
```

输出结果为：

```
=== RUN   TestHSet
redis.Student
{1231 1231231}
--- PASS: TestHSet (0.00s)
PASS
```

在java中，Java的反射（reflection）机制是指在程序的运行状态中，可以构造任意一个类的对象，可以了解任意一个对象所属的类，可以了解任意一个类的成员变量和方法，可以调用任意一个对象的属性和方法。这种动态获取程序信息以及动态调用对象的功能称为Java语言的反射机制。反射被视为动态语言的关键。当然go是一种编译型的静态语言

reflect包提供了较为完善的机制来支持使用反射的特性，如 Type 和 Value 都提供了 Kind()⽅法⽤来获取其属于的 Kind 常量。(如上述例子，stu.Kind就是struct)

对⽤户⾃定义的不同结构体⽽⾔，其 reflect.Type 不⼀样。reflect.Type 之间的相互⽐较，会循环递归保证内部所有域确保⼀致。⽤户⾃定义的结构体和内置类型同样凑效。

那么提问2个问题

1. 上⽂提到的每种类型都有对应⼀种编码函数 encoderFunc，究竟是对不同的 Kind，还是对应不同的 Type？
2. 下⽂会讲到序列化的缓存中间结果，那么缓存是针对不同的 Kind 还是针对不同的 Type 来缓存？

第一个问题对于深入源码之后可以看到，它对应的是kind进行的编码函数



### 对于json序列化的解析流程

json.Marshal的流程，有两条主线，

一条是利用golang的反射原理，使用递归的解析结构体内所有的字段，生成字节序列，有两处递归，令外一条主线是尽可能的缓存可复用的中间状态结果，以提高性能，有三处缓存

对valueEncoder内部会进行预处理typeEncoder，而在预处理typeEncoder中，会调用newTypeEncoder，生成每种类型对应的encoderFunc，会对每种类型调用相对应的encoderFuce执行具体序列。

```go
func (e *encodeState) reflectValue(v reflect.Value, opts encOpts) {
	valueEncoder(v)(e, v, opts)
}

func valueEncoder(v reflect.Value) encoderFunc {
	if !v.IsValid() {
		return invalidValueEncoder
	}
	return typeEncoder(v.Type())
}
```

```go
/生成对应的encoderFunc
func typeEncoder(t reflect.Type) encoderFunc {
	if fi, ok := encoderCache.Load(t); ok {
		return fi.(encoderFunc)
	}

	// To deal with recursive types, populate the map with an
	// indirect func before we build it. This type waits on the
	// real func (f) to be ready and then calls it. This indirect
	// func is only used for recursive types.
	var (
		wg sync.WaitGroup
		f  encoderFunc
	)
	wg.Add(1)
	fi, loaded := encoderCache.LoadOrStore(t, encoderFunc(func(e *encodeState, v reflect.Value, opts encOpts) {
		wg.Wait()
		f(e, v, opts)
	}))
	if loaded {
		return fi.(encoderFunc)
	}

	// Compute the real encoder and replace the indirect func with it.
	f = newTypeEncoder(t, true)
	wg.Done()
	encoderCache.Store(t, f)
	return f
}


/该方法就是对每种类型调用的相对应的encoderFunc执行具体序列
// newTypeEncoder constructs an encoderFunc for a type.
// The returned encoder only checks CanAddr when allowAddr is true.
func newTypeEncoder(t reflect.Type, allowAddr bool) encoderFunc {
	// If we have a non-pointer value whose type implements
	// Marshaler with a value receiver, then we're better off taking
	// the address of the value - otherwise we end up with an
	// allocation as we cast the value to an interface.
	if t.Kind() != reflect.Pointer && allowAddr && reflect.PointerTo(t).Implements(marshalerType) {
		return newCondAddrEncoder(addrMarshalerEncoder, newTypeEncoder(t, false))
	}
	if t.Implements(marshalerType) {
		return marshalerEncoder
	}
	if t.Kind() != reflect.Pointer && allowAddr && reflect.PointerTo(t).Implements(textMarshalerType) {
		return newCondAddrEncoder(addrTextMarshalerEncoder, newTypeEncoder(t, false))
	}
	if t.Implements(textMarshalerType) {
		return textMarshalerEncoder
	}

	switch t.Kind() {
	case reflect.Bool:
		return boolEncoder
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return intEncoder
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return uintEncoder
	case reflect.Float32:
		return float32Encoder
	case reflect.Float64:
		return float64Encoder
	case reflect.String:
		return stringEncoder
	case reflect.Interface:
		return interfaceEncoder
	case reflect.Struct:
		return newStructEncoder(t)
	case reflect.Map:
		return newMapEncoder(t)
	case reflect.Slice:
		return newSliceEncoder(t)
	case reflect.Array:
		return newArrayEncoder(t)
	case reflect.Pointer:
		return newPtrEncoder(t)
	default:
		return unsupportedTypeEncoder
	}
}
```

<!-- valueEncoder的⽬的是输出每种 Kind 的 encoderFunc-->

#### 主线一：两处递归

1. 递归遍历结构体树状结构，对内部节点生成器对应的类型编码器encoderFunc，或者开发者自定义的编码器，递归的结束条件是最终递归之基本类型，生成基本类型编码器。
2. 启动类型编码器调用，依赖类型编码器函数内部递归，从根结底啊安一次调用整棵树的序列化函数，递归的结束条件是递归至基本类型编码器，生成字符编码。

递归是为了处理复杂数据类型，例如结构体，ptr指针，slice，array，map。Ptr 编码器函数 通过 t.Elem() 递归调⽤typeEncoder；Array/Slice 编码器函数通过 t.Elem() 递归调⽤ typeEncoder；Map 稍微复杂，不但通过 t.Elem() 递归调⽤ typeEncoder，其额外的操作是的对其 Key 进⾏处理，通过判断其 Key 类型

若key为非string类型且没有实现TextMarshal则输出不支持错误

例子：

```go
type Student struct {
    Name string `json:"name"`
    Age int
    CurrentClass *Class
    Classes []Class
    Friends map[string]Student
}
```

递归解析Student：

首先对结构体进行调用，遍历struct的属性，而后对于属性进行encode解析，如果属性是指针那么对该指针的元素类型进行获取（type.elem()），如果为嵌套结构体那么意味着递归操作，直到字段属性为空。



#### 主线二：三处缓存

##### sync.poll存储encodeState

sync.poll 是 Golang 官⽅提供⽤来缓存分配的内存的对象，以降低分配和 GC 压⼒。 序 列化中，encodeState 的⾸要作⽤是存储字符编码，其内部包含了 bytes.Buffer，由于 在 json.Marshal 在 IO 密集的业务程序中，通常会被⼤量的调⽤，如果不断的释放⽣成 新的 bytes.Buffer，会降低性能。 官⽅包的源码可以看到， encodeState 结构体被放 进 sync.poll 内==（var encodeStatePool ）==，来保存和复⽤临时对象，减少内存分配， 降低 GC 压力。

sync.poll 内部的 bytes.Buffer 提供可扩容的字节缓冲区，其实质是对切⽚的封装，结 构中包含⼀个 64 字节的⼩切⽚，避免⼩内存分配，并可以依据使⽤情况⾃动扩充。⽽ 且，其空切⽚ buf[:0] 在该场合下⾮常有⽤，是直接在原内存上进⾏操作，⾮常⾼效， 每次开始序列化之处，会将 Reset()。

##### sync.map存储encodeFunc


