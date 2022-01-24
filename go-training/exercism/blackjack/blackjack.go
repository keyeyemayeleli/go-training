package blackjack

// ParseCard returns the integer value of a card following blackjack ruleset.
func ParseCard(card string) int {
	cards := make(map[string]int)
	cards["ace"] = 11
	cards["two"] = 2
	cards["three"] = 3
	cards["four"] = 4
	cards["five"] = 5
	cards["six"] = 6
	cards["seven"] = 7
	cards["eight"] = 8
	cards["nine"] = 9
	cards["ten"] = 10
	cards["jack"] = 10
	cards["queen"] = 10
	cards["king"] = 10
	cards["other"] = 0

	return cards[card]
}

// IsBlackjack returns true if the player has a blackjack, false otherwise.
func IsBlackjack(card1, card2 string) bool {
	if ParseCard(card1)+ParseCard(card2) == 21 {
		return true
	} else {
		return false
	}
}

// LargeHand implements the decision tree for hand scores larger than 20 points.
func LargeHand(isBlackjack bool, dealerScore int) string {
	var o string
	switch {
	case isBlackjack && dealerScore < 10:
		o = "W"
	case isBlackjack && dealerScore >= 10:
		o = "S"
	default:
		o = "P"
	}
	return o
}

// SmallHand implements the decision tree for hand scores with less than 21 points.
func SmallHand(handScore, dealerScore int) string {
	var o string
	switch {
	case handScore >= 17:
		o = "S"
	case handScore <= 11:
		o = "H"
	case handScore >= 12 && handScore <= 16 && dealerScore >= 7:
		o = "H"
	case handScore >= 12 && handScore <= 16 && dealerScore < 7:
		o = "S"
	}
	return o
}
