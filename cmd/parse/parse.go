package parse

import (
	"github.com/cosmos/cosmos-sdk/x/genutil/types"
	parse "github.com/forbole/juno/v6/cmd/parse/types"
	"github.com/forbole/juno/v6/modules"
	"github.com/forbole/juno/v6/types/utils"
	"github.com/spf13/cobra"

	parseblocks "github.com/forbole/juno/v6/cmd/parse/blocks"
	nodeconfig "github.com/forbole/juno/v6/node/config"

	parsetransaction "github.com/forbole/juno/v6/cmd/parse/transactions"
	parsecmdtypes "github.com/forbole/juno/v6/cmd/parse/types"

	parseauth "github.com/forbole/callisto/v4/cmd/parse/auth"
	parsebank "github.com/forbole/callisto/v4/cmd/parse/bank"
	parsedistribution "github.com/forbole/callisto/v4/cmd/parse/distribution"
	parsefeegrant "github.com/forbole/callisto/v4/cmd/parse/feegrant"
	parsegov "github.com/forbole/callisto/v4/cmd/parse/gov"
	parsemint "github.com/forbole/callisto/v4/cmd/parse/mint"
	parsepricefeed "github.com/forbole/callisto/v4/cmd/parse/pricefeed"
	parsestaking "github.com/forbole/callisto/v4/cmd/parse/staking"
)

// NewParseCmd returns the Cobra command allowing to parse some chain data without having to re-sync the whole database
func NewParseCmd(parseCfg *parse.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:               "parse",
		Short:             "Parse some data without the need to re-syncing the whole database from scratch",
		PersistentPreRunE: runPersistentPreRuns(parse.ReadConfigPreRunE(parseCfg)),
	}

	cmd.AddCommand(
		parseauth.NewAuthCmd(parseCfg),
		parsebank.NewBankCmd(parseCfg),
		parseblocks.NewBlocksCmd(parseCfg),
		parsedistribution.NewDistributionCmd(parseCfg),
		parsefeegrant.NewFeegrantCmd(parseCfg),
		NewGenesisCmd(parseCfg),
		parsegov.NewGovCmd(parseCfg),
		parsemint.NewMintCmd(parseCfg),
		parsepricefeed.NewPricefeedCmd(parseCfg),
		parsestaking.NewStakingCmd(parseCfg),
		parsetransaction.NewTransactionsCmd(parseCfg),
	)

	return cmd
}

func runPersistentPreRuns(preRun func(_ *cobra.Command, _ []string) error) func(_ *cobra.Command, _ []string) error {
	return func(cmd *cobra.Command, args []string) error {
		if root := cmd.Root(); root != nil {
			if root.PersistentPreRunE != nil {
				err := root.PersistentPreRunE(root, args)
				if err != nil {
					return err
				}
			}
		}

		return preRun(cmd, args)
	}
}

// NewGenesisCmd returns the Cobra command allowing to parse the genesis file
func NewGenesisCmd(parseConfig *parsecmdtypes.Config) *cobra.Command {
	var flagPath = "genesis-file-path"
	cmd := &cobra.Command{
		Use:   "genesis-file",
		Short: "Parse the genesis file",
		Long: `
Parse the genesis file only.
Note that the modules built will NOT have access to the node as they are only supposed to deal with the genesis
file itself and not the on-chain data.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Read the configuration
			cfg, err := parsecmdtypes.ReadConfig(parseConfig)
			if err != nil {
				return err
			}

			// Set the node to be of type None so that the node won't be built
			cfg.Node.Type = nodeconfig.TypeNone

			// Build the parsing context
			parseCtx, err := parsecmdtypes.GetParserContext(cfg, parseConfig)
			if err != nil {
				return err
			}

			// Get the file path
			genesisFilePath := cfg.Parser.GenesisFilePath
			customPath, _ := cmd.Flags().GetString(flagPath)
			if customPath != "" {
				genesisFilePath = customPath
			}

			// Read the genesis file
			genesis, err := types.AppGenesisFromFile(genesisFilePath)
			if err != nil {
				return err
			}

			// Convert it to genesis doc
			genDoc, err := genesis.ToGenesisDoc()
			if err != nil {
				return err
			}

			// Get the genesis state
			genState, err := utils.GetGenesisState(genDoc)
			if err != nil {
				return err
			}

			for _, module := range parseCtx.Modules {
				if module, ok := module.(modules.GenesisModule); ok {
					err = module.HandleGenesis(genDoc, genState)
					if err != nil {
						return err
					}
				}
			}

			return nil
		},
	}

	cmd.Flags().String(flagPath, "", "Path to the genesis file to be used. If empty, the path will be taken from the config file")

	return cmd
}
