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
		if errors.Is(err, usecase.ErrUnsafePassword) {
			return gen.PostApiAuth400JSONResponse{Errors: pkg.PointerTo(err.Error())}, nil
		}
		if errors.Is(err, usecase.ErrLongPassword) {
			return gen.PostApiAuth400JSONResponse{Errors: pkg.PointerTo(err.Error())}, nil
		}

		return gen.PostApiAuth500JSONResponse{Errors: pkg.PointerTo(err.Error())}, nil
	}

	return gen.PostApiAuth200JSONResponse(gen.AuthResponse{Token: (*string)(token)}), nil
}

func (r *server) GetApiBuyItem(ctx context.Context, request gen.GetApiBuyItemRequestObject) (gen.GetApiBuyItemResponseObject, error) {
	item := entity.Item(request.Item)
	//username := ctx.Value("username")

	var username any = "T"
	usrStr := username.(string)

	req := entity.ItemRequest{
		Username: usrStr,
		Item:     string(item),
	}

	err := r.userUsecase.Buy(ctx, req)
	if err != nil {
		if errors.Is(err, usecase.ErrInsufficientBalance) {
			return gen.GetApiBuyItem400JSONResponse{Errors: pkg.PointerTo(err.Error())}, nil
		}
		if errors.Is(err, usecase.ErrItemNotExists) {
			return gen.GetApiBuyItem400JSONResponse{Errors: pkg.PointerTo(err.Error())}, nil
		}

		return gen.GetApiBuyItem500JSONResponse{Errors: pkg.PointerTo(err.Error())}, nil
	}

	return gen.GetApiBuyItem200Response{}, nil

}

func (r *server) GetApiInfo(ctx context.Context, request gen.GetApiInfoRequestObject) (gen.GetApiInfoResponseObject, error) {

	infoReq := entity.InfoRequest{
		Username: "T",
	}

	info, err := r.userUsecase.Info(ctx, infoReq)
	if err != nil {
		return gen.GetApiInfo500JSONResponse{Errors: pkg.PointerTo(err.Error())}, nil
	}

	return gen.GetApiInfo200JSONResponse(convertInfoToInfoResponse(info)), nil
}

func convertInfoToInfoResponse(info *entity.InfoResponse) gen.InfoResponse {
	var received *[]struct {
		Amount   *int    `json:"amount,omitempty"`
		FromUser *string `json:"fromUser,omitempty"`
	}

	if len(info.CoinHistory.Received) > 0 {
		rec := make([]struct {
			Amount   *int    `json:"amount,omitempty"`
			FromUser *string `json:"fromUser,omitempty"`
		}, len(info.CoinHistory.Received))

		for i := range info.CoinHistory.Received {
			rec[i] =
				struct {
					Amount   *int    `json:"amount,omitempty"`
					FromUser *string `json:"fromUser,omitempty"`
				}{
					Amount:   &info.CoinHistory.Received[i].Amount,
					FromUser: &info.CoinHistory.Received[i].FromUser,
				}
		}

		received = &rec
	}

	var sent *[]struct {
		Amount *int    `json:"amount,omitempty"`
		ToUser *string `json:"toUser,omitempty"`
	}

	if len(info.CoinHistory.Sent) > 0 {
		s := make([]struct {
			Amount *int    `json:"amount,omitempty"`
			ToUser *string `json:"toUser,omitempty"`
		}, len(info.CoinHistory.Sent))

		for i := range info.CoinHistory.Sent {
			s[i] =
				struct {
					Amount *int    `json:"amount,omitempty"`
					ToUser *string `json:"toUser,omitempty"`
				}{
					Amount: &info.CoinHistory.Sent[i].Amount,
					ToUser: &info.CoinHistory.Sent[i].ToUser,
				}
		}

		sent = &s
	}

	var inventory *[]struct {
		Quantity *int    `json:"quantity,omitempty"`
		Type     *string `json:"type,omitempty"`
	}

	if len(info.Inventory) > 0 {
		inv := make([]struct {
			Quantity *int    `json:"quantity,omitempty"`
			Type     *string `json:"type,omitempty"`
		}, len(info.Inventory))

		for i := range info.Inventory {
			inv[i] =
				struct {
					Quantity *int    `json:"quantity,omitempty"`
					Type     *string `json:"type,omitempty"`
				}{
					Quantity: &info.Inventory[i].Quantity,
					Type:     &info.Inventory[i].Type,
				}
		}

		inventory = &inv
	}

	return gen.InfoResponse{
		CoinHistory: &struct {
			Received *[]struct {
				Amount   *int    "json:\"amount,omitempty\""
				FromUser *string "json:\"fromUser,omitempty\""
			} "json:\"received,omitempty\""
			Sent *[]struct {
				Amount *int    "json:\"amount,omitempty\""
				ToUser *string "json:\"toUser,omitempty\""
			} "json:\"sent,omitempty\""
		}{
			Received: received,
			Sent:     sent,
		},
		Inventory: inventory,
		Coins:     &info.Coins,
	}
}

func (r *server) PostApiSendCoin(ctx context.Context, request gen.PostApiSendCoinRequestObject) (gen.PostApiSendCoinResponseObject, error) {
	var username any = "T"
	usrStr := username.(string)

	req := entity.SendCoinRequest{
		Amount:   request.Body.Amount,
		FromUser: usrStr,
		ToUser:   request.Body.ToUser,
	}

	err := r.userUsecase.Send(ctx, req)
	if err != nil {
		if errors.Is(err, usecase.ErrWrongCoinAmount) {
			return gen.PostApiSendCoin400JSONResponse{Errors: pkg.PointerTo(err.Error())}, nil
		}

		return gen.PostApiSendCoin500JSONResponse{Errors: pkg.PointerTo(err.Error())}, nil
	}

	return gen.PostApiSendCoin200Response{}, nil
}

func (r *server) Start() {
	srv := gen.NewStrictHandler(r, nil)
	handler := gen.Handler(srv)

	http.ListenAndServe(r.address, handler)
}
