package gamerule

import (
	"github.com/YWJSonic/ServerUtility/foundation"
	"github.com/YWJSonic/ServerUtility/foundation/math"
	"github.com/YWJSonic/ServerUtility/gameplate"
)

type result struct {
	Normalresult   map[string]interface{}
	Otherdata      map[string]interface{}
	Normaltotalwin int64
	FreeResult     []map[string]interface{}
	FreeTotalWin   int64
	FreeGameCount  int
}

// // Result att 0: freecount
// func (r *Rule) newlogicResult(betMoney int64) result {
// 	normalresult, otherdata, normaltotalwin := r.outputGame(betMoney)

// 	if otherdata["isrespin"].(int) != 1 {
// 		return result{
// 			Normalresult:   normalresult,
// 			Otherdata:      otherdata,
// 			Normaltotalwin: normaltotalwin,
// 		}
// 	}
// 	respinresult, respintotalwin := r.outRespin(betMoney)
// 	return result{
// 		Normalresult:   normalresult,
// 		Otherdata:      otherdata,
// 		Normaltotalwin: normaltotalwin,
// 		Respinresult:   respinresult,
// 		Respintotalwin: respintotalwin,
// 	}
// }

// Result ...
func (r *Rule) newlogicResult(betMoney int64) result {
	option := gameplate.PlateOption{
		Scotter: []int{r.scotter1(), r.scotter2()},
		Wild:    []int{r.wild()},
	}

	normalresult, otherdata, normaltotalwin := r.outputGame(betMoney, option)
	// fmt.Println("----normalresult----", normalresult)
	// fmt.Println("----otherdata----", otherdata)
	// fmt.Println("----normaltotalwin----", normaltotalwin)
	FreeGameCount := foundation.InterfaceToInt(otherdata["freegamecount"])
	// result["freegamecount"] = FreeGameCount
	// result["normalresult"] = normalresult
	// result["isfreegame"] = 0
	// totalWin += normaltotalwin

	if FreeGameCount > 0 {
		freeresult, _, freetotalwin := r.outputFreeGame(betMoney, FreeGameCount, option)
		// fmt.Println("----freeresult----", freeresult)
		// fmt.Println("----freetotalwin----", freetotalwin)
		// result["freeresult"] = freeresult
		// result["isfreegame"] = 1
		// totalWin += freetotalwin
		return result{
			Normalresult:   normalresult,
			Otherdata:      otherdata,
			Normaltotalwin: normaltotalwin,
			FreeResult:     freeresult,
			FreeTotalWin:   freetotalwin,
			FreeGameCount:  FreeGameCount,
		}
	}

	// result["totalwinscore"] = totalWin
	return result{
		Normalresult:   normalresult,
		Otherdata:      otherdata,
		Normaltotalwin: normaltotalwin,
	}

}

// outputGame out put normal game result, mini game status, totalwin
func (r *Rule) outputGame(betMoney int64, option gameplate.PlateOption) (map[string]interface{}, map[string]interface{}, int64) {
	var totalScores int64
	normalResult := make(map[string]interface{})
	otherdata := make(map[string]interface{})

	normalResult, otherdata, totalScores = r.aRound(betMoney, r.normalReel(), option, 1)
	return normalResult, otherdata, totalScores
}

func (r *Rule) outputFreeGame(betMoney int64, freeCount int, option gameplate.PlateOption) ([]map[string]interface{}, map[string]interface{}, int64) {
	var totalScores int64
	// var wildCount, bonusRate int
	otherdata := make(map[string]interface{})
	var freeResult []map[string]interface{}

	for i, imax := 0, freeCount; i < imax; i++ {
		tmpResult, _, tmpTotalScores := r.aRound(betMoney, r.freeReel(), option, 2)
		totalScores += tmpTotalScores
		freeResult = append(freeResult, tmpResult)
	}
	return freeResult, otherdata, totalScores
}

func (r *Rule) aRound(betMoney int64, scorll [][]int, option gameplate.PlateOption, gameType int) (map[string]interface{}, map[string]interface{}, int64) {

	var winLineInfo = make([]interface{}, 0)
	var totalScores int64
	var freeGameCount, scotterCount int
	// var paylinestr string
	var isLink bool
	var scotterLineSymbol, scotterLinePoint [][]int
	var plateSymbolCollectResult map[string]interface{}
	result := map[string]interface{}{
		"bonusrate":  int64(0),
		"bonusscore": int64(0),
	}
	otherdata := map[string]interface{}{
		"isfreegame":    0,
		"freegamecount": 0,
	}

	plateIndex, plateSymbol := gameplate.NewPlate2D(r.NormalReelSize, scorll)
	// plateSymbol = [][]int{
	// 	{11, 0, 9},
	// 	{0, 11, 6},
	// 	{13, 4, 12},
	// 	{0, 14, 12},
	// 	{12, 10, 1},
	// }
	plateLineMap := gameplate.PlateToLinePlate(plateSymbol, r.LineMap)

	for lineIndex, plateLine := range plateLineMap {
		newLine := gameplate.CutSymbolLink(plateLine, option) // cut line to win line point
		mulityLine := gameplate.LineMulitResult(newLine, option)

		if len(mulityLine) > 1 {
			isLink = false
			infoLine := gameplate.NewInfoLine()
			for _, winLine := range mulityLine {
				for _, payLine := range r.ItemResults[len(winLine)] {
					if r.isWin(winLine, payLine, option) {
						isLink = true
						tmpline := r.winResult(betMoney, lineIndex, newLine, payLine, option, gameType)
						if tmpline.Score > infoLine.Score {
							infoLine = tmpline
							infoLine.LineWinIndex = lineIndex
						}
					}
				}
			}
			if isLink {
				totalScores += infoLine.Score
				winLineInfo = append(winLineInfo, infoLine)
			}
		} else {
			for _, payLine := range r.ItemResults[len(newLine)] { // win line result group
				if r.isWin(newLine, payLine, option) { // win result check
					infoLine := r.winResult(betMoney, lineIndex, newLine, payLine, option, gameType)
					infoLine.LineWinIndex = lineIndex
					totalScores += infoLine.Score
					winLineInfo = append(winLineInfo, infoLine)
				}
			}
		}
	}

	// scotter 1 handle
	plateSymbolCollectResult = gameplate.PlateSymbolCollect(r.scotter1(), plateSymbol, option, map[string]interface{}{"isincludewild": false, "isseachallplate": true})
	scotterCount = foundation.InterfaceToInt(plateSymbolCollectResult["targetsymbolcount"])
	scotterCount = math.ClampInt(scotterCount, 0, len(r.FreeGameCountArray))
	scotterLineSymbol = plateSymbolCollectResult["symbolnumcollation"].([][]int)
	scotterLinePoint = plateSymbolCollectResult["symbolpointcollation"].([][]int)

	if scotterCount >= r.scotter1GameLimit() {
		infoLine := gameplate.NewInfoLine()

		for i, max := 0, len(scotterLineSymbol); i < max; i++ {
			if len(scotterLineSymbol[i]) > 0 {
				infoLine.AddNewLine(scotterLineSymbol[i], scotterLinePoint[i], option)
			} else {
				infoLine.AddEmptyPoint()
			}
		}

		infoLine.LineWinRate = r.scotter1LineRate(scotterCount)
		infoLine.Score = int64(infoLine.LineWinRate) * betMoney
		totalScores += infoLine.Score
		winLineInfo = append(winLineInfo, infoLine)

		freeGameCount = r.FreeGameCountArray[scotterCount]
		otherdata["freegamecount"] = freeGameCount
		otherdata["isfreegame"] = 1

	}

	// scotter 2 handle
	plateSymbolCollectResult = gameplate.PlateSymbolCollect(r.scotter2(), plateSymbol, option, map[string]interface{}{"isincludewild": false, "isseachallplate": true})
	scotterCount = foundation.InterfaceToInt(plateSymbolCollectResult["targetsymbolcount"])
	scotterCount = math.ClampInt(scotterCount, 0, len(r.FreeGameCountArray))
	scotterLineSymbol = plateSymbolCollectResult["symbolnumcollation"].([][]int)
	scotterLinePoint = plateSymbolCollectResult["symbolpointcollation"].([][]int)

	if scotterCount >= r.scotter2GameLimit() {
		infoLine := gameplate.NewInfoLine()

		for i, max := 0, len(scotterLineSymbol); i < max; i++ {
			if len(scotterLineSymbol[i]) > 0 {
				infoLine.AddNewLine(scotterLineSymbol[i], scotterLinePoint[i], option)
			} else {
				infoLine.AddEmptyPoint()
			}
		}

		infoLine.LineWinRate = r.scotter2LineRate(scotterCount)
		infoLine.Score = int64(infoLine.LineWinRate) * betMoney
		totalScores += infoLine.Score
		winLineInfo = append(winLineInfo, infoLine)

		bonusrate := r.BonusRate[foundation.RangeRandom(r.getBonusWeightings())]
		result["bonusrate"] = bonusrate
		result["bonusscore"] = bonusrate * betMoney
		totalScores += bonusrate * betMoney

	}

	if len(winLineInfo) > 0 {
		result = foundation.AppendMap(result, gameplate.ResultMapLine(plateIndex, plateSymbol, winLineInfo))
	} else {
		result = foundation.AppendMap(result, gameplate.ResultMapLine(plateIndex, plateSymbol, []interface{}{}))
	}

	result["gameresult"] = winLineInfo
	result["scores"] = totalScores
	return result, otherdata, totalScores
}

// isWin symbol line compar parline is win
func (r *Rule) isWin(lineSymbol []int, payLineSymbol []int, option gameplate.PlateOption) bool {
	targetSymbol := 0
	isWin := true
	EmptyNum := option.EmptyNum()
	mainSymbol := EmptyNum

	for lineIndex, max := 0, len(payLineSymbol)-1; lineIndex < max; lineIndex++ {
		targetSymbol = lineSymbol[lineIndex]

		if isWild, _ := option.IsWild(targetSymbol); isWild {
			if mainSymbol == EmptyNum {
				mainSymbol = targetSymbol
			}
			continue
		}

		switch payLineSymbol[lineIndex] {
		case targetSymbol:
			mainSymbol = targetSymbol
		default:
			isWin = false
			return isWin
		}
	}

	if mainSymbol != payLineSymbol[0] {
		return false
	}

	return isWin
}

func (r *Rule) winResult(betMoney int64, lineIndex int, newLine, payLine []int, option gameplate.PlateOption, gameType int) gameplate.InfoLine {
	mainSymbol := payLine[0]
	infoLine := gameplate.NewInfoLine()

	for i, max := 0, len(payLine)-1; i < max; i++ {
		infoLine.AddNewPoint(newLine[i], r.LineMap[lineIndex][i], option)
	}

	if isScotter, _ := option.IsScotter(mainSymbol); isScotter {
		infoLine.LineWinRate = payLine[len(payLine)-1]
		infoLine.Score = int64(infoLine.LineWinRate) * betMoney
	} else {
		infoLine.LineWinRate = payLine[len(payLine)-1]
		infoLine.Score = int64(infoLine.LineWinRate) * (betMoney / int64(r.BetLine))
	}
	return infoLine
}
