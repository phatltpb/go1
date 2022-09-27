package model

type AuthRequest struct {
	ID          uint64 `json:"id"`
	UID         uint64 `json:"uid"`
	PhoneNumber string `json:"phone_number"`
	Status      string `json:"status"`
	Expired_at  int64  `json:"expired_at"`
}
