package source

import codectypes "github.com/cosmos/cosmos-sdk/codec/types"

type Source interface {
	GetAllAnyAccounts(height int64) ([]*codectypes.Any, error)
}
