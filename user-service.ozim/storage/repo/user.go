package repo

import (
	pb "github.com/najimovmashhurbek/project-api/user-service.ozim/genproto"
)

//UserStorageI ...
type UserStorageI interface {
	CreateUser(*pb.User) (*pb.CreatePostRes, error)
	DeleteUser(*pb.DeleteById) (*pb.DeleteUserRes, error)
	UpdateUser(*pb.User) (*pb.UpdateUserRes, error)
	GetAllUser(*pb.GetAllById) (*pb.User, error)
	ListUsers(limit, page int64) ([]*pb.User, int64, error)
}
