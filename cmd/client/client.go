package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/HmmerHead/go-grpc/pb"
	"google.golang.org/grpc"
)

func main() {
	connection, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Could not connect to GRPC Server: %v", err)
	}

	defer connection.Close()

	cliente := pb.NewUserServiceClient(connection)

	//AddUser(cliente)
	//AddUserVerbose(cliente)
	// AddUsers(cliente)
	AddUserStreamBoth(cliente)
}

func AddUser(client pb.UserServiceClient) {
	request := &pb.User{
		Id:    "69",
		Name:  "Ze",
		Email: "ze@ze.com",
	}

	response, err := client.AddUser(context.Background(),
		request)

	if err != nil {
		log.Fatalf("Could not make GRPC Request: %v", err)
	}

	fmt.Println(response)
}

func AddUserVerbose(client pb.UserServiceClient) {
	request := &pb.User{
		Id:    "69",
		Name:  "Ze",
		Email: "ze@ze.com",
	}

	response, err := client.AddUserVerbose(context.Background(),
		request)

	if err != nil {
		log.Fatalf("Could not make GRPC Request: %v", err)
	}

	for {
		stream, err := response.Recv()

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatalf("Could not receive the message: %v", err)
		}

		fmt.Println("Status: ", stream.Status)
	}
}

func AddUsers(client pb.UserServiceClient) {
	request := []*pb.User{
		&pb.User{
			Id:    "g1",
			Name:  "gg",
			Email: "G@g.com",
		},
		&pb.User{
			Id:    "g2",
			Name:  "gg2",
			Email: "G2@g.com",
		},
		&pb.User{
			Id:    "g3",
			Name:  "gg3",
			Email: "G3@g.com",
		},
		&pb.User{
			Id:    "g4",
			Name:  "gg4",
			Email: "G4@g.com",
		},
		&pb.User{
			Id:    "g9",
			Name:  "gg5",
			Email: "G5@g.com",
		},
		&pb.User{
			Id:    "g8",
			Name:  "gg6",
			Email: "G6@g.com",
		},
		&pb.User{
			Id:    "g8",
			Name:  "gg7",
			Email: "G7@g.com",
		},
	}

	stream, err := client.AddUsers(context.Background())

	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	for _, request := range request {
		stream.Send(request)
		time.Sleep(time.Second * 2)
	}

	response, err := stream.CloseAndRecv()

	if err != nil {
		log.Fatalf("Error receving response: %v", err)
	}

	fmt.Println(response)
}

func AddUserStreamBoth(client pb.UserServiceClient) {

	request := []*pb.User{
		&pb.User{
			Id:    "g1",
			Name:  "gg",
			Email: "G@g.com",
		},
		&pb.User{
			Id:    "g2",
			Name:  "gg2",
			Email: "G2@g.com",
		},
		&pb.User{
			Id:    "g3",
			Name:  "gg3",
			Email: "G3@g.com",
		},
		&pb.User{
			Id:    "g4",
			Name:  "gg4",
			Email: "G4@g.com",
		},
		&pb.User{
			Id:    "g9",
			Name:  "gg5",
			Email: "G5@g.com",
		},
		&pb.User{
			Id:    "g8",
			Name:  "gg6",
			Email: "G6@g.com",
		},
		&pb.User{
			Id:    "g8",
			Name:  "gg7",
			Email: "G7@g.com",
		},
	}

	stream, err := client.AddUserStreamBoth(context.Background())

	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	wait := make(chan int)

	go func() {
		for _, requests := range request {
			fmt.Println("Sending user: ", requests.Name)
			stream.Send(requests)
			time.Sleep(time.Second * 2)
		}
		stream.CloseSend()
	}()

	go func() {
		for {
			response, err := stream.Recv()

			if err != io.EOF {
				break
			}

			if err != nil {
				log.Fatalf("Error creating request: %v", err)
				break
			}

			fmt.Printf("Recebendo user %v com status: %v", response.GetUser().GetName(), response.GetStatus())
		}

		close(wait)
	}()

	<-wait
}
