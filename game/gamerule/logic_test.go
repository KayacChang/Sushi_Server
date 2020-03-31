package gamerule

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/YWJSonic/ServerUtility/foundation/fileload"
	"github.com/YWJSonic/ServerUtility/igame"
)

func TestNew(t *testing.T) {
	gamejsStr := fileload.Load("../../file/gameconfig.json")
	var gameRule = &Rule{}
	if err := json.Unmarshal([]byte(gamejsStr), &gameRule); err != nil {
		panic(errors.New("gameconfig error: "))
	}
	fmt.Println(gameRule.newlogicResult(0))
}

func TestGameRequest(t *testing.T) {
	gamejsStr := fileload.Load("../../file/gameconfig.json")
	var gameRule = &Rule{}
	if err := json.Unmarshal([]byte(gamejsStr), &gameRule); err != nil {
		panic(errors.New("gameconfig error: "))
	}

	for i, count := 0, 200; i < count; i++ {
		result := gameRule.GameRequest(&igame.RuleRequest{
			BetIndex: 0,
		})
		if respin, ok := result.GameResult["isfreegame"]; ok && respin.(int) == 1 {
			fmt.Println(result)
		}
	}
}
