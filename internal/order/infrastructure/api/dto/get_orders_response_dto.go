package dto

type GetOrdersResponseDto struct {
	Orders []*GetOrderResponseDto `json:"orders"`
}
