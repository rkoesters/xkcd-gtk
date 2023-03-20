package cache

// ViewRefreshWither provides methods for refreshing a view of cache statistics
// with the provided stat. All methods must silently accept a nil receiver.
type ViewRefreshWither interface {
	RefreshMetadataWith(Stat)
	RefreshImagesWith(Stat)
}

// ViewRefresher provides methods for refreshing a view of cache statistics.
// All methods must silently accept a nil receiver.
type ViewRefresher interface {
	ViewRefreshWither
	// RefreshMetadata queries StatMetadata then calls RefreshMetadataWith.
	RefreshMetadata()
	// RefreshImages queries StatImages then calls RefreshImagesWith.
	RefreshImages()
}

// ViewRefresherGetter returns a ViewRefresher. Useful for lazily passing the
// ViewRefresher as an argument.
type ViewRefresherGetter func() ViewRefresher

// ViewRefreshWitherGetter returns a ViewRefreshWither. Useful for lazily
// passing the ViewRefreshWither as an argument.
type ViewRefreshWitherGetter func() ViewRefreshWither

type nullRefresher struct{}

func (_ *nullRefresher) RefreshMetadata()           {}
func (_ *nullRefresher) RefreshMetadataWith(_ Stat) {}
func (_ *nullRefresher) RefreshImages()             {}
func (_ *nullRefresher) RefreshImagesWith(_ Stat)   {}

var nilRefresher *nullRefresher = nil
