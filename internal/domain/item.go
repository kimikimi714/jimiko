package domain

// Category is a category of item in shopping list.
type Category string

const (
	// Foods is 食品
	Foods          Category = "食品"
	// HouseholdGoods is 日用品
	HouseholdGoods Category = "日用品"
)

// Item is one item in shopping list.
type Item struct {
	Category Category `json:"category"`
	Amount   int      `json:"amount"`
	Name     string   `json:"name"`
}
