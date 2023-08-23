package ggRPC

import (
	"common/ggConfig"
	"common/ggDiscovery"
	"common/ggIDL/relation"
	"common/ggIDL/user"
	"common/ggLog"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"
)

func GetUserClient() user.UserClient {
	conn := initClient(ggConfig.Config.UserServer.Name)
	return user.NewUserClient(conn)
}

func GetRelationClient() relation.RelationClient {
	conn := initClient(ggConfig.Config.UserServer.Name)
	return relation.NewRelationClient(conn)
}

func GetVideoClient() {

}

func GetChatClient() {

}

func initClient(name string) *grpc.ClientConn {
	etcdRegister := ggDiscovery.NewResolver(ggConfig.Config.Etcd.Addrs)
	resolver.Register(etcdRegister)

	target := fmt.Sprintf("etcd:///%s", name)

	conn, err := grpc.Dial(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		ggLog.Fatal(err)
	}
	return conn
}
