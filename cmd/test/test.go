package main

import (
	"context"
	"fmt"
	"log"
	"qd-email-api/pb/gen/go/pb_email"
	"time"

	pkgLogger "github.com/gustavo-m-franco/qd-common/pkg/log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func contextWithCorrelationID(correlationID string) context.Context {
	md := metadata.New(map[string]string{
		pkgLogger.CorrelationIDKey: correlationID,
	})
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	return ctx
}

func main() {
	// Set up a connection to the gRPC server.
	conn, err := grpc.Dial("localhost:8081", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	// Create a gRPC client for the EmailService.
	client := pb_email.NewEmailServiceClient(conn)

	// Define the SendEmailRequest.
	request := &pb_email.SendEmailRequest{
		To:      "gusfran17@gmail.com",
		Subject: "Hello Gus, gRPC!",
		Body:    "This is a test email sent via gRPC. By me, GUS!!!!!!",
	}

	// Set a context with a timeout.
	ctx, cancel := context.WithTimeout(contextWithCorrelationID("1111111111111111111"), 10*time.Second)
	defer cancel()

	// Call the SendEmail method on the gRPC server.
	response, err := client.SendEmail(ctx, request)
	if err != nil {
		log.Fatalf("SendEmail failed: %v", err)
	}

	// Handle the response.
	if response.Success {
		fmt.Printf("Email sent successfully. Response message: %s\n", response.Message)
	} else {
		fmt.Printf("Email sending failed. Response message: %s\n", response.Message)
	}
}
