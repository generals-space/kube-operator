# kube-operator

参考文章

1. [net/http2: http.Flusher doesn't send data to the client](https://github.com/golang/go/issues/18510#issuecomment-270319411)
2. [bugs/golang/18510/server.go](https://github.com/odeke-em/bugs/blob/master/golang/18510/server.go)
    - 这个server示例用到了`fmt.Fprintf()`
3. [Send a chunked HTTP response from a Go server](https://stackoverflow.com/questions/26769626/send-a-chunked-http-response-from-a-go-server)
    - 高票答案也用到了`fmt.Fprintf()`
4. [Go：http Transfer-Encoding chunked 实时读写](https://blog.csdn.net/test1280/article/details/116273611)
    - client 侧代码示例
    - client 代码示例已经提到要使用'\n'进行分隔.

本示例模拟了 apiserver 中的 watch 接口实现(其实就是一个长连接示例), 非常简单, 无需要任何外部库.

在 apiserver 中, handler 实现可见[kubernetes apiserver源码](https://github.com/kubernetes/kubernetes/blob/v1.17.2/staging/src/k8s.io/apiserver/pkg/endpoints/handlers/watch.go#L206)

启动 server 端, ta只有一个`/watch`接口. 当客户端访问时, server 会每隔1s向客户端输出当前时间.

```console
$ go run server.go
2022/10/30 21:46:18 http server started at :8897...
```

直接使用 curl 访问, 无需添加任何额外参数, 就是一个简单的 Get 请求.

```console
$ curl localhost:8897/watch
2022-10-30 20:59:50
2022-10-30 21:00:28
2022-10-30 21:00:46
2022-10-30 21:00:55
2022-10-30 21:01:02
2022-10-30 21:01:07
```

也可以执行 client.go 文件进行模拟.

```console
$ go run client.go
2022/10/30 22:19:39 response: *http.Response
2022/10/30 22:19:39 response.Body: *http.bodyEOFSignal
2022/10/30 22:19:39 2022-10-30 22:19:39
2022/10/30 22:19:40 2022-10-30 22:19:40
2022/10/30 22:19:41 2022-10-30 22:19:41
^Cexit status 2
```

## FAQ

刚开始写 server 端代码时, 遇到了一个问题.

curl 访问时无法实时打印时间消息, 只能等到 watch handler 函数返回时(可以用某种方式(如计数器) return), 或者客户端自动断开, 出现如下错误时, 会把积攒的输入全部打印出来.

```console
$ curl localhost:8897/watch
curl: (18) transfer closed with outstanding read data remaining
2022-10-30 21:11:592022-10-30 21:12:002022-10-30 21:12:012022-10-30 21:12:02
```

当时代码中是这么写的

```go
    log.Printf("info: %s\n", info)
    resp.Write([]byte(info))
    flusher.Flush()
```

kubernetes中, watch 方法是调用了 apimachinery 库实现的, 不过最终也还是调用了`response.Write()`方法.

我换了很多方法, 比如想着可能是`response.Write()`其实需要接收一个`bytes.Buffer{}`中的数据, 需要在`flusher.Flush()`后再调用一把`buf.Reset()`...

最终找到了参考文章1, 其中的 issue 回答者提到了参考文章2, 其中的 server 示例用的是`fmt.Fprintf()`.

```go
    log.Printf("info: %s\n", info)
    fmt.Fprintf(resp, "%s\n", info)
    flusher.Flush()
```

我想了2个小时也没想明白为什么这个函数可以, 其他就不行, 明明最终都是调用的`response.Write()`.

最后偶然发现, 原来写的的数据要改`\n`结尾...
