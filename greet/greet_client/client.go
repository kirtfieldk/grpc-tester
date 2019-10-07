package main

import (
	"context"
	"fmt"
	"log"

	"github.com/keithkfield/grpc-go-course-tester/greet/greetpb"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect %v", err)
	}
	defer conn.Close()
	c := greetpb.NewGreetServiceClient(conn)
	doUniary(c)

}

func doUniary(c greetpb.GreetServiceClient) {
	fmt.Printf("Working")
	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Keith",
			LastName:  "Kirtfield",
		},
	}
	response, err := c.Greet(context.Background(), req)
	if err != nil {
		log.Fatalf("Error calling Greet RPC: %v", err)
	}
	log.Printf("Response from Greet: %v", response.Result)
}
