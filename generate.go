package echo_swagger

import (
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/ruanlianjun/gutils"
)

func generate(options *swagger) error {
	parser := options.parser

	if err := parser.ParseAPI(options.rootPath, options.mainFileName, 10); err != nil {
		_, err := color.New(color.FgRed).Println(err)
		return err
	}

	bytes, err := parser.GetSwagger().MarshalJSON()
	if err != nil {
		_, err := color.New(color.FgRed).Println(err)
		return err
	}

	dir := filepath.Join(options.rootPath, options.filename)
	if err = gutils.MkdirAll(dir); err != nil {
		_, err := color.New(color.FgRed).Println(err)
		return err
	}

	if err = os.WriteFile(dir, bytes, 0655); err != nil {
		_, err := color.New(color.FgRed).Println(err)
		return err
	}
	return nil
}
