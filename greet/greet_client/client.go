package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/keithkfield/grpc-go-course-tester/greet/greetpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	// doCLientStreaming(c)
	// doBidiStreaming(c)
	doUniaryWithDeadiline(c, 5*time.Second)
	doUniaryWithDeadiline(c, 1*time.Second)

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

// BIDI Streaming
func doBidiStreaming(c greetpb.GreetServiceClient) {
	fmt.Print("Started bidi streaming")
	req := []*greetpb.GreetEveryoneRequest{
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Keith",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "doung",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Jack",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Bernie",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Cam",
			},
		},
	}
	stream, err := c.GreetEveryone(context.Background())
	if err != nil {
		log.Fatalf("Error Sending data: %V", err)
	}
	waitc := make(chan struct{})

	go func() {
		// Func to send alot of data
		for _, request := range req {
			fmt.Printf("Sending: %v\n", request)
			stream.Send(request)
			time.Sleep(1000 * time.Millisecond)
		}
		stream.CloseSend()
	}()
	go func() {
		// Recieve a bunch of messages
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Error with Recieving: %v\n", err)
				break
			}
			fmt.Printf("Recieved: %v", res.GetResult())
		}
	}()
	<-waitc
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
func doUniaryWithDeadiline(c greetpb.GreetServiceClient, seconds time.Duration) {
	fmt.Printf("Working")
	req := &greetpb.GreetWithDeadlineRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Keith",
			LastName:  "Kirtfield",
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), seconds)
	defer cancel()
	response, err := c.GreetWithDeadline(ctx, req)
	if err != nil {
		statusErr, ok := status.FromError(err)
		if ok {
			if statusErr.Code() == codes.DeadlineExceeded {
				fmt.Println("Time Exceeded")
			} else {
				fmt.Printf("Unexpected error: %v", statusErr)
			}
			return
		}

		log.Fatalf("Error calling Greet RPC: %v", err)
	}
	log.Printf("Response from Greet: %v", response.Response)
}
