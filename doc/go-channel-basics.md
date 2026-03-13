# Go Channel 入门

这份文档专门解释 Go 里的 `channel` 是什么，以及下面三种类型分别是什么意思：

- `chan T` = 可读可写通道
- `<-chan T` = 只读通道
- `chan<- T` = 只写通道

如果你刚开始学 Go 并发，最容易卡住的就是这里。

## 1. channel 是什么

`channel` 可以理解成 goroutine 之间传递数据和同步的“管道”。

最简单的例子：

```go
ch := make(chan int)

go func() {
    ch <- 10
}()

x := <-ch
fmt.Println(x) // 10
```

这段代码的意思是：

- 创建一个 `int` 类型的 channel
- 后台 goroutine 往里面发送 `10`
- 主流程从 channel 里接收这个值

所以你可以先记住：

- `ch <- 10` 表示发送
- `x := <-ch` 表示接收

## 2. `chan T`

`chan T` 表示一个普通 channel，可读可写。

比如：

```go
var ch chan int
```

意思是：

- `ch` 是一个 channel
- 里面传的是 `int`
- 这个 channel 可以发送，也可以接收

例子：

```go
ch := make(chan int)

go func() {
    ch <- 42
}()

v := <-ch
fmt.Println(v)
```

一句话：

`chan T` = 双向通道。

## 3. `<-chan T`

`<-chan T` 表示一个只读 channel。

比如：

```go
var ch <-chan int
```

意思是：

- 这是一个“只能接收、不能发送”的 channel

可以这样用：

```go
func consume(ch <-chan int) {
    v := <-ch
    fmt.Println(v)
}
```

这里 `consume()` 只能从 `ch` 里读数据，不能往里面写。

错误示例：

```go
func consume(ch <-chan int) {
    ch <- 10 // 编译报错
}
```

一句话：

`<-chan T` = 只读通道。

## 4. `chan<- T`

`chan<- T` 表示一个只写 channel。

比如：

```go
var ch chan<- int
```

意思是：

- 这是一个“只能发送、不能接收”的 channel

可以这样用：

```go
func produce(ch chan<- int) {
    ch <- 10
}
```

错误示例：

```go
func produce(ch chan<- int) {
    v := <-ch // 编译报错
    fmt.Println(v)
}
```

一句话：

`chan<- T` = 只写通道。

## 5. 三者对照

### `chan T`

```go
var ch chan int
ch <- 1
v := <-ch
```

能发，也能收。

### `<-chan T`

```go
var ch <-chan int
v := <-ch
```

只能收，不能发。

### `chan<- T`

```go
var ch chan<- int
ch <- 1
```

只能发，不能收。

## 6. 为什么要分只读和只写

这样做的好处是：函数职责更清楚，类型更安全。

比如：

```go
func producer(out chan<- int) {
    out <- 1
    out <- 2
}

func consumer(in <-chan int) {
    fmt.Println(<-in)
    fmt.Println(<-in)
}
```

这里一眼就能看出来：

- `producer` 负责写
- `consumer` 负责读

也避免在函数里误把读写写反。

## 7. 你项目里的例子

你代码里有这种函数签名：

```go
func workerPool(ctx context.Context, jobs <-chan Job, n int, sem semaphore) []<-chan Result
```

拆开理解：

- `jobs <-chan Job`
  表示 `jobs` 是只读任务通道，worker pool 只能从里面拿任务

- `[]<-chan Result`
  表示返回值是一个切片，切片里每个元素都是“只读结果通道”

也就是：

- 输入：一个任务 channel
- 输出：多个结果 channel

## 8. `[]<-chan Result` 为什么这么怪

因为它其实是两层类型叠在一起：

- `[]T` = `T` 的切片
- `T` 这里是 `<-chan Result`

所以：

```go
[]<-chan Result
```

就是：

“一个切片，里面每个元素都是只读的 `Result` channel”

你可以脑补成这样：

```go
[
    resultCh1,
    resultCh2,
    resultCh3,
]
```

每个 `resultCh` 都是 `<-chan Result`。

## 9. `jobs` 不是数组，而是共享任务队列

很多人第一次看到下面这种代码时，会误以为 `jobs` 是数组：

```go
job, ok := <-jobs
```

但这里的 `jobs` 不是 `[]Job`，而是：

```go
jobs <-chan Job
```

也就是说，`jobs` 是一个只读 channel，不是数组。

所以这句代码的意思不是：

```go
job = jobs[i]
```

而是：

“从共享任务队列里取出下一个任务”

更像 Python 里的：

```python
job = queue.get()
```

而不是：

```python
job = jobs[3]
```

## 10. `job, ok := <-jobs` 到底是什么意思

这句可以拆成两部分看：

- `job`
  取到的任务值
- `ok`
  这个 channel 是否还开着

例子：

```go
job, ok := <-jobs
if !ok {
    return
}
```

意思是：

- 如果 `jobs` 里还有任务，就取一个出来放到 `job`
- 如果 `jobs` 已经关闭了，`ok` 就是 `false`
- `ok == false` 通常表示“没有更多任务了，可以退出了”

## 11. 为什么不会被多个 worker 重复拿到

因为 channel 的接收是“消费型”的。

也就是说：

```go
job := <-jobs
```

不是“看一眼”，而是“拿走一个”。

一旦某个 worker 从 `jobs` 里取走了一个任务，这个任务就不在 channel 里了，别的 worker 不会再拿到同一个任务。

你可以把多个 worker 理解成多个工人一起从同一个任务队列取活：

```text
jobs channel -> worker1
             -> worker2
             -> worker3
```

谁先取到，谁就处理；同一个任务不会同时发给两个人。

所以你通常不需要提前知道“这是第几个任务”或者“下标是多少”，你只需要知道：

- 它是当前轮到这个 worker 取到的下一个任务

## 12. `go func() { ... }()` 是什么

这是 Go 里非常常见的写法，可以拆开理解：

- `func() { ... }`
  定义一个匿名函数
- 后面的 `()`
  立刻调用这个匿名函数
- 前面的 `go`
  不是普通调用，而是启动一个 goroutine 异步执行

比如：

```go
go func() {
    fmt.Println("hello")
}()
```

意思是：

- 启动一个后台 goroutine
- 让它去执行这段匿名函数代码

可以粗略类比成 Python：

```python
threading.Thread(target=run).start()
```

或者：

```python
asyncio.create_task(run())
```

## 13. `defer` 是什么

`defer` 的意思是：

“先记下来，等当前函数快结束的时候再执行”

比如：

```go
defer close(out)
```

意思不是立刻关闭 `out`，而是：

- 先继续执行后面的逻辑
- 等这个函数返回前，再自动执行 `close(out)`

它很适合做收尾动作，比如：

- 关闭文件
- 关闭 channel
- 解锁
- 回收资源

可以粗略类比成 Python 的 `finally`：

```python
try:
    ...
finally:
    close_out()
```

## 14. buffered channel 和 unbuffered channel

### 无缓冲

```go
ch := make(chan int)
```

发送和接收要彼此配合，常用于同步。

### 有缓冲

```go
ch := make(chan int, 3)
```

表示这个 channel 最多可以先放 3 个值。

比如：

```go
ch := make(chan int, 2)
ch <- 1
ch <- 2
// 再发一个就会阻塞，直到有人接收
```

## 15. 为什么 `semaphore` 用的是 `chan struct{}`

你项目里有：

```go
type semaphore chan struct{}
```

意思是拿 channel 当“并发名额控制器”。

往里面发一个空结构体：

```go
s <- struct{}{}
```

表示占用一个名额。

从里面取一个：

```go
<-s
```

表示释放一个名额。

这里用 `struct{}` 是因为它不携带业务数据，只表示“一个信号”，而且几乎不占空间。

## 16. 最常见的 channel 心智模型

可以先把 channel 理解成下面这三种角色：

- `chan T`
  我既能放东西进去，也能拿东西出来

- `<-chan T`
  我只能拿，不能放

- `chan<- T`
  我只能放，不能拿

## 17. 最短总结

记住这三句就够了：

- `chan T` = 可读可写通道
- `<-chan T` = 只读通道
- `chan<- T` = 只写通道

再记一个发送接收语法：

- `ch <- v` = 发送
- `v := <-ch` = 接收

如果这两条记住了，Go 的很多并发代码就没那么吓人了。
