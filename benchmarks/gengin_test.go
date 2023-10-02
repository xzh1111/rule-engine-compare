package benchmarks

import (
	"context"
	"fmt"
	"math/rand"
	"testing"

	"github.com/cookedsteak/gengine/engine"
	"github.com/stretchr/testify/assert"
	cengine "github.com/xzh1111/rule-engine-compare/engine"
)

//业务规则
const service_rules string = `
rule "1" "1"
begin
	resp.At = room.GetAttention()
	println("rule 1...")
end 

rule "2" "2"
begin
	resp.Num = room.GetNum()
	println("rule 2...")
end
`

//业务接口
type MyService struct {
	//gengine pool
	Pool *engine.GenginePool

	//other params
}

//request
type Request struct {
	Rid       int64
	RuleNames []string
	//other params
}

//resp
type Response struct {
	At  int64
	Num int64
	//other params
}

//特定的场景服务
type Room struct {
}

func (r *Room) GetAttention( /*params*/ ) int64 {
	// logic
	return 100
}

func (r *Room) GetNum( /*params*/ ) int64 {
	//logic
	return 111
}

//初始化业务服务
//apiOuter这里最好仅注入一些无状态函数，方便应用中的状态管理
func NewMyService(poolMinLen, poolMaxLen int64, em int, rulesStr string, apiOuter map[string]interface{}) *MyService {
	pool, e := engine.NewGenginePool(poolMinLen, poolMaxLen, em, rulesStr, apiOuter)
	if e != nil {
		panic(fmt.Sprintf("初始化gengine失败，err:%+v", e))
	}

	myService := &MyService{Pool: pool}
	return myService
}

//service
func (ms *MyService) Service(req *Request) (*Response, error) {

	resp := &Response{}

	//基于需要注入接口或数据,data这里最好仅注入与本次请求相关的结构体或数据，便于状态管理
	data := make(map[string]interface{})
	data["req"] = req
	data["resp"] = resp

	//模块化业务逻辑,api
	room := &Room{}
	data["room"] = room

	//
	e, _ := ms.Pool.ExecuteSelectedRules(data, req.RuleNames)
	if e != nil {
		println(fmt.Sprintf("pool execute rules error: %+v", e))
		return nil, e
	}

	return resp, nil
}

//模拟调用
func Test_run(t *testing.T) {

	//初始化
	//注入api，请确保注入的API属于并发安全
	apis := make(map[string]interface{})
	apis["println"] = fmt.Println
	msr := NewMyService(10, 20, 1, service_rules, apis)

	//调用
	req := &Request{
		Rid:       123,
		RuleNames: []string{"1", "2"},
	}
	response, e := msr.Service(req)
	if e != nil {
		println(fmt.Sprintf("service err:%+v", e))
		return
	}

	println("resp result = ", response.At, response.Num)
}

func BenchmarkGenginExecute(b *testing.B) {

	err := cengine.SyncRules(cengine.LoadGenginRules)
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

	knowRules := []string{"rule1.erl","rule2.erl"}
	// 循环运行测试函数
	for i := 0; i < b.N; i++ {
		idx := rand.Intn(len(ufs))
		uf := ufs[idx]
		for _, knowRule := range knowRules {
			err := cengine.GenginExecute(ctx, map[string]interface{}{"UserFact": uf}, knowRule)
			if err != nil {
				b.Errorf("Error executing GruleExecute: %s %v", knowRule, err)
			}
		}
	}
}


func TestGenginExecute(t *testing.T) {
	type args struct {
		ctx           context.Context
		facts         map[string]interface{}
		KnowledgeName string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := cengine.GenginExecute(tt.args.ctx, tt.args.facts, tt.args.KnowledgeName); (err != nil) != tt.wantErr {
				t.Errorf("GenginExecute() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
