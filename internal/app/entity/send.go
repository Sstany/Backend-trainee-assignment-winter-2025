package entity

type SendCoinRequest struct {
	// Amount Количество монет, которые необходимо отправить.
	Amount   int `json:"amount"`
	FromUser string
	// ToUser Имя пользователя, которому нужно отправить монеты.
	ToUser string `json:"toUser"`
}
