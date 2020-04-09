package cache

import (
	"fmt"
	"time"

	"github.com/YWJSonic/ServerUtility/cacheinfo"
	"github.com/YWJSonic/ServerUtility/code"
	"github.com/YWJSonic/ServerUtility/messagehandle"
	"github.com/gomodule/redigo/redis"
)

// Setting ...
type Setting struct {
	ConnectTimeout, ReadTimeout, WriteTimeout, CacheDeleteTime time.Duration
	URL                                                        string
}

// GameCache ICache
type GameCache struct {
	cachePool *redis.Pool
	Setting   Setting
}

// GetCachePool ...
func (c *GameCache) GetCachePool() *redis.Pool {
	if c.cachePool == nil {
		c.cachePool = &redis.Pool{
			MaxIdle:     50,
			IdleTimeout: 240 * time.Second,
			MaxActive:   50,
			Wait:        true,
			Dial: func() (redis.Conn, error) {
				c, err := redis.Dial("tcp", c.Setting.URL,
					redis.DialConnectTimeout(c.Setting.ConnectTimeout),
					redis.DialReadTimeout(c.Setting.ReadTimeout),
					redis.DialWriteTimeout(c.Setting.WriteTimeout))
				if err != nil {
					messagehandle.ErrorLogPrintln("newCachePool-1", c, err)
					return nil, fmt.Errorf("redis connection error: %s", err)
				}
				//验证redis密码
				// if _, authErr := c.Do("AUTH", RedisPassword); authErr != nil {
				// 	return nil, fmt.Errorf("redis auth password error: %s", authErr)
				// }
				return c, nil
			},
			TestOnBorrow: func(c redis.Conn, t time.Time) error {
				_, err := c.Do("PING")
				if err != nil {
					return fmt.Errorf("ping redis error: %s", err)
				}
				return nil
			},
		}
	}
	return c.cachePool
}

// GetToken ...
func (c *GameCache) GetToken(gameAccount string) string {
	value, err := cacheinfo.GetString(c.GetCachePool(), fmt.Sprintf("TOK%s", gameAccount))

	if err != nil {
		return ""
	}
	return value
}

// SetToken ...
func (c *GameCache) SetToken(gameAccount, Token string) {
	now := time.Now()
	lastHour := 23 - now.Hour()
	lastMinute := 59 - now.Minute()
	lastsecod := 60 - now.Second()
	lasttime := time.Duration(lastHour*60*60+lastMinute*60+lastsecod) * time.Second
	cacheinfo.RunSet(c.GetCachePool(), fmt.Sprintf("TOK%s", gameAccount), Token, lasttime)
}

// GetAccountInfo Get Account Struct
func (c *GameCache) GetAccountInfo(gameAccount string) (interface{}, messagehandle.ErrorMsg) {
	err := messagehandle.New()
	info, errMsg := cacheinfo.Get(c.GetCachePool(), fmt.Sprintf("ACC%s", gameAccount))

	if errMsg != nil {
		err.ErrorCode = code.FailedPrecondition
		err.Msg = fmt.Sprintln(errMsg)
		messagehandle.ErrorLogPrintln("GetAccountInfo-1", errMsg, gameAccount)
		return nil, err
	}

	return info, err
}

// SetAccountInfo Set Account Struct
func (c *GameCache) SetAccountInfo(gameAccount string, Value interface{}) {
	cacheinfo.RunSet(c.GetCachePool(), fmt.Sprintf("ACC%s", gameAccount), Value, c.Setting.CacheDeleteTime)
}

// GetPlayerInfo Get PlayerInfo Struct
func (c *GameCache) GetPlayerInfo(playerid int64) (interface{}, messagehandle.ErrorMsg) {
	err := messagehandle.New()
	info, errMsg := cacheinfo.Get(c.GetCachePool(), fmt.Sprintf("ID%dJS", playerid))

	if errMsg != nil {
		err.ErrorCode = code.FailedPrecondition
		err.Msg = fmt.Sprintln(errMsg)
		messagehandle.ErrorLogPrintln("GetPlayerInfo-1", errMsg, playerid)
		return nil, err
	}

	return info, err
}

// SetPlayerInfo Set PlayerInfo Struct
func (c *GameCache) SetPlayerInfo(playerid int64, Value interface{}) {
	cacheinfo.RunSet(c.GetCachePool(), fmt.Sprintf("ID%dJS", playerid), Value, c.Setting.CacheDeleteTime)
}

// ClearPlayerCache ...
func (c *GameCache) ClearPlayerCache(playerid int64, gameAccount string) {
	cacheinfo.Del(c.GetCachePool(), fmt.Sprintf("ID%dJS", playerid))
	cacheinfo.Del(c.GetCachePool(), fmt.Sprintf("ACC%s", gameAccount))
	cacheinfo.Del(c.GetCachePool(), fmt.Sprintf("TOK%s", gameAccount))
}

// ClearAllCache ...
func (c *GameCache) ClearAllCache() {
	cacheinfo.RunFlush(c.GetCachePool())
}

//------ third party request -------

// SetULGInfo Set ULG info
func (c *GameCache) SetULGInfo(playerid int64, value interface{}) {
	key := fmt.Sprintf("ULG%d", playerid)
	cacheinfo.RunSet(c.GetCachePool(), key, value, c.Setting.CacheDeleteTime)
}

// GetULGInfoCache Get ULG info
func (c *GameCache) GetULGInfoCache(playerid int64) interface{} {
	err := messagehandle.New()
	key := fmt.Sprintf("ULG%d", playerid)
	info, errMsg := cacheinfo.Get(c.GetCachePool(), key)

	if errMsg != nil {
		err.ErrorCode = code.FailedPrecondition
		err.Msg = fmt.Sprintln(errMsg)
		messagehandle.ErrorLogPrintln("GetULGInfoCache-1", errMsg, playerid)
		return nil
	}

	return info
}

//------ game info per each player -----

// SetAttach ...
func (c *GameCache) SetAttach(playerid int64, value interface{}) {
	key := fmt.Sprintf("attach%d", playerid)
	cacheinfo.RunSet(c.GetCachePool(), key, value, c.Setting.CacheDeleteTime)
}

// GetAttach game data request
func (c *GameCache) GetAttach(playerid int64) interface{} {
	err := messagehandle.New()
	key := fmt.Sprintf("attach%d", playerid)
	info, errMsg := cacheinfo.Get(c.GetCachePool(), key)

	if errMsg != nil {
		err.ErrorCode = code.FailedPrecondition
		err.Msg = fmt.Sprintln(errMsg)
		messagehandle.ErrorLogPrintln("GetAttach-1", errMsg, key)
		return nil
	}

	return info
}
