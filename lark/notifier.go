package lark

import (
	"fmt"
	"github.com/nuominmin/notifier"
)

const (
	WEBHOOK_URL_FORMAT = "https://open.larksuite.com/open-apis/bot/v2/hook/%s"
	TEXT_BODY_FORMAT   = `{"msg_type":"text","content":{"text":"%s"}}`
)

func NewNotifier(token string) notifier.Notifier {
	return notifier.NewNotifier(fmt.Sprintf(WEBHOOK_URL_FORMAT, token), TEXT_BODY_FORMAT)
}

func NewDelayNotifier(token string) notifier.Notifier {
	return notifier.NewDelayNotifier(fmt.Sprintf(WEBHOOK_URL_FORMAT, token), TEXT_BODY_FORMAT)
}
