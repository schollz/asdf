package jsparse

import (
	"strings"

	"github.com/dop251/goja"
)

func Parse(code string) (variables []string, values []string, err error) {
	vm := goja.New()
	_, err = vm.RunString(code)
	if err != nil {
		return nil, nil, err
	}

	varNames := vm.GlobalObject().Keys()
	for _, name := range varNames {
		value := vm.Get(name)
		strValue, ok := value.Export().(string)
		if ok {
			variables = append(variables, name)
			values = append(values, strings.TrimSpace(strValue))
		}
	}

	return variables, values, nil
}
