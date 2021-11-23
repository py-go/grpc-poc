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
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"

	"google.golang.org/grpc"

	pb "github.com/py-go/grpc-poc/pkg/filesvc"
	"github.com/spf13/cobra"
)

const (
	address     = "localhost:9000"
	defaultName = "dr-who"
)

// clientCmd represents the client command
var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "Query the gRPC server",

	Run: func(cmd *cobra.Command, args []string) {
		var conn *grpc.ClientConn
		conn, err := grpc.Dial(address, grpc.WithInsecure())
		if err != nil {
			log.Fatalf("did not connect: %s", err)
		}
		defer conn.Close()

		client := pb.NewFileSVCClient(conn)

		var name string

		// Contact the server and print out its response.
		// name := defaultName
		if len(os.Args) > 2 {
			name = os.Args[2]
		}
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
		defer cancel()
		// r, err := client.GetFileSVC(ctx, &pb.FileSVCRequest{Name: name})
		r, err := client.GetFileSVCByte(ctx, &pb.FileSVCRequest{Name: name})

		if err != nil {
			log.Fatalf("could not greet: %v", err)
		}
		tmpfile, err := ioutil.TempFile("./images", "image.*.png")
		destName := tmpfile.Name()
		if err != nil {
			log.Fatal(err)
		}
		// defer os.Remove(destName) // clean up
		_, err = tmpfile.Write(r.GetMessage())
		if err != nil {
			log.Fatalf("could not save image: %v", err)
		}
		if err := tmpfile.Close(); err != nil {
			log.Fatal(err)
		}
		fullpath, _ := filepath.Abs(destName)
		log.Printf("image of %s is saved here: %s\n", name, fullpath)

	},
}

func init() {
	rootCmd.AddCommand(clientCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// clientCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// clientCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
