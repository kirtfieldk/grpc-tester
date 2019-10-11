package main

import (
	"context"
	"fmt"
	"log"

	"github.com/keithkfield/grpc-go-course-tester/blog/blogpb"
	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Starting blog client")
	opts := grpc.WithInsecure()
	cc, err := grpc.Dial("localHost:50051", opts)
	if err != nil {
		log.Fatalf("Error trying to connect with server: %v", err)
	}
	defer cc.Close()
	c := blogpb.NewBlogServiceClient(cc)
	// Createing blog
	// blog := &blogpb.Blog{
	// 	AuthorId: "Keith Kirtfield",
	// 	Title:    "Miser and a tree",
	// 	Content:  "A simple book",
	// }
	// res, err := c.CreateBlog(context.Background(), &blogpb.CreateBlogRequest{Blog: blog})
	// if err != nil {
	// 	log.Fatalf("Error adding blog: %v", err)
	// }
	// fmt.Printf("Added blog: %v", res)

	// read Blog
	fmt.Println("Envoke Read blog Func")

	_, err2 := c.ReadBlog(context.Background(), &blogpb.ReadBlogRequest{BlogId: "weifhewfhefoijewfe"})
	if err2 != nil {
		fmt.Printf("Error: %v", err2)
	}
	readBlog, err3 := c.ReadBlog(context.Background(), &blogpb.ReadBlogRequest{BlogId: "5d9fd5110b4e6da8a2f6dbef"})
	if err3 != nil {
		fmt.Printf("Error: %v", err3)
	}

	fmt.Printf("Blog retived: %v", readBlog)
}
