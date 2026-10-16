package ivs

type IVS struct {
	Channels  []Channel  `tree:"channels"`
	ChatRooms []ChatRoom `tree:"chat_rooms"`
}
