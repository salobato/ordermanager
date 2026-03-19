package usecase_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"log/slog"
	"strings"
	"testing"

	"github.com/salobato/ordermanager/internal/core/entity"
	"github.com/salobato/ordermanager/internal/core/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ─── Helpers ─────────────────────────────────────────────────────────────────

// testLogger returns a logger that writes JSON lines to a buffer so tests
// can assert on what was actually logged.
func testLogger(buf *bytes.Buffer) *slog.Logger {
	return slog.New(slog.NewJSONHandler(buf, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
}

// logLines parses every JSON line written to buf into a slice of maps.
func logLines(buf *bytes.Buffer) []map[string]interface{} {
	var lines []map[string]interface{}
	for _, raw := range strings.Split(strings.TrimSpace(buf.String()), "\n") {
		if raw == "" {
			continue
		}
		var entry map[string]interface{}
		if err := json.Unmarshal([]byte(raw), &entry); err == nil {
			lines = append(lines, entry)
		}
	}
	return lines
}

func stubOrder() *entity.Order {
	return &entity.Order{
		ID:          "order-1",
		OrderNumber: "ORD-2024-000001",
		Status:      entity.OrderCreated,
		CustomerID:  "customer-1",
		Total:       100.0,
	}
}

// ─── WithPlaceOrderLogging ───────────────────────────────────────────────────

func TestWithPlaceOrderLogging_Success(t *testing.T) {
	var buf bytes.Buffer
	decorated := usecase.WithPlaceOrderLogging(
		func(input usecase.PlaceOrderInput) (*entity.Order, error) {
			return stubOrder(), nil
		},
		testLogger(&buf),
	)

	order, err := decorated(usecase.PlaceOrderInput{CustomerID: "customer-1", Total: 100.0})
	require.NoError(t, err)
	require.NotNil(t, order)

	lines := logLines(&buf)
	require.Len(t, lines, 2, "expected exactly two log lines: starting + completed")

	start := lines[0]
	assert.Equal(t, "INFO", start["level"])
	assert.Equal(t, "starting", start["msg"])
	assert.Equal(t, "PlaceOrder", start["use_case"])
	assert.Equal(t, "customer-1", start["customer_id"])
	assert.Equal(t, 100.0, start["total"])

	done := lines[1]
	assert.Equal(t, "INFO", done["level"])
	assert.Equal(t, "completed", done["msg"])
	assert.Equal(t, "order-1", done["order_id"])
	assert.Equal(t, "ORD-2024-000001", done["order_number"])
	assert.Equal(t, string(entity.OrderCreated), done["status"])
	assert.NotNil(t, done["duration"])
}

func TestWithPlaceOrderLogging_Error(t *testing.T) {
	var buf bytes.Buffer
	decorated := usecase.WithPlaceOrderLogging(
		func(input usecase.PlaceOrderInput) (*entity.Order, error) {
			return nil, errors.New("validation failed")
		},
		testLogger(&buf),
	)

	order, err := decorated(usecase.PlaceOrderInput{CustomerID: "customer-1", Total: 100.0})
	assert.Error(t, err)
	assert.Nil(t, order)

	lines := logLines(&buf)
	require.Len(t, lines, 2, "expected exactly two log lines: starting + failed")

	assert.Equal(t, "ERROR", lines[1]["level"])
	assert.Equal(t, "failed", lines[1]["msg"])
	assert.Equal(t, "validation failed", lines[1]["error"])
	assert.NotNil(t, lines[1]["duration"])
}

func TestWithPlaceOrderLogging_PassesThroughInput(t *testing.T) {
	var captured usecase.PlaceOrderInput
	var buf bytes.Buffer
	decorated := usecase.WithPlaceOrderLogging(
		func(input usecase.PlaceOrderInput) (*entity.Order, error) {
			captured = input
			return stubOrder(), nil
		},
		testLogger(&buf),
	)

	input := usecase.PlaceOrderInput{CustomerID: "customer-99", Total: 49.99}
	_, err := decorated(input)
	require.NoError(t, err)
	assert.Equal(t, input, captured)
}

// ─── WithUpdateOrderStatusLogging ────────────────────────────────────────────

func TestWithUpdateOrderStatusLogging_Success(t *testing.T) {
	var buf bytes.Buffer
	decorated := usecase.WithUpdateOrderStatusLogging(
		func(input usecase.UpdateOrderStatusInput) (*entity.Order, error) {
			return &entity.Order{
				ID:          input.OrderID,
				OrderNumber: "ORD-2024-000001",
				Status:      entity.OrderProcessing,
			}, nil
		},
		testLogger(&buf),
	)

	order, err := decorated(usecase.UpdateOrderStatusInput{
		OrderID: "order-1",
		Status:  string(entity.OrderProcessing),
	})
	require.NoError(t, err)
	require.NotNil(t, order)

	lines := logLines(&buf)
	require.Len(t, lines, 2)

	start := lines[0]
	assert.Equal(t, "starting", start["msg"])
	assert.Equal(t, "UpdateOrderStatus", start["use_case"])
	assert.Equal(t, "order-1", start["order_id"])
	assert.Equal(t, string(entity.OrderProcessing), start["requested_status"])

	done := lines[1]
	assert.Equal(t, "completed", done["msg"])
	assert.Equal(t, string(entity.OrderProcessing), done["status"])
	assert.NotNil(t, done["duration"])
}

func TestWithUpdateOrderStatusLogging_Error(t *testing.T) {
	var buf bytes.Buffer
	decorated := usecase.WithUpdateOrderStatusLogging(
		func(input usecase.UpdateOrderStatusInput) (*entity.Order, error) {
			return nil, errors.New("transição inválida")
		},
		testLogger(&buf),
	)

	_, err := decorated(usecase.UpdateOrderStatusInput{OrderID: "order-1", Status: "invalido"})
	assert.Error(t, err)

	lines := logLines(&buf)
	require.Len(t, lines, 2)
	assert.Equal(t, "ERROR", lines[1]["level"])
	assert.Equal(t, "transição inválida", lines[1]["error"])
}

// ─── WithFindByIDLogging ──────────────────────────────────────────────────────

func TestWithFindByIDLogging_Success(t *testing.T) {
	var buf bytes.Buffer
	decorated := usecase.WithFindByIDLogging(
		func(id string) (*entity.Order, error) {
			return &entity.Order{
				ID:          id,
				OrderNumber: "ORD-2024-000001",
				Status:      entity.OrderCreated,
			}, nil
		},
		testLogger(&buf),
	)

	order, err := decorated("order-1")
	require.NoError(t, err)
	require.NotNil(t, order)

	lines := logLines(&buf)
	require.Len(t, lines, 2)

	start := lines[0]
	assert.Equal(t, "starting", start["msg"])
	assert.Equal(t, "FindByID", start["use_case"])
	assert.Equal(t, "order-1", start["order_id"])

	done := lines[1]
	assert.Equal(t, "completed", done["msg"])
	assert.Equal(t, "ORD-2024-000001", done["order_number"])
	assert.NotNil(t, done["duration"])
}

func TestWithFindByIDLogging_NotFound(t *testing.T) {
	var buf bytes.Buffer
	decorated := usecase.WithFindByIDLogging(
		func(id string) (*entity.Order, error) {
			return nil, errors.New("pedido não encontrado")
		},
		testLogger(&buf),
	)

	_, err := decorated("non-existent")
	assert.Error(t, err)

	lines := logLines(&buf)
	require.Len(t, lines, 2)
	assert.Equal(t, "ERROR", lines[1]["level"])
	assert.Equal(t, "pedido não encontrado", lines[1]["error"])
	assert.Equal(t, "non-existent", lines[1]["order_id"])
}

func TestWithFindByIDLogging_PassesThroughID(t *testing.T) {
	var capturedID string
	var buf bytes.Buffer
	decorated := usecase.WithFindByIDLogging(
		func(id string) (*entity.Order, error) {
			capturedID = id
			return stubOrder(), nil
		},
		testLogger(&buf),
	)

	_, err := decorated("order-42")
	require.NoError(t, err)
	assert.Equal(t, "order-42", capturedID)
}
