package database

import (
	flag "github.com/spf13/pflag"

	"github.com/iotaledger/hornet/pkg/database"
	"github.com/iotaledger/hornet/pkg/node"
)

const (
	// the used database engine (pebble/rocksdb/mapdb).
	CfgDatabaseEngine = "db.engine"
	// the path to the database folder.
	CfgDatabasePath = "db.path"
	// whether to automatically start revalidation on startup if the database is corrupted.
	CfgDatabaseAutoRevalidation = "db.autoRevalidation"
	// ignore the check for corrupted databases (should only be used for debug reasons).
	CfgDatabaseDebug = "db.debug"
	// whether to check if the ledger state matches the total supply on startup
	CfgCheckLedgerStateOnStartup = "db.checkLedgerStateOnStartup"
)

var params = &node.PluginParams{
	Params: map[string]*flag.FlagSet{
		"nodeConfig": func() *flag.FlagSet {
			fs := flag.NewFlagSet("", flag.ContinueOnError)
			fs.String(CfgDatabaseEngine, string(database.EngineRocksDB), "the used database engine (pebble/rocksdb/mapdb)")
			fs.String(CfgDatabasePath, "mainnetdb", "the path to the database folder")
			fs.Bool(CfgDatabaseAutoRevalidation, false, "whether to automatically start revalidation on startup if the database is corrupted")
			fs.Bool(CfgDatabaseDebug, false, "ignore the check for corrupted databases (should only be used for debug reasons)")
			fs.Bool(CfgCheckLedgerStateOnStartup, false, "whether to check if the ledger state matches the total supply on startup")
			return fs
		}(),
	},
	Masked: nil,
}
