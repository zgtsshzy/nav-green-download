package main

import (
	"context"
	"nav-green-download/pkg/conf"
	"nav-green-download/pkg/manage"
	"nav-green-download/pkg/server"
	"syscall"

	"github.com/sirupsen/logrus"
)

func init() {
	c := conf.New()
	c.Show()
}

func main() {
	smoc := server.NewSMOCDownloader()
	mfwam := server.NewMFWAMDownloader()

	manager := manage.New(
		"气象源数据处理",
		manage.Server(smoc, mfwam),
		manage.BeforeStart(BeforeStartFunc),
		manage.AfterStop(AfterStopFunc),
		manage.Signal(syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT),
	)

	defer func() {
		if err := manager.Stop(); err != nil {
			logrus.Errorf("停止所有的服务失败: %v", err)
		}
	}()

	if err := manager.Run(); err != nil {
		logrus.Errorf("启动所有的服务失败: %v", err)
	}
}

func BeforeStartFunc(ctx context.Context) error {

	return nil
}

func AfterStopFunc(ctx context.Context) error {

	return nil
}
