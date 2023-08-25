gen: gen_user gen_video

gen_user:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative  ./user/user.proto

gen_video:
	protoc --proto_path=./ --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative  ./video/video.proto