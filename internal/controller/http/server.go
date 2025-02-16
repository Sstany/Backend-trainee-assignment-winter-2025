package handler

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"shop/internal/app/entity"
	"shop/internal/app/usecase"
	"shop/internal/controller/http/gen"
	"shop/internal/pkg"

	"github.com/oapi-codegen/runtime/strictmiddleware/nethttp"
	"go.uber.org/zap"
)

var _ gen.StrictServerInterface = (*Server)(nil)

type contextKey string

const usernameContextKey contextKey = "username"

type Server struct {
	address     string
	userUsecase usecase.UserUseCase
	authUsecase usecase.AuthUseCase
	logger      *zap.Logger
}

func NewServer(
	logger *zap.Logger,
	userUsecase usecase.UserUseCase,
	authUsecase usecase.AuthUseCase,
	address string,
) *Server {
	return &Server{
		address:     address,
		userUsecase: userUsecase,
		authUsecase: authUsecase,
		logger:      logger,
	}
}

func (r Server) PostApiAuth(ctx context.Context, request gen.PostApiAuthRequestObject) (gen.PostApiAuthResponseObject, error) {
	if request.Body == nil {
		return gen.PostApiAuth400JSONResponse{Errors: pkg.PointerTo("empty body")}, nil
	}

	login := entity.Login(*request.Body)

	token, err := r.authUsecase.Auth(ctx, login)
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

func (r *Server) GetApiBuyItem(ctx context.Context, request gen.GetApiBuyItemRequestObject) (gen.GetApiBuyItemResponseObject, error) {
	item := entity.Item(request.Item)
	username := ctx.Value(usernameContextKey)
	usrStr, ok := username.(string)
	if !ok {
		return gen.GetApiBuyItem500JSONResponse{Errors: pkg.PointerTo("username  not string")}, nil
	}

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
		if errors.Is(err, usecase.ErrReveicerNotExists) {
			return gen.GetApiBuyItem400JSONResponse{Errors: pkg.PointerTo(err.Error())}, nil
		}

		return gen.GetApiBuyItem500JSONResponse{Errors: pkg.PointerTo(err.Error())}, nil
	}

	return gen.GetApiBuyItem200Response{}, nil

}

func (r *Server) GetApiInfo(ctx context.Context, request gen.GetApiInfoRequestObject) (gen.GetApiInfoResponseObject, error) {
	username := ctx.Value(usernameContextKey)
	usrStr, ok := username.(string)
	if !ok {
		return gen.GetApiInfo500JSONResponse{Errors: pkg.PointerTo("username  not string")}, nil
	}

	infoReq := entity.InfoRequest{
		Username: usrStr,
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

func (r *Server) PostApiSendCoin(ctx context.Context, request gen.PostApiSendCoinRequestObject) (gen.PostApiSendCoinResponseObject, error) {
	username := ctx.Value(usernameContextKey)
	usrStr, ok := username.(string)
	if !ok {
		return gen.PostApiSendCoin500JSONResponse{Errors: pkg.PointerTo("username  not string")}, nil
	}

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

func (r *Server) NewAuthMiddleware() nethttp.StrictHTTPMiddlewareFunc {
	return func(f nethttp.StrictHTTPHandlerFunc, operationID string) nethttp.StrictHTTPHandlerFunc {
		return func(
			ctx context.Context,
			w http.ResponseWriter,
			req *http.Request,
			request interface{},
		) (response interface{}, err error) {
			if operationID == "PostApiAuth" {
				return f(ctx, w, req, request)
			}

			authHeader := req.Header.Get("Authorization")
			if authHeader == "" {
				responseErr(w, "no authorization header", http.StatusUnauthorized)
				return nil, nil
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			if tokenString == authHeader {
				responseErr(w, "Invalid token format", http.StatusUnauthorized)
				return nil, nil
			}

			username, valid, err := r.authUsecase.AuthenticateWithAccessToken(tokenString)
			if err != nil {
				responseErr(w, "authentication with access token failed", http.StatusUnauthorized)
				return nil, nil
			}
			if !valid {
				responseErr(w, "Invalid token", http.StatusUnauthorized)
				return nil, nil
			}

			uCtx := context.WithValue(ctx, usernameContextKey, username)

			return f(uCtx, w, req, request)
		}
	}
}

func requestErrorHandler(w http.ResponseWriter, r *http.Request, err error) {
	responseErr(w, err.Error(), http.StatusInternalServerError)
}

func responseErrorHandlerFunc(w http.ResponseWriter, r *http.Request, err error) {
	responseErr(w, err.Error(), http.StatusInternalServerError)
}

func responseErr(w http.ResponseWriter, errStr string, status int) {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	resp := gen.ErrorResponse{
		Errors: &errStr,
	}

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	_, err = w.Write(jsonResp)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
}

func (r *Server) Start() {
	srv := gen.NewStrictHandlerWithOptions(
		r,
		[]gen.StrictMiddlewareFunc{r.NewAuthMiddleware()},
		gen.StrictHTTPServerOptions{
			RequestErrorHandlerFunc:  requestErrorHandler,
			ResponseErrorHandlerFunc: responseErrorHandlerFunc,
		},
	)

	handler := gen.Handler(srv)

	s := http.Server{
		Addr:              r.address,
		Handler:           handler,
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       2 * time.Minute,
	}

	log.Fatal(s.ListenAndServe())
}
