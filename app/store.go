package app

// Store ...
type Store interface {
	Close() error
	Migrate(collection, id string) error
	GetViews_(collection, id string) (int64, error)
	IncView_(collection, id string) error
	GetViews(id string) (int64, error)
	IncViews(id string) error
}
