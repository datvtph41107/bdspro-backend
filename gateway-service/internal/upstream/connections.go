package upstream

import (
	"common/rpc"
	"errors"
	"fmt"
	gatewayconfig "gateway/config"
	"gateway/internal/tqdtransport"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const backoffMaxDelay = 5 * time.Second

// Connections is the explicit Gateway-owned outbound topology. Shared RPC
// owns channel mechanics; Gateway owns which services exist and when they close.
type Connections struct {
	Notification *grpc.ClientConn
	Payment      *grpc.ClientConn
	BDSPro       *grpc.ClientConn
	CRM          *grpc.ClientConn
	Organization *grpc.ClientConn
	Social       *grpc.ClientConn
	User         *grpc.ClientConn
	Chat         *grpc.ClientConn
	Auth         *grpc.ClientConn
	Assistant    *grpc.ClientConn
	Hub          *grpc.ClientConn
	TQD          *grpc.ClientConn
}

func Open(cfg gatewayconfig.Runtime) (*Connections, error) {
	connections := &Connections{}
	opened := []*grpc.ClientConn{}
	open := func(name string, callOptions ...grpc.CallOption) (*grpc.ClientConn, error) {
		target, err := cfg.Endpoint(name)
		if err != nil {
			return nil, err
		}

		conn, err := rpc.NewClient(rpc.ClientConfig{
			Target:             target,
			BackoffMaxDelay:    backoffMaxDelay,
			Credentials:        insecure.NewCredentials(),
			Transport:          cfg.RPCTransport,
			DefaultCallOptions: callOptions,
		})
		if err != nil {
			return nil, fmt.Errorf("open %s gRPC connection: %w", name, err)
		}

		opened = append(opened, conn)
		return conn, nil
	}

	var err error
	if connections.Notification, err = open("notification"); err != nil {
		closeAll(opened)
		return nil, err
	}
	if connections.Payment, err = open("payment"); err != nil {
		closeAll(opened)
		return nil, err
	}
	if connections.BDSPro, err = open("bdspro"); err != nil {
		closeAll(opened)
		return nil, err
	}
	if connections.CRM, err = open("crm"); err != nil {
		closeAll(opened)
		return nil, err
	}
	if connections.Organization, err = open("organization"); err != nil {
		closeAll(opened)
		return nil, err
	}
	if connections.Social, err = open("social"); err != nil {
		closeAll(opened)
		return nil, err
	}
	if connections.User, err = open("user"); err != nil {
		closeAll(opened)
		return nil, err
	}
	if connections.Chat, err = open("chat"); err != nil {
		closeAll(opened)
		return nil, err
	}
	if connections.Auth, err = open("auth"); err != nil {
		closeAll(opened)
		return nil, err
	}
	if connections.Assistant, err = open("assistant"); err != nil {
		closeAll(opened)
		return nil, err
	}
	if connections.Hub, err = open("hub"); err != nil {
		closeAll(opened)
		return nil, err
	}
	if connections.TQD, err = open("tqd",
		grpc.MaxCallRecvMsgSize(tqdtransport.MaxGRPCMessageBytes),
		grpc.MaxCallSendMsgSize(tqdtransport.MaxGRPCMessageBytes),
	); err != nil {
		closeAll(opened)
		return nil, err
	}
	return connections, nil
}

// LIFO — Last In, First Out. Reverse order khả năng phụ thuộc vào resource được tạo trước.
func (c *Connections) Close() error {
	if c == nil {
		return nil
	}
	return errors.Join(
		closeConn(c.TQD), closeConn(c.Hub), closeConn(c.Assistant),
		closeConn(c.Auth), closeConn(c.Chat), closeConn(c.User),
		closeConn(c.Social), closeConn(c.Organization), closeConn(c.CRM),
		closeConn(c.BDSPro), closeConn(c.Payment), closeConn(c.Notification),
	)
}

func closeConn(conn *grpc.ClientConn) error {
	if conn == nil {
		return nil
	}
	return conn.Close()
}

func closeAll(conns []*grpc.ClientConn) {
	for i := len(conns) - 1; i >= 0; i-- {
		_ = conns[i].Close()
	}
}
