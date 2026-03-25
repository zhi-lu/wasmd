package types

import (
	"fmt"

	"cosmossdk.io/errors"
)

var (
	ErrInvalidCreator         = errors.Register(ModuleName, 1100, "invalid creator address")
	ErrInvalidID              = errors.Register(ModuleName, 1101, "invalid id")
	ErrInvalidName            = errors.Register(ModuleName, 1102, "invalid name")
	ErrInvalidURL             = errors.Register(ModuleName, 1103, "invalid url")
	ErrInvalidCreatorForEqual = errors.Register(ModuleName, 1104, "creator address mismatch")
	ErrInvalidDeleteModel     = errors.Register(ModuleName, 1105, "unauthorized model deletion")
	ErrInvalidOperator        = errors.Register(ModuleName, 1106, "invalid operator address")
	ErrInvalidModelRef        = errors.Register(ModuleName, 1107, "referenced model does not exist")
	ErrAgentNotActive         = errors.Register(ModuleName, 1108, "agent is not active")
	ErrInsufficientFee        = errors.Register(ModuleName, 1109, "fee below agent minimum")
	ErrInvalidTaskStatus      = errors.Register(ModuleName, 1110, "invalid task status for this operation")
	ErrUnauthorized           = errors.Register(ModuleName, 1111, "unauthorized")
	ErrInvalidInputHash       = errors.Register(ModuleName, 1112, "invalid input hash")
	ErrInvalidFee             = errors.Register(ModuleName, 1113, "invalid fee")
)

type ErrModelAlreadyExists struct{ ID string }
type ErrModelNotFound struct{ ID string }
type ErrAgentAlreadyExistsError struct{ ID string }
type ErrAgentNotFoundError struct{ ID string }
type ErrTaskNotFoundError struct{ ID string }

var (
	_ error = (*ErrModelAlreadyExists)(nil)
	_ error = (*ErrModelNotFound)(nil)
	_ error = (*ErrAgentAlreadyExistsError)(nil)
	_ error = (*ErrAgentNotFoundError)(nil)
	_ error = (*ErrTaskNotFoundError)(nil)
)

func (e ErrModelAlreadyExists) Error() string { return fmt.Sprintf("model %s already exists", e.ID) }
func (e ErrModelNotFound) Error() string      { return fmt.Sprintf("model %s not found", e.ID) }
func (e ErrAgentAlreadyExistsError) Error() string {
	return fmt.Sprintf("agent %s already exists", e.ID)
}
func (e ErrAgentNotFoundError) Error() string { return fmt.Sprintf("agent %s not found", e.ID) }
func (e ErrTaskNotFoundError) Error() string  { return fmt.Sprintf("task %s not found", e.ID) }
