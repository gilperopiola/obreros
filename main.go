package main

import (
	"os"
	"os/signal"
	"syscall"
)

/* ~-~--~-~-~-~-~-~-~-~ */
/*      - Obreros -     */
/* ~-~--~-~-~-~-~-~-~-~ */

func main() {
	runApp, cleanUp := NewApp()
	defer cleanUp()

	runApp()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-ch
}
