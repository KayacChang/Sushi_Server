package game

import (
	"errors"
	"fmt"
	"strings"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"gitlab.fbk168.com/gamedevjp/backend-utility/utility/igame"
	"gitlab.fbk168.com/gamedevjp/backend-utility/utility/iserver"
	_ "gitlab.fbk168.com/gamedevjp/backend-utility/utility/mysql"
	"gitlab.fbk168.com/gamedevjp/backend-utility/utility/playerinfo"
	"gitlab.fbk168.com/gamedevjp/backend-utility/utility/restfult"
	"gitlab.fbk168.com/gamedevjp/backend-utility/utility/socket"
	"gitlab.fbk168.com/gamedevjp/backend-utility/utility/user"
	"gitlab.fbk168.com/gamedevjp/sushi/server/game/cache"
	"gitlab.fbk168.com/gamedevjp/sushi/server/game/protoc"
)

var version string = "v1"

// AuthUserURL ...
const AuthUserURL string = "%s/%s/users/%s"

// NewOrderURL ...
const NewOrderURL string = "%s/%s/orders"

// GetOrderURL ...
const GetOrderURL string = "%s/%s/orders/%s"

// Game ...
type Game struct {
	Server    *iserver.Service
	Cache     *cache.GameCache
	IGameRule igame.ISlotRule

	// ProtocolMap map[string]func(r *http.Request) protocol.IProtocol
}

// RESTfulURLs ...
func (g *Game) RESTfulURLs() []restfult.Setting {
	return []restfult.Setting{
		restfult.Setting{
			RequestType: "POST",
			URL:         "game/init",
			Fun:         g.gameinit,
			ConnType:    restfult.Client,
		},
		restfult.Setting{
			RequestType: "POST",
			URL:         "game/result",
			Fun:         g.gameresult,
			ConnType:    restfult.Client,
		},
	}
}

// SocketURLs ...
func (g *Game) SocketURLs() []socket.Setting {
	return []socket.Setting{
		socket.Setting{
			URL: "lobby/createNewSocket",
			Fun: g.createNewSocket,
		},
	}
}

// NewUser *Not Use
func (g *Game) NewUser(token, gameAccount string) *user.Info {
	return &user.Info{}
}

// GetUser ...
func (g *Game) GetUser(userToken string) (*user.Info, *protoc.Error, error) {
	if g.Server.Setting.ServerMod == "dev" {
		return &user.Info{
			UserServerInfo: &playerinfo.AccountInfo{},
			UserGameInfo: &playerinfo.Info{
				IDStr:  "devtest",
				Money:  10000000,
				MoneyU: 10000000,
			},
		}, nil, nil
	}

	tokens := strings.Split(userToken, " ")
	if len(tokens) < 2 {
		return nil, nil, errors.New("token error")
	}

	res, err := g.Server.Transfer.AuthUser(fmt.Sprintf(AuthUserURL, g.Server.Transfer.Path, version, tokens[1]))
	if err != nil {
		if res != nil {
			errorProto := &protoc.Error{}
			if jserr := errorProto.XXX_Unmarshal(res); jserr != nil {
				return nil, nil, jserr
			}
			return nil, errorProto, err
		}
		return nil, nil, err
	}

	userProto := &protoc.User{}
	if jserr := userProto.XXX_Unmarshal(res); jserr != nil {
		return nil, nil, jserr
	}

	return &user.Info{
		UserServerInfo: &playerinfo.AccountInfo{},
		UserGameInfo: &playerinfo.Info{
			IDStr:  userProto.GetUserId(),
			Money:  int64(userProto.GetBalance()),
			MoneyU: userProto.GetBalance(),
		},
	}, nil, nil
}

// NewOrder ...
func (g *Game) NewOrder(token, userIDStr string, betMoney int64) (*protoc.Order, *protoc.Error, error) {
	if g.Server.Setting.ServerMod == "dev" {
		return &protoc.Order{
			UserId:  userIDStr,
			GameId:  g.IGameRule.GetGameTypeID(),
			Bet:     uint64(betMoney),
			OrderId: "testOrder",
		}, nil, nil
	}

	orderProto := &protoc.Order{
		UserId: userIDStr,
		GameId: g.IGameRule.GetGameTypeID(),
		Bet:    uint64(betMoney),
	}
	payload, err := proto.Marshal(orderProto)
	if err != nil {
		return nil, nil, err
	}

	res, err := g.Server.Transfer.NewOrder(fmt.Sprintf(NewOrderURL, g.Server.Transfer.Path, version), token, payload)
	if err != nil {
		if res != nil {
			errorProto := &protoc.Error{}
			if jserr := errorProto.XXX_Unmarshal(res); jserr != nil {
				return nil, nil, jserr
			}
			return nil, errorProto, err
		}
		return nil, nil, err
	}

	if jserr := orderProto.XXX_Unmarshal(res); jserr != nil {
		return nil, nil, jserr
	}
	return orderProto, nil, nil

}

// EndOrder ...
func (g *Game) EndOrder(token string, orderProto *protoc.Order) (*protoc.Order, *protoc.Error, error) {
	orderProto.CompletedAt = ptypes.TimestampNow()
	if g.Server.Setting.ServerMod == "dev" {
		return orderProto, nil, nil
	}

	payload, err := proto.Marshal(orderProto)
	if err != nil {
		return nil, nil, err
	}

	res, err := g.Server.Transfer.EndOrder(fmt.Sprintf(GetOrderURL, g.Server.Transfer.Path, version, orderProto.GetOrderId()), token, payload)
	if err != nil {
		if res != nil {
			errorProto := &protoc.Error{}
			if jserr := errorProto.XXX_Unmarshal(res); jserr != nil {
				return nil, nil, jserr
			}
			return nil, errorProto, err
		}
		return nil, nil, err
	}

	if jserr := orderProto.XXX_Unmarshal(res); jserr != nil {
		return nil, nil, jserr
	}
	return orderProto, nil, nil
}

// GetUserByGameID ...
func (g *Game) GetUserByGameID(token string, userID int64) (*user.Info, error) {
	return &user.Info{}, nil
}

// CheckGameType *Not Use
func (g *Game) CheckGameType(clientGameTypeID string) bool {
	return true
}

// CheckToken *Not Use
func (g *Game) CheckToken(token string) error {
	if serverToken, err := g.getToken(); err != nil {
		return errors.New("getToken error: ")
	} else if serverToken != token {
		return errors.New("token not equal: ")
	}
	return nil
}

func (g *Game) getToken() (string, error) {
	token, err := g.Server.HTTPConn.HTTPPostRawRequest("", nil)
	if err != nil {
		return "", err
	}
	return string(token), nil
}
