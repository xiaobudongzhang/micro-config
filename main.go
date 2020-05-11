package main

import (
	"context"
	"net"
	"sync"

	"github.com/micro/go-micro/v2/config"
	"github.com/micro/go-micro/v2/config/source/file"
	"github.com/micro/go-micro/v2/util/log"
	proto "github.com/micro/go-plugins/config/source/grpc/proto"
	grpc2 "google.golang.org/grpc"
)

var (
	mux        sync.RWMutex
	configMaps = make(map[string]*proto.ChangeSet)
	apps       = []string("micro")
)

type Service struct{}

func main() {
	defer func() {
		if r := recover(); r != nil {
			log.Logf("recover %v", r)
		}
	}()

	err := loadAndWatchConfigFile()
	if err != nil {
		log.Fatal(err)
	}

	service := grpc2.NewServer()
	proto.RegisterSourceServer(service, new(Service))
	ts, err := net.Listen("tcp", ":9600")
	if err != nil {
		log.Fatal(err)
	}

	log.Logf("configServer started")

	err = service.Server(ts)
	if err != nil {
		log.Fatal(err)
	}
}

func (s Service) Read(ctx context.Context, req *proto.ReadRequest) (rsp *proto.ReadResponse, err error) {
	appName := parsePath(req.Path)

	rsp = &proto.ReadResponse{
		ChangeSet: getConfig(appName)
	}
	return
}

func (s Service) Watch(req *proto.WatchRequest, server proto.Source_WatchServer) (err error) {
	appName := parsePath(req.Path)
	rsp := &proto.WatchReponse{
		ChangeSet:getConfig(appName)
	}

	if err = server.Send(rsp); err != nil {
		log.Logf("侦听处理异常 %s", err)
		return err
	}

	return
}


func loadAndWatchConfigFile() (err error) {

	for _,app := range apps {
		if err := config.Load(file.NewSource(
			file.WithPath("./conf" + app + ".yml")
		)); err != nil {
			log.Fatal("加载应用配置文件异常 %s", err)
			return err
		}
	}



	watcher, err := config.Watch()

	

}