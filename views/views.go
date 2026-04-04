package views

import "embed"

//go:embed *.html
var fs embed.FS

func MustRead(name string) []byte {
	data, err := fs.ReadFile(name)
	if err != nil {
		panic("views: " + err.Error())
	}
	return data
}
