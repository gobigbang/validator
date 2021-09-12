package types

import "context"

func Validate(ctx context.Context, input interface{}, rules interface{}, config ValidatorConfig) error {
	var err error

	switch rules.(type) {
	case Validatable:
		err = rules.(Validatable).Validate(ctx, input, config)
	}

	return err
}
