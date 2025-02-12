package handler

import (
	"context"
	"errors"
	"net/http"

	"shop/internal/app/entity"
	"shop/internal/app/usecase"
	"shop/internal/controller/http/gen"
	"shop/internal/pkg"

	"go.uber.org/zap"
)

var _ gen.StrictServerInterface = (*server)(nil)

type server struct {
	address     string
	userUsecase usecase.UserUseCase
	logger      *zap.Logger
}

func NewServer(logger *zap.Logger, userUsecase usecase.UserUseCase, address string) *server {
	return &server{
		address:     address,
		userUsecase: userUsecase,
		logger:      logger,
	}
}

func (r *server) PostApiAuth(ctx context.Context, request gen.PostApiAuthRequestObject) (gen.PostApiAuthResponseObject, error) {
	if request.Body == nil {
		return gen.PostApiAuth400JSONResponse{Errors: pkg.PointerTo("empty body")}, nil
	}

	login := entity.Login(*request.Body)

	token, err := r.userUsecase.Auth(ctx, login)
	if err != nil {
		if errors.Is(err, usecase.ErrInvalidPassword) {
			return gen.PostApiAuth400JSONResponse{Errors: pkg.PointerTo(err.Error())}, nil
		}

		return gen.PostApiAuth500JSONResponse{Errors: pkg.PointerTo(err.Error())}, nil
	}

	return gen.PostApiAuth200JSONResponse(gen.AuthResponse{Token: (*string)(token)}), nil
}

func (r *server) GetApiBuyItem(ctx context.Context, request gen.GetApiBuyItemRequestObject) (gen.GetApiBuyItemResponseObject, error) {
	return nil, nil
}

func (r *server) GetApiInfo(ctx context.Context, request gen.GetApiInfoRequestObject) (gen.GetApiInfoResponseObject, error) {
	return nil, nil
}

func (r *server) PostApiSendCoin(ctx context.Context, request gen.PostApiSendCoinRequestObject) (gen.PostApiSendCoinResponseObject, error) {
	return nil, nil
}

func (r *server) Start() {
	srv := gen.NewStrictHandler(r, nil)
	handler := gen.Handler(srv)

	http.ListenAndServe(r.address, handler)
}
