package qyweixin

import (
	"bytes"
	"context"
	"fmt"
	"github.com/nuominmin/notifier"
	"io"
)

type notify struct {
	n notifier.Notifier
}

func NewNotifier(token string) notifier.Notifier {
	webhookURL := fmt.Sprintf("https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=%s", token)
	return &notify{
		n: notifier.New(webhookURL, func(text string) io.Reader {
			return bytes.NewBufferString(fmt.Sprintf(`{"msgtype": "text", "text": {"content": "%s"}}`, text))
		}),
	}
}

func (n *notify) SendMessage(ctx context.Context, message string) error {
	return n.n.SendMessage(ctx, message)
}
