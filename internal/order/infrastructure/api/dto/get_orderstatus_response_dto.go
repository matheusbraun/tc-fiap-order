package dto

type GetOrderStatusResponseDto struct {
	ID                       uint   `json:"id"`
	CreatedAt                string `json:"created_at"`
	CurrentStatus            uint   `json:"current_status"`
	CurrentStatusDescription string `json:"current_status_description"`
	OrderId                  uint   `json:"order_id"`
}
