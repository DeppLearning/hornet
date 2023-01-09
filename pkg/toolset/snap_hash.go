package toolset

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	flag "github.com/spf13/pflag"

	coreDatabase "github.com/iotaledger/hornet/core/database"
	"github.com/iotaledger/hornet/pkg/database"
	"github.com/iotaledger/hornet/pkg/model/storage"
	"github.com/iotaledger/hornet/pkg/snapshot"
)

func snapshotHash(args []string) error {

	fs := flag.NewFlagSet("", flag.ContinueOnError)
	fullSnapshotPathFlag := fs.String(FlagToolSnapshotPathFull, "snapshots/mainnet/full_snapshot.bin", "the path to the full snapshot file")
	deltaSnapshotPathFlag := fs.String(FlagToolSnapshotPathDelta, "snapshots/mainnet/delta_snapshot.bin", "the path to the delta snapshot file (optional)")
	outputJSONFlag := fs.Bool(FlagToolOutputJSON, false, FlagToolDescriptionOutputJSON)

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", ToolSnapHash)
		fs.PrintDefaults()
		println(fmt.Sprintf("\nexample: %s --%s %s --%s %s",
			ToolSnapHash,
			FlagToolSnapshotPathFull,
			"snapshots/mainnet/full_snapshot.bin",
			FlagToolSnapshotPathDelta,
			"snapshots/mainnet/delta_snapshot.bin"))
	}

	if err := parseFlagSet(fs, args); err != nil {
		return err
	}

	if len(*fullSnapshotPathFlag) == 0 {
		return fmt.Errorf("'%s' not specified", FlagToolSnapshotPathFull)
	}

	fullPath := *fullSnapshotPathFlag
	deltaPath := *deltaSnapshotPathFlag

	targetEngine, err := database.DatabaseEngineAllowed(database.EnginePebble)
	if err != nil {
		return err
	}

	tempDir, err := os.MkdirTemp("", "snapHash")
	if err != nil {
		return fmt.Errorf("can't create temp dir: %w", err)
	}

	tangleStore, err := database.StoreWithDefaultSettings(filepath.Join(tempDir, coreDatabase.TangleDatabaseDirectoryName), true, targetEngine)
	if err != nil {
		return fmt.Errorf("%s database initialization failed: %w", coreDatabase.TangleDatabaseDirectoryName, err)
	}

	utxoStore, err := database.StoreWithDefaultSettings(filepath.Join(tempDir, coreDatabase.UTXODatabaseDirectoryName), true, targetEngine)
	if err != nil {
		return fmt.Errorf("%s database initialization failed: %w", coreDatabase.UTXODatabaseDirectoryName, err)
	}

	// clean up temp db
	defer func() {
		_ = tangleStore.Close()
		_ = utxoStore.Close()

		_ = os.RemoveAll(tempDir)
	}()

	dbStorage, err := storage.New(tangleStore, utxoStore)
	if err != nil {
		return err
	}

	_, _, err = snapshot.LoadSnapshotFilesToStorage(context.Background(), dbStorage, fullPath, deltaPath)
	if err != nil {
		return err
	}

	return calculateDatabaseLedgerHash(dbStorage, *outputJSONFlag)
}
