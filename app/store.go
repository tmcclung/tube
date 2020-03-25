package app

// Store ...
type Store interface {
	Close() error
	GetViews(collection, id string) (int64, error)
	IncView(collection, id string) error
}
