package health

import "github.com/gin-gonic/gin"

type HealthChecks struct {
	Mongo    func() error
	RabbitMQ func() error
}

type Handler struct {
	health HealthChecks
}

func NewHandler(health HealthChecks) *Handler {
	return &Handler{
		health: health,
	}
}

func (h *Handler) Health(c *gin.Context) {
	type serviceStatus struct {
		Status string `json:"status"`
		Error  string `json:"error,omitempty"`
	}

	response := map[string]interface{}{
		"status": "ok",
		"checks": map[string]serviceStatus{},
	}

	checks := response["checks"].(map[string]serviceStatus)

	// Mongo
	if err := h.health.Mongo(); err != nil {
		checks["mongodb"] = serviceStatus{
			Status: "down",
			Error:  err.Error(),
		}
		response["status"] = "degraded"
	} else {
		checks["mongodb"] = serviceStatus{Status: "up"}
	}

	// RabbitMQ
	if err := h.health.RabbitMQ(); err != nil {
		checks["rabbitmq"] = serviceStatus{
			Status: "down",
			Error:  err.Error(),
		}
		response["status"] = "degraded"
	} else {
		checks["rabbitmq"] = serviceStatus{Status: "up"}
	}

	// HTTP status
	if response["status"] == "ok" {
		c.JSON(200, response)
	} else {
		c.JSON(503, response)
	}
}
