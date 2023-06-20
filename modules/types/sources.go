package types

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/simapp/params"
	"github.com/forbole/juno/v4/node/remote"

	markettypes "github.com/akash-network/akash-api/go/node/market/v1beta3"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/forbole/juno/v4/node/local"

	nodeconfig "github.com/forbole/juno/v4/node/config"

	providertypes "github.com/akash-network/akash-api/go/node/provider/v1beta3"
	banksource "github.com/forbole/bdjuno/v4/modules/bank/source"
	remotebanksource "github.com/forbole/bdjuno/v4/modules/bank/source/remote"
	distrsource "github.com/forbole/bdjuno/v4/modules/distribution/source"
	remotedistrsource "github.com/forbole/bdjuno/v4/modules/distribution/source/remote"
	govsource "github.com/forbole/bdjuno/v4/modules/gov/source"
	remotegovsource "github.com/forbole/bdjuno/v4/modules/gov/source/remote"
	marketsource "github.com/forbole/bdjuno/v4/modules/market/source"
	remotemarketsource "github.com/forbole/bdjuno/v4/modules/market/source/remote"
	mintsource "github.com/forbole/bdjuno/v4/modules/mint/source"
	remotemintsource "github.com/forbole/bdjuno/v4/modules/mint/source/remote"
	providersource "github.com/forbole/bdjuno/v4/modules/provider/source"
	remoteprovidersource "github.com/forbole/bdjuno/v4/modules/provider/source/remote"
	slashingsource "github.com/forbole/bdjuno/v4/modules/slashing/source"
	remoteslashingsource "github.com/forbole/bdjuno/v4/modules/slashing/source/remote"
	stakingsource "github.com/forbole/bdjuno/v4/modules/staking/source"
	remotestakingsource "github.com/forbole/bdjuno/v4/modules/staking/source/remote"
)

type Sources struct {
	BankSource     banksource.Source
	DistrSource    distrsource.Source
	GovSource      govsource.Source
	MarketSource   marketsource.Source
	MintSource     mintsource.Source
	ProviderSource providersource.Source
	SlashingSource slashingsource.Source
	StakingSource  stakingsource.Source
}

func BuildSources(nodeCfg nodeconfig.Config, encodingConfig *params.EncodingConfig) (*Sources, error) {
	switch cfg := nodeCfg.Details.(type) {
	case *remote.Details:
		return buildRemoteSources(cfg)
	case *local.Details:
		return nil, fmt.Errorf("local node is currently not supported: %T", cfg)

	default:
		return nil, fmt.Errorf("invalid configuration type: %T", cfg)
	}
}

// func buildLocalSources(cfg *local.Details, encodingConfig *params.EncodingConfig) (*Sources, error) {
// 	source, err := local.NewSource(cfg.Home, encodingConfig)
// 	if err != nil {
// 		return nil, err
// 	}

// 	app := simapp.NewSimApp(
// 		log.NewTMLogger(log.NewSyncWriter(os.Stdout)), source.StoreDB, nil, true, map[int64]bool{},
// 		cfg.Home, 0, simapp.MakeTestEncodingConfig(), simapp.EmptyAppOptions{},
// 	)

// 	// For MarketSource & ProviderSource
// 	akashapp := akashapp.NewApp(
// 		log.NewTMLogger(log.NewSyncWriter(os.Stdout)), source.StoreDB, nil, true, 0, map[int64]bool{},
// 		cfg.Home, simapp.EmptyAppOptions{},
// 	)
// 	escrowKeeper := escrowKeeper.NewKeeper(encodingConfig.Marshaler, sdk.NewKVStoreKey(akashprovider.StoreKey), app.BankKeeper, nil)

// 	sources := &Sources{
// 		BankSource:     localbanksource.NewSource(source, banktypes.QueryServer(app.BankKeeper)),
// 		DistrSource:    localdistrsource.NewSource(source, distrtypes.QueryServer(app.DistrKeeper)),
// 		GovSource:      localgovsource.NewSource(source, govtypes.QueryServer(app.GovKeeper)),
// 		MarketSource:   localmarketsource.NewSource(source, akashmarket.NewKeeper(encodingConfig.Marshaler, sdk.NewKVStoreKey(akashprovider.StoreKey), akashapp.GetSubspace(markettypes.ModuleName), escrowKeeper).NewQuerier()),
// 		MintSource:     localmintsource.NewSource(source, minttypes.QueryServer(app.MintKeeper)),
// 		ProviderSource: localprovidersource.NewSource(source, akashprovider.NewKeeper(encodingConfig.Marshaler, sdk.NewKVStoreKey(akashprovider.StoreKey)).NewQuerier()),
// 		SlashingSource: localslashingsource.NewSource(source, slashingtypes.QueryServer(app.SlashingKeeper)),
// 		StakingSource:  localstakingsource.NewSource(source, stakingkeeper.Querier{Keeper: app.StakingKeeper}),
// 	}

// 	// Mount and initialize the stores
// 	err = source.MountKVStores(app, "keys")
// 	if err != nil {
// 		return nil, err
// 	}

// 	err = source.MountTransientStores(app, "tkeys")
// 	if err != nil {
// 		return nil, err
// 	}

// 	err = source.MountMemoryStores(app, "memKeys")
// 	if err != nil {
// 		return nil, err
// 	}

// 	err = source.InitStores()
// 	if err != nil {
// 		return nil, err
// 	}

// 	return sources, nil
// }

func buildRemoteSources(cfg *remote.Details) (*Sources, error) {
	source, err := remote.NewSource(cfg.GRPC)
	if err != nil {
		return nil, fmt.Errorf("error while creating remote source: %s", err)
	}

	return &Sources{
		BankSource:     remotebanksource.NewSource(source, banktypes.NewQueryClient(source.GrpcConn)),
		DistrSource:    remotedistrsource.NewSource(source, distrtypes.NewQueryClient(source.GrpcConn)),
		GovSource:      remotegovsource.NewSource(source, govtypes.NewQueryClient(source.GrpcConn)),
		MarketSource:   remotemarketsource.NewSource(source, markettypes.NewQueryClient(source.GrpcConn)),
		MintSource:     remotemintsource.NewSource(source, minttypes.NewQueryClient(source.GrpcConn)),
		ProviderSource: remoteprovidersource.NewSource(source, providertypes.NewQueryClient(source.GrpcConn)),
		SlashingSource: remoteslashingsource.NewSource(source, slashingtypes.NewQueryClient(source.GrpcConn)),
		StakingSource:  remotestakingsource.NewSource(source, stakingtypes.NewQueryClient(source.GrpcConn)),
	}, nil
}
