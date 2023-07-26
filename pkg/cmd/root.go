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

func (r *RootCmd) GetTimeIntervalFlag(cmd *cobra.Command) int {
	port, _ := cmd.Flags().GetInt(constant.TimeIntervalFlag)
	return port
}

func (r *RootCmd) AddReceiverFlag() {
	r.Command.Flags().StringSliceP(constant.ReceiverFlag, "r", []string{}, "receiver id list")
}

func (r *RootCmd) GetAddReceiverFlag(cmd *cobra.Command) []string {
	receivers, _ := cmd.Flags().GetStringSlice(constant.ReceiverFlag)
	return receivers
}

func (r *RootCmd) AddSenderFlag() {
	r.Command.Flags().StringSliceP(constant.SenderFlag, "s", []string{}, "sender id list")
}

func (r *RootCmd) GetSenderFlag(cmd *cobra.Command) []string {
	senders, _ := cmd.Flags().GetStringSlice(constant.SenderFlag)
	return senders
}

func (r *RootCmd) AddMessageNumberFlag() {
	r.Command.Flags().IntP(constant.MessageNumberFlag, "m", 0, "single sender send message Number")
}

func (r *RootCmd) GetMessageNumberFlag(cmd *cobra.Command) int {
	senders, _ := cmd.Flags().GetInt(constant.MessageNumberFlag)
	return senders
}

func (r *RootCmd) Execute() error {
	return r.Command.Execute()
}
