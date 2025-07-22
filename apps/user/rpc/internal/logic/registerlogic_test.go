package logic

import (
	"PaiPai/apps/user/rpc/user"
	"context"
	"testing"
)

// 单例测试

func TestRegisterLogic_Register(t *testing.T) {
	type args struct {
		in *user.RegisterReq
	}
	tests := []struct {
		name      string
		args      args
		wantPrint bool
		wantErr   bool
	}{
		{
			"1", args{in: &user.RegisterReq{
				Phone:    "11111111111",
				Nickname: "PaiCloud",
				Password: "123456",
				Avatar:   "pai.jpg",
				Sex:      1,
			}}, true, false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewRegisterLogic(context.Background(), svcCtx)
			got, err := l.Register(tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("Register() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantPrint {
				t.Log(tt.name, got)
			}
		})
	}
}
