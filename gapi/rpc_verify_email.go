package gapi

import (
	"context"

	"github.com/rs/zerolog/log"
	db "github.com/terenjit/simplebank/db/sqlc"
	"github.com/terenjit/simplebank/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) VerifyEmail(ctx context.Context, req *pb.VerifyEmailRequest) (*pb.VerifyEmailResponse, error) {
	log.Info().
		Int64("email_id", req.GetEmailId()).
		Str("secret_code", req.GetSecretCode()).
		Msg("Verifying email")
	txResult, err := server.store.VerifyEmailTx(ctx, db.VerifyEmailTxParams{
		EmailId:    req.GetEmailId(),
		SecretCode: req.GetSecretCode(),
	})

	if err != nil {
		log.Error().Err(err).Msg("VerifyEmailTx failed")
		return nil, status.Errorf(codes.Internal, "failed to verify email")
	}

	rsp := &pb.VerifyEmailResponse{
		IsVerified: txResult.User.IsEmailVerified,
	}
	return rsp, nil
}
