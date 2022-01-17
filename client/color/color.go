package color

import (
	"fmt"
	"github.com/fatih/color"
)

const Clearln = "\r\x1b[2K"

func PrintlnYellow(format string) {
	color.Set(color.FgYellow)
	fmt.Println(Clearln + fmt.Sprintf(format))
	color.Unset()
}
