package errs

import (
	"fmt"
)

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*          - Worker Errors -          */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

type WorkerErr struct {
	Err      error
	Metadata string // Holds extra info.
}

func (werr WorkerErr) Error() string {
	if werr.Metadata == "" {
		return werr.Unwrap().Error()
	}
	return fmt.Sprintf("%s -> %v", werr.Metadata, werr.Unwrap())
}

func (werr WorkerErr) Unwrap() error {
	if werr.Err != nil {
		return werr.Err
	}
	return fmt.Errorf(werr.Metadata)
}
