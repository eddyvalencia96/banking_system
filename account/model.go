package account

type Account struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	InitialBalance int    `json:"initial_balance"`
	Balance        int    `json:"balance"`
}

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
