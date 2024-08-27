package lark

import (
	"context"
	"github.com/nuominmin/notifier"
	"testing"
)

func TestNotifier(t *testing.T) {
	n := NewNotifier(notifier.LARK_TOKEN_A, notifier.LARK_TOKEN_B)
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
	//n := NewDelayNotifier(notifier.LARK_TOKEN)
	n := NewDelayNotifier(notifier.LARK_TOKEN_A, notifier.LARK_TOKEN_B)
	n.SetIdentity("test")

	_ = n.SendMessage(context.Background(), "testtesttestttesttesttestesttesttesttest")
	_ = n.SendMessage(context.Background(), "testtesttestt1esttesttesttesttesttesttest")
	_ = n.SendMessage(context.Background(), "testtesttestt2esttesttesttesttesttesttest")
	_ = n.SendMessage(context.Background(), "testtesttest3testtesttesttesttesttesttest")
	_ = n.SendMessage(context.Background(), "testtesttestt4esttesttesttesttesttesttest")
	_ = n.SendMessage(context.Background(), "testtesttestte5sttesttesttesttesttesttest")
	_ = n.SendMessage(context.Background(), "testtesttestte6sttesttesttesttesttesttest")
	_ = n.SendMessage(context.Background(), "testtesttesttest7testtesttesttesttesttest")
	n.Close()
}
