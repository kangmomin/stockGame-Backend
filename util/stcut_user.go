package util

type UserInfo struct {
	DiscordId    string `json:"discord_id"`
	Name         string `json:"name"`
	Coin         int    `json:"coin"`
	Bank         int    `json:"bank"`
	Tax          int    `json:"tax"`
	GambleTicket int    `json:"gamble_ticket"`
	Stock        []struct {
		Name string `json:"name"`

		// 원가
		Cost  int `json:"cost"`
		Count int `json:"count"`
	} `json:"stock"`
}
