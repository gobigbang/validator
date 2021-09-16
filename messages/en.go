package messages

var En = map[string]interface{}{
	`no_validation_input`: `Nothing to validate`,
	`key_not_allowed`:     `The field "{Field}" is not allowed here.`,
	`invalid_rule`:        `The rule {Rule} is invalid.`,
	`required`:            `The field "{Field}" is required.`,
	`not_zero`:            `The field "{Field}" is required.`,
	`is_array`:            `The field "{Field}" must be an array.`,
	`is_string`:           `The field "{Field}" must be a string.`,
	`is_numeric`:          `The field "{Field}" must be a number.`,
	`is_integer`:          `The field "{Field}" must be an integer.`,
	`is_map`:              `The field "{Field}" must be a key:value map.`,
	`is_struct`:           `The field "{Field}" must be an object.`,
	`eqfield`:             `The fields "{Field}" and {Arg0} must match.`,
	`in`:                  `The field "{Field}" is invalid`,
	`not_in`:              `The field "{Field}" is invalid`,
	`max_length`: map[string]interface{}{
		`string`: `The field "{Field}" may not be greater than {Arg0} characters long.`,
		`array`:  `The field "{Field}" may not have more than {Arg0} items.`,
		`map`:    `The field "{Field}" may not have more than {Arg0} items.`,
	},
	`min_length`: map[string]interface{}{
		`string`: `The field "{Field}" must be at least {Arg0} characters long.`,
		`array`:  `The field "{Field}" must have at least {Arg0} items.`,
		`map`:    `The field "{Field}" must have at least {Arg0} items.`,
	},
}
