package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/keithkfield/grpc-go-course-tester/calculator/calculatorpb"
	"google.golang.org/grpc"
)

type server struct{}

func (*server) Sum(ctx context.Context, req *calculatorpb.SumRequest) (*calculatorpb.SumResponse, error) {
	fmt.Printf("Recieved: %v", req)
	firstNumber := req.FirstNumber
	secondNumber := req.SecondNumber
	sum := firstNumber + secondNumber
	res := &calculatorpb.SumResponse{
		SeunResult: sum,
	}
	return res, nil
}
func main() {
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to Listen: %v", err)
	}
	s := grpc.NewServer()
	calculatorpb.RegisterCalculatorServiceServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve %v", err)
	}
	fmt.Print("working")
}
