package types

import (
	"fmt"

	"cosmossdk.io/errors"
)

var (
	ErrInvalidCreator         = errors.Register(ModuleName, 1100, "invalid creator, creator address cannot be empty.")
	ErrInvalidID              = errors.Register(ModuleName, 1101, "invalid id, model id cannot be empty.")
	ErrInvalidName            = errors.Register(ModuleName, 1102, "invalid name, model name cannot be empty.")
	ErrInvalidURL             = errors.Register(ModuleName, 1103, "invalid url, model url cannot be empty.")
	ErrInvalidCreatorForEqual = errors.Register(ModuleName, 1104, "invalid creator, creator address must match the given address.")
	ErrInvalidDeleteModel     = errors.Register(ModuleName, 1105, "invalid creator, creator cannot delete the model.")
)

type ErrModelAlreadyExists struct {
	ErrorContent string
}

type ErrModelNotFound struct {
	ErrorContent string
}

// judge whether the error is impl error interface
var (
	_ error = (*ErrModelAlreadyExists)(nil)
	_ error = (*ErrModelNotFound)(nil)
)

func (err ErrModelAlreadyExists) Error() string {
	return fmt.Sprintf("model %s already exists", err.ErrorContent)
}

func (err ErrModelNotFound) Error() string {
	return fmt.Sprintf("model %s not found", err.ErrorContent)
}
