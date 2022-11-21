package domain

type BidRepository interface {
	GetBid(bidId string) *Bid
	GetBidsByUserId(userId string) *[]*Bid
	GetBidsByItemId(itemId string) *[]*Bid
	SaveBid(bid *Bid)
	SaveBids(bids *[]*Bid)
	DeleteBid(bidId string)
	NextBidId() string
}
