package order

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/salobato/ordermanager/internal/core/entity"
	"github.com/salobato/ordermanager/internal/core/usecase"
)

type UseCases struct {
	PlaceOrder        func(usecase.PlaceOrderInput) (*entity.Order, error)
	UpdateOrderStatus func(usecase.UpdateOrderStatusInput) (*entity.Order, error)
	FindByID          func(string) (*entity.Order, error)
}

type Handler struct {
	uc UseCases
}

func NewHandler(uc UseCases) *Handler {
	return &Handler{uc: uc}
}

func (h *Handler) PlaceOrder(c *gin.Context) {
	var req PlaceOrderRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	input := usecase.PlaceOrderInput{
		CustomerID: req.CustomerID,
		Total:      req.Total,
	}

	order, err := h.uc.PlaceOrder(input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, toOrderResponse(order))
}

func (h *Handler) UpdateOrderStatus(c *gin.Context) {
	id := c.Param("id")

	var req UpdateStatusRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	input := usecase.UpdateOrderStatusInput{
		OrderID: id,
		Status:  req.Status,
	}

	order, err := h.uc.UpdateOrderStatus(input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, toOrderResponse(order))
}

func (h *Handler) FindByID(c *gin.Context) {
	id := c.Param("id")

	order, err := h.uc.FindByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, toOrderResponse(order))
}
