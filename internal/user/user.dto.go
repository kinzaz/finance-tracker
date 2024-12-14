package user

type UserProfileResponseDto struct {
	Email   string  `json:"email"`
	Name    string  `json:"name"`
	Balance float64 `json:"balance"`
}
