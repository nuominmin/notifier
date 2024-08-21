package qyweixin

import (
	"fmt"
	"github.com/nuominmin/notifier"
)

const (
	WEBHOOK_URL_FORMAT = "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=%s"
	TEXT_BODY_FORMAT   = `{"msgtype": "text", "text": {"content": "%s"}}`
)

func NewNotifier(token string) notifier.Notifier {
	return notifier.NewNotifier(fmt.Sprintf(WEBHOOK_URL_FORMAT, token), TEXT_BODY_FORMAT)
}

func NewDelayNotifier(token string) notifier.Notifier {
	return notifier.NewDelayNotifier(fmt.Sprintf(WEBHOOK_URL_FORMAT, token), TEXT_BODY_FORMAT)
}
