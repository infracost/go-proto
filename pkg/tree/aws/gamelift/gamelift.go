package gamelift

import "strings"

type GameLift struct {
	Builds           []Build           `tree:"builds"`
	Fleets           []Fleet           `tree:"fleets"`
	GameServerGroups []GameServerGroup `tree:"game_server_groups"`
	StreamGroups     []StreamGroup     `tree:"stream_groups"`
}

func (g *GameLift) PostProcess() {
	for i := range g.Fleets {
		g.Fleets[i].Relationships = FleetRelationships{}
	}
	for i, fleet := range g.Fleets {
		if fleet.BuildID.IsEmpty() {
			continue
		}
		for j := range g.Builds {
			build := &g.Builds[j]
			if fleet.BuildID.Value() == build.ID || strings.Contains(fleet.BuildID.Value(), build.ID) ||
				(!build.Name.IsEmpty() && fleet.BuildID.Value() == build.Name.Value()) {
				g.Fleets[i].Relationships.Build = build
				break
			}
		}
	}
}
