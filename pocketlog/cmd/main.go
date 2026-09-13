package main

import (
	"go-pocket-sized-projects/pocketlog"
	"os"
)

func main() {
	logger := pocketlog.New(pocketlog.LevelInfo, pocketlog.WithOutput(os.Stdout))
	logger.Infof("A little copying is better than a little dependency.")
	logger.Errorf("Errors are values. Documentation is for %s.", "users")
	logger.Debugf("Make the zero (%d) value useful.", 0)
}
