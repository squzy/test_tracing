package main

import (
	"context"
	api "github.com/squzy/squzy_generated/generated/proto/v1"
	"github.com/squzy/squzy_go/core"
	sGrpc "github.com/squzy/squzy_go/integrations/grpc"
	sHttp "github.com/squzy/squzy_go/integrations/http"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
	service "test_grpc/generated"
)

type server struct {
	application *core.Application
}

func (s *server) Echo(ctx context.Context, msg *service.EchoMsg) (*service.EchoMsg, error) {
	trx := core.GetTransactionFromContext(ctx).CreateTransaction("calcaulate time", api.TransactionType_TRANSACTION_TYPE_INTERNAL, nil)

	client := http.Client{Transport: sHttp.NewRoundTripper(s.application, nil)}

	req, err := http.NewRequest("GET", "https://api.exchangeratesapi.io/latest?base=USD", nil)

	if err != nil {
		return nil, err
	}
	_, err = client.Do(sHttp.NewRequest(trx, req))

	trx.End(nil)
	return &service.EchoMsg{}, nil
}

func main() {
	squzy, err := core.CreateApplication(nil, &core.Options{
		ApiHost:         "http://localhost:8080",
		ApplicationName: "Go app test grpc",
	})

	if err != nil {
		log.Fatal(err)
	}
	lis, err := net.Listen("tcp", ":7879")
	if err != nil {
		log.Fatal(err)
	}

	s := grpc.NewServer(grpc.ChainUnaryInterceptor(sGrpc.NewServerUnaryInterceptor(squzy)))

	service.RegisterEchoServiceServer(s, &server{
		application: squzy,
	})

	log.Fatal(s.Serve(lis))
}
