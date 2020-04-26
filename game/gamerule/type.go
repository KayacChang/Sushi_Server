package gamerule

import (
	"github.com/YWJSonic/ServerUtility/igame"
)

// Rule ...
type Rule struct {
	BetLine             int       `json:"BetLine"`
	BetRate             []int64   `json:"BetRate"`
	BetRateDefaultIndex int64     `json:"BetRateDefaultIndex"`
	BetRateLinkIndex    []int64   `json:"BetRateLinkIndex"`
	BonusRate           []int64   `json:"BonusRate"`
	BonusWeightings     [][]int   `json:"BonusWeightings"`
	FreeGameCountArray  []int     `json:"FreeGameCountArray"`
	FreeReel            [][][]int `json:"FreeReel"`
	GameIndex           int64     `json:"GameIndex"`
	GameTypeID          string    `json:"GameTypeID"`
	IsAttachInit        bool      `json:"IsAttachInit"`
	ItemResults         [][][]int `json:"ItemResults"`
	Items               []int     `json:"Items"`
	LineMap             [][]int   `json:"LineMap"`
	NormalReel          [][][]int `json:"NormalReel"`
	NormalReelSize      []int     `json:"NormalReelSize"`
	RTPSetting          []int     `json:"RTPSetting"`
	ScotterItemIndex    []int     `json:"ScotterItemIndex"`
	ScotterGameLimit    []int     `json:"ScotterGameLimit"`
	ScotterLineRate     [][]int   `json:"scotterLineRate"`
	ServerDayPayDefault int64     `json:"ServerDayPayDefault"`
	ServerDayPayLimit   int64     `json:"ServerDayPayLimit"`
	Version             string    `json:"Version"`
	WildsItemIndex      []int     `json:"WildsItemIndex"`
	WinBetRateLimit     int64     `json:"WinBetRateLimit"`
	WinScoreLimit       int64     `json:"WinScoreLimit"`
}

// GetGameIndex ...
func (r *Rule) GetGameIndex() int64 {
	return r.GameIndex
}

// GetGameTypeID ...
func (r *Rule) GetGameTypeID() string {
	return r.GameTypeID
}

// GetBetMoney ...
func (r *Rule) GetBetMoney(index int64) int64 {
	return r.BetRate[index]
}

// GetReel ...
func (r *Rule) GetReel() map[string][][]int {
	scrollmap := map[string][][]int{
		"normalreel": r.normalReel(),
		"freereel":   r.freeReel(),
	}
	return scrollmap
}

// GetBetSetting ...
func (r *Rule) GetBetSetting() map[string]interface{} {
	tmp := make(map[string]interface{})
	tmp["betrate"] = r.BetRate                         //betRate
	tmp["betratelinkindex"] = r.BetRateLinkIndex       //betRateLinkIndex
	tmp["betratedefaultindex"] = r.BetRateDefaultIndex //betRateDefaultIndex
	return tmp
}

// CheckGameType ...
func (r *Rule) CheckGameType(gameTypeID string) bool {
	if r.GameTypeID != gameTypeID {
		return false
	}
	return true
}

func (r *Rule) normalReel() [][]int {
	return r.NormalReel[r.normalRTP()]
}
func (r *Rule) freeReel() [][]int {
	return r.FreeReel[r.freeRTP()]
}
func (r *Rule) normalRTP() int {
	return r.RTPSetting[0]
}
func (r *Rule) freeRTP() int {
	return r.RTPSetting[1]
}
func (r *Rule) getBonusWeightings() []int {
	return r.BonusWeightings[r.RTPSetting[0]]
}

func (r *Rule) wild() int {
	return r.WildsItemIndex[0]
}
func (r *Rule) scotter1() int {
	return r.ScotterItemIndex[0]
}

func (r *Rule) scotter1GameLimit() int {
	return r.ScotterGameLimit[0]
}

func (r *Rule) scotter1LineRate(index int) int {
	return r.ScotterLineRate[0][index]
}

func (r *Rule) scotter2() int {
	return r.ScotterItemIndex[1]
}

func (r *Rule) scotter2GameLimit() int {
	return r.ScotterGameLimit[1]
}
func (r *Rule) scotter2LineRate(index int) int {
	return r.ScotterLineRate[1][index]
}

// GameRequest ...
func (r *Rule) GameRequest(config *igame.RuleRequest) *igame.RuleRespond {
	betMoney := r.GetBetMoney(config.BetIndex)
	result := make(map[string]interface{})
	otherData := make(map[string]interface{})
	var totalWin int64

	gameResult := r.newlogicResult(betMoney)

	result["normalresult"] = gameResult.Normalresult
	result["isfreegame"] = 0
	result["freegamecount"] = gameResult.FreeGameCount
	totalWin += gameResult.Normaltotalwin

	if gameResult.FreeResult != nil {
		result["freeresult"] = gameResult.FreeResult
		result["isfreegame"] = 1
		totalWin += gameResult.FreeTotalWin
	}

	result["totalwinscore"] = totalWin

	// fmt.Println(foundation.JSONToString(result))
	return &igame.RuleRespond{
		BetMoney:      betMoney,
		Totalwinscore: totalWin,
		GameResult:    result,
		OtherData:     otherData,
	}
}
