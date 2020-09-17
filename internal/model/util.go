package model

import (
	"bytes"
	"fmt"
)

// PrintFiles print files info to console
func PrintFiles(files []DirItem) []byte {
	b := bytes.NewBuffer([]byte{})
	b.WriteString(fmt.Sprintf("total %d\n", len(files)))
	for _, f := range files {
		var t string = "f"
		var s = f.Size
		if f.IsDir {
			t = "d"
			s = 0
		} else if f.Error != "" {
			t = "e"
			s = 0
		}
		if f.Error != "" {
			b.WriteString(fmt.Sprintf("%s %15d %30s %s(%s)\n", t, s, "", f.Name, f.Error))
		} else {
			b.WriteString(fmt.Sprintf("%s %15d %30s %s\n", t, s, f.Modified, f.Name))
		}
	}
	return b.Bytes()
}
