package comparer

import (
	"testing"

	"github.com/illuminarean-labs/yaml-diff-reporter/domain"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	type args struct {
		config Config
	}
	tests := []struct {
		name string
		args args
		want Comparer
	}{
		{
			name: "성공",
			args: args{
				config: Config{
					IgnoredKeys: []string{"hello"},
					Modes:       domain.CompareModes{Type, Value},
				},
			},
			want: comparer{
				results: &domain.ErrorResults{},
				config: Config{
					IgnoredKeys: []string{"hello"},
					Modes:       domain.CompareModes{Type, Value},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := New(tt.args.config)
			assert.Equalf(t, got, tt.want, "New() = %v, want %v", got, tt.want)
		})
	}
}

func Test_sliceKey(t *testing.T) {
	type args struct {
		parent string
		index  int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "성공",
			args: args{
				parent: "hello.world",
				index:  3,
			},
			want: "hello.world[3]",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := sliceKey(tt.args.parent, tt.args.index)
			assert.Equalf(t, got, tt.want, "sliceKey() = %v, want %v", got, tt.want)
		})
	}
}

func Test_mapKey(t *testing.T) {
	type args struct {
		parent string
		key    string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "parent가 빈 문자열인 경우",
			args: args{
				parent: "",
				key:    "key",
			},
			want: "key",
		},
		{
			name: "parent, key 모두 있는 경우",
			args: args{
				parent: "hello",
				key:    "world",
			},
			want: "hello.world",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := mapKey(tt.args.parent, tt.args.key)
			assert.Equalf(t, got, tt.want, "mapKey() = %v, want %v", got, tt.want)
		})
	}
}

func Test_comparer_Compare(t *testing.T) {
	type fields struct {
		Results *domain.ErrorResults
		Config  Config
	}
	type args struct {
		parent string
		lhs    any
		rhs    any
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name:   "",
			fields: fields{},
			args:   args{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := comparer{
				results: tt.fields.Results,
				config:  tt.fields.Config,
			}
			c.Compare(tt.args.parent, tt.args.lhs, tt.args.rhs)
		})
	}
}
