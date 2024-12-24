package errs

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*              - Errors -             */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

// NOTE -> These are just strings, error messages, NOT actual errors.

const (

	/* -~-~-~-~-~-~ Fatal error messages (init/shutdown) ~-~-~-~-~-~- */

	FailedToCreateLogger = "Failed to create Logger: %v"
	FailedToConnectToDB  = "Failed to connect to the DB: %v"

	/* -~-~-~-~-~ Non-Fatal error messages (init/shutdown) ~-~-~-~-~- */

	FailedToGetSQLDB   = "Failed to get SQL DB connection: %v"
	FailedToCloseSQLDB = "Failed to close SQL DB connection: %v"
)

const (

	// DB Errors
	DBNoQueryOpts      = "db error -> no query options"
	DBCreatingWebpage  = "db error -> creating webpage"
	DBCountingWebpages = "db error -> counting webpages"
)
