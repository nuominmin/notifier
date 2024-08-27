package qyweixin

import (
	"github.com/nuominmin/notifier"
)

const (
	WEBHOOK_URL_FORMAT = "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=%s"
	TEXT_BODY_FORMAT   = `{"msgtype": "text", "text": {"content": "%s"}}`
)

func NewNotifier(token string) notifier.Notifier {
	return notifier.NewNotifier(WEBHOOK_URL_FORMAT, TEXT_BODY_FORMAT, token)
}

func NewDelayNotifier(token string) notifier.DelayNotifier {
	return notifier.NewDelayNotifier(WEBHOOK_URL_FORMAT, TEXT_BODY_FORMAT, token)
}
