package entity_test

import (
	"testing"
	"time"

	"github.com/salobato/ordermanager/internal/core/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewOrderNumber_Format(t *testing.T) {
	on := entity.NewOrderNumber(1)
	currentYear := time.Now().Year()

	assert.True(t, on.IsValid(), "O número de pedido novo deve ser válido")
	assert.Equal(t, currentYear, on.Year(), "O ano extraído deve ser igual ao ano atual")
	assert.Equal(t, 1, on.Sequence(), "A sequência extraída deve ser igual à sequência recebida por parâmetro")
}

func TestNewOrderNumber_SequencePadding(t *testing.T) {
	cases := []struct {
		seq      int64
		expected int
	}{
		{1, 1},
		{42, 42},
		{999999, 999999},
	}

	for _, tc := range cases {
		on := entity.NewOrderNumber(tc.seq)
		assert.True(t, on.IsValid())
		assert.Equal(t, tc.expected, on.Sequence())
	}
}

func TestOrderNumber_String(t *testing.T) {
	on := entity.NewOrderNumber(7)
	assert.NotEmpty(t, on.String())
	assert.Equal(t, string(on), on.String())
}

func TestOrderNumber_IsValid(t *testing.T) {
	year := time.Now().Year()

	validCases := []entity.OrderNumber{
		entity.NewOrderNumber(1),
		entity.NewOrderNumber(999999),
	}
	for _, v := range validCases {
		assert.True(t, v.IsValid(), "expected %q to be valid", v)
	}

	invalidCases := []entity.OrderNumber{
		"",
		"ORD-",
		"ORD-ABCD-000001",  // Ano não numérico
		"ORD-2024-12345",   // Sequência muito curta (5 dígitos)
		"ORD-2024-1234567", // Sequência muito longa (7 dígitos)
		"ord-2024-000001",  // Prefixo em letras minúsculas
		entity.OrderNumber("ORD-" + string(rune(year)) + "-000001"), // Ano em formato incorreto
	}
	for _, v := range invalidCases {
		assert.False(t, v.IsValid(), "Esperado que %q seja inválido", v)
	}
}

func TestOrderNumber_Year(t *testing.T) {
	on := entity.NewOrderNumber(1)
	assert.Equal(t, time.Now().Year(), on.Year())
}

func TestOrderNumber_Sequence(t *testing.T) {
	cases := []int64{1, 100, 999999}
	for _, seq := range cases {
		on := entity.NewOrderNumber(seq)
		assert.Equal(t, int(seq), on.Sequence())
	}
}

func TestOrderStatus_IsValid_AllStatuses(t *testing.T) {
	valid := []entity.OrderStatus{
		entity.OrderCreated,
		entity.OrderProcessing,
		entity.OrderShipped,
		entity.OrderDelivered,
	}
	for _, s := range valid {
		assert.True(t, s.IsValid(), "Esperado que %q seja válido", s)
	}
}

func TestOrderStatus_IsValid_Invalid(t *testing.T) {
	invalid := []entity.OrderStatus{"", "invalido", "CRIADO", "shipped"}
	for _, s := range invalid {
		assert.False(t, s.IsValid(), "Esperado que %q seja inválido", s)
	}
}

func TestOrderStatus_CanTransitionTo(t *testing.T) {
	cases := []struct {
		from     entity.OrderStatus
		to       entity.OrderStatus
		expected bool
	}{
		// Transições válidas
		{entity.OrderCreated, entity.OrderProcessing, true},
		{entity.OrderProcessing, entity.OrderShipped, true},
		{entity.OrderShipped, entity.OrderDelivered, true},

		// Transições inválidas (pulando etapas)
		{entity.OrderCreated, entity.OrderShipped, false},
		{entity.OrderCreated, entity.OrderDelivered, false},
		{entity.OrderProcessing, entity.OrderDelivered, false},

		// Transições inválidas (voltando etapas)
		{entity.OrderProcessing, entity.OrderCreated, false},
		{entity.OrderShipped, entity.OrderCreated, false},
		{entity.OrderShipped, entity.OrderProcessing, false},
		{entity.OrderDelivered, entity.OrderShipped, false},

		// "entregue" é terminal
		{entity.OrderDelivered, entity.OrderCreated, false},
		{entity.OrderDelivered, entity.OrderProcessing, false},
		{entity.OrderDelivered, entity.OrderDelivered, false},
	}

	for _, tc := range cases {
		result := tc.from.CanTransitionTo(tc.to)
		assert.Equal(t, tc.expected, result, "%s -> %s", tc.from, tc.to)
	}
}

func TestOrder_Validate_Valid(t *testing.T) {
	o := entity.Order{
		CustomerID: "customer-1",
		Total:      99.99,
	}
	assert.NoError(t, o.Validate())
}

func TestOrder_Validate_EmptyCustomerID(t *testing.T) {
	o := entity.Order{
		CustomerID: "",
		Total:      99.99,
	}
	err := o.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cliente")
}

func TestOrder_Validate_ZeroTotal(t *testing.T) {
	o := entity.Order{
		CustomerID: "customer-1",
		Total:      0,
	}
	err := o.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "total")
}

func TestOrder_Validate_NegativeTotal(t *testing.T) {
	o := entity.Order{
		CustomerID: "customer-1",
		Total:      -10.0,
	}
	err := o.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "total")
}

func TestNewOrder_Success(t *testing.T) {
	before := time.Now().Truncate(time.Second)
	order, err := entity.NewOrder("customer-1", 150.0, 1)
	after := time.Now().Add(time.Second)

	require.NoError(t, err)
	require.NotNil(t, order)

	assert.NotEmpty(t, order.CustomerID)
	assert.Equal(t, "customer-1", order.CustomerID)
	assert.Equal(t, 150.0, order.Total)
	assert.Equal(t, entity.OrderCreated, order.Status)
	assert.True(t, order.OrderNumber.IsValid())

	assert.False(t, order.PlacedAt.Before(before), "PlacedAt deveria ser >= que before")
	assert.False(t, order.PlacedAt.After(after), "PlacedAt deveria ser <= que after")
	assert.False(t, order.UpdatedAt.Before(before))
	assert.False(t, order.UpdatedAt.After(after))
}

func TestNewOrder_EmptyCustomerID(t *testing.T) {
	order, err := entity.NewOrder("", 100.0, 1)
	assert.Error(t, err)
	assert.Nil(t, order)
}

func TestNewOrder_ZeroTotal(t *testing.T) {
	order, err := entity.NewOrder("customer-1", 0, 1)
	assert.Error(t, err)
	assert.Nil(t, order)
}

func TestNewOrder_NegativeTotal(t *testing.T) {
	order, err := entity.NewOrder("customer-1", -50.0, 1)
	assert.Error(t, err)
	assert.Nil(t, order)
}

func TestOrder_ChangeStatus_FullLifecycle(t *testing.T) {
	order, err := entity.NewOrder("customer-1", 100.0, 1)
	require.NoError(t, err)

	transitions := []entity.OrderStatus{
		entity.OrderProcessing,
		entity.OrderShipped,
		entity.OrderDelivered,
	}

	for _, next := range transitions {
		prevUpdatedAt := order.UpdatedAt
		time.Sleep(time.Millisecond)

		err := order.ChangeStatus(next)
		assert.NoError(t, err, "Transição para %s deve ser bem sucedida", next)
		assert.Equal(t, next, order.Status)
		assert.True(t, order.UpdatedAt.After(prevUpdatedAt), "UpdatedAt deve ser atualizado")
	}
}

func TestOrder_ChangeStatus_InvalidTransition(t *testing.T) {
	order, err := entity.NewOrder("customer-1", 100.0, 1)
	require.NoError(t, err)

	err = order.ChangeStatus(entity.OrderDelivered)
	assert.Error(t, err)
	assert.Equal(t, entity.OrderCreated, order.Status, "Status não deve mudar se erro")
}

func TestOrder_ChangeStatus_TerminalState(t *testing.T) {
	order, err := entity.NewOrder("customer-1", 100.0, 1)
	require.NoError(t, err)

	require.NoError(t, order.ChangeStatus(entity.OrderProcessing))
	require.NoError(t, order.ChangeStatus(entity.OrderShipped))
	require.NoError(t, order.ChangeStatus(entity.OrderDelivered))

	for _, s := range []entity.OrderStatus{
		entity.OrderCreated,
		entity.OrderProcessing,
		entity.OrderShipped,
		entity.OrderDelivered,
	} {
		err := order.ChangeStatus(s)
		assert.Error(t, err, "Nenhuma transição deve ser permitida de OrderDelivered para %s", s)
	}
}

func TestOrder_ChangeStatus_UpdatesTimestamp(t *testing.T) {
	order, err := entity.NewOrder("customer-1", 100.0, 1)
	require.NoError(t, err)

	original := order.UpdatedAt
	time.Sleep(time.Millisecond)

	require.NoError(t, order.ChangeStatus(entity.OrderProcessing))
	assert.True(t, order.UpdatedAt.After(original))
}
