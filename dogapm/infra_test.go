package dogapm

import (
	"context"
	"io"
	"net/http"
	"protos"
	"testing"
	"time"
)

func TestInit(t *testing.T){
	Infra.Init(
		InfraDbOption("root:password@tcp(localhost:3307)/ordersvc"),
		InfraRdbOption("localhost:6380"),
	)
}

func TestNewHttpServer(t *testing.T) {
	server := NewHttpServer(":8080")
	server.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	if server == nil {
		t.Fatal("Expected NewHttpServer to return a non-nil server")
	}
	if err := server.Start(); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := server.Close(); err != nil {
			t.Fatal(err)
		}
	})

	var lastErr error
	for i := 0; i < 20; i++ {
		resp, err := http.Get("http://127.0.0.1:8080/test")
		if err != nil {
			lastErr = err
			time.Sleep(50 * time.Millisecond)
			continue
		}

		body, err := io.ReadAll(resp.Body)
		_ = resp.Body.Close()
		if err != nil {
			t.Fatal(err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("unexpected status code: got %d", resp.StatusCode)
		}
		if string(body) != "OK" {
			t.Fatalf("unexpected response body: got %q", string(body))
		}
		return
	}

	t.Fatalf("server did not become ready: %v", lastErr)
}

type helloSvc struct{
	protos.UnimplementedHelloServiceServer
}

func (h *helloSvc) Receive(ctx context.Context, in *protos.HelloMsg) (*protos.HelloMsg, error) {
	return &protos.HelloMsg{Msg: "Hello, " + in.Msg + "!"}, nil
}

func TestNewGrpcServer(t *testing.T) {
	server := NewGrpcServer(":50051")
	protos.RegisterHelloServiceServer(server.Server, &helloSvc{})
	if err := server.Start(); err != nil {
		t.Fatal(err)
	}
	defer server.Stop()

	var lastErr error
	for i := 0; i < 20; i++ {
		client, err := NewGrpcClient("localhost:50051")
		if err != nil {
			lastErr = err
			time.Sleep(50 * time.Millisecond)
			continue
		}

		res, err := protos.NewHelloServiceClient(client).Receive(context.TODO(), &protos.HelloMsg{Msg: "world"})
		if err != nil {
			_ = client.Close()
			lastErr = err
			time.Sleep(50 * time.Millisecond)
			continue
		}
		_ = client.Close()
		if res.Msg != "Hello, world!" {
			t.Fatalf("unexpected response: got %q", res.Msg)
		}
		t.Log(res.Msg)
		return
	}

	if lastErr != nil {
		t.Fatalf("server did not become ready: %v", lastErr)
	}
	t.Fatal("server did not become ready")
}