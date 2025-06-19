package main

import (
	"context"
	"log"
	"net"

	"todocli/internal/model"
	"todocli/internal/service"
	"todocli/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct {
	pb.UnimplementedTaskServiceServer
}

func (s *server) CreateTask(ctx context.Context, in *pb.Task) (*pb.TaskID, error) {
	t, err := model.NewTask(service.GenerateTaskID(), in.Title, in.Description, in.Status)
	if err != nil {
		return nil, err
	}
	service.AddTask(t)
	return &pb.TaskID{Id: int32(t.ID)}, nil
}

func (s *server) GetTask(ctx context.Context, in *pb.TaskID) (*pb.Task, error) {
	t := service.FindTaskByID(int(in.Id))
	if t == nil {
		return nil, grpc.Errorf(404, "task not found")
	}
	return &pb.Task{
		Id:          int32(t.ID),
		Title:       t.Title,
		Description: t.Description,
		Status:      t.StatusRaw,
	}, nil
}

func (s *server) ListTasks(ctx context.Context, _ *pb.Empty) (*pb.TaskList, error) {
	tasks := service.GetAllTasks()
	var list []*pb.Task
	for _, t := range tasks {
		list = append(list, &pb.Task{
			Id:          int32(t.ID),
			Title:       t.Title,
			Description: t.Description,
			Status:      t.StatusRaw,
		})
	}
	return &pb.TaskList{Tasks: list}, nil
}

func (s *server) UpdateTask(ctx context.Context, in *pb.Task) (*pb.Task, error) {
	t := service.FindTaskByID(int(in.Id))
	if t == nil {
		return nil, grpc.Errorf(404, "task not found")
	}
	t.Title = in.Title
	t.Description = in.Description
	t.StatusRaw = in.Status
	t.Normalize()
	service.AddTask(t)
	return in, nil
}

func (s *server) DeleteTask(ctx context.Context, in *pb.TaskID) (*pb.Empty, error) {
	err := service.DeleteTask(int(in.Id))
	if err != nil {
		return nil, grpc.Errorf(404, "task not found")
	}
	return &pb.Empty{}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(err)
	}
	s := grpc.NewServer()
	pb.RegisterTaskServiceServer(s, &server{})
	reflection.Register(s)
	log.Println("gRPC server listening on :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
