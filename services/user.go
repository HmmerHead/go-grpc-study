package services

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/HmmerHead/go-grpc/pb"
)

// type UserServiceServer interface {
// 	AddUser(context.Context, *User) (*User, error)
//  AddUserVerbose(ctx context.Context, in *User, opts ...grpc.CallOption) (UserService_AddUserVerboseClient, error)
// 	AddUsers(ctx context.Context, opts ...grpc.CallOption) (UserService_AddUsersClient, error)
//	AddUserStreamBoth(ctx context.Context, opts ...grpc.CallOption) (UserService_AddUserStreamBothClient, error)
//  mustEmbedUnimplementedUserServiceServer()
// }

type UserService struct {
	pb.UnimplementedUserServiceServer
}

func NewUserService() *UserService {
	return &UserService{}
}

func (*UserService) AddUser(context context.Context, request *pb.User) (*pb.User, error) {

	return &pb.User{
		Id:    "121",
		Name:  request.GetName(),
		Email: request.GetEmail(),
	}, nil
}

func (*UserService) AddUserVerbose(request *pb.User, stream pb.UserService_AddUserVerboseServer) error {

	stream.Send(&pb.UserResultStream{
		Status: "Init",
		User:   &pb.User{},
	})

	time.Sleep(time.Second * 3)

	stream.Send(&pb.UserResultStream{
		Status: "Inserting",
		User:   &pb.User{},
	})

	time.Sleep(time.Second * 3)

	stream.Send(&pb.UserResultStream{
		Status: "User has been nserte",
		User: &pb.User{
			Id:    "123",
			Name:  request.GetName(),
			Email: request.GetEmail(),
		},
	})

	time.Sleep(time.Second * 3)

	stream.Send(&pb.UserResultStream{
		Status: "Completed",
		User: &pb.User{
			Id:    "123",
			Name:  request.GetName(),
			Email: request.GetEmail(),
		},
	})

	time.Sleep(time.Second * 3)

	return nil
}

func (*UserService) AddUsers(stream pb.UserService_AddUsersServer) error {

	users := []*pb.User{}

	for {

		request, err := stream.Recv()

		if err == io.EOF {
			return stream.SendAndClose(&pb.Users{
				User: users,
			})
		}

		if err != nil {
			log.Fatalf("Error receiving stream: %v", err)
		}

		users = append(users, &pb.User{
			Id:    request.GetId(),
			Name:  request.GetName(),
			Email: request.GetEmail(),
		})

		fmt.Println("Adding ", request.GetName())
	}

}

func (*UserService) AddUserStreamBoth(stream pb.UserService_AddUserStreamBothServer) error {

	for {
		request, err := stream.Recv()

		if err == io.EOF {
			return nil
		}

		if err != nil {
			log.Fatalf("Error receving stream from client: %v", err)
		}

		err = stream.Send(&pb.UserResultStream{
			Status: "Added",
			User:   request,
		})

		if err != nil {
			log.Fatalf("Error senng stream to the client: %v", err)
		}


	}

}
