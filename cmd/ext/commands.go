package ext

import "mbtv2/cmd/mta/models"

type Cmd struct {
	Info    string
	Command []string
}

// ExeCmd - Get the build operation's
//Todo - get from external resources
func ExeCmd(m models.Modules) []Cmd {

	switch m.Type {
	case "html5":
		// TODO get the maps from external source
		return []Cmd{
			{"# installing module dependencies & execute grunt & remove dev dependencies",
				[]string{"npm install", "grunt", "npm prune --production"}},
		}
	case "nodejs":
		return []Cmd{{"# TODO build for node.js",
			[]string{"Not supported yet"}},
		}
	default:
		return []Cmd{{"# New module type",
			[]string{"Not supported yet"}}}
	}

}
