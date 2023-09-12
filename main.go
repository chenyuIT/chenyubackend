// Code generated by hertz generator.

package main

import (
	"chenyucms/bootstrap"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/goravel/framework/facades"
	"github.com/hertz-contrib/logger/accesslog"
)

func main() {
	//程序异常恢复
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("程序出现错误:", r)
		}
	}()
	//加载系统功能
	bootstrap.Boot()

	// Start GRPC server
	go func() {
		if err := facades.Grpc().Run(); err != nil {
			facades.Log().Errorf("Run grpc error: %+v", err)
		}
	}()

	// Start HTTP service with Hertz
	h := server.Default(
		server.WithHostPorts(facades.Config().GetString("APP_HOST") + ":" + facades.Config().GetString("APP_PORT")),
	)

	h.Name = "ChenyuCMS"
	h.Use(accesslog.New())
	register(h)
	h.Spin()
}
