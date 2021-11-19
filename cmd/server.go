/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/spf13/cobra"

	pb "github.com/py-go/grpc-poc/pkg/filesvc"

	"google.golang.org/grpc"
)

const (
	port = ":9000"
)

// server is used to implement filesvc.FileSVCServer.
type Server struct {
	pb.UnimplementedFileSVCServer
}

// server is used to implement filesvc.FileSVCServer.
type ServerByte struct {
	pb.UnimplementedFileSVCServer
}

type APIResponse struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Starts the Schema gRPC server",

	Run: func(cmd *cobra.Command, args []string) {
		lis, err := net.Listen("tcp", port)
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}

		grpcServer := grpc.NewServer()

		// Register services
		pb.RegisterFileSVCServer(grpcServer, &Server{})

		log.Printf("GRPC server listening on %v", lis.Addr())

		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}

// GetFileSVC implements filesvc.FileSVCServer
func (s *Server) GetFileSVC(ctx context.Context, req *pb.FileSVCRequest) (res *pb.FileSVCReply, err error) {
	res = &pb.FileSVCReply{}

	// Check request
	if req == nil {
		fmt.Println("request must not be nil")
		err = fmt.Errorf("request must not be nil")
		return
	}

	if req.Name == "" {
		err = fmt.Errorf("name must not be empty in the request")
		return
	}

	log.Printf("Received: %v", req.GetName())
	//Call Image API in order to get Dog Images's URL
	url, err := getRandomImageURL(req.GetName())
	if err != nil {
		err = fmt.Errorf("failed to call ImageAPI: %w", err)
		log.Println(err)
		return
	}
	log.Printf("Responded: %v", url)
	res.Message = url
	return res, nil
}

// GetFileSVCByte implements filesvc.FileSVCServer
func (s *Server) GetFileSVCByte(ctx context.Context, req *pb.FileSVCRequest) (res *pb.FileSVCByteReply, err error) {
	res = &pb.FileSVCByteReply{}

	// Check request
	if req == nil {
		fmt.Println("request must not be nil")
		err = fmt.Errorf("request must not be nil")
		return
	}

	if req.Name == "" {
		err = fmt.Errorf("name must not be empty in the request")
		return
	}

	log.Printf("Received: %v", req.GetName())

	//Call Image API in order to get Dog Images's URL
	url, err := getRandomImageURL(req.GetName())
	log.Println(url, err)

	if err != nil {
		err = fmt.Errorf("failed to call ImageAPI: %w", err)
		log.Println(err)
		return
	}
	imageData, _, err := getData(url)
	if err != nil {
		err = fmt.Errorf("can't get the Dog Images: %w", err)
		log.Println(err)
		return
	}
	log.Printf("Responded image data of : %v", url)
	res.Message = imageData
	return
}
