package godownloader

import (
	context "context"
	"fmt"
	"log"
	"time"

	"github.com/asim/go-micro/plugins/server/grpc/v4"
	"github.com/geebytes/go-downloader/pb"
	"go-micro.dev/v4"
	"go-micro.dev/v4/registry"
)

var downloader *Downloader = NewDownloader("/data/work/download")

type ServerInterface interface {
	Callback(src string, err error) (string, error)
	// Download(ctx context.Context, in *DownloadRequest) (*DownloadResponse, error)
	pb.DownloaderHandler
}

type DefaultServer struct {
	// pb.UnimplementedDownloaderServer
}

func (s *DefaultServer) Callback(src string, err error) (string, error) {
	return src, err
}
func (s *DefaultServer) Download(ctx context.Context, in *pb.DownloadRequest, out *pb.DownloadResponse) error {
	dst, err := downloader.Download(in.Request, in.FileName)
	dst, err = s.Callback(dst, err)
	if err != nil {
		fmt.Println(err.Error())
		out.Dst = ""
		out.Msg = err.Error()
		out.Status = 500

		return err

	}
	out.Dst = dst
	out.Status = 200
	out.Msg = "success"
	return nil

}

func StartServer(server ServerInterface) {
	grpcServer := grpc.NewServer()
	service := micro.NewService(
		micro.Server(grpcServer),
		micro.Name("downloader.service"),
		micro.Address("0.0.0.0:9527"),
		micro.RegisterTTL(time.Second*30),
		micro.RegisterInterval(time.Second*10),
		micro.Registry(registry.NewRegistry(registry.Addrs("127.0.0.1:2379"))),
	)

	// optionally setup command line usage
	service.Init()

	// Register Handlers
	pb.RegisterDownloaderHandler(service.Server(), server)

	// Run server
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
