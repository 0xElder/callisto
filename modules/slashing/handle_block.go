package slashing

import (
	"fmt"

	juno "github.com/forbole/juno/v6/types"

	tmctypes "github.com/cometbft/cometbft/rpc/core/types"
	"github.com/rs/zerolog/log"
)

// HandleBlock implements BlockModule
func (m *Module) HandleBlock(
	block *tmctypes.ResultBlock, results *tmctypes.ResultBlockResults, _ []*juno.Transaction, _ *tmctypes.ResultValidators,
) error {
	// Update the signing infos
	err := m.updateSigningInfo(block.Block.Height)
	if err != nil {
		return fmt.Errorf("error while updating signing info: %s", err)
	}

	return nil
}

// updateSigningInfo reads from the LCD the current staking pool and stores its value inside the database
func (m *Module) updateSigningInfo(height int64) error {
	log.Debug().Str("module", "slashing").Int64("height", height).Msg("updating signing info")

	signingInfos, err := m.getSigningInfos(height)
	if err != nil {
		return err
	}

	return m.db.SaveValidatorsSigningInfos(signingInfos)
}
