package schema

type UserDealLink struct {
	UserId int `db:"user_id"`
	DealId int `db:"deal_id"`
}
