package collectors

import "context"

const (
	START     = "start"
	STOP      = "stop"
	RESTART   = "restart"
	STATUS    = "status"
	RELOAD    = "reload"
	INVENTORY = "inventory"
)

func runCommands(ctx context.Context, a Action) {
	for _, command := range a.Commands {
		// load collector inventory
		switch command.Command {
		case INVENTORY:
			// re-download inventory
		case START:
			// start collector
		case STOP:
			// stop collector
		case RESTART:
			// restart collector
		case STATUS:
			// get collector status
		}
	}
}
