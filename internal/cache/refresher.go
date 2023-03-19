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

func (_ *nullRefresher) RefreshMetadata()           {}
func (_ *nullRefresher) RefreshMetadataWith(_ Stat) {}
func (_ *nullRefresher) RefreshImages()             {}
func (_ *nullRefresher) RefreshImagesWith(_ Stat)   {}

var nilRefresher *nullRefresher = nil
