package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"time"

	"github.com/keithkfield/grpc-go-course-tester/blog/blogpb"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var collection *mongo.Collection

type server struct {
}

type item struct {
	// ID       objectid.ObjectID
	AuthorID string `bson: "author_id"`
	Content  string `bson: "content"`
	Title    string `bson: "title"`
}

func main() {
	// If crash return file and line number
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	fmt.Println("RUnning blog service")
	port := "0.0.0.0:50051"
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to Listen: %v", err)
	}
	opt := []grpc.ServerOption{}
	s := grpc.NewServer(opt...)
	blogpb.RegisterBlogServiceServer(s, &server{})
	// Connectiing to Mongo
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb+srv://keithkfield:Icecat12!@goconnect-glenv.mongodb.net/test?retryWrites=true&w=majority"))
	if err != nil {
		log.Fatalf("Error connecting to Mongo Database: %v", err)
	}
	collection = client.Database("blogs").Collection("blog-post")
	// END

	go func() {
		fmt.Printf("Listening on port: %v", port)
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Failed to serve %v", err)
		}
		fmt.Print("working")
	}()
	// Wait for control c to exit
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	// Block untill signal recieved
	<-ch
	fmt.Println("Stopping Server")
	s.Stop()
	fmt.Println("CLoesing Listner")
	lis.Close()
	fmt.Println("Cloesing mongo COnnection")
	client.Disconnect(context.TODO())
	fmt.Println("Cloesing Program")

}

// GRPC FUNCTION
func (*server) CreateBlog(ctx context.Context, req *blogpb.CreateBlogRequest) (*blogpb.CreateBlogResponse, error) {
	blog := req.GetBlog()
	data := item{
		AuthorID: blog.GetAuthorId(),
		Title:    blog.GetTitle(),
		Content:  blog.GetContent(),
	}
	res, err := collection.InsertOne(context.Background(), data)
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Internal Error: %v", err))
	}
	fmt.Printf("Added: %v", res)

	return &blogpb.CreateBlogResponse{
		Blog: &blogpb.Blog{
			AuthorId: blog.GetAuthorId(),
			Title:    blog.GetTitle(),
			Content:  blog.GetContent(),
		},
	}, nil
}
