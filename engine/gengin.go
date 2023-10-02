package engine

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/cookedsteak/gengine/engine"
)

var genginPoolMap map[string]*engine.GenginePool = make(map[string]*engine.GenginePool)

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
		pool, err := engine.NewGenginePool(2, 100, 1, string(content), nil)
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
