package toolbox

import (
	"github.com/gilperopiola/obreros/core"
	"github.com/gilperopiola/obreros/toolbox/api_clients"
	"github.com/gilperopiola/obreros/toolbox/db_tool/sqldb"
)

var _ core.Toolbox = (*Toolbox)(nil)

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*            - Toolbox -              */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~ v1 */

// 🛠️ Things that perform actions 🛠️
type Toolbox struct {
	core.APIs            // -> API Clients.
	core.DBTool          // -> Storage (DB, Cache, etc).
	core.FileManager     // -> Creates folders and files.
	core.Retrier         // -> Executes a fn and retries if it fails.
	core.ShutdownJanitor // -> Cleans up and frees resources on application shutdown.
}

func Setup(cfg *core.Config) *Toolbox {
	toolbox := Toolbox{}

	toolbox.Retrier = NewRetrier(&cfg.RetrierCfg)
	sqlDB := sqldb.NewSqlDB(&cfg.DBCfg, toolbox.Retrier)
	toolbox.DBTool = sqldb.NewDBTool(sqlDB)
	toolbox.APIs = api_clients.NewAPIClients()
	toolbox.FileManager = NewFileManager("etc/data/")
	toolbox.ShutdownJanitor = NewShutdownJanitor()

	return &toolbox
}
