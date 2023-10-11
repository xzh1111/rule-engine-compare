package engine

import (
	"context"
	"flag"
	"fmt"
	"log"
	"strconv"
	"sync"

	"github.com/hyperjumptech/grule-rule-engine/ast"
	"github.com/hyperjumptech/grule-rule-engine/builder"
	"github.com/hyperjumptech/grule-rule-engine/engine"
	"github.com/hyperjumptech/grule-rule-engine/pkg"
)

var rule_dir = flag.String("rule_dir", "./data/config", "rule Directory")

var grulePoolMap map[string]*sync.Pool
var rb *builder.RuleBuilder
var klib *ast.KnowledgeLibrary

func gruleInit() {
	klib = ast.NewKnowledgeLibrary()
	rb = builder.NewRuleBuilder(klib)
	grulePoolMap = make(map[string]*sync.Pool)
}

func LoadGRules() error {
	once.Do(gruleInit)
	ruleFiles, err := loadRuleFiles(*rule_dir, ".grl")
	if err != nil {
		return fmt.Errorf("error loadGrlFiles: %v", err)
	}

	for _, ruleFile := range ruleFiles {
		if !ruleFile.Change {
			continue
		}
		err := rb.BuildRuleFromResource(ruleFile.FileName, strconv.Itoa(ruleFile.Version), pkg.NewFileResource(*rule_dir+"/"+ruleFile.FileName))
		if err != nil {
			log.Fatalf("BuildRuleFromResource, ruleFile: %+v , error: %v\n", ruleFile, err)
			return err
		}
		if _, ok := grulePoolMap[ruleFile.FileName]; ok {
			log.Printf("grulePoolMap update for %s", ruleFile.FileName)
		} else {
			log.Printf("grulePoolMap create for %s", ruleFile.FileName)
		}
		mu.Lock()
		grulePoolMap[ruleFile.FileName] = &sync.Pool{
			New: func() interface{} {
				kb, _ := klib.NewKnowledgeBaseInstance(ruleFile.FileName, strconv.Itoa(ruleFile.Version))
				return kb
			},
		}
		mu.Unlock()

	}
	return nil
}

func GruleExecute(ctx context.Context, facts map[string]interface{}, KnowledgeName string) error {
	e := engine.NewGruleEngine()
	dataCtx := ast.NewDataContext()
	var err error
	for k, v := range facts {
		err = dataCtx.Add(k, v)
		if err != nil {
			log.Printf("DataContext, Add %s fact failed , error: %v\n", k, err)
			return err
		}
	}
	mu.RLock()
	kbPool, ok := grulePoolMap[KnowledgeName]
	mu.RUnlock()
	if !ok {
		return fmt.Errorf("KnowledgeBase %s not found", KnowledgeName)
	} else {
		kb := kbPool.Get().(*ast.KnowledgeBase)
		err = e.Execute(dataCtx, kb)
		defer kbPool.Put(kb)
	}

	return err
}
