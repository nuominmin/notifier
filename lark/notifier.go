package lark

import (
	"github.com/nuominmin/notifier"
)

const (
	WEBHOOK_URL_FORMAT = "https://open.larksuite.com/open-apis/bot/v2/hook/%s"
	TEXT_BODY_FORMAT   = `{"msg_type":"text","content":{"text":"%s"}}`
)

func NewNotifier(tokens ...string) notifier.Notifier {
	return notifier.NewNotifier(WEBHOOK_URL_FORMAT, TEXT_BODY_FORMAT, tokens...)
}

func NewDelayNotifier(tokens ...string) notifier.DelayNotifier {
	return notifier.NewDelayNotifier(WEBHOOK_URL_FORMAT, TEXT_BODY_FORMAT, tokens...)
}
