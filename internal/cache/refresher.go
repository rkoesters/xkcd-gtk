package cache

// ViewRefresher provides methods for refreshing the view of cache statistics.
// All methods must silently accept a nil receiver.
type ViewRefresher interface {
	RefreshMetadata()
	RefreshMetadataWith(Stat)
	RefreshImages()
	RefreshImagesWith(Stat)
}

// ViewRefresherGetter returns a ViewRefresher. Useful for lazily passing the
// ViewRefresher as an argument.
type ViewRefresherGetter func() ViewRefresher

type nullRefresher struct{}

func newNullRefresher() ViewRefresher              { return nullRefresher{} }
func (r nullRefresher) RefreshMetadata()           {}
func (r nullRefresher) RefreshMetadataWith(_ Stat) {}
func (r nullRefresher) RefreshImages()             {}
func (r nullRefresher) RefreshImagesWith(_ Stat)   {}
