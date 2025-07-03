package gapi

import (
	db "github.com/terenjit/simplebank/db/sqlc"
	"github.com/terenjit/simplebank/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func CovertUser(User db.User) *pb.User {
	return &pb.User{
		Username:          User.Username,
		FullName:          User.FullName,
		Email:             User.Email,
		PasswordChangedAt: timestamppb.New(User.PasswordChangedAt),
		CreatedAt:         timestamppb.New(User.CreatedAt),
	}
}
