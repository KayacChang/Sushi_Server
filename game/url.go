package game

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/YWJSonic/GameServer/sushi/game/protocol"
	"github.com/YWJSonic/ServerUtility/code"
	"github.com/YWJSonic/ServerUtility/httprouter"
	"github.com/YWJSonic/ServerUtility/igame"
	"github.com/YWJSonic/ServerUtility/messagehandle"
	"github.com/YWJSonic/ServerUtility/socket"
	"github.com/gorilla/websocket"
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
	// var result = make(map[string]interface{})
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

	// get game result
	RuleRequest := &igame.RuleRequest{
		BetIndex: proto.BetIndex,
		UserID:   user.UserGameInfo.ID,
	}
	result := g.IGameRule.GameRequest(RuleRequest)
	user.UserGameInfo.SumMoney(result.Totalwinscore - result.BetMoney)
	resultMap := result.GameResult
	// resultMap["totalwinscore"] = result.Totalwinscore
	// resultMap["playermoney"] = user.UserGameInfo.GetMoney()
	// resultMap["normalresult"] = result.GameResult["normalresult"]
	// resultMap["freegamecount"] = result.OtherData["freegamecount"]

	g.EndOrder(proto.Token, order)
	g.Server.HTTPResponse(w, resultMap, messagehandle.New())
}
