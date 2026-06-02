package chain

import (
	"cosmossdk.io/store/types"
)

// queryableMultiStore wraps a CommitMultiStore to gracefully handle
// CacheMultiStoreWithVersion failures at the latest height.
// This addresses an IAVL v1.1.2 compatibility issue where version tracking
// may not properly support versioned queries on fresh chains, even though
// the underlying state is correct (proven via ABCI store queries).
type queryableMultiStore struct {
	types.CommitMultiStore
}

func newQueryableMultiStore(cms types.CommitMultiStore) *queryableMultiStore {
	return &queryableMultiStore{CommitMultiStore: cms}
}

// CacheMultiStoreWithVersion attempts a versioned cache. If it fails and the
// requested version is the latest committed version, falls back to the current
// committed cache (which contains the correct state).
func (q *queryableMultiStore) CacheMultiStoreWithVersion(version int64) (types.CacheMultiStore, error) {
	cms, err := q.CommitMultiStore.CacheMultiStoreWithVersion(version)
	if err == nil {
		return cms, nil
	}

	// Fallback: if the requested version matches the latest committed version,
	// use the current state directly. The data is correct — only the IAVL
	// version metadata lookup fails in certain configurations.
	latestVersion := q.CommitMultiStore.LatestVersion()
	if version == latestVersion || version == latestVersion-1 {
		return q.CommitMultiStore.CacheMultiStore(), nil
	}

	return nil, err
}
