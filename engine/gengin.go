package engine

import (
	"context"
	"fmt"
	"log"
	"math"
	"os"
	"strings"
	"time"

	"github.com/cookedsteak/gengine/engine"
)

var genginPoolMap map[string]*engine.GenginePool = make(map[string]*engine.GenginePool)

var ApiOuters map[string]interface{} = make(map[string]interface{})

func AddApiOuters(apiOuters map[string]interface{}) error {
	for k, v := range apiOuters {
		ApiOuters[k] = v
	}
	return nil
}

func AddApiOuter(name string, api interface{}) error {
	ApiOuters[name] = api
	return nil
}

func init() {
	apiOuters := map[string]interface{}{
		"Println":        fmt.Println,
		"Now":            time.Now(),
		"StringContains": strings.Contains,
		"Abs":            math.Abs,
		"MathLog":        math.Log,
		"Log10":          math.Log10,
		"Log1p":          math.Log1p,
		"Log2":           math.Log2,
		"Mod":            math.Mod,
		"Pow":            math.Pow,
		"Pow10":          math.Pow10,
		"Round":          math.Round,
	}

	AddApiOuters(apiOuters)
}

func LoadGenginRules() error {
	ruleFiles, err := loadRuleFiles(*rule_dir, ".erl")
	if err != nil {
		return fmt.Errorf("error loadGrlFiles: %v", err)
	}

	for _, ruleFile := range ruleFiles {
		if !ruleFile.Change {
			continue
		}
		content, err := os.ReadFile(*rule_dir + "/" + ruleFile.FileName)
		if err != nil {
			return err
		}
		if _, ok := genginPoolMap[ruleFile.FileName]; ok {
			log.Printf("genginPoolMap update for %s", ruleFile.FileName)
		} else {
			log.Printf("genginPoolMap create for %s", ruleFile.FileName)
		}
		mu.Lock()
		pool, err := engine.NewGenginePool(2, 100, 1, string(content), ApiOuters)
		if err != nil {
			panic(fmt.Sprintf("init gengine failed, err:%+v", err))
		}
		genginPoolMap[ruleFile.FileName] = pool
		mu.Unlock()

	}
	return nil
}

func GenginExecute(ctx context.Context, facts map[string]interface{}, KnowledgeName string) error {
	var err error
	mu.RLock()
	pool, ok := genginPoolMap[KnowledgeName]
	mu.RUnlock()
	if !ok {
		return fmt.Errorf("KnowledgeBase %s not found", KnowledgeName)
	} else {
		err, _ = pool.Execute(facts, false)
	}

	return err
}
