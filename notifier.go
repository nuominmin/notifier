package notifier

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"sync/atomic"
)

// Notifier 消息通知接口
type Notifier interface {
	SendMessage(ctx context.Context, message string) error
}

type Notify struct {
	client           *http.Client
	method           string
	headers          map[string]string
	webhookUrlFormat string   // 机器人钩子地址格式
	textBodyFormat   string   // 消息转换为 JSON 的函数模板
	tokens           []string // 钩子地址所需要的 token
	currTokenIdx     uint64   // 当前使用 token 的索引
}

// NewNotifier 创建一个新的通知实例
func NewNotifier(webhookUrlFormat string, textBodyFormat string, tokens ...string) Notifier {
	return &Notify{
		client:           http.DefaultClient,
		method:           http.MethodPost, // 默认方法为 POST
		headers:          map[string]string{"Content-Type": "application/json"},
		webhookUrlFormat: webhookUrlFormat,
		textBodyFormat:   textBodyFormat,
		tokens:           tokens,
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

func (n *Notify) getWebhookUrl() string {
	totalToken := uint64(len(n.tokens))
	if totalToken == 0 {
		return ""
	}
	idx := atomic.AddUint64(&n.currTokenIdx, 1) - 1
	if idx >= totalToken*100 {
		atomic.StoreUint64(&n.currTokenIdx, 0) // 重置为 0，防止溢出
	}

	return fmt.Sprintf(n.webhookUrlFormat, n.tokens[idx%totalToken])
}

// SendMessage 发送消息到指定的 Webhook URL
func (n *Notify) SendMessage(ctx context.Context, message string) error {
	webhookUrl := n.getWebhookUrl()
	if webhookUrl == "" {
		return fmt.Errorf("no webhook url")
	}

	req, err := http.NewRequestWithContext(ctx, n.method, webhookUrl, bytes.NewBufferString(fmt.Sprintf(n.textBodyFormat, message)))
	if err != nil {
		return fmt.Errorf("create request error: %v", err)
	}

	for key, value := range n.headers {
		req.Header.Set(key, value)
	}

	var resp *http.Response
	if resp, err = n.client.Do(req); err != nil {
		return fmt.Errorf("send request error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send message, status code: %d", resp.StatusCode)
	}

	return nil
}
