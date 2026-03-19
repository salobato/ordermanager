package order_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	api "github.com/salobato/ordermanager/internal/adapter/api/gin/order"
	"github.com/salobato/ordermanager/internal/core/entity"
	"github.com/salobato/ordermanager/internal/core/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupRouter(uc api.UseCases) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := api.NewHandler(uc)
	r.POST("/orders", h.PlaceOrder)
	r.PATCH("/orders/:id/status", h.UpdateOrderStatus)
	r.GET("/orders/:id", h.FindByID)
	return r
}

func stubOrder() *entity.Order {
	return &entity.Order{
		ID:          "order-1",
		OrderNumber: "ORD-2024-000001",
		Status:      entity.OrderCreated,
		CustomerID:  "customer-1",
		Total:       100.0,
		PlacedAt:    time.Now(),
		UpdatedAt:   time.Now(),
	}
}

func doRequest(router *gin.Engine, method, path string, body interface{}) *httptest.ResponseRecorder {
	var buf bytes.Buffer
	if body != nil {
		json.NewEncoder(&buf).Encode(body)
	}
	req, _ := http.NewRequest(method, path, &buf)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}

func TestPlaceOrder_Success(t *testing.T) {
	order := stubOrder()
	router := setupRouter(api.UseCases{
		PlaceOrder: func(input usecase.PlaceOrderInput) (*entity.Order, error) {
			return order, nil
		},
	})

	w := doRequest(router, http.MethodPost, "/orders", map[string]interface{}{
		"customer_id": "customer-1",
		"total":       100.0,
	})

	assert.Equal(t, http.StatusCreated, w.Code)

	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, "order-1", resp["id"])
	assert.Equal(t, "ORD-2024-000001", resp["order_number"])
	assert.Equal(t, string(entity.OrderCreated), resp["status"])
}

func TestPlaceOrder_InvalidJSON(t *testing.T) {
	router := setupRouter(api.UseCases{})

	req, _ := http.NewRequest(http.MethodPost, "/orders", bytes.NewBufferString("not-json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Contains(t, resp, "error")
}

func TestPlaceOrder_MissingRequiredFields(t *testing.T) {
	router := setupRouter(api.UseCases{
		PlaceOrder: func(input usecase.PlaceOrderInput) (*entity.Order, error) {
			return nil, errors.New("O ID do cliente não pode ser vazio")
		},
	})

	w := doRequest(router, http.MethodPost, "/orders", map[string]interface{}{
		"total": 100.0,
	})

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Contains(t, resp, "error")
}

func TestPlaceOrder_UseCaseError(t *testing.T) {
	router := setupRouter(api.UseCases{
		PlaceOrder: func(input usecase.PlaceOrderInput) (*entity.Order, error) {
			return nil, errors.New("use case failure")
		},
	})

	w := doRequest(router, http.MethodPost, "/orders", map[string]interface{}{
		"customer_id": "customer-1",
		"total":       100.0,
	})

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, "use case failure", resp["error"])
}

func TestUpdateOrderStatus_Success(t *testing.T) {
	router := setupRouter(api.UseCases{
		UpdateOrderStatus: func(input usecase.UpdateOrderStatusInput) (*entity.Order, error) {
			return &entity.Order{
				ID:     input.OrderID,
				Status: entity.OrderProcessing,
			}, nil
		},
	})

	w := doRequest(router, http.MethodPatch, "/orders/order-1/status", map[string]interface{}{
		"status": string(entity.OrderProcessing),
	})

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, string(entity.OrderProcessing), resp["status"])
}

func TestUpdateOrderStatus_InvalidJSON(t *testing.T) {
	router := setupRouter(api.UseCases{})

	req, _ := http.NewRequest(http.MethodPatch, "/orders/order-1/status", bytes.NewBufferString("not-json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdateOrderStatus_UseCaseError(t *testing.T) {
	router := setupRouter(api.UseCases{
		UpdateOrderStatus: func(input usecase.UpdateOrderStatusInput) (*entity.Order, error) {
			return nil, errors.New("Transição de status inválida")
		},
	})

	w := doRequest(router, http.MethodPatch, "/orders/order-1/status", map[string]interface{}{
		"status": string(entity.OrderDelivered),
	})

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, "Transição de status inválida", resp["error"])
}

func TestUpdateOrderStatus_ForwardsIDFromURL(t *testing.T) {
	var capturedInput usecase.UpdateOrderStatusInput
	router := setupRouter(api.UseCases{
		UpdateOrderStatus: func(input usecase.UpdateOrderStatusInput) (*entity.Order, error) {
			capturedInput = input
			return &entity.Order{ID: input.OrderID, Status: entity.OrderProcessing}, nil
		},
	})

	w := doRequest(router, http.MethodPatch, "/orders/order-99/status", map[string]interface{}{
		"status": string(entity.OrderProcessing),
	})

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "order-99", capturedInput.OrderID)
}

// ─── FindByID ────────────────────────────────────────────────────────────────

func TestFindByID_Success(t *testing.T) {
	order := stubOrder()
	router := setupRouter(api.UseCases{
		FindByID: func(id string) (*entity.Order, error) {
			return order, nil
		},
	})

	w := doRequest(router, http.MethodGet, "/orders/order-1", nil)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, "order-1", resp["id"])
}

func TestFindByID_NotFound(t *testing.T) {
	router := setupRouter(api.UseCases{
		FindByID: func(id string) (*entity.Order, error) {
			return nil, errors.New("Pedido não encontrado")
		},
	})

	w := doRequest(router, http.MethodGet, "/orders/non-existent", nil)

	assert.Equal(t, http.StatusNotFound, w.Code)

	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, "Pedido não encontrado", resp["error"])
}

func TestFindByID_ForwardsIDFromURL(t *testing.T) {
	var capturedID string
	router := setupRouter(api.UseCases{
		FindByID: func(id string) (*entity.Order, error) {
			capturedID = id
			return &entity.Order{ID: id}, nil
		},
	})

	doRequest(router, http.MethodGet, "/orders/order-42", nil)
	assert.Equal(t, "order-42", capturedID)
}
