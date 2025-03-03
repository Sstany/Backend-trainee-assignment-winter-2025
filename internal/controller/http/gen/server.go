// Package gen provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/oapi-codegen/oapi-codegen/v2 version v2.4.1 DO NOT EDIT.
package gen

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/oapi-codegen/runtime"
	strictnethttp "github.com/oapi-codegen/runtime/strictmiddleware/nethttp"
)

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Аутентификация и получение JWT-токена. При первой аутентификации пользователь создается автоматически.
	// (POST /api/auth)
	PostApiAuth(w http.ResponseWriter, r *http.Request)
	// Купить предмет за монеты.
	// (GET /api/buy/{item})
	GetApiBuyItem(w http.ResponseWriter, r *http.Request, item string)
	// Получить информацию о монетах, инвентаре и истории транзакций.
	// (GET /api/info)
	GetApiInfo(w http.ResponseWriter, r *http.Request)
	// Отправить монеты другому пользователю.
	// (POST /api/sendCoin)
	PostApiSendCoin(w http.ResponseWriter, r *http.Request)
}

// Unimplemented server implementation that returns http.StatusNotImplemented for each endpoint.

type Unimplemented struct{}

// Аутентификация и получение JWT-токена. При первой аутентификации пользователь создается автоматически.
// (POST /api/auth)
func (_ Unimplemented) PostApiAuth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

// Купить предмет за монеты.
// (GET /api/buy/{item})
func (_ Unimplemented) GetApiBuyItem(w http.ResponseWriter, r *http.Request, item string) {
	w.WriteHeader(http.StatusNotImplemented)
}

// Получить информацию о монетах, инвентаре и истории транзакций.
// (GET /api/info)
func (_ Unimplemented) GetApiInfo(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

// Отправить монеты другому пользователю.
// (POST /api/sendCoin)
func (_ Unimplemented) PostApiSendCoin(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler            ServerInterface
	HandlerMiddlewares []MiddlewareFunc
	ErrorHandlerFunc   func(w http.ResponseWriter, r *http.Request, err error)
}

type MiddlewareFunc func(http.Handler) http.Handler

// PostApiAuth operation middleware
func (siw *ServerInterfaceWrapper) PostApiAuth(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	ctx = context.WithValue(ctx, BearerAuthScopes, []string{})

	r = r.WithContext(ctx)

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.PostApiAuth(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

// GetApiBuyItem operation middleware
func (siw *ServerInterfaceWrapper) GetApiBuyItem(w http.ResponseWriter, r *http.Request) {

	var err error

	// ------------- Path parameter "item" -------------
	var item string

	err = runtime.BindStyledParameterWithOptions("simple", "item", chi.URLParam(r, "item"), &item, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "item", Err: err})
		return
	}

	ctx := r.Context()

	ctx = context.WithValue(ctx, BearerAuthScopes, []string{})

	r = r.WithContext(ctx)

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetApiBuyItem(w, r, item)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

// GetApiInfo operation middleware
func (siw *ServerInterfaceWrapper) GetApiInfo(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	ctx = context.WithValue(ctx, BearerAuthScopes, []string{})

	r = r.WithContext(ctx)

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetApiInfo(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

// PostApiSendCoin operation middleware
func (siw *ServerInterfaceWrapper) PostApiSendCoin(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	ctx = context.WithValue(ctx, BearerAuthScopes, []string{})

	r = r.WithContext(ctx)

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.PostApiSendCoin(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

type UnescapedCookieParamError struct {
	ParamName string
	Err       error
}

func (e *UnescapedCookieParamError) Error() string {
	return fmt.Sprintf("error unescaping cookie parameter '%s'", e.ParamName)
}

func (e *UnescapedCookieParamError) Unwrap() error {
	return e.Err
}

type UnmarshalingParamError struct {
	ParamName string
	Err       error
}

func (e *UnmarshalingParamError) Error() string {
	return fmt.Sprintf("Error unmarshaling parameter %s as JSON: %s", e.ParamName, e.Err.Error())
}

func (e *UnmarshalingParamError) Unwrap() error {
	return e.Err
}

type RequiredParamError struct {
	ParamName string
}

func (e *RequiredParamError) Error() string {
	return fmt.Sprintf("Query argument %s is required, but not found", e.ParamName)
}

type RequiredHeaderError struct {
	ParamName string
	Err       error
}

func (e *RequiredHeaderError) Error() string {
	return fmt.Sprintf("Header parameter %s is required, but not found", e.ParamName)
}

func (e *RequiredHeaderError) Unwrap() error {
	return e.Err
}

type InvalidParamFormatError struct {
	ParamName string
	Err       error
}

func (e *InvalidParamFormatError) Error() string {
	return fmt.Sprintf("Invalid format for parameter %s: %s", e.ParamName, e.Err.Error())
}

func (e *InvalidParamFormatError) Unwrap() error {
	return e.Err
}

type TooManyValuesForParamError struct {
	ParamName string
	Count     int
}

func (e *TooManyValuesForParamError) Error() string {
	return fmt.Sprintf("Expected one value for %s, got %d", e.ParamName, e.Count)
}

// Handler creates http.Handler with routing matching OpenAPI spec.
func Handler(si ServerInterface) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{})
}

type ChiServerOptions struct {
	BaseURL          string
	BaseRouter       chi.Router
	Middlewares      []MiddlewareFunc
	ErrorHandlerFunc func(w http.ResponseWriter, r *http.Request, err error)
}

// HandlerFromMux creates http.Handler with routing matching OpenAPI spec based on the provided mux.
func HandlerFromMux(si ServerInterface, r chi.Router) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{
		BaseRouter: r,
	})
}

func HandlerFromMuxWithBaseURL(si ServerInterface, r chi.Router, baseURL string) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{
		BaseURL:    baseURL,
		BaseRouter: r,
	})
}

// HandlerWithOptions creates http.Handler with additional options
func HandlerWithOptions(si ServerInterface, options ChiServerOptions) http.Handler {
	r := options.BaseRouter

	if r == nil {
		r = chi.NewRouter()
	}
	if options.ErrorHandlerFunc == nil {
		options.ErrorHandlerFunc = func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	}
	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: options.Middlewares,
		ErrorHandlerFunc:   options.ErrorHandlerFunc,
	}

	r.Group(func(r chi.Router) {
		r.Post(options.BaseURL+"/api/auth", wrapper.PostApiAuth)
	})
	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/api/buy/{item}", wrapper.GetApiBuyItem)
	})
	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/api/info", wrapper.GetApiInfo)
	})
	r.Group(func(r chi.Router) {
		r.Post(options.BaseURL+"/api/sendCoin", wrapper.PostApiSendCoin)
	})

	return r
}

type PostApiAuthRequestObject struct {
	Body *PostApiAuthJSONRequestBody
}

type PostApiAuthResponseObject interface {
	VisitPostApiAuthResponse(w http.ResponseWriter) error
}

type PostApiAuth200JSONResponse AuthResponse

func (response PostApiAuth200JSONResponse) VisitPostApiAuthResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	return json.NewEncoder(w).Encode(response)
}

type PostApiAuth400JSONResponse ErrorResponse

func (response PostApiAuth400JSONResponse) VisitPostApiAuthResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(400)

	return json.NewEncoder(w).Encode(response)
}

type PostApiAuth401JSONResponse ErrorResponse

func (response PostApiAuth401JSONResponse) VisitPostApiAuthResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(401)

	return json.NewEncoder(w).Encode(response)
}

type PostApiAuth500JSONResponse ErrorResponse

func (response PostApiAuth500JSONResponse) VisitPostApiAuthResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(500)

	return json.NewEncoder(w).Encode(response)
}

type GetApiBuyItemRequestObject struct {
	Item string `json:"item"`
}

type GetApiBuyItemResponseObject interface {
	VisitGetApiBuyItemResponse(w http.ResponseWriter) error
}

type GetApiBuyItem200Response struct {
}

func (response GetApiBuyItem200Response) VisitGetApiBuyItemResponse(w http.ResponseWriter) error {
	w.WriteHeader(200)
	return nil
}

type GetApiBuyItem400JSONResponse ErrorResponse

func (response GetApiBuyItem400JSONResponse) VisitGetApiBuyItemResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(400)

	return json.NewEncoder(w).Encode(response)
}

type GetApiBuyItem401JSONResponse ErrorResponse

func (response GetApiBuyItem401JSONResponse) VisitGetApiBuyItemResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(401)

	return json.NewEncoder(w).Encode(response)
}

type GetApiBuyItem500JSONResponse ErrorResponse

func (response GetApiBuyItem500JSONResponse) VisitGetApiBuyItemResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(500)

	return json.NewEncoder(w).Encode(response)
}

type GetApiInfoRequestObject struct {
}

type GetApiInfoResponseObject interface {
	VisitGetApiInfoResponse(w http.ResponseWriter) error
}

type GetApiInfo200JSONResponse InfoResponse

func (response GetApiInfo200JSONResponse) VisitGetApiInfoResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	return json.NewEncoder(w).Encode(response)
}

type GetApiInfo400JSONResponse ErrorResponse

func (response GetApiInfo400JSONResponse) VisitGetApiInfoResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(400)

	return json.NewEncoder(w).Encode(response)
}

type GetApiInfo401JSONResponse ErrorResponse

func (response GetApiInfo401JSONResponse) VisitGetApiInfoResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(401)

	return json.NewEncoder(w).Encode(response)
}

type GetApiInfo500JSONResponse ErrorResponse

func (response GetApiInfo500JSONResponse) VisitGetApiInfoResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(500)

	return json.NewEncoder(w).Encode(response)
}

type PostApiSendCoinRequestObject struct {
	Body *PostApiSendCoinJSONRequestBody
}

type PostApiSendCoinResponseObject interface {
	VisitPostApiSendCoinResponse(w http.ResponseWriter) error
}

type PostApiSendCoin200Response struct {
}

func (response PostApiSendCoin200Response) VisitPostApiSendCoinResponse(w http.ResponseWriter) error {
	w.WriteHeader(200)
	return nil
}

type PostApiSendCoin400JSONResponse ErrorResponse

func (response PostApiSendCoin400JSONResponse) VisitPostApiSendCoinResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(400)

	return json.NewEncoder(w).Encode(response)
}

type PostApiSendCoin401JSONResponse ErrorResponse

func (response PostApiSendCoin401JSONResponse) VisitPostApiSendCoinResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(401)

	return json.NewEncoder(w).Encode(response)
}

type PostApiSendCoin500JSONResponse ErrorResponse

func (response PostApiSendCoin500JSONResponse) VisitPostApiSendCoinResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(500)

	return json.NewEncoder(w).Encode(response)
}

// StrictServerInterface represents all server handlers.
type StrictServerInterface interface {
	// Аутентификация и получение JWT-токена. При первой аутентификации пользователь создается автоматически.
	// (POST /api/auth)
	PostApiAuth(ctx context.Context, request PostApiAuthRequestObject) (PostApiAuthResponseObject, error)
	// Купить предмет за монеты.
	// (GET /api/buy/{item})
	GetApiBuyItem(ctx context.Context, request GetApiBuyItemRequestObject) (GetApiBuyItemResponseObject, error)
	// Получить информацию о монетах, инвентаре и истории транзакций.
	// (GET /api/info)
	GetApiInfo(ctx context.Context, request GetApiInfoRequestObject) (GetApiInfoResponseObject, error)
	// Отправить монеты другому пользователю.
	// (POST /api/sendCoin)
	PostApiSendCoin(ctx context.Context, request PostApiSendCoinRequestObject) (PostApiSendCoinResponseObject, error)
}

type StrictHandlerFunc = strictnethttp.StrictHTTPHandlerFunc
type StrictMiddlewareFunc = strictnethttp.StrictHTTPMiddlewareFunc

type StrictHTTPServerOptions struct {
	RequestErrorHandlerFunc  func(w http.ResponseWriter, r *http.Request, err error)
	ResponseErrorHandlerFunc func(w http.ResponseWriter, r *http.Request, err error)
}

func NewStrictHandler(ssi StrictServerInterface, middlewares []StrictMiddlewareFunc) ServerInterface {
	return &strictHandler{ssi: ssi, middlewares: middlewares, options: StrictHTTPServerOptions{
		RequestErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		},
		ResponseErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		},
	}}
}

func NewStrictHandlerWithOptions(ssi StrictServerInterface, middlewares []StrictMiddlewareFunc, options StrictHTTPServerOptions) ServerInterface {
	return &strictHandler{ssi: ssi, middlewares: middlewares, options: options}
}

type strictHandler struct {
	ssi         StrictServerInterface
	middlewares []StrictMiddlewareFunc
	options     StrictHTTPServerOptions
}

// PostApiAuth operation middleware
func (sh *strictHandler) PostApiAuth(w http.ResponseWriter, r *http.Request) {
	var request PostApiAuthRequestObject

	var body PostApiAuthJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		sh.options.RequestErrorHandlerFunc(w, r, fmt.Errorf("can't decode JSON body: %w", err))
		return
	}
	request.Body = &body

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (interface{}, error) {
		return sh.ssi.PostApiAuth(ctx, request.(PostApiAuthRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "PostApiAuth")
	}

	response, err := handler(r.Context(), w, r, request)

	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(PostApiAuthResponseObject); ok {
		if err := validResponse.VisitPostApiAuthResponse(w); err != nil {
			sh.options.ResponseErrorHandlerFunc(w, r, err)
		}
	} else if response != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("unexpected response type: %T", response))
	}
}

// GetApiBuyItem operation middleware
func (sh *strictHandler) GetApiBuyItem(w http.ResponseWriter, r *http.Request, item string) {
	var request GetApiBuyItemRequestObject

	request.Item = item

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (interface{}, error) {
		return sh.ssi.GetApiBuyItem(ctx, request.(GetApiBuyItemRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "GetApiBuyItem")
	}

	response, err := handler(r.Context(), w, r, request)

	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(GetApiBuyItemResponseObject); ok {
		if err := validResponse.VisitGetApiBuyItemResponse(w); err != nil {
			sh.options.ResponseErrorHandlerFunc(w, r, err)
		}
	} else if response != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("unexpected response type: %T", response))
	}
}

// GetApiInfo operation middleware
func (sh *strictHandler) GetApiInfo(w http.ResponseWriter, r *http.Request) {
	var request GetApiInfoRequestObject

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (interface{}, error) {
		return sh.ssi.GetApiInfo(ctx, request.(GetApiInfoRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "GetApiInfo")
	}

	response, err := handler(r.Context(), w, r, request)

	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(GetApiInfoResponseObject); ok {
		if err := validResponse.VisitGetApiInfoResponse(w); err != nil {
			sh.options.ResponseErrorHandlerFunc(w, r, err)
		}
	} else if response != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("unexpected response type: %T", response))
	}
}

// PostApiSendCoin operation middleware
func (sh *strictHandler) PostApiSendCoin(w http.ResponseWriter, r *http.Request) {
	var request PostApiSendCoinRequestObject

	var body PostApiSendCoinJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		sh.options.RequestErrorHandlerFunc(w, r, fmt.Errorf("can't decode JSON body: %w", err))
		return
	}
	request.Body = &body

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (interface{}, error) {
		return sh.ssi.PostApiSendCoin(ctx, request.(PostApiSendCoinRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "PostApiSendCoin")
	}

	response, err := handler(r.Context(), w, r, request)

	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(PostApiSendCoinResponseObject); ok {
		if err := validResponse.VisitPostApiSendCoinResponse(w); err != nil {
			sh.options.ResponseErrorHandlerFunc(w, r, err)
		}
	} else if response != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("unexpected response type: %T", response))
	}
}
