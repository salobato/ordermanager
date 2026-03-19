package mongo_test

import (
	"testing"

	"github.com/salobato/ordermanager/internal/adapter/repository/mongo"

	"github.com/stretchr/testify/assert"
)

func TestCounterRepository_FirstSequence(t *testing.T) {
	db := setupTestDB(t)
	repo := mongo.NewCounterRepository(db)

	seq, err := repo.GetNextSequence("order_number")

	assert.NoError(t, err)
	assert.Equal(t, int64(1), seq)
}

func TestCounterRepository_IncrementSequence(t *testing.T) {
	db := setupTestDB(t)
	repo := mongo.NewCounterRepository(db)

	seq1, _ := repo.GetNextSequence("order_number")
	seq2, _ := repo.GetNextSequence("order_number")
	seq3, _ := repo.GetNextSequence("order_number")

	assert.Equal(t, int64(1), seq1)
	assert.Equal(t, int64(2), seq2)
	assert.Equal(t, int64(3), seq3)
}

func TestCounterRepository_GetCurrentSequence(t *testing.T) {
	db := setupTestDB(t)
	repo := mongo.NewCounterRepository(db)

	_, _ = repo.GetNextSequence("order_number")
	_, _ = repo.GetNextSequence("order_number")

	current, err := repo.GetCurrentSequence("order_number")

	assert.NoError(t, err)
	assert.Equal(t, int64(2), current)
}

func TestCounterRepository_GetCurrentSequence_NotFound(t *testing.T) {
	db := setupTestDB(t)
	repo := mongo.NewCounterRepository(db)

	current, err := repo.GetCurrentSequence( "unknown_counter")

	assert.NoError(t, err)
	assert.Equal(t, int64(0), current)
}

func TestCounterRepository_ConcurrentAccess(t *testing.T) {
	db := setupTestDB(t)
	repo := mongo.NewCounterRepository(db)

	

	total := 10
	results := make(chan int64, total)

	for i := 0; i < total; i++ {
		go func() {
			seq, _ := repo.GetNextSequence("order_number")
			results <- seq
		}()
	}

	seen := make(map[int64]bool)

	for i := 0; i < total; i++ {
		seq := <-results
		seen[seq] = true
	}

	assert.Len(t, seen, total)
}
