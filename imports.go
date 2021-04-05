// +build aopgen

package goaop

import (
	`github.com/Just-maple/xtoolinternal/gocommand`
	`github.com/Just-maple/xtoolinternal/imports`
	imports2 "golang.org/x/tools/imports"
)

var (
	opt2   = &imports2.Options{Comments: true, TabIndent: true, TabWidth: 8}
	intopt = &imports.Options{
		Env: &imports.ProcessEnv{
			GocmdRunner: &gocommand.Runner{},
		},
		LocalPrefix: imports2.LocalPrefix,
		AllErrors:   opt2.AllErrors,
		Comments:    opt2.Comments,
		FormatOnly:  opt2.FormatOnly,
		Fragment:    opt2.Fragment,
		TabIndent:   opt2.TabIndent,
		TabWidth:    opt2.TabWidth,
	}
)

func Process(filename string, src []byte, _ interface{}) (ret []byte, err error) {
	return imports.Process("", src, intopt)
}

