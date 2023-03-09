package echo_swagger

import (
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/ruanlianjun/gutils"
)

func generate(options *swagger) {
	parser := options.parser

	if err := parser.ParseAPI(options.rootPath, options.mainFileName, 10); err != nil {
		color.New(color.FgRed).Println(err)
		return
	}

	bytes, err := parser.GetSwagger().MarshalJSON()
	if err != nil {
		color.New(color.FgRed).Println(err)
		return
	}

	dir := filepath.Join(options.rootPath, options.filename)
	if err = gutils.MkdirAll(dir); err != nil {
		color.New(color.FgRed).Println(err)
		return
	}

	if err = os.WriteFile(dir, bytes, 0655); err != nil {
		color.New(color.FgRed).Println(err)
		return
	}
}
