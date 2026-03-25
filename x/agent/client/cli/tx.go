package cli

import (
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/CosmWasm/wasmd/x/agent/types"
)

func GetTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Agent transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
		SilenceUsage:               true,
	}
	txCmd.AddCommand(
		registerModelCmd(),
		deleteModelCmd(),
		registerAgentCmd(),
		updateAgentCmd(),
		deactivateAgentCmd(),
		createTaskCmd(),
		submitTaskResultCmd(),
		cancelTaskCmd(),
	)
	return txCmd
}

// ── Model ──

func registerModelCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "register-model [id] [name] [url]",
		Short: "Register a new AI model",
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

// ── Agent ──

func registerAgentCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "register-agent [name] [model-id] [description] [fee-amount]",
		Short: "Register a new AI agent (e.g., register-agent my-agent gpt4 \"GPT-4 Agent\" 100stake)",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			fee, err := sdk.ParseCoinsNormalized(args[3])
			if err != nil {
				return err
			}
			msg := &types.MsgRegisterAgent{
				Operator:    clientCtx.GetFromAddress().String(),
				Name:        args[0],
				ModelId:     args[1],
				Description: args[2],
				FeePerTask:  fee,
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func updateAgentCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-agent [agent-id] [name] [description] [fee-amount]",
		Short: "Update an existing agent",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			fee, err := sdk.ParseCoinsNormalized(args[3])
			if err != nil {
				return err
			}
			msg := &types.MsgUpdateAgent{
				Operator:    clientCtx.GetFromAddress().String(),
				AgentId:     args[0],
				Name:        args[1],
				Description: args[2],
				FeePerTask:  fee,
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func deactivateAgentCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deactivate-agent [agent-id]",
		Short: "Deactivate an agent (set status to INACTIVE)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			msg := &types.MsgDeactivateAgent{
				Operator: clientCtx.GetFromAddress().String(),
				AgentId:  args[0],
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// ── Task ──

func createTaskCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-task [agent-id] [input-hash] [fee-amount]",
		Short: "Create a task for an agent (fee goes to escrow)",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			fee, err := sdk.ParseCoinsNormalized(args[2])
			if err != nil {
				return err
			}
			msg := &types.MsgCreateTask{
				Creator:   clientCtx.GetFromAddress().String(),
				AgentId:   args[0],
				InputHash: args[1],
				Fee:       fee,
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func submitTaskResultCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "submit-task-result [task-id] [output-hash] [result-url]",
		Short: "Submit result for a task (operator only, releases escrow)",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			msg := &types.MsgSubmitTaskResult{
				Operator:   clientCtx.GetFromAddress().String(),
				TaskId:     args[0],
				OutputHash: args[1],
				ResultUrl:  args[2],
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func cancelTaskCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cancel-task [task-id]",
		Short: "Cancel a pending task (refunds fee to creator)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			msg := &types.MsgCancelTask{
				Creator: clientCtx.GetFromAddress().String(),
				TaskId:  args[0],
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}
