package flags

import (
	"bytes"
	"fmt"
	flag "github.com/spf13/pflag"
	"golang.org/x/term"
	"os"
	"strings"
)

var (
	usage Config
)

type CustomCommandUsage interface {
	// Usage returns the usage generated from the config provided to NewCustomCommandUsage and the provided FlagSet
	Usage(*flag.FlagSet) string
}

type commandUsageGenerator struct {
	config Config
}

// NewCustomCommandUsage returns a CustomCommandUsage that can be used to generate a custom usage
// string based on the provided config and flagset
func NewCustomCommandUsage(ccUsage Config) CustomCommandUsage {
	return &commandUsageGenerator{ccUsage}
}

func (cug *commandUsageGenerator) Usage(flagset *flag.FlagSet) string {
	buf := new(bytes.Buffer)

	maxlen := 0
	lines := []string{}
	cmd := os.Args[0]
	cols, _, err := term.GetSize(int(os.Stdin.Fd()))
	if err != nil { // just disable line wrapping if we can't do it for whatever reason
		cols = 0
		err = nil
	}

	fmt.Fprintf(buf, "Usage of %s\n", cmd)
	fmt.Fprintf(buf, "Author: %s\n", usage.Author)
	if len(usage.DescriptionElements) > 0 {
		fmt.Fprintf(buf, "Description:\n")

		for index, elem := range usage.DescriptionElements {
			switch elem.Kind {
			case Paragraph:
				fmt.Fprintf(buf, "%s\n\n", wrap(2, cols, elem.Content))
			case BulletPoint:
				var fmtString string
				if index+1 < len(usage.DescriptionElements) && usage.DescriptionElements[index+1].Kind == BulletPoint {
					fmtString = "  • %s\n"
				} else {
					fmtString = "  • %s\n\n"
				}

				fmt.Fprintf(buf, fmtString, wrap(4, cols, elem.Content))
			}
		}
	}

	// Assemble flag synopssis lines
	fmt.Fprintf(buf, "Flag Synopsis:\n")
	flagset.VisitAll(func(fl *flag.Flag) {
		if fl.Hidden {
			return
		}

		line := ""
		if fl.Shorthand != "" && fl.ShorthandDeprecated == "" {
			line = fmt.Sprintf("  -%s, --%s", fl.Shorthand, fl.Name)
		} else {
			line = fmt.Sprintf("      --%s", fl.Name)
		}

		varname, usage := flag.UnquoteUsage(fl)
		if varname != "" {
			line += " " + varname
		}

		// This special character will be replaced with spacing once the
		// correct alignment is calculated
		line += "\x00"
		if len(line) > maxlen {
			maxlen = len(line)
		}

		line += usage

		line += fmt.Sprintf(" (Default: '%v')", fl.DefValue)

		if len(fl.Deprecated) != 0 {
			line += fmt.Sprintf(" (DEPRECATED: %s)", fl.Deprecated)
		}

		line += fmt.Sprintln()
		lines = append(lines, line)
	})

	// wrap flag synopsis lines
	for _, line := range lines {
		sidx := strings.Index(line, "\x00")
		spacing := strings.Repeat(" ", maxlen-sidx)
		// maxlen + 2 comes from + 1 for the \x00 and + 1 for the (deliberate) off-by-one in maxlen-sidx
		fmt.Fprintln(buf, line[:sidx], spacing, wrap(maxlen+2, cols, line[sidx+1:]))
	}

	if len(cug.config.Examples) > 0 {
		fmt.Fprintf(buf, "\n\n")
		fmt.Fprintf(buf, "Examples:\n")

		for _, example := range cug.config.Examples {
			fmt.Fprintf(buf, "  %s\n", wrap(2, cols, example.Example))
			fmt.Fprintln(buf, "     ", wrap(6, cols, example.Description))
			fmt.Fprintln(buf)
		}
	} else {
		fmt.Fprintf(buf, "\n")
	}

	return buf.String()
}

func wrap(i, w int, s string) string {
	if w == 0 {
		return strings.Replace(s, "\n", "\n"+strings.Repeat(" ", i), -1)
	}

	// space between indent i and end of line width w into which
	// we should wrap the text.
	wrap := w - i

	var r, l string

	// Not enough space for sensible wrapping. Wrap as a block on
	// the next line instead.
	if wrap < 24 {
		i = 16
		wrap = w - i
		r += "\n" + strings.Repeat(" ", i)
	}
	// If still not enough space then don't even try to wrap.
	if wrap < 24 {
		return strings.Replace(s, "\n", r, -1)
	}

	// Try to avoid short orphan words on the final line, by
	// allowing wrapN to go a bit over if that would fit in the
	// remainder of the line.
	slop := 5
	wrap = wrap - slop

	// Handle first line, which is indented by the caller (or the
	// special case above)
	l, s = wrapN(wrap, slop, s)
	r = r + strings.Replace(l, "\n", "\n"+strings.Repeat(" ", i), -1)

	// Now wrap the rest
	for s != "" {
		var t string

		t, s = wrapN(wrap, slop, s)
		r = r + "\n" + strings.Repeat(" ", i) + strings.Replace(t, "\n", "\n"+strings.Repeat(" ", i), -1)
	}

	return r

}

func wrapN(i, slop int, s string) (string, string) {
	if i+slop > len(s) {
		return s, ""
	}

	w := strings.LastIndexAny(s[:i], " \t\n")
	if w <= 0 {
		return s, ""
	}
	nlPos := strings.LastIndex(s[:i], "\n")
	if nlPos > 0 && nlPos < w {
		return s[:nlPos], s[nlPos+1:]
	}
	return s[:w], s[w+1:]
}
