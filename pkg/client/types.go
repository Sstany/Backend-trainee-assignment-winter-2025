// Package client provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/oapi-codegen/oapi-codegen/v2 version v2.4.1 DO NOT EDIT.
package client

const (
	BearerAuthScopes = "BearerAuth.Scopes"
)

// AuthRequest defines model for AuthRequest.
type AuthRequest struct {
	// Password Пароль для аутентификации.
	Password string `json:"password"`

	// Username Имя пользователя для аутентификации.
	Username string `json:"username"`
}

// AuthResponse defines model for AuthResponse.
type AuthResponse struct {
	// Token JWT-токен для доступа к защищенным ресурсам.
	Token *string `json:"token,omitempty"`
}

// ErrorResponse defines model for ErrorResponse.
type ErrorResponse struct {
	// Errors Сообщение об ошибке, описывающее проблему.
	Errors *string `json:"errors,omitempty"`
}

// InfoResponse defines model for InfoResponse.
type InfoResponse struct {
	CoinHistory *struct {
		Received *[]struct {
			// Amount Количество полученных монет.
			Amount *int `json:"amount,omitempty"`

			// FromUser Имя пользователя, который отправил монеты.
			FromUser *string `json:"fromUser,omitempty"`
		} `json:"received,omitempty"`
		Sent *[]struct {
			// Amount Количество отправленных монет.
			Amount *int `json:"amount,omitempty"`

			// ToUser Имя пользователя, которому отправлены монеты.
			ToUser *string `json:"toUser,omitempty"`
		} `json:"sent,omitempty"`
	} `json:"coinHistory,omitempty"`

	// Coins Количество доступных монет.
	Coins     *int `json:"coins,omitempty"`
	Inventory *[]struct {
		// Quantity Количество предметов.
		Quantity *int `json:"quantity,omitempty"`

		// Type Тип предмета.
		Type *string `json:"type,omitempty"`
	} `json:"inventory,omitempty"`
}

// SendCoinRequest defines model for SendCoinRequest.
type SendCoinRequest struct {
	// Amount Количество монет, которые необходимо отправить.
	Amount int `json:"amount"`

	// ToUser Имя пользователя, которому нужно отправить монеты.
	ToUser string `json:"toUser"`
}

// PostApiAuthJSONRequestBody defines body for PostApiAuth for application/json ContentType.
type PostApiAuthJSONRequestBody = AuthRequest

// PostApiSendCoinJSONRequestBody defines body for PostApiSendCoin for application/json ContentType.
type PostApiSendCoinJSONRequestBody = SendCoinRequest
