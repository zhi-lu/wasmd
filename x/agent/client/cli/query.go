package cli

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"

	"github.com/CosmWasm/wasmd/x/agent/types"
)

// GetQueryCmd returns the query commands for this module.
func GetQueryCmd() *cobra.Command {
	queryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Agent query subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
		SilenceUsage:               true,
	}
	for _, cmd := range []*cobra.Command{modelQueryCmd(), modelsQueryCmd()} {
		queryCmd.AddCommand(cmd)
	}
	return queryCmd
}

func modelQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "model [id]",
		Short: "Query a model by id",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			qc := types.NewQueryClient(clientCtx)
			resp, err := qc.Model(context.Background(), &types.QueryModelRequest{Id: args[0]})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(resp)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func modelsQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "models",
		Short: "Query all models",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			qc := types.NewQueryClient(clientCtx)
			resp, err := qc.Models(context.Background(), &types.QueryModelsRequest{})
			if err != nil {
				return err
			}
			if resp == nil {
				return fmt.Errorf("empty response")
			}
			return clientCtx.PrintProto(resp)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}
