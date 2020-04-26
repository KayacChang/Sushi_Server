package main

import (
	"encoding/json"
	"log"

	"github.com/joho/godotenv"
	"gitlab.fbk168.com/gamedevjp/sushi/server/env"
	"gitlab.fbk168.com/gamedevjp/sushi/server/game"
)

type ENV struct {
	Maintain             bool   `json:"Maintain"`
	MaintainStartTime    string `json:"MaintainStartTime"`
	MaintainFinishTime   string `json:"MaintainFinishTime"`
	MaintainCheckoutTime string `json:"ULGMaintainCheckoutTime"`

	ServerIP         string `json:"IP"`
	ServerHTTPPort   string `json:"PORT"`
	ServerSocketPort string `json:"SocketPORT"`

	HTTPS bool   `json:"Https"`
	Cert  string `json:"Cert"`
	Key   string `json:"Key"`

	DBIP       string `json:"DBIP"`
	DBPort     string `json:"DBPORT"`
	DBUser     string `json:"DBUser"`
	DBPassword string `json:"DBPassword"`

	RedisURL string `json:"RedisURL"`

	APIURL string `json:"TransferURL"`

	AccountEncode string `json:"AccountEncodeStr"`
	ServerMod     string `json:"ServerMod"`
	GameTypeID    string `json:"GameTypeID"`
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Panicf("No [ .env ] file found...\n")
	}

	config := ENV{
		Maintain:             env.GetEnvAsBool("MAINTAIN"),
		MaintainStartTime:    env.GetEnvAsString("MAINTAIN_START_TIME"),
		MaintainFinishTime:   env.GetEnvAsString("MAINTAIN_FINISH_TIME"),
		MaintainCheckoutTime: env.GetEnvAsString("MAINTAIN_CHECKOUT_TIME"),

		ServerIP:         env.GetEnvAsString("SERVER_IP"),
		ServerHTTPPort:   env.GetEnvAsString("SERVER_HTTP_PORT"),
		ServerSocketPort: env.GetEnvAsString("SERVER_SOCKET_PORT"),

		HTTPS: env.GetEnvAsBool("SERVER_HTTPS"),
		Cert:  env.GetEnvAsString("SERVER_CERT"),
		Key:   env.GetEnvAsString("SERVER_KEY"),

		DBIP:       env.GetEnvAsString("DB_IP"),
		DBPort:     env.GetEnvAsString("DB_PORT"),
		DBUser:     env.GetEnvAsString("DB_USER"),
		DBPassword: env.GetEnvAsString("DB_PASSWORD"),

		RedisURL: env.GetEnvAsString("REDIS_URL"),

		APIURL: env.GetEnvAsString("API_URL"),

		AccountEncode: env.GetEnvAsString("ACCOUNT_ENCODE"),
		ServerMod:     env.GetEnvAsString("SERVER_MOD"),
		GameTypeID:    env.GetEnvAsString("GAME_TYPEID"),
	}

	jsonbyte, err := json.Marshal(config)
	if err != nil {
		log.Panicf("error: %s", err.Error())
	}

	game.NewGameServer(string(jsonbyte))
}
