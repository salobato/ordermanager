package order

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.Engine, h *Handler) {
	orders := r.Group("/orders")

	orders.POST("", h.PlaceOrder)
	orders.PATCH("/:id/status", h.UpdateOrderStatus)
	orders.GET("/:id", h.FindByID)
}
