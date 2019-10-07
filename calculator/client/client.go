package main

import (
	"context"
	"fmt"
	"log"

	"github.com/keithkfield/grpc-go-course-tester/calculator/calculatorpb"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect %v", err)
	}
	defer conn.Close()
	c := calculatorpb.NewCalculatorServiceClient(conn)
	doUniary(c)

}

func doUniary(c calculatorpb.CalculatorServiceClient) {
	fmt.Printf("Working")
	req := &calculatorpb.SumRequest{

		FirstNumber:  9,
		SecondNumber: 13,
	}
	response, err := c.Sum(context.Background(), req)
	if err != nil {
		log.Fatalf("Error calling Greet RPC: %v", err)
	}
	log.Printf("Response from Sum: %v", response.SeunResult)
}
