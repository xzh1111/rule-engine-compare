package engine

import (
	"context"
	"fmt"
	"log"
	"math"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/cookedsteak/gengine/engine"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
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

func Printf(format string, a ...any) {
	fmt.Printf(format, a...)
	fmt.Print("\n")
}

func SliceLen(s []string) int {
	return len(s)
}

func StringLen(s string) int {
	return len(s)
}

func Len(i interface{}) int {
	v := reflect.ValueOf(i)

	switch v.Kind() {
	case reflect.String:
		return len(v.String())
	case reflect.Slice:
		return v.Len()
	default:
		fmt.Printf("customLen does not support this type: %s\n", v.Kind())
		return -1
	}
}

func FindRange(val float64, ranges ...int64) []int64 {
	for i := 0; i < len(ranges)-1; i++ {
		if float64(ranges[i]) < val && val <= float64(ranges[i+1]) {
			return []int64{ranges[i], ranges[i+1]}
		}
	}
	return []int64{}
}

var sugar *zap.SugaredLogger

func init() {
	config := zap.Config{
		Encoding:         "console",
		Level:            zap.NewAtomicLevelAt(zap.InfoLevel),
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig: zapcore.EncoderConfig{
			EncodeTime: zapcore.ISO8601TimeEncoder,
		},
	}
	logger, _ := config.Build()
	sugar = logger.Sugar()
	apiOuters := map[string]interface{}{
		"ZapLogf":        sugar.Infof,
		"Len":            Len,
		"Sprintf":        fmt.Sprintf,
		"Printf":         Printf,
		"Println":        fmt.Println,
		"Now":            time.Now,
		"Abs":            math.Abs,
		"MathLog":        math.Log,
		"Log10":          math.Log10,
		"Log1p":          math.Log1p,
		"Log2":           math.Log2,
		"Mod":            math.Mod,
		"Pow":            math.Pow,
		"Pow10":          math.Pow10,
		"Round":          math.Round,
		"ToLower":        strings.ToLower,
		"StringContains": strings.Contains,
		"Atoi":           strconv.Atoi,
		"Itoa":           strconv.Itoa,
		"StringSplit":    strings.Split,
		"SliceLen":       SliceLen,
		"StringLen":      StringLen,
		"FindRange":      FindRange,
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
