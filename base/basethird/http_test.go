package basethird

import (
	"fmt"
	"reflect"
	"testing"
)

func TestCopyMap(t *testing.T) {
	type args struct {
		original map[string]string
	}
	tests := []struct {
		name string
		args args
		want map[string]string
	}{
		{
			name: "test1",
			args: args{
				original: map[string]string{
					"key1": "value1",
					"key2": "value2",
				},
			},
			want: map[string]string{
				"key1": "value1",
				"key2": "value2",
			},
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CopyMap(tt.args.original)
			for k, v := range got {
				got[k] = fmt.Sprintf("%s-new", v)
			}
			t.Logf("%v", got)
			if !reflect.DeepEqual(tt.args.original, tt.want) {
				t.Errorf("CopyMap() = %v, want %v", tt.args.original, tt.want)
			}
		})
	}
}
