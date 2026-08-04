package cmd

import (
	_operationbinding "common/operationbinding"
	_operationbindingcatalog "common/operationbindingcatalog"
	_operationcatalog "common/operationcatalog"
	_operationregistry "common/operationregistry"
	_operationruntime "common/operationruntime"

	"fmt"
)

/** newTQDOperationRuntime builds approved operation selection. */
func newTQDOperationRuntime() (*_operationruntime.Runtime, error) {
	operations, err := _operationregistry.New(_operationcatalog.TQDWorkspaceV1())
	if err != nil {
		return nil, fmt.Errorf("build operation registry: %w", err)
	}

	bindings, err := _operationbinding.New(
		_operationbindingcatalog.TQDWorkspaceV1(),
		operations,
	)
	if err != nil {
		return nil, fmt.Errorf("build operation bindings: %w", err)
	}

	runtime, err := _operationruntime.New(operations, bindings, nil)
	if err != nil {
		return nil, fmt.Errorf("build operation runtime: %w", err)
	}
	if err := runtime.RequireReady(
		_operationbindingcatalog.TQDWorkspaceRequiredTransportsV1(),
	); err != nil {
		return nil, fmt.Errorf("validate operation runtime: %w", err)
	}
	return runtime, nil
}
