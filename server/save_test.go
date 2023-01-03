package server

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_write(t *testing.T) {
	type args struct {
		server    Server
		overwrite bool
		serverDir string
	}
	tests := []struct {
		name    string
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.wantErr(t, write(tt.args.server, tt.args.overwrite, tt.args.serverDir), fmt.Sprintf("write(%v, %v, %v)", tt.args.server, tt.args.overwrite, tt.args.serverDir))
		})
	}
}
