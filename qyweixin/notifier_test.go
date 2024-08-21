package qyweixin

import (
	"context"
	"github.com/nuominmin/notifier"
	"testing"
)

func TestNotifier(t *testing.T) {
	n := NewNotifier(notifier.QYWEIXIN_TOKEN)
	_ = n.SendMessage(context.Background(), "test")
	_ = n.SendMessage(context.Background(), "test1")
}
