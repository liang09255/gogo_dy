package ggRPC

import (
	"common/ggConfig"
	"common/ggDiscovery"
	"common/ggIDL/chat"
	"common/ggIDL/relation"
	"common/ggIDL/user"
	"common/ggIDL/video"
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

func GetVideoClient() video.VideoServiceClient {
	conn := initClient(ggConfig.Config.VideoServer.Name)
	return video.NewVideoServiceClient(conn)
}

func GetRelationClient() relation.RelationClient {
	conn := initClient(ggConfig.Config.UserServer.Name)
	return relation.NewRelationClient(conn)
}

func GetChatClient() chat.ChatClient {
	conn := initClient(ggConfig.Config.ChatServer.Name)
	return chat.NewChatClient(conn)
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
