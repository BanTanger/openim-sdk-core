package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"open_im_sdk/pkg/constant"
)

type RootCmd struct {
	Command cobra.Command
	name    string
}

func NewRootCmd(name string) RootCmd {
	c := cobra.Command{
		Short: fmt.Sprintf(`start %s`, name),
		Long:  fmt.Sprintf(`start %s`, name),
	}
	return RootCmd{
		Command: c,
		name:    name,
	}
}

func (r *RootCmd) AddTimeIntervalFlag() {
	r.Command.Flags().IntP(constant.TimeIntervalFlag, "t", 0, "interval time mill second")
}

func (r *RootCmd) getTimeIntervalFlag(cmd *cobra.Command) int {
	port, _ := cmd.Flags().GetInt(constant.TimeIntervalFlag)
	return port
}
