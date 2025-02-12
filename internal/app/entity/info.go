package entity

// type InfoResponse struct {
// 	CoinHistory *struct {
// 		Received *[]struct {
// 			// Amount Количество полученных монет.
// 			Amount *int `json:"amount,omitempty"`

// 			// FromUser Имя пользователя, который отправил монеты.
// 			FromUser *string `json:"fromUser,omitempty"`
// 		} `json:"received,omitempty"`
// 		Sent *[]struct {
// 			// Amount Количество отправленных монет.
// 			Amount *int `json:"amount,omitempty"`

// 			// ToUser Имя пользователя, которому отправлены монеты.
// 			ToUser *string `json:"toUser,omitempty"`
// 		} `json:"sent,omitempty"`
// 	} `json:"coinHistory,omitempty"`

// 	// Coins Количество доступных монет.
// 	Coins     *int `json:"coins,omitempty"`
// 	Inventory *[]struct {
// 		// Quantity Количество предметов.
// 		Quantity *int `json:"quantity,omitempty"`

// 		// Type Тип предмета.
// 		Type *string `json:"type,omitempty"`
// 	} `json:"inventory,omitempty"`
// }

type UserInfo struct {
	CoinHistory *CoinHistory
	Coins       *[]int
	Inventory   *Inventory
}

type CoinHistory struct {
}

type Inventory struct {
}
