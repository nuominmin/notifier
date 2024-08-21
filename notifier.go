package notifier

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

// Notifier 是一个通用的消息通知接口，支持泛型 T
type Notifier interface {
	SendMessage(ctx context.Context, message string) error
}

type Notify struct {
	client     *http.Client
	webhookURL string
	method     string
	headers    map[string]string
	formatFunc func(message string) io.Reader // 泛型消息转换为 JSON 的函数
}

// New 创建一个新的通知实例，支持泛型 T
func New(webhookURL string, formatFunc func(text string) io.Reader) Notifier {
	return &Notify{
		client:     http.DefaultClient,
		webhookURL: webhookURL,
		method:     http.MethodPost, // 默认方法为 POST
		headers:    map[string]string{"Content-Type": "application/json"},
		formatFunc: formatFunc,
	}
}

// SetClient 设置 HTTP 客户端
func (n *Notify) SetClient(client *http.Client) {
	n.client = client
}

// SetRequestConfig 设置 HTTP 请求方法和头部信息
func (n *Notify) SetRequestConfig(method string, headers map[string]string) {
	n.method = method
	n.headers = headers
}

// SendMessage 发送消息到指定的 Webhook URL
func (n *Notify) SendMessage(ctx context.Context, message string) error {
	req, err := http.NewRequestWithContext(ctx, n.method, n.webhookURL, n.formatFunc(message))
	if err != nil {
		return err
	}

	for key, value := range n.headers {
		req.Header.Set(key, value)
	}

	var resp *http.Response
	if resp, err = n.client.Do(req); err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send message, status code: %d", resp.StatusCode)
	}

	return nil
}
