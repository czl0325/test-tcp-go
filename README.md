# test-tcp-go
测试go语言的tcp连接

先执行命令编译 `go build -o server main.go server.go user.go`，
编译后生成server文件，在运行`./server`启动服务器。
另外开一个终端，执行`nc 127.0.0.1 8888`可以连接服务器