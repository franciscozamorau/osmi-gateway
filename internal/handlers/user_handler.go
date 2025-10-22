package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	pb "github.com/franciscozamorau/osmi-gateway/gen"
	"google.golang.org/grpc"
)

func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	var req pb.UserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	conn, err := grpc.Dial("osmi-server:50051", grpc.WithInsecure())
	if err != nil {
		http.Error(w, "Failed to connect to gRPC", http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	client := pb.NewOsmiServiceClient(conn)
	resp, err := client.CreateUser(context.Background(), &req)
	if err != nil {
		http.Error(w, "gRPC error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(resp)
}
