package entity

type Response struct {
	CoinHistory *struct {
		Received *[]struct {
			// Amount Количество полученных монет.
			Amount int

			// FromUser Имя пользователя, который отправил монеты.
			FromUser *string
		}
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

type InfoResponse struct {
	CoinHistory CoinHistory
	Coins       int
	Inventory   []Inventory
}

type InfoRequest struct {
	Username string
}
type CoinHistory struct {
	Received []Received
	Sent     []Sent
}

type Received struct {
	Amount   int
	FromUser string
}
type Sent struct {
	Amount int
	ToUser string
}

type Coins int
type Inventory struct {
	Quantity int
	Type     string
}
