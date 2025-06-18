package core

import "errors"

var (
	ErrChainDatabaseClosed = errors.New("blockchain leveldb database is closed")
)

var (
	ErrBlockChainInsertFailed = errors.New("failed to insert block into the chain")
)
