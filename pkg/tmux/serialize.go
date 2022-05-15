package tmux

import (
	"reflect"

	"github.com/sergei-dyshel/claug/internal/utils"

	"github.com/iancoleman/strcase"
)

func Serialize(cmd any) (args []string) {
	value := reflect.ValueOf(cmd)
	utils.Assertf(
		value.Kind() == reflect.Pointer,
		"command must be pointer to struct (got %s)",
		value.Kind(),
	)
	value = value.Elem()
	utils.Assertf(
		value.Kind() == reflect.Struct,
		"command must be pointer to struct (got pointer to %s)",
		value.Kind(),
	)
	type_ := value.Type()
	args = []string{strcase.ToKebab(type_.Name())}
	seenPos := false
	for i := 0; i < type_.NumField(); i++ {
		field := type_.Field(i)
		val := value.Field(i)
		tag := field.Tag
		opt, hasOpt := tag.Lookup("opt")
		utils.Assertf(!seenPos || !hasOpt, "optional field comes after positional (%s)", field.Name)
		seenPos = seenPos || !hasOpt
		appendOpt := func(params ...string) { args = append(args, append([]string{"-" + opt}, params...)...) }
		switch field.Type.Kind() {
		case reflect.Bool:
			utils.Assertf(hasOpt, "bool field must be optional (%s)", field.Name)
			if val.Bool() {
				appendOpt()
			}
		case reflect.String:
			if hasOpt {
				if !val.IsZero() {
					appendOpt(val.String())
				}
			} else {
				args = append(args, val.String())
			}
		default:
			utils.Panicf("unsupported field type %s", field.Type.Kind())
		}
	}
	return
}
