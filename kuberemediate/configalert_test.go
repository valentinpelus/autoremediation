package kuberemediate

import "testing"

func TestLoadConfAlert(t *testing.T) {
	type args struct {
		confPath string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.

		{
			name: "Testcase 1",
			args: args{
				confPath: "../config.yaml",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			LoadConfAlert(tt.args.confPath)
		})
	}
}
