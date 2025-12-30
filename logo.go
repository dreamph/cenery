package cenery

import "io"

const LogoANSI = "" +
	"\x1b[38;5;214m      __\x1b[38;5;202m====\x1b[38;5;214m__\x1b[0m\n" +
	"\x1b[38;5;214m  ___/  \x1b[38;5;202m====\x1b[38;5;214m  \\___\x1b[0m\n" +
	"\x1b[38;5;214m /  _   \x1b[38;5;202m====\x1b[38;5;214m   _  \\ \x1b[0m\n" +
	"\x1b[38;5;214m|  ( )  \x1b[38;5;202m====\x1b[38;5;214m  ( )  |\x1b[0m\n" +
	"\x1b[38;5;214m|   _   \x1b[38;5;202m====\x1b[38;5;214m   _   |\x1b[0m\n" +
	"\x1b[38;5;214m \\__\x1b[38;5;202mFIRE\x1b[38;5;214m__\x1b[38;5;202mSTART\x1b[38;5;214m__/ \x1b[0m\n" +
	"\x1b[38;5;45m     C E N E R Y\x1b[0m\n"

// PrintLogo writes the ANSI-colored cenery logo to w.
func PrintLogo(w io.Writer) error {
	if w == nil {
		return nil
	}
	_, err := io.WriteString(w, LogoANSI)
	return err
}
