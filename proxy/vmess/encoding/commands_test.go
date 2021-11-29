package encoding_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/qazz-shyper/website/common"
	"github.com/qazz-shyper/website/common/buf"
	"github.com/qazz-shyper/website/common/protocol"
	"github.com/qazz-shyper/website/common/uuid"
	. "github.com/qazz-shyper/website/proxy/vmess/encoding"
)

func TestSwitchAccount(t *testing.T) {
	sa := &protocol.CommandSwitchAccount{
		Port:     1234,
		ID:       uuid.New(),
		AlterIds: 1024,
		Level:    128,
		ValidMin: 16,
	}

	buffer := buf.New()
	common.Must(MarshalCommand(sa, buffer))

	cmd, err := UnmarshalCommand(1, buffer.BytesFrom(2))
	common.Must(err)

	sa2, ok := cmd.(*protocol.CommandSwitchAccount)
	if !ok {
		t.Fatal("failed to convert command to CommandSwitchAccount")
	}
	if r := cmp.Diff(sa2, sa); r != "" {
		t.Error(r)
	}
}
