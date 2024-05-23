package loan

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"

	modulev1 "loan/api/loan/loan"
)

// AutoCLIOptions implements the autocli.HasAutoCLIConfig interface.
func (am AppModule) AutoCLIOptions() *autocliv1.ModuleOptions {
	return &autocliv1.ModuleOptions{
		Query: &autocliv1.ServiceCommandDescriptor{
			Service: modulev1.Query_ServiceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "Params",
					Use:       "params",
					Short:     "Shows the parameters of the module",
				},
				{
					RpcMethod:      "GetAllLoan",
					Use:            "get-all-loan [type --- 0 all,1 pending,2 approved,3 canceled,4 liquidate]",
					Short:          "Query get-all-loan",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "type"}},
				},

				{
					RpcMethod:      "GetLoan",
					Use:            "get-loan [loan-id]",
					Short:          "Query get-loan",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "loanId"}},
				},

				// this line is used by ignite scaffolding # autocli/query
			},
		},
		Tx: &autocliv1.ServiceCommandDescriptor{
			Service:              modulev1.Msg_ServiceDesc.ServiceName,
			EnhanceCustomCommand: true, // only required if you want to use the custom command
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "UpdateParams",
					Skip:      true, // skipped because authority gated
				},
				{
					RpcMethod:      "CreateLoan",
					Use:            "create-loan [amount] [apy] [collateral] [borrow-time]",
					Short:          "Send a create-loan tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "amount"}, {ProtoField: "apy"}, {ProtoField: "collateral"}, {ProtoField: "borrowTime"}},
				},
				{
					RpcMethod:      "ApprovedLoan",
					Use:            "approved-loan [loan-id]",
					Short:          "Send a approved-loan tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "loanId"}},
				},
				// this line is used by ignite scaffolding # autocli/tx
			},
		},
	}
}
