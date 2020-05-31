package testutil

import (
	"context"
	"fmt"
	"io/ioutil"
	"net"
	"sync"

	"github.com/phayes/freeport"
	"github.com/pkg/errors"
)

type MockServer struct {
	Address  string
	In       chan string
	listener net.Listener
	work     sync.WaitGroup
}

func RunMockServer(ctx context.Context, network string) (*MockServer, error) {
	port, err := freeport.GetFreePort()
	if err != nil {
		return nil, errors.Wrap(err, "find free port")
	}
	address := fmt.Sprintf("localhost:%d", port)
	listener, err := new(net.ListenConfig).Listen(ctx, network, address)
	if err != nil {
		return nil, err
	}
	mock := &MockServer{
		Address:  address,
		In:       make(chan string, 100),
		listener: listener,
	}
	mock.work.Add(1)
	go mock.listen()
	return mock, nil
}

func (m *MockServer) listen() {
	defer func() {
		close(m.In)
		m.work.Done()
	}()

	for {
		conn, err := m.listener.Accept()
		if err != nil {
			return
		}

		data, err := ioutil.ReadAll(conn)
		_ = conn.Close()
		if err != nil {
			continue
		}

		m.In <- string(data)
	}
}

func (m *MockServer) Close() {
	_ = m.listener.Close()
	m.work.Wait()
}
