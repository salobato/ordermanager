package repository

type CounterRepository interface {
	GetNextSequence(counterName string) (int64, error)
	GetCurrentSequence(counterName string) (int64, error)
}
