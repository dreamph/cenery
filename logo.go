package cenery

import (
	"fmt"
	"io"
	"strings"
)

// PrintLogo writes the cenery logo and runtime info to w.
func PrintLogo(w io.Writer, engineName string, addr string) error {
	if w == nil {
		return nil
	}
	if engineName == "" {
		engineName = "Unknown"
	}

	if strings.HasPrefix(addr, ":") {
		addr = "0.0.0.0" + addr
	}
	const (
		colorOrange = "\x1b[38;5;214m"
		colorRed    = "\x1b[38;5;202m"
		colorCyan   = "\x1b[38;5;45m"
		colorReset  = "\x1b[0m"
	)
	_, err := fmt.Fprintf(w, ""+
		colorOrange+"      __"+colorRed+"===="+colorOrange+"__"+colorReset+"\n"+
		colorOrange+"  ___/  "+colorRed+"===="+colorOrange+"  \\___"+colorReset+"\n"+
		colorOrange+" /  _   "+colorRed+"===="+colorOrange+"   _  \\ "+colorReset+"\n"+
		colorOrange+"|  ( )  "+colorRed+"===="+colorOrange+"  ( )  |"+colorReset+"\n"+
		colorOrange+"|   _   "+colorRed+"===="+colorOrange+"   _   |"+colorReset+"\n"+
		colorOrange+" \\__"+colorRed+"FIRE"+colorOrange+"__"+colorRed+"START"+colorOrange+"__/ "+colorReset+"\n"+
		colorCyan+"     C E N E R Y"+colorReset+"\n"+
		colorOrange+"Engine"+colorReset+" : %s\n"+
		colorOrange+"Port"+colorReset+" : %s\n", engineName, addr)
	return err
}
