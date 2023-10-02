package benchmarks

import (
	"context"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xzh1111/rule-engine-compare/engine"
)

func BenchmarkGruleExecute(b *testing.B) {

	err := engine.SyncRules(engine.LoadGRules)
	assert.Equal(b, err, nil)

	ufs := []*UserFact{
		{
			UserKey: "playone",
			Itier:   16,
			RoundNo: 1,
		},
		{
			UserKey: "playone",
			Itier:   8,
			RoundNo: 1,
		},
		{
			UserKey: "noplayone",
			Itier:   8,
			RoundNo: 1,
		},
		{
			UserKey: "playone",
			Itier:   17,
			RoundNo: 2,
		},
	}

	// 创建测试上下文
	ctx := context.Background()

	knowRules := []string{"rule1.grl", "rule2.grl"}
	// 循环运行测试函数
	for i := 0; i < b.N; i++ {
		idx := rand.Intn(len(ufs))
		uf := ufs[idx]
		for _, knowRule := range knowRules {
			err := engine.GruleExecute(ctx, map[string]interface{}{"UserFact": uf}, knowRule)
			if err != nil {
				b.Errorf("Error executing GruleExecute: %s %v", knowRule, err)
			}
		}
	}
}
