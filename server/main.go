package main

import (
	"context"
	"log"
	"net"

	pb "github.com/toyozaki/grpc_sample/gen"
	"google.golang.org/grpc"
)

var address = "localhost:8080"

type server struct {
	pb.UnimplementedGreetServiceServer
}

func (s *server) UnaryGreet(ctx context.Context, req *pb.UnaryGreetRequest) (*pb.UnaryGreetReply, error) {
	log.Printf("Received: %v", req.GetName())
	return &pb.UnaryGreetReply{Message: "Hello " + req.GetName()}, nil
}

func (s *server) ClientStreamGreet(stream pb.GreetService_ClientStreamGreetServer) error {
	count := 0
	for {
		req, err := stream.Recv()
		if err != nil {
			if err.Error() == "EOF" {
				return nil
			}

			return err
		}
		log.Printf("Received: %v", req.GetName())

		if count == 10 {
			if err := stream.SendAndClose(&pb.ClientStreamGreetReply{Message: "Goodby " + req.GetName()}); err != nil {
				return err
			}
		}

		count++
	}
}

func (s *server) ServerStreamGreet(req *pb.ServerStreamGreetRequest, stream pb.GreetService_ServerStreamGreetServer) error {
	for i := 0; i < 5; i++ {
		if err := stream.Send(&pb.ServerStreamGreetReply{Message: "Hello " + req.GetName()}); err != nil {
			return err
		}
	}
	return nil
}

func (s *server) BidirectionalStreamGreet(stream pb.GreetService_BidirectionalStreamGreetServer) error {
	count := 0
	for {
		req, err := stream.Recv()
		if err != nil {
			return err
		}
		log.Printf("Received: %v", req.GetName())

		if count > 5 {
			if err := stream.Send(&pb.BidirectionalStreamGreetReply{Message: "Goodby " + req.GetName()}); err != nil {
				return err
			}
			return nil
		} else {
			if err := stream.Send(&pb.BidirectionalStreamGreetReply{Message: "Hello " + req.GetName()}); err != nil {
				return err
			}
		}

		count++
	}
}

func main() {
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal("Failed to listen: ", err)
	}

	log.Println("Launch server on ", address)
	s := grpc.NewServer()

	pb.RegisterGreetServiceServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatal("Failed to serve: ", err)
	}
}
