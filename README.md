# notifier

一个轻量级的、可扩展的 Go 语言包，旨在为多种消息通知服务提供统一的接口。通过实现 Notifier 接口，用户可以轻松地集成多个第三方消息通知服务（例如 Lark）到他们的应用程序中
## 安装
你可以使用' go get '来安装这个包:

```sh
go get github.com/nuominmin/notifier
```

## 示例
```go
n := lark.NewNotifier(notifier.LARK_TOKEN)
_ = n.SendMessage(context.Background(), "test")
_ = n.SendMessage(context.Background(), "test1")
```