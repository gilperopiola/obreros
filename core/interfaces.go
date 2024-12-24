package core

import (
	"github.com/gilperopiola/god"
	"github.com/gilperopiola/obreros/core/models"
)

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*            - Interfaces -           */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

/* -~-~-~-~- Main Interfaces -~-~-~-~- */

type Obrero interface {
	Laburar()
}

// With this you can avoid importing the toolbox pkg.
// Remember to add new tools on the app.go file as well.
type Toolbox interface {
	APIs
	DBTool
	FileManager
	ShutdownJanitor
	Retrier
}

// Unifies our SQL and Mongo AnyDB Interfaces.
type AnyDB interface {
	GetInnerDB() any // *gorm.DB or *mongo.Client
}

/* -~-~-~-~- Toolbox: Tools -~-~-~-~- */

type (
	FileManager interface {
		CreateFolder(path string) error
		CreateFolders(paths ...string) error
	}
	ShutdownJanitor interface {
		AddCleanupFunc(fn func())
		AddCleanupFuncWithErr(fn func() error)
		Cleanup()
	}
	Retrier interface {
		TryTo(doThis func() error, nTimes int) error
		TryToOrElse(doThis func() (any, error), orElse func(), nTimes int) (any, error)
	}
	DBTool interface {
		GetDB() AnyDB
		CloseDB()
		IsNotFound(err error) bool

		// Webpages
		InsertWebpage(ctx god.Ctx, url, title, content string) (*models.Webpage, error)
	}
)

type (
	APIs interface {
		InternalAPIs
		ExternalAPIs
	}
	InternalAPIs interface{}
	ExternalAPIs interface{}
)

/* -~-~-~- SQL DB ~-~-~- */

// Low-level API for our SQL Database.
// It's an adapter for Gorm. Concrete types sql.sqlAdapter and mocks.Gorm implement this.
type (
	SqlDB interface {
		AnyDB
		AddError(err error) error
		AutoMigrate(dst ...any) error
		Association(column string) SqlDBAssociation
		Close()
		Count(value *int64) SqlDB
		Create(value any) SqlDB
		Debug() SqlDB
		Delete(value any, where ...any) SqlDB
		Error() error
		Find(out any, where ...any) SqlDB
		First(out any, where ...any) SqlDB
		FirstOrCreate(out any, where ...any) SqlDB
		Group(query string) SqlDB
		Joins(query string, args ...any) SqlDB
		Limit(value int) SqlDB
		Model(value any) SqlDB
		Offset(value int) SqlDB
		Or(query any, args ...any) SqlDB
		Order(value string) SqlDB
		Pluck(column string, value any) SqlDB
		Raw(sql string, values ...any) SqlDB
		Row() SqlRow
		Rows() (SqlRows, error)
		RowsAffected() int64
		Save(value any) SqlDB
		Scan(dest any) SqlDB
		Scopes(funcs ...func(SqlDB) SqlDB) SqlDB
		WithContext(ctx god.Ctx) SqlDB
		Where(query any, args ...any) SqlDB
	}

	SqlDBOpt func(SqlDB) // Optional functions to apply to a query
)

// Used to avoid importing the gorm and sql pkgs: *gorm.Association, *sql.Row, *sql.Rows.
type (
	SqlDBAssociation interface {
		Append(values ...interface{}) error
	}
	SqlRow interface {
		Scan(dest ...any) error
	}
	SqlRows interface {
		Next() bool
		Scan(dest ...any) error
		Close() error
	}
)
