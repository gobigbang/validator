package rules

// func isInClosure(ctx context.Context, value interface{}, params types.ValidateRequest) bool {
// 	in := params.Self.(types.RuleWithArgs).GetArgs()
// 	isIn := false
// 	if value != nil {

// 		switch k := params.ValueDescriptor.RKind; {
// 		case k == reflect.String:
// 			isIn = func() bool {
// 				for _, a := range in {
// 					if value == fmt.Sprint(a) {
// 						return true
// 					}
// 				}
// 				return false
// 			}()
// 		case (k <= reflect.Float64) && (k >= reflect.Int):
// 			isIn = func() bool {
// 				for _, a := range in {
// 					if reflect.TypeOf(a).Kind() == k {
// 						if value == a {
// 							return true
// 						}
// 					} else {
// 						if fmt.Sprint(value) == fmt.Sprint(a) {
// 							return true
// 						}
// 					}
// 				}
// 				return false
// 			}()
// 		default:
// 			isIn = true
// 		}

// 	}
// 	return isIn
// }

// var In = types.RuleWithArgsCtr(func(args ...interface{}) types.RuleWithArgs {
// 	return types.ClosureWithArgs(func(ctx context.Context, value interface{}, params types.ValidateRequest) (e error) {
// 		if value == nil {
// 			return
// 		}
// 		in := isInClosure(ctx, value, params)
// 		if !in {
// 			e = errors.New(params.Message)
// 		}
// 		return
// 	}).Message("The field value is invalid!").Code("in").Args(args...)
// })

// var NotIn = types.RuleWithArgsCtr(func(args ...interface{}) types.RuleWithArgs {
// 	return types.ClosureWithArgs(func(ctx context.Context, value interface{}, params types.ValidateRequest) (e error) {
// 		if value == nil {
// 			return
// 		}
// 		in := isInClosure(ctx, value, params)
// 		if in {
// 			e = errors.New(params.Message)
// 		}
// 		return
// 	}).Message("The field value is invalid!").Code("not_in").Args(args...)
// })
