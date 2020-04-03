package game

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/YWJSonic/GameServer/sushi/game/db"
	"github.com/YWJSonic/GameServer/sushi/game/protocol"
	"github.com/gorilla/websocket"
	"gitlab.fbk168.com/gamedevjp/backend-utility/utility/code"
	"gitlab.fbk168.com/gamedevjp/backend-utility/utility/foundation"
	"gitlab.fbk168.com/gamedevjp/backend-utility/utility/httprouter"
	"gitlab.fbk168.com/gamedevjp/backend-utility/utility/igame"
	"gitlab.fbk168.com/gamedevjp/backend-utility/utility/messagehandle"
	"gitlab.fbk168.com/gamedevjp/backend-utility/utility/socket"
	"gitlab.fbk168.com/gamedevjp/cyberpunk/server/game/constants"
)

func (g *Game) createNewSocket(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("token")

	if err := g.CheckToken(token); err != nil {
		log.Fatal("createNewSocket: not this token\n")
		return
	}

	c, err := g.Server.Socket.Upgrade(w, r, r.Header)
	if err != nil {
		log.Println("upgrade:", err)
		return
	}
	g.Server.Socket.AddNewConn("f", c, func(msg socket.Message) error {
		fmt.Println("#-- socket --#", msg)
		return nil
	})
	// g.Server.Socket.AddNewConn(user.GetGameInfo().GameAccount, c, g.SocketMessageHandle)

	time.Sleep(time.Second * 3)
	g.Server.Socket.ConnMap["f"].Send(websocket.CloseMessage, []byte{})
}

// SocketMessageHandle ...
func (g *Game) SocketMessageHandle(msg socket.Message) error {
	fmt.Println("#-- socket --#", msg)
	return nil
}

func (g *Game) gameinit(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var proto protocol.InitRequest
	proto.InitData(r)

	// get user
	user, _, err := g.GetUser(proto.Token)
	if err != nil {
		err := messagehandle.New()
		err.ErrorCode = code.NoThisPlayer
		g.Server.HTTPResponse(w, "", err)
		return
	}

	result := map[string]interface{}{
		"betrate": g.IGameRule.GetBetSetting(),
		"player": map[string]interface{}{
			"gameaccount": g.IGameRule.GetGameTypeID(),
			"id":          user.UserGameInfo.IDStr,
			"money":       user.UserGameInfo.GetMoney(),
		},
		"reel": g.IGameRule.GetReel(),
	}
	g.Server.HTTPResponse(w, result, messagehandle.New())
}

func (g *Game) gameresult(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var proto protocol.GameRequest
	proto.InitData(r)

	user, _, err := g.GetUser(proto.Token)
	if err != nil {
		err := messagehandle.New()
		err.Msg = "GameTypeError"
		err.ErrorCode = code.GameTypeError
		// messagehandle.ErrorLogPrintln("GetPlayerInfoByPlayerID-2", err, token, betIndex, betMoney)
		g.Server.HTTPResponse(w, "", err)
		return
	}
	if user.UserGameInfo.Money < g.IGameRule.GetBetMoney(proto.BetIndex) {
		err := messagehandle.New()
		err.Msg = "NoMoneyToBet"
		err.ErrorCode = code.NoMoneyToBet
		g.Server.HTTPResponse(w, "", err)
		return
	}

	order, errproto, err := g.NewOrder(proto.Token, user.UserGameInfo.IDStr, g.IGameRule.GetBetMoney(proto.BetIndex))

	if errproto != nil {
		err := messagehandle.New()
		err.Msg = errproto.String()
		err.ErrorCode = code.ULGInfoFormatError
		g.Server.HTTPResponse(w, "", err)
		return
	}
	if err != nil {
		err := messagehandle.New()
		err.Msg = "ULGInfoFormatError"
		err.ErrorCode = code.ULGInfoFormatError
		g.Server.HTTPResponse(w, "", err)
		return
	}

	oldMoney := user.UserGameInfo.Money
	// get game result
	RuleRequest := &igame.RuleRequest{
		BetIndex: proto.BetIndex,
		UserID:   user.UserGameInfo.ID,
	}
	result := g.IGameRule.GameRequest(RuleRequest)
	user.UserGameInfo.SumMoney(result.Totalwinscore - result.BetMoney)
	resultMap := result.GameResult
	resultMap["playermoney"] = user.UserGameInfo.GetMoney()

	msg := foundation.JSONToString(result.GameResult)
	msg = strings.ReplaceAll(msg, "\"", "\\\"")
	errMsg := db.SetLog(
		g.Server.DBConn("logdb"),
		user.UserGameInfo.IDStr,
		0,
		time.Now().Unix(),
		constants.ActionGameResult,
		oldMoney,
		user.UserGameInfo.Money,
		result.Totalwinscore,
		"",
		"",
		"",
		msg,
	)
	if errMsg.ErrorCode != code.OK {
		g.Server.HTTPResponse(w, resultMap, errMsg)
		return
	}

	g.EndOrder(proto.Token, order)
	g.Server.HTTPResponse(w, resultMap, messagehandle.New())
}
