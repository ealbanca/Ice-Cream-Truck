package models

type CartItem struct {
	ID            int
	ProductName   string
	SizeLabel     string
	FlavorLabels  string
	ToppingLabels string
	Quantity      int
	TotalPrice    float64
}

// Example: In-memory cart for demonstration (replace with DB logic as needed)
var cartItems []CartItem

func GetCartItems(userID int) []CartItem {
	// TODO: Replace with DB fetch by userID
	return cartItems
}

func AddToCart(item CartItem) {
	cartItems = append(cartItems, item)
}

func RemoveFromCart(itemID int) {
	for i, item := range cartItems {
		if item.ID == itemID {
			cartItems = append(cartItems[:i], cartItems[i+1:]...)
			break
		}
	}
}

func GetCartTotal() float64 {
	total := 0.0
	for _, item := range cartItems {
		total += item.TotalPrice * float64(item.Quantity)
	}
	return total
}
