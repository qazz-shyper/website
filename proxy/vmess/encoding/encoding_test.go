package encoding_test

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/qazz-shyper/website/common"
	"github.com/qazz-shyper/website/common/buf"
	"github.com/qazz-shyper/website/common/net"
	"github.com/qazz-shyper/website/common/protocol"
	"github.com/qazz-shyper/website/common/uuid"
	"github.com/qazz-shyper/website/proxy/vmess"
	. "github.com/qazz-shyper/website/proxy/vmess/encoding"
)

func toAccount(a *vmess.Account) protocol.Account {
	account, err := a.AsAccount()
	common.Must(err)
	return account
}

func TestRequestSerialization(t *testing.T) {
	user := &protocol.MemoryUser{
		Level: 0,
		Email: "test@example.com",
	}
	id := uuid.New()
	account := &vmess.Account{
		Id:      id.String(),
		AlterId: 0,
	}
	user.Account = toAccount(account)

	expectedRequest := &protocol.RequestHeader{
		Version:  1,
		User:     user,
		Command:  protocol.RequestCommandTCP,
		Address:  net.DomainAddress("www.example.com"),
		Port:     net.Port(443),
		Security: protocol.SecurityType_AES128_GCM,
	}

	buffer := buf.New()
	client := NewClientSession(context.TODO(), true, protocol.DefaultIDHash)
	common.Must(client.EncodeRequestHeader(expectedRequest, buffer))

	buffer2 := buf.New()
	buffer2.Write(buffer.Bytes())

	sessionHistory := NewSessionHistory()
	defer common.Close(sessionHistory)

	userValidator := vmess.NewTimedUserValidator(protocol.DefaultIDHash)
	userValidator.Add(user)
	defer common.Close(userValidator)

	server := NewServerSession(userValidator, sessionHistory)
	actualRequest, err := server.DecodeRequestHeader(buffer, false)
	common.Must(err)

	if r := cmp.Diff(actualRequest, expectedRequest, cmp.AllowUnexported(protocol.ID{})); r != "" {
		t.Error(r)
	}

	_, err = server.DecodeRequestHeader(buffer2, false)
	// anti replay attack
	if err == nil {
		t.Error("nil error")
	}
}

func TestInvalidRequest(t *testing.T) {
	user := &protocol.MemoryUser{
		Level: 0,
		Email: "test@example.com",
	}
	id := uuid.New()
	account := &vmess.Account{
		Id:      id.String(),
		AlterId: 0,
	}
	user.Account = toAccount(account)

	expectedRequest := &protocol.RequestHeader{
		Version:  1,
		User:     user,
		Command:  protocol.RequestCommand(100),
		Address:  net.DomainAddress("www.example.com"),
		Port:     net.Port(443),
		Security: protocol.SecurityType_AES128_GCM,
	}

	buffer := buf.New()
	client := NewClientSession(context.TODO(), true, protocol.DefaultIDHash)
	common.Must(client.EncodeRequestHeader(expectedRequest, buffer))

	buffer2 := buf.New()
	buffer2.Write(buffer.Bytes())

	sessionHistory := NewSessionHistory()
	defer common.Close(sessionHistory)

	userValidator := vmess.NewTimedUserValidator(protocol.DefaultIDHash)
	userValidator.Add(user)
	defer common.Close(userValidator)

	server := NewServerSession(userValidator, sessionHistory)
	_, err := server.DecodeRequestHeader(buffer, false)
	if err == nil {
		t.Error("nil error")
	}
}

func TestMuxRequest(t *testing.T) {
	user := &protocol.MemoryUser{
		Level: 0,
		Email: "test@example.com",
	}
	id := uuid.New()
	account := &vmess.Account{
		Id:      id.String(),
		AlterId: 0,
	}
	user.Account = toAccount(account)

	expectedRequest := &protocol.RequestHeader{
		Version:  1,
		User:     user,
		Command:  protocol.RequestCommandMux,
		Security: protocol.SecurityType_AES128_GCM,
		Address:  net.DomainAddress("v1.mux.cool"),
	}

	buffer := buf.New()
	client := NewClientSession(context.TODO(), true, protocol.DefaultIDHash)
	common.Must(client.EncodeRequestHeader(expectedRequest, buffer))

	buffer2 := buf.New()
	buffer2.Write(buffer.Bytes())

	sessionHistory := NewSessionHistory()
	defer common.Close(sessionHistory)

	userValidator := vmess.NewTimedUserValidator(protocol.DefaultIDHash)
	userValidator.Add(user)
	defer common.Close(userValidator)

	server := NewServerSession(userValidator, sessionHistory)
	actualRequest, err := server.DecodeRequestHeader(buffer, false)
	common.Must(err)

	if r := cmp.Diff(actualRequest, expectedRequest, cmp.AllowUnexported(protocol.ID{})); r != "" {
		t.Error(r)
	}
}
