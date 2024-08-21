package lark

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
	webhookURL := fmt.Sprintf("https://open.larksuite.com/open-apis/bot/v2/hook/%s", token)
	return &notify{
		n: notifier.New(webhookURL, func(text string) io.Reader {
			return bytes.NewBufferString(fmt.Sprintf(`{"msg_type":"text","content":{"text":"%s"}}`, text))
		}),
	}
}

func (n *notify) SendMessage(ctx context.Context, message string) error {
	return n.n.SendMessage(ctx, message)
}
