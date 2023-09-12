package routes

import (
	"chenyucms/app/grpc/controllers"
	proto "github.com/goravel/example-proto"
	"github.com/goravel/framework/facades"
)

func Grpc() {
	proto.RegisterUserServiceServer(facades.Grpc().Server(), controllers.NewUserController())
}
