package dto

type TransferInquiryReq struct {
	AccountNumber string  `json:"account_number"`
	Amount        float64 `json:"amount"`
	UserID        int64   `json:"user_id,omitempty"`
}
