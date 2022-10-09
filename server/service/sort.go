package service

import "github.com/rmarken5/ffdp/protobuf/proto_files/player_proto"

type ByLastNameAsc []*player_proto.Player

func (b ByLastNameAsc) Len() int      { return len(b) }
func (b ByLastNameAsc) Swap(i, j int) { b[i], b[j] = b[j], b[i] }
func (b ByLastNameAsc) Less(i, j int) bool {
	if b[i].LastName == b[j].LastName {
		return b[i].FirstName > b[j].FirstName
	}
	return b[i].LastName > b[j].LastName
}

type ByLastNameDesc []*player_proto.Player

func (b ByLastNameDesc) Len() int      { return len(b) }
func (b ByLastNameDesc) Swap(i, j int) { b[i], b[j] = b[j], b[i] }
func (b ByLastNameDesc) Less(i, j int) bool {
	if b[i].LastName == b[j].LastName {
		return b[i].FirstName < b[j].FirstName
	}
	return b[i].LastName < b[j].LastName
}

type ByFirstNameAsc []*player_proto.Player

func (b ByFirstNameAsc) Len() int      { return len(b) }
func (b ByFirstNameAsc) Swap(i, j int) { b[i], b[j] = b[j], b[i] }
func (b ByFirstNameAsc) Less(i, j int) bool {
	return b[i].FirstName > b[j].FirstName
}

type ByFirstNameDesc []*player_proto.Player

func (b ByFirstNameDesc) Len() int           { return len(b) }
func (b ByFirstNameDesc) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }
func (b ByFirstNameDesc) Less(i, j int) bool { return b[i].FirstName < b[j].FirstName }

type ByValueDesc []*player_proto.Player

func (b ByValueDesc) Len() int           { return len(b) }
func (b ByValueDesc) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }
func (b ByValueDesc) Less(i, j int) bool { return b[i].Value < b[j].Value }

type ByValueAsc []*player_proto.Player

func (b ByValueAsc) Len() int           { return len(b) }
func (b ByValueAsc) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }
func (b ByValueAsc) Less(i, j int) bool { return b[i].Value > b[j].Value }
