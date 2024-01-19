// Copyright © 2023 Circonus, Inc. <support@circonus.com>
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//

package agents

type Platforms map[string]Architectures

type Architectures map[string]ProcessManagers

type ProcessManagers map[string]string

var commands = Platforms{
	"linux": Architectures{
		"amd64": ProcessManagers{
			"systemd": "/usr/bin/systemctl",
		},
		"arm64": ProcessManagers{
			"systemd": "/usr/bin/systemctl",
		},
	},
	"darwin": Architectures{
		"amd64": ProcessManagers{
			"brew":      "/usr/local/bin/brew",
			"launchctl": "/bin/launchctl",
		},
		"arm64": ProcessManagers{
			"brew":      "/opt/homebrew/bin/brew",
			"launchctl": "/bin/launchctl",
		},
	},
	"windows": Architectures{
		"amd64": ProcessManagers{
			"powershell": "powershell.exe",
		},
	},
}

/*
brew services start <name>
brew services restart <name>
brew services kill <name> // don't use stop - it will unregister it from starting up at login/boot
brew services info --json <name> // status (search for `"running": true`)
*/

/*
powershell.exe Restart-Service -Name <name>
powershell.exe Get-Service -Name <name> // for status `Status field should be "running"`
*/

// brew restart []string{"services","restart"}
// windows restart   []string{"Restart-Service", "-Name"},
