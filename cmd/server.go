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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/spf13/cobra"
	"golang.org/x/xerrors"

	pb "github.com/py-go/grpc-poc/pkg/filesvc"
	"google.golang.org/grpc"
)

const (
	port    = ":9000"
	BaseAPI = "https://dog.ceo/api"
)

// server is used to implement filesvc.FileSVCServer.
type Server struct {
	pb.UnimplementedFileSVCServer
}

type APIResponse struct {
	Message []string `json: "message"`
	Status  string   `json: "status"`
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
func (s *Server) GetFileSVC(ctx context.Context, req *pb.FileSVCRequest) (*pb.FileSVCReply, error) {
	res := &pb.FileSVCReply{}

	// Check request
	if req == nil {
		fmt.Println("request must not be nil")
		return res, xerrors.Errorf("request must not be nil")
	}

	if req.Name == "" {
		fmt.Println("name must not be empty in the request")
		return res, xerrors.Errorf("name must not be empty in the request")
	}

	log.Printf("Received: %v", req.GetName())

	//Call KuteGo API in order to get Dog Images's URL
	response, err := http.Get(BaseAPI + fmt.Sprintf("/breed/%s/images", req.GetName()))
	if err != nil {
		log.Fatalf("failed to call KuteGoAPI: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode == 200 {
		// Transform our response to a []byte
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Fatalf("failed to read response body: %v", err)
		}

		// Put only needed informations of the JSON document in our array of Dog Images
		var data APIResponse
		err = json.Unmarshal(body, &data)
		if err != nil {
			log.Fatalf("failed to unmarshal JSON: %v", err)
		}

		// Create a string with all of the Dog Images's name and a blank line as separator
		var stringData strings.Builder
		for _, gopher := range data.Message {
			stringData.WriteString(gopher + "\n")
		}
		res.Message = stringData.String()
	} else {
		log.Println("Can't get the Dog Images :-(")
	}

	return res, nil
}
