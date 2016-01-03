package model

type Camera struct {
	ID        int    `json:"id"`
	MinChange int    `json:"min_change"`
	Name      string `json:"name"`
	Password  string `json:"password"`
	Threshold int    `json:"threshold"`
	URL       string `json:"url"`
	Username  string `json:"username"`
}
