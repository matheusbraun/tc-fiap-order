package controller

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	orderController "github.com/viniciuscluna/tc-fiap-50/internal/order/controller"
	"github.com/viniciuscluna/tc-fiap-50/internal/order/infrastructure/api/dto"
)

type orderApiController struct {
	controller orderController.OrderController
}

func NewOrderController(controller orderController.OrderController) *orderApiController {
	return &orderApiController{
		controller: controller,
	}
}

func (c *orderApiController) RegisterRoutes(r chi.Router) {
	prefix := "/v1/order"
	r.Post(prefix, c.Add)
	r.Get(prefix+"/{orderId}", c.GetOrder)
	r.Get(prefix, c.GetOrders)
	r.Get(prefix+"/{orderId}/status", c.GetOrderStatus)
	r.Put(prefix+"/{orderId}/status", c.UpdateOrderStatus)
}

// @Summary     Add order
// @Description Add order
// @Tags        Order
// @Accept      json
// @Produce     json
// @Param       body body dto.AddOrderDto true "Body"
// @Success     201  {object} dto.GetOrderResponseDto
// @Router      /v1/order [post]
func (c *orderApiController) Add(w http.ResponseWriter, r *http.Request) {
	var orderRequest dto.AddOrderDto

	if err := json.NewDecoder(r.Body).Decode(&orderRequest); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
	}

	orderId, err := c.controller.Add(&orderRequest)

	if err != nil {
		http.Error(w, "Error processing request", http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(orderId)
}

// @Summary     Get order
// @Description Get order
// @Tags        Order
// @Accept      json
// @Produce     json
// @Param       orderId path uint true "Order ID"
// @Success     200  {object} dto.GetOrderResponseDto
// @Router      /v1/order/{orderId} [get]
func (c *orderApiController) GetOrder(w http.ResponseWriter, r *http.Request) {
	orderId, err := getOrderIDFromPath(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	order, err := c.controller.GetOrder(orderId)

	if err != nil {
		http.Error(w, "Error processing request", http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(order)
}

// @Summary     Get orders
// @Description Get orders
// @Tags        Order
// @Accept      json
// @Produce     json
// @Success     200  {object} dto.GetOrdersResponseDto
// @Router      /v1/order [get]
func (c *orderApiController) GetOrders(w http.ResponseWriter, r *http.Request) {
	orders, err := c.controller.GetOrders()

	if err != nil {
		http.Error(w, "Error processing request", http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(orders)
}

// @Summary     Get order status
// @Description Get order status
// @Tags        Order
// @Accept      json
// @Produce     json
// @Param       orderId path uint true "Order ID"
// @Success     200  {object} dto.GetOrderStatusResponseDto
// @Router      /v1/order/{orderId}/status [get]
func (c *orderApiController) GetOrderStatus(w http.ResponseWriter, r *http.Request) {
	orderId, err := getOrderIDFromPath(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	status, err := c.controller.GetOrderStatus(orderId)

	if err != nil {
		http.Error(w, "Error processing request", http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(status)
}

// @Summary     Update order status
// @Description Update order status
// @Tags        Order
// @Accept      json
// @Produce     json
// @Param       orderId path uint true "Order ID"
// @Param       status body dto.UpdateOrderStatusRequestDto true "Status"
// @Success     200
// @Router      /v1/order/{orderId}/status [put]
func (c *orderApiController) UpdateOrderStatus(w http.ResponseWriter, r *http.Request) {
	orderId, err := getOrderIDFromPath(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var statusRequest dto.UpdateOrderStatusRequestDto

	if err := json.NewDecoder(r.Body).Decode(&statusRequest); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
	}

	err = c.controller.UpdateOrderStatus(orderId, &statusRequest)

	if err != nil {
		http.Error(w, "Error processing request", http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
}

func getOrderIDFromPath(r *http.Request) (uint, error) {
	vars := chi.URLParam(r, "orderId")
	id, err := strconv.ParseUint(vars, 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(id), nil
}
