package lark

import (
	"context"
	"github.com/nuominmin/notifier"
	"testing"
)

func TestNotifier(t *testing.T) {
	n := NewNotifier(notifier.LARK_TOKEN)
	_ = n.SendMessage(context.Background(), "test")
	_ = n.SendMessage(context.Background(), "test1")
	_ = n.SendMessage(context.Background(), "test1")
	_ = n.SendMessage(context.Background(), "test1")
	_ = n.SendMessage(context.Background(), "test1")
	_ = n.SendMessage(context.Background(), "test1")
	_ = n.SendMessage(context.Background(), "test1")
	_ = n.SendMessage(context.Background(), "test1")
	_ = n.SendMessage(context.Background(), "test1")
}

func TestDelayNotifier(t *testing.T) {
	n := NewDelayNotifier(notifier.LARK_TOKEN)
	_ = n.SendMessage(context.Background(), "testtesttestttesttesttestesttesttesttest")
	_ = n.SendMessage(context.Background(), "testtesttesttesttesttesttesttesttesttest")
	select {}
}
