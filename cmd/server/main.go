package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"time"

	hellopb "grpc_study/pkg/grpc"

	// "google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	// "google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
	// "google.golang.org/grpc/status"
)

type myServer struct {
	hellopb.UnimplementedGreetingServiceServer
}

// Unary RPC の通信終了時
func (s *myServer) Hello(ctx context.Context, req *hellopb.HelloRequest) (*hellopb.HelloResponse, error) {
	// (なにか処理をしてエラーが発生した)
	// stat := status.New(codes.Unknown, "unknown error occurred")
	// stat, _ = stat.WithDetails(&errdetails.DebugInfo{
	// 	Detail: "detail reason of err",
	// })
	// err := stat.Err()
	// return nil, err
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		log.Println(md)
	}
	headMD := metadata.New(map[string]string{"type": "unary", "from": "server", "in": "header"})
	if err := grpc.SetHeader(ctx, headMD); err != nil {
		return nil, err
	}
	trailerMD := metadata.New(map[string]string{"type": "unary", "from": "server", "in": "trailer"})
	if err := grpc.SetTrailer(ctx, trailerMD); err != nil {
		return nil, err
	}
	return &hellopb.HelloResponse{
		Message: fmt.Sprintf("Hello, %s!", req.GetName()),
	}, nil
}

// Server stream RPC の通信終了時
func (s *myServer) HelloServerStream(req *hellopb.HelloRequest, stream hellopb.GreetingService_HelloServerStreamServer) error {
	resCount := 5
	for i := 0; i < resCount; i++ {
		if err := stream.Send(&hellopb.HelloResponse{
			Message: fmt.Sprintf("[%d] Hello, %s", i, req.GetName()),
		}); err != nil {
			return err
		}
		time.Sleep(time.Second * 1)
	}
	return nil
}

// Client Stream RPCがリクエストを受け取るところ
func (s *myServer) HelloClientStream(stream hellopb.GreetingService_HelloClientStreamServer) error {
	nameList := make([]string, 0)
	for {
		// stream のRecvメソッドを呼び出してリクエスト内容を取得する
		req, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			message := fmt.Sprintf("Hello, %v!", nameList)
			return stream.SendAndClose(&hellopb.HelloResponse{
				Message: message,
			})
		}
		if err != nil {
			return err
		}
		nameList = append(nameList, req.GetName())
	}
}

// 双方向ストリーミングの場合
func (s *myServer) HelloBiStreams(stream hellopb.GreetingService_HelloBiStreamsServer) error {
	if md, ok := metadata.FromIncomingContext(stream.Context()); ok {
		log.Println(md)
	}
	// (パターン1)すぐにヘッダーを送信したいならばこちら
	headerMD := metadata.New(map[string]string{"type": "stream", "from": "server", "in": "header"})
	if err := stream.SendHeader(headerMD); err != nil {
		return err
	}
	// (パターン2)本来ヘッダーを送るタイミングで送りたいならばこちら
	if err := stream.SetHeader(headerMD); err != nil {
		return err
	}

	trailerMD := metadata.New(map[string]string{"type": "stream", "from": "server", "in": "trailer"})
	stream.SetTrailer(trailerMD)

	for {
		// リクエスト受信
		req, err := stream.Recv()
		// 得られたエラーがio.EOFならばもうリクエストは送られてこない
		if errors.Is(err, io.EOF) {
			return nil
		}
		if err != nil {
			return err
		}
		message := fmt.Sprintf("Hello, %v!", req.GetName())
		if err := stream.Send(&hellopb.HelloResponse{
			Message: message,
		}); err != nil {
			return err
		}
	}
}

// 自作サービス構造体のコンストラクタを定義
func NewMyServer() *myServer {
	return &myServer{}
}

func main() {
	// 8080番ポートのListenerを作成
	port := 8080
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		panic(err)
	}

	// gRPCサーバーを作成
	s := grpc.NewServer(
	// grpc.StreamInterceptor(myStreamServerInterceptor1),
	// grpc.ChainUnaryInterceptor(
	// 	myUnaryServerInterceptor1,
	// 	myUnaryServerInterceptor2,
	// ),

	// grpc.ChainStreamInterceptor(
	// 	myStreamServerInterceptor1,
	// 	myStreamServerInterceptor2,
	// ),
	)

	// gRPCサーバーにGreetingServiceを登録
	hellopb.RegisterGreetingServiceServer(s, NewMyServer())

	// サーバーリフレクションの設定
	reflection.Register(s)

	//作成したgRPCサーバーを、8080番ポートで稼働させる
	go func() {
		log.Printf("start gRPC server port: %v", port)
		s.Serve(listener)
	}()

	// Ctrl+C が入力されたらGraceful shutdown されるようにする
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("stopping gRPC server ...")
	s.GracefulStop()
}
