package cli

import (
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"

	"github.com/CosmWasm/wasmd/x/agent/types"
)

// GetTxCmd returns the transaction commands for this module.
func GetTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Agent transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
		SilenceUsage:               true,
	}
	for _, cmd := range []*cobra.Command{registerModelCmd(), deleteModelCmd()} {
		txCmd.AddCommand(cmd)
	}
	return txCmd
}

func registerModelCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "register-model [id] [name] [url]",
		Short: "Register a model",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			msg := &types.MsgRegisterModel{
				Creator: clientCtx.GetFromAddress().String(),
				Model: &types.Model{
					Id:      args[0],
					Name:    args[1],
					Url:     args[2],
					Creator: clientCtx.GetFromAddress().String(),
				},
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func deleteModelCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete-model [id]",
		Short: "Delete a model by id",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			msg := &types.MsgDeleteModel{
				Creator: clientCtx.GetFromAddress().String(),
				ModelId: args[0],
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}
