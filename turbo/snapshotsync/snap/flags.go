package snap

import (
	"fmt"

	"github.com/ledgerwatch/erigon-lib/kv"
	"github.com/ledgerwatch/erigon/eth/ethconfig"
)

var (
	blockSnapshotEnabledKey = []byte("blocksSnapshotEnabled")
)

func Enabled(tx kv.Getter) (bool, error) {
	return kv.GetBool(tx, kv.DatabaseInfo, blockSnapshotEnabledKey)
}

// makes sure that erigon is on the same syncmode used previously
func EnsureNotChanged(tx kv.GetPut, cfg ethconfig.Snapshot) error {
	ok, v, err := kv.EnsureNotChangedBool(tx, kv.DatabaseInfo, blockSnapshotEnabledKey, cfg.Enabled)
	if err != nil {
		return err
	}
	if !ok {
		if v {
			return fmt.Errorf("we recently changed default of --syncmode flag, or you forgot to set --syncmode flag, please add flag --syncmode=snap")
		} else {
			return fmt.Errorf("we recently changed default of --syncmode flag, or you forgot to set --syncmode flag, please add flag --syncmode=full")
		}
	}
	return nil
}

// ForceSetFlags - if you know what you are doing
func ForceSetFlags(tx kv.GetPut, cfg ethconfig.Snapshot) error {
	if cfg.Enabled {
		if err := tx.Put(kv.DatabaseInfo, blockSnapshotEnabledKey, []byte{1}); err != nil {
			return err
		}
	} else {
		if err := tx.Put(kv.DatabaseInfo, blockSnapshotEnabledKey, []byte{0}); err != nil {
			return err
		}
	}
	return nil
}
