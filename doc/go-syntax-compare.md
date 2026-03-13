# Go 和 Java/Python/JS 基本语法对照表

这份文档主要给刚接触 Go 的人快速建立语法映射，重点放在最常用、最容易和 Java / Python / JavaScript 搞混的部分。

## 变量

### Go

```go
x := 10
var y int = 20
var z int
```

### Java

```java
int x = 10;
int y = 20;
int z;
```

### Python

```python
x = 10
y = 20
z = None
```

### JavaScript

```javascript
let x = 10;
let y = 20;
let z;
```

记忆点：

- Go 的 `:=` 是“声明并赋值”
- Go 的 `=` 是“给已有变量重新赋值”
- Go 的类型通常写在变量名后面

## 常量

### Go

```go
const Pi = 3.14
```

### Java

```java
final double PI = 3.14;
```

### Python

```python
PI = 3.14
```

### JavaScript

```javascript
const PI = 3.14;
```

## if

### Go

```go
if x > 0 {
    fmt.Println("positive")
} else {
    fmt.Println("non-positive")
}
```

### Java

```java
if (x > 0) {
    System.out.println("positive");
} else {
    System.out.println("non-positive");
}
```

### Python

```python
if x > 0:
    print("positive")
else:
    print("non-positive")
```

### JavaScript

```javascript
if (x > 0) {
  console.log("positive");
} else {
  console.log("non-positive");
}
```

记忆点：

- Go 的 `if` 没有括号
- Go 的 `if` 必须带 `{}` 块

Go 里还有一个很常见的写法：

```go
if err != nil {
    return err
}
```

也可以在 `if` 里顺手声明变量：

```go
if v, err := f(); err != nil {
    return err
} else {
    fmt.Println(v)
}
```

## for

Go 只有 `for`，没有单独的 `while`。

### Go

```go
for i := 0; i < 10; i++ {
    fmt.Println(i)
}

for x < 10 {
    x++
}

for {
    break
}
```

### Java

```java
for (int i = 0; i < 10; i++) {
    System.out.println(i);
}

while (x < 10) {
    x++;
}
```

### Python

```python
for i in range(10):
    print(i)

while x < 10:
    x += 1
```

### JavaScript

```javascript
for (let i = 0; i < 10; i++) {
  console.log(i);
}

while (x < 10) {
  x++;
}
```

记忆点：

- Go 用一个 `for` 同时覆盖 `for`、`while`、死循环

## 遍历数组 / 切片 / 列表

### Go

```go
nums := []int{10, 20, 30}

for i, v := range nums {
    fmt.Println(i, v)
}
```

如果不要下标：

```go
for _, v := range nums {
    fmt.Println(v)
}
```

### Java

```java
int[] nums = {10, 20, 30};
for (int i = 0; i < nums.length; i++) {
    System.out.println(i + " " + nums[i]);
}
```

### Python

```python
nums = [10, 20, 30]
for i, v in enumerate(nums):
    print(i, v)
```

### JavaScript

```javascript
const nums = [10, 20, 30];
nums.forEach((v, i) => {
  console.log(i, v);
});
```

记忆点：

- Go 的 `range` 很常用
- `_` 表示“这个值我不要”

## 函数

### Go

```go
func add(a int, b int) int {
    return a + b
}
```

简写：

```go
func add(a, b int) int {
    return a + b
}
```

### Java

```java
int add(int a, int b) {
    return a + b;
}
```

### Python

```python
def add(a, b):
    return a + b
```

### JavaScript

```javascript
function add(a, b) {
  return a + b;
}
```

记忆点：

- Go 参数类型写在后面
- Go 返回类型也写在后面

## 多返回值

Go 非常常见的一种写法是返回“结果 + 错误”。

### Go

```go
func div(a, b int) (int, error) {
    if b == 0 {
        return 0, fmt.Errorf("divide by zero")
    }
    return a / b, nil
}
```

调用：

```go
v, err := div(10, 2)
if err != nil {
    fmt.Println(err)
    return
}
fmt.Println(v)
```

记忆点：

- Go 经常返回 `(value, error)`
- 所以你会经常看到 `if err != nil`

## 数组和切片

### Go

```go
arr := [3]int{1, 2, 3}  // 数组，长度固定
s := []int{1, 2, 3}     // 切片，更常用
s = append(s, 4)
```

### Java

```java
int[] arr = {1, 2, 3};
```

### Python

```python
arr = [1, 2, 3]
arr.append(4)
```

### JavaScript

```javascript
const arr = [1, 2, 3];
arr.push(4);
```

记忆点：

- Go 日常里更常用的是切片 `slice`
- `append()` 是非常高频的内置函数

## map / 字典 / 对象

### Go

```go
m := map[string]int{
    "a": 1,
    "b": 2,
}

fmt.Println(m["a"])
```

判断 key 是否存在：

```go
v, ok := m["c"]
if ok {
    fmt.Println(v)
}
```

### Java

```java
Map<String, Integer> m = new HashMap<>();
m.put("a", 1);
```

### Python

```python
m = {"a": 1, "b": 2}
print(m["a"])
```

### JavaScript

```javascript
const m = { a: 1, b: 2 };
console.log(m.a);
```

记忆点：

- Go 常用 `v, ok := m[key]` 判断某个 key 在不在

## 结构体 / 类 / 对象

### Go

```go
type User struct {
    Name string
    Age  int
}

u := User{Name: "Tom", Age: 18}
fmt.Println(u.Name)
```

### Java

```java
class User {
    String name;
    int age;
}
```

### Python

```python
class User:
    def __init__(self, name, age):
        self.name = name
        self.age = age
```

### JavaScript

```javascript
const user = { name: "Tom", age: 18 };
```

记忆点：

- Go 有 `struct`
- Go 不走传统面向对象里那种重继承路线，更偏组合

## 方法

### Go

```go
type User struct {
    Name string
}

func (u User) SayHello() {
    fmt.Println("hello", u.Name)
}
```

### Java

```java
class User {
    String name;

    void sayHello() {
        System.out.println("hello " + name);
    }
}
```

记忆点：

- `func (u User) SayHello()` 里的 `(u User)` 是 receiver
- 可以理解成“给某个类型绑定方法”

## 指针

### Go

```go
x := 10
p := &x
fmt.Println(*p)
*p = 20
```

### Java / Python / JavaScript

这些语言通常没有显式指针语法。

记忆点：

- `&x` 是取地址
- `*p` 是解引用
- Go 有指针，但不像 C 那样支持随意做指针运算

## 错误处理

### Go

```go
v, err := doSomething()
if err != nil {
    return err
}
fmt.Println(v)
```

### Java

```java
try {
    var v = doSomething();
} catch (Exception e) {
    e.printStackTrace();
}
```

### Python

```python
try:
    v = do_something()
except Exception as e:
    print(e)
```

### JavaScript

```javascript
try {
  const v = doSomething();
} catch (e) {
  console.error(e);
}
```

记忆点：

- Go 倾向于显式返回错误
- Go 不靠异常做主流程控制

## switch

### Go

```go
switch x {
case 1:
    fmt.Println("one")
case 2:
    fmt.Println("two")
default:
    fmt.Println("other")
}
```

记忆点：

- Go 的 `switch` 默认不会自动落到下一个 `case`
- 所以一般不需要像 Java / JS 那样频繁写 `break`

## 可见性

Go 没有 `public` / `private` 关键字，而是靠首字母大小写区分。

```go
type User struct {
    Name string // public
    age  int    // private
}
```

记忆点：

- 大写开头：可导出
- 小写开头：包内可见

## 包和导入

### Go

```go
package main

import "fmt"
```

使用：

```go
fmt.Println("hello")
```

### Java

```java
import java.util.List;
```

### Python

```python
import os
```

### JavaScript

```javascript
import fs from "fs";
```

## 最值得先适应的 Go 模板

如果你刚开始学 Go，先把这几个写顺：

```go
x := 10
```

```go
if err != nil {
    return err
}
```

```go
for _, v := range items {
    fmt.Println(v)
}
```

```go
func add(a, b int) int {
    return a + b
}
```

```go
v, ok := m["key"]
```

## 一句话总结

Go 和 Java / Python / JavaScript 最大的不同，不是“功能少”，而是：

- 语法更硬
- 规则更统一
- 写法更模板化
- 错误处理更显式
- 并发支持是语言级别的

如果把这些常见模板写顺了，Go 的“别扭感”会下降很多。
