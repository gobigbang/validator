package rules

import (
	"context"

	"github.com/gobigbang/validator/messages"
	"github.com/gobigbang/validator/types"
)

func Required(ctx context.Context, params types.RuleCheckerParameters) error {
	var err error
	passes := params.Value != nil
	if !passes {
		err = types.NewValidationError("required", messages.GetMessage(params.Messages, "required"), params)
	}
	return err
}
