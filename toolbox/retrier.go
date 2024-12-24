package toolbox

import (
	"fmt"
	"math"
	"time"

	"github.com/gilperopiola/obreros/core"
)

var _ core.Retrier = (*retrier)(nil)

/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */
/*             - Retrier -             */
/* -~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~-~- */

type retrier struct {
	config *core.RetrierCfg // Empty for now
}

func NewRetrier(cfg *core.RetrierCfg) core.Retrier {
	return &retrier{cfg}
}

func (r retrier) TryTo(doThis func() error, nTries int) error {
	adapter := func() (any, error) {
		// Used to match the signature of TryToOrElse() first param
		return nil, doThis()
	}
	_, err := r.TryToOrElse(adapter, func() {}, nTries)
	return err
}

func (r retrier) TryToOrElse(doThis func() (any, error), orElse func(), nTries int) (any, error) {
	var result any
	var err error

	for nTry := 1; nTry <= nTries; nTry++ {
		if result, err = doThis(); err == nil {
			break
		}

		core.LogUnexpectedError(fmt.Errorf("try %d/%d failed: %v", nTry, nTries, err))

		if nTry == nTries { // Don't sleep on the last try
			break
		}

		sleepFor := math.Pow(float64(nTry), 2) // 2, 4, 8, 16, 32...
		time.Sleep(time.Second * time.Duration(sleepFor))
		orElse()
	}

	return result, nil
}
