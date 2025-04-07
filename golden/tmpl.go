package golden

import (
	"text/template"
)

type (
	// TmplFuncFactoryVars is a variables for TmplFuncFactory used to create functions collection.
	TmplFuncFactoryVars struct{}

	// TmplFuncFactory is a text/template functions collection factory.
	TmplFuncFactory func(TestingT, TmplFuncFactoryVars) template.FuncMap
)
