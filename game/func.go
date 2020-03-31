package game

import (
	"encoding/json"
	"errors"

	"github.com/YWJSonic/GameServer/sushi/game/gamerule"
	"github.com/YWJSonic/ServerUtility/foundation"
	"github.com/YWJSonic/ServerUtility/foundation/fileload"
	"github.com/YWJSonic/ServerUtility/iserver"
	"github.com/YWJSonic/ServerUtility/restfult"
	"github.com/YWJSonic/ServerUtility/socket"
)

// NewGameServer ...
func NewGameServer() {

	jsStr := fileload.Load("./file/config.json")
	config := foundation.StringToJSON(jsStr)
	baseSetting := iserver.NewSetting()
	baseSetting.SetData(config)

	gamejsStr := fileload.Load("./file/gameconfig.json")
	var gameRule = &gamerule.Rule{}
	if err := json.Unmarshal([]byte(gamejsStr), &gameRule); err != nil {
		panic(errors.New("gameconfig error: "))
	}

	var gameserver = iserver.NewService()
	var game = &Game{
		IGameRule: gameRule,
		Server:    gameserver,
		// ProtocolMap: protocol.NewProtocolMap(),
	}
	gameserver.Restfult = restfult.NewRestfultService()
	gameserver.Socket = socket.NewSocket()
	gameserver.IGame = game

	// start Server
	gameserver.Launch(baseSetting)

	// start DB service
	setting := gameserver.Setting.DBSetting()
	gameserver.LaunchDB("gameDB", setting)
	gameserver.LaunchDB("logDB", setting)
	gameserver.LaunchDB("payDB", setting)

	// start restful service
	go gameserver.LaunchRestfult(game.RESTfulURLs())
	go gameserver.LaunchSocket(game.SocketURLs())

	<-gameserver.ShotDown
}
