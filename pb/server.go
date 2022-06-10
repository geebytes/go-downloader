package pb

import (
	context "context"
	"fmt"
	"net"

	grpc "google.golang.org/grpc"
)

var downloader *Downloader = NewDownloader("/data/work/download")

type ServerInterface interface {
	Callback(src string, err error) (string, error)
	// Download(ctx context.Context, in *DownloadRequest) (*DownloadResponse, error)
	DownloaderServer
}

type DefaultServer struct {
	UnimplementedDownloaderServer
}

func (s *DefaultServer) Callback(src string, err error) (string, error) {
	return src, err
}
func (s *DefaultServer) Download(ctx context.Context, in *DownloadRequest) (*DownloadResponse, error) {
	dst, err := downloader.Download(in.Request, in.FileName)
	dst, err = s.Callback(dst, err)
	if err != nil {
		fmt.Println(err.Error())
		return &DownloadResponse{Dst: dst, Status: 500, Msg: err.Error()}, err

	}
	return &DownloadResponse{Dst: dst, Status: 200}, nil

}

func StartServer(server ServerInterface) {
	lis, err := net.Listen("tcp", ":9527")
	if err != nil {
		fmt.Printf("failed to listen: %v", err)
		return
	}
	s := grpc.NewServer()               // 创建gRPC服务器
	RegisterDownloaderServer(s, server) // 在gRPC服务端注册服务
	err = s.Serve(lis)
	if err != nil {
		panic(err)
	}
}
