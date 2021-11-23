# grpc-poc
It's test project to show the basic implementation of grpc using golang.

This project consists of Server as well client examples.


## GRPC Golang Server
The server is responsible to make the external API call to retrieve the image URL and sending the byte string of image data back to the client.

To start the server you can simply run the below command:

```go run main.go server```

## GRPC Golang Client
From the client, you can send the dog breed name and it will download the image for you in your local folder.

### Command

To run below command to call the client with breed ```hound```.

```go run main.go client hound```

### Response

```2021/11/19 22:50:23 image of hound is saved here: /home/nroot/grpc-poc/images/image.2878467758.png```
