package main

import (
	"context"
	"errors"
	"log"
	"net"
	"os"
	"strings"

	users "github.com/grpc/practical_go/user-service/service"
	"google.golang.org/grpc"
)

type userService struct {
	users.UnimplementedUsersServer
}

func (s userService) GetUser(ctx context.Context, in *users.UserGetRequest) (*users.UserGetReply, error) {
	log.Printf("Received request for user with Email: %s Id: %s\n", in.Email, in.Id)
	components := strings.Split(in.Email, "@")
	log.Println(components)
	if len(components) != 2 {
		return nil, errors.New("invalid email address")
	}
	u := users.User{
		Id:        in.Id,
		FirstName: components[0],

		LastName: components[1],
		Age:      36,
	}
	return &users.UserGetReply{User: &u}, nil
}

func registerService(s *grpc.Server) {
	users.RegisterUsersServer(s, &userService{})
}

func startServer(s *grpc.Server, l net.Listener) error {
	return s.Serve(l)
}
func main() {
	listenAddr := os.Getenv("LISTEN_ADDR")
	if len(listenAddr) == 0 {
		listenAddr = ":50051"
	}
	lis, err := net.Listen("tcp", listenAddr)
	if err != nil {
		log.Fatal(err)
	}
	s := grpc.NewServer()
	registerService((s))
	log.Fatal(startServer(s, lis))
}
