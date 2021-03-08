package commands

import (
	//"fmt"
	"github.com/jfrog/jfrog-cli-core/plugins/components"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestConfigCmdArguments(t *testing.T) {
	tests := []struct {
		name             string
		ctx              *components.Context
		wantErrMsgPrefix string
	}{
		{
			name: "zero argument  (without interactive menu)",
			ctx: &components.Context{
				Arguments: []string{},
			},
			wantErrMsgPrefix: "wrong number of arguments",
		},
		{
			name: "one argument",
			ctx: &components.Context{
				Arguments: []string{"a"},
			},
			wantErrMsgPrefix: "wrong number of arguments",
		},
		{
			name: "two argument",
			ctx: &components.Context{
				Arguments: []string{"a", "b"},
			},
			wantErrMsgPrefix: "product id",
		},
		{
			name: "three argument",
			ctx: &components.Context{
				Arguments: []string{"a", "b", "c"},
			},
			wantErrMsgPrefix: "wrong number of arguments",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := configCmd(tt.ctx)
			assert.NotNil(t, err)
			assert.True(t, strings.HasPrefix(err.Error(), tt.wantErrMsgPrefix))
		})
	}
}



//TODO: create a context mock to trigger the command manually
