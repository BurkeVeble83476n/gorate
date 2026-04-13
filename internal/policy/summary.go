package policy

import (
	"fmt"
	"io"
	"text/tabwriter"
)

// Summary holds a human-readable description of a single policy.
type Summary struct {
	Name     string
	Endpoint string
	Method   string
	Limit    int
	Window   int
}

// ToSummary converts a Policy into a Summary.
func ToSummary(p Policy) Summary {
	return Summary{
		Name:     p.Name,
		Endpoint: p.Endpoint,
		Method:   p.Method,
		Limit:    p.Limit,
		Window:   p.Window,
	}
}

// PrintTable writes a formatted table of policy summaries to w.
func PrintTable(w io.Writer, policies []Policy) {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "NAME\tENDPOINT\tMETHOD\tLIMIT\tWINDOW(s)")
	fmt.Fprintln(tw, "----\t--------\t------\t-----\t---------")
	for _, p := range policies {
		s := ToSummary(p)
		fmt.Fprintf(tw, "%s\t%s\t%s\t%d\t%d\n",
			s.Name, s.Endpoint, s.Method, s.Limit, s.Window)
	}
	tw.Flush()
}
