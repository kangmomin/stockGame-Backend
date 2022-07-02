package util

import (
	"log"
	"net/http"
)

type UserInfo struct {
	DiscordId    string `json:"discord_id"`
	Name         string `json:"name"`
	Coin         int    `json:"coin"`
	Bank         int    `json:"bank"`
	Tax          int    `json:"tax"`
	GambleTicket int    `json:"gamble_ticket"`
}

type UserStock struct {
	UserId string `json:"user_id"`
	Name   string `json:"name"`
	Cost   int    `json:"cost"`
	Count  int    `json:"count"`
}

type SignUp struct {
	UserId string `json:"user_id"`
	Name   string `json:"name"`
	Token  string `json:"token"`
}

type DiscordUserInfo struct {
	ID            string      `json:"id"`
	Username      string      `json:"username"`
	Avatar        interface{} `json:"avatar"`
	Discriminator string      `json:"discriminator"`
	PublicFlags   int         `json:"public_flags"`
	Banner        interface{} `json:"banner"`
	BannerColor   interface{} `json:"banner_color"`
	AccentColor   interface{} `json:"accent_color"`
}

// 유저가 존재하는지 디스코드 api를 통해 확인
func IsValidUserId(t string, userId string) bool {
	req, err := http.NewRequest("GET", "https://discord.com/api/v9/users/"+userId, nil)
	if err != nil {
		log.Println(err)
	}

	req.Header.Add("Authorization", "Bot "+t)
	c := http.Client{}
	resp, _ := c.Do(req)

	return resp.Status == "200 OK"
}
