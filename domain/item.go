package domain

type Category string

const (
	Foods          Category = "食品"
	HouseholdGoods Category = "日用品"
)

type Item struct {
	Category Category `json:"category"`
	Amount   int      `json:"amount"`
	Name     string   `json:"name"`
}
