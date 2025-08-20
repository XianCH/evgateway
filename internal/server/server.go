package server

import (
	"fmt"
	"net"
	"time"

	"github.com/x14n/evgateway/internal/config"
	"github.com/x14n/evgateway/internal/gateway"
	"github.com/x14n/evgateway/internal/handlers"
	"github.com/x14n/evgateway/internal/protocol"
	"github.com/x14n/evgateway/utils"
	"github.com/x14n/evgateway/version"
)

type Server struct {
	Addr       string
	Gateway    *gateway.Gateway
	Dispatcher *gateway.Dispatcher
	Workerpool *WorkerPool
}

func NewServer(addr string, gw *gateway.Gateway, dispatcher *gateway.Dispatcher, wp *WorkerPool) *Server {
	return &Server{
		Addr:       addr,
		Gateway:    gw,
		Dispatcher: dispatcher,
		Workerpool: wp,
	}
}

func (s *Server) ListenAndServer() error {
	ln, err := net.Listen("tcp", s.Addr)
	if err != nil {
		fmt.Printf("tcp listen error: %v\n", err)
		return err
	}
	fmt.Println("TCPServer listen at :", s.Addr)

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Printf("tcp accept error %v", err)
			continue
		}
		fmt.Printf("New connection from %s\n", conn.RemoteAddr().String())

		session := &gateway.Session{
			ID:       "",
			Addr:     conn.RemoteAddr().String(),
			Conn:     conn,
			Lastseen: time.Now(),
		}

		s.Gateway.AddSession(session)

	}
}

func handleConnect(conn net.Conn, session *gateway.Session, srv *Server) {
	defer func() {
		conn.Close()
		srv.Gateway.RemoveSession(session.ID)
		fmt.Printf("Connection closed for session %s\n", session.ID)
	}()

	parser := protocol.NewParser(conn)
	parser.Start()
	defer parser.Stop()

	for {
		select {
		case frame, ok := <-parser.Frames():
			if !ok {
				fmt.Printf("may parser closed for session %s\n", session.ID)
				return
			}

			srv.Workerpool.Submit(func() {
				srv.Dispatcher.Dispatch(srv.Gateway, session, frame)
			})

		case err := <-parser.Errors():
			fmt.Printf("connect parser error %v", err)
		}
	}
}

func Run() {

	fmt.Printf("[gateway] starting EV Gateway v%s\n", version.Version)

	cfg := config.LoadConfig()

	gw := gateway.NewGateway()

	dispatcher := gateway.NewDispatcher()
	handlers.RegisterAllHandlers(dispatcher)

	wp := NewWorkerPool(cfg.WorkerPoolSize)
	wp.Start(cfg.WorkerPoolSize)
	defer wp.Stop()

	// 启动定时清理过期会话
	utils.StartSessionCleaner(gw, cfg.HeatbeatTTL)

	srv := NewServer(cfg.Addr, gw, dispatcher, wp)

	if err := srv.ListenAndServer(); err != nil {
		fmt.Printf("server error: %v\n", err)
	}
}
