package printer

import (
	"fmt"
	"github.com/rmarken5/ffdp/protobuf/proto_files/player_proto"
)

const (
	PlayerFmt = "|%-4s|%-6s|%-35s|%-8s|%-4s|%-3s|%-5s|%-4d|%-5f|\n"
)

func Print(player *player_proto.Player) string {
	return fmt.Sprintf(PlayerFmt, player.Rank, player.ID, player.FullName, player.Position, player.Team, player.Bye, player.ADP, player.PointTotal, player.Value)
}
