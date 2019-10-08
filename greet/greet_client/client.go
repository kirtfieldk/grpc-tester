package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

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
	// doUniary(c)
	// doServerStreaming(c)
	doCLientStreaming(c)

}

func doServerStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Starting Stream")

	req := &greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Keith",
			LastName:  "Kirtfield",
		},
	}
	stream, err := c.GreetManyTimes(context.Background(), req)
	if err != nil {
		log.Fatalf("Error at : %v", err)
	}
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			// End Of Stream
			break
		}
		if err != nil {
			log.Fatalf("Error while reading stream : %v", err)
		}
		log.Printf("Message: %v", msg.GetResult())
	}
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

func doCLientStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Starting client streaming")

	req := []*greetpb.LongGreetRequest{
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Keith",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "doung",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Jack",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Bernie",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Cam",
			},
		},
	}

	stream, err := c.LongGreet(context.Background())
	if err != nil {
		log.Fatalf("Error reading stream: %v", err)
	}
	// Iterate and send each interval
	for _, re := range req {
		fmt.Printf("Sending: %v", re)
		stream.Send(re)
		time.Sleep(1000 * time.Millisecond)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error while recienving resopnse: %v", err)
	}
	fmt.Printf("LongGreet Response: %v\n", res)

}
