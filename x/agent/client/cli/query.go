package cli

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"

	"github.com/CosmWasm/wasmd/x/agent/types"
)

func GetQueryCmd() *cobra.Command {
	queryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Agent query subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
		SilenceUsage:               true,
	}
	queryCmd.AddCommand(
		modelQueryCmd(),
		modelsQueryCmd(),
		agentQueryCmd(),
		agentsQueryCmd(),
		agentsByModelQueryCmd(),
		taskQueryCmd(),
		tasksQueryCmd(),
		tasksByAgentQueryCmd(),
		tasksByCreatorQueryCmd(),
	)
	return queryCmd
}

// ── Model ──

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

// ── Agent ──

func agentQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "agent [id]",
		Short: "Query an agent by id",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			qc := types.NewQueryClient(clientCtx)
			resp, err := qc.Agent(context.Background(), &types.QueryAgentRequest{Id: args[0]})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(resp)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func agentsQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "agents",
		Short: "Query all agents",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			qc := types.NewQueryClient(clientCtx)
			resp, err := qc.Agents(context.Background(), &types.QueryAgentsRequest{})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(resp)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func agentsByModelQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "agents-by-model [model-id]",
		Short: "Query agents by model id",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			qc := types.NewQueryClient(clientCtx)
			resp, err := qc.AgentsByModel(context.Background(), &types.QueryAgentsByModelRequest{ModelId: args[0]})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(resp)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// ── Task ──

func taskQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "task [id]",
		Short: "Query a task by id",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			qc := types.NewQueryClient(clientCtx)
			resp, err := qc.Task(context.Background(), &types.QueryTaskRequest{Id: args[0]})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(resp)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func tasksQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tasks",
		Short: "Query all tasks",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			qc := types.NewQueryClient(clientCtx)
			resp, err := qc.Tasks(context.Background(), &types.QueryTasksRequest{})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(resp)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func tasksByAgentQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tasks-by-agent [agent-id]",
		Short: "Query tasks by agent id",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			qc := types.NewQueryClient(clientCtx)
			resp, err := qc.TasksByAgent(context.Background(), &types.QueryTasksByAgentRequest{AgentId: args[0]})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(resp)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func tasksByCreatorQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tasks-by-creator [creator-address]",
		Short: "Query tasks by creator address",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			qc := types.NewQueryClient(clientCtx)
			resp, err := qc.TasksByCreator(context.Background(), &types.QueryTasksByCreatorRequest{Creator: args[0]})
			if err != nil {
				return err
			}
			return clientCtx.PrintProto(resp)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}
