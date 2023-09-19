package gapi

import (
	"context"

	"github.com/lib/pq"
	db "github.com/nphattai/go-simple-bank/db/sqlc"
	"github.com/nphattai/go-simple-bank/pb"
	"github.com/nphattai/go-simple-bank/util"
	"github.com/nphattai/go-simple-bank/val"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	if violation := validateCreateUserRequest(req); violation != nil {
		return nil, invalidArgumentError(violation)
	}

	hashedPassword, err := util.HashPassword(req.GetPassword())

	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to hash password")
	}

	arg := db.CreateUserParams{
		Username:       req.GetUsername(),
		HashedPassword: hashedPassword,
		FullName:       req.GetFullName(),
		Email:          req.GetEmail(),
	}

	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok {
			if pgErr.Code.Name() == "unique_violation" {
				return nil, status.Errorf(codes.AlreadyExists, "username already exists: %s", err)
			}
			return nil, status.Errorf(codes.Internal, "failed to create user: %s", err)
		}
	}

	rsp := &pb.CreateUserResponse{
		User: converter(user),
	}

	return rsp, nil
}

func validateCreateUserRequest(req *pb.CreateUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := val.ValidateUsername(req.Username); err != nil {
		violations = append(violations, fieldViolation("username", err))
	}
	if err := val.ValidateFullName(req.FullName); err != nil {
		violations = append(violations, fieldViolation("full_name", err))
	}
	if err := val.ValidatePassword(req.Password); err != nil {
		violations = append(violations, fieldViolation("password", err))
	}
	if err := val.ValidateEmail(req.Email); err != nil {
		violations = append(violations, fieldViolation("email", err))
	}

	return violations
}
