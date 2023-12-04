// Ref: https://droctothorpe.github.io/posts/2020/07/leveled-logs-with-cobra-and-logrus/

package internal

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var Debug bool

func ToggleDebug(cmd *cobra.Command, args []string) {
	if Debug {
		log.SetLevel(log.DebugLevel)
		log.SetFormatter(&log.TextFormatter{})
	}
}
