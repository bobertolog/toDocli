package main

import (
	"context"
	"log"
	"net"

	"todocli/internal/repository"
	"todocli/internal/service"
	"todocli/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

var taskService service.TaskService

type server struct {
	pb.UnimplementedTaskServiceServer
}

func (s *server) CreateTask(ctx context.Context, in *pb.Task) (*pb.TaskID, error) {
	task, err := taskService.Create(in.Title, in.Description, in.Status)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to create task: %v", err)
	}
	return &pb.TaskID{Id: int32(task.ID)}, nil
}

func (s *server) GetTask(ctx context.Context, in *pb.TaskID) (*pb.Task, error) {
	task, err := taskService.GetByID(int(in.Id))
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "task not found")
	}
	return &pb.Task{
		Id:          int32(task.ID),
		Title:       task.Title,
		Description: task.Description,
		Status:      task.Status.String(),
	}, nil
}

func (s *server) ListTasks(ctx context.Context, _ *pb.Empty) (*pb.TaskList, error) {
	tasks := taskService.GetAll()
	var list []*pb.Task
	for _, t := range tasks {
		list = append(list, &pb.Task{
			Id:          int32(t.ID),
			Title:       t.Title,
			Description: t.Description,
			Status:      t.Status.String(),
		})
	}
	return &pb.TaskList{Tasks: list}, nil
}

func (s *server) UpdateTask(ctx context.Context, in *pb.Task) (*pb.Task, error) {
	err := taskService.Update(int(in.Id), in.Title, in.Description, in.Status)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "task not found or update failed: %v", err)
	}
	return in, nil
}

func (s *server) DeleteTask(ctx context.Context, in *pb.TaskID) (*pb.Empty, error) {
	err := taskService.Delete(int(in.Id))
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "task not found")
	}
	return &pb.Empty{}, nil
}

func main() {
	repo := repository.NewInMemoryRepository()
	taskService = service.NewTaskService(repo)

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
