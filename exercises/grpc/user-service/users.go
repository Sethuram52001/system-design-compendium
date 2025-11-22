package main

import (
	"errors"

	pb "github.com/Sethuram52001/system-design-compendium/exercises/grpc/user-service/userpb"
)

type UserStore struct {
	users map[string]*pb.User
}

func NewUserStore() *UserStore {
	store := &UserStore{
		users: make(map[string]*pb.User),
	}

	store.users["1"] = &pb.User{
		Id:    "1",
		Name:  "John Doe",
		Email: "john@example.com",
		Age:   30,
	}
	store.users["2"] = &pb.User{
		Id:    "2",
		Name:  "Jane Smith",
		Email: "jane@example.com",
		Age:   25,
	}

	return store
}

func (s *UserStore) GetUser(id string) (*pb.User, error) {
	user, exists := s.users[id]
	if !exists {
		return nil, errors.New("user not found")
	}
	return user, nil
}
