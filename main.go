package main

import (
	"flag"
	"github.com/tal-tech/go-zero/core/conf"
	"github.com/tal-tech/go-zero/core/service"
	"github.com/tal-tech/go-zero/core/threading"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"job/internal/config"
	"job/internal/handler"
	"job/internal/svc"
)

/**
* @Description TODO
* @Version 1.0
**/
var configFile = flag.String("f", "/Users/seven/Developer/goenv/未命名文件夹/job/etc/job.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	group := service.NewServiceGroup()
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	threading.GoSafe(func() {
		for s := range ch {
			switch s {
			case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
				group.Stop()
			}
		}
	})

	handler.RegisterJob(ctx,group)

	//阻塞直至有信号传入
	s := <-ch
	fmt.Println("退出job..", s)
}