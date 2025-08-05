package main

import (
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	"strings"
	"testing"
)

func Test_Jaeger(t *testing.T) {
	cfg := jaegercfg.Configuration{
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:          true,
			CollectorEndpoint: fmt.Sprintf("http://%s/api/traces", "x.x.x.x:14268"),
		},
	}
	Jaeger, err := cfg.InitGlobalTracer("client test", jaegercfg.Logger(jaeger.StdLogger))
	if err != nil {
		t.Log(err)
		return
	}
	defer Jaeger.Close()

	// 执行任务
	tracer := opentracing.GlobalTracer()

	// 任务节点定义span
	parentSpan := tracer.StartSpan("A")
	defer parentSpan.Finish()

	B(tracer, parentSpan)
}

func B(tracer opentracing.Tracer, parentSpan opentracing.Span) {
	childSpan := tracer.StartSpan("B", opentracing.ChildOf(parentSpan.Context()))
	defer childSpan.Finish()

}

func Test_reverseWords(t *testing.T) {
	fmt.Println([]byte("pai cloud"))
	tests := []struct {
		input  string
		output string
	}{
		{"pai cloud", "cloud pai"}, {" Hello paicloud ", "paicloud Hello"}, {"good luck to you", "you to luck good"}, {"", ""}, {" a ", "a"},
	}
	for _, tt := range tests {
		result := reverseWords(tt.input)
		if result != tt.output {
			t.Errorf("expected:%s\t,result:%s\n", tt.output, result)
		} else {
			fmt.Printf("expected:%s\t,result:%s\n", tt.output, result)
		}
	}
}
func reverseWords(s string) string {
	s = strings.TrimSpace(s)
	sb := []byte(s)
	left, right := 0, 0
	for right < len(s) {
		for 0 < right && right < len(s) && s[right] == ' ' && s[right-1] == ' ' {
			right++
			continue
		}
		sb[left] = sb[right]
		right++
		left++
	}
	sb = sb[0:left]
	s = string(sb)
	ans := make([]string, 0, len(sb))
	j := 0
	for i := 0; i <= len(s); i++ {
		if i == len(s) {
			ans = append(ans, s[j:i])
			break
		}
		if sb[i] == ' ' {
			ans = append(ans, s[j:i])
			j = i + 1
		}
	}
	a := strings.Builder{}
	for i := len(ans) - 1; i >= 0; i-- {
		a.WriteString(ans[i])
		if i != 0 {
			a.WriteString(" ")
		}
	}
	return a.String()
}
