package common

import (
	"bytes"
	"text/tabwriter"
)

type TabWriterWrapper struct {
	writer *tabwriter.Writer
}

var output bytes.Buffer
var TabWriter *TabWriterWrapper

func InitTabWriter(padding ...int) {
	// Init the global tab writer, for prettiness :)
	pad := 0
	if len(padding) > 0 {
		pad = padding[0]
	}
	TabWriter = &TabWriterWrapper{
		writer: tabwriter.NewWriter(&output, 0, 8, pad, '\t', tabwriter.AlignRight),
	}
}

func (w *TabWriterWrapper) Write(raw *[]string) string {
	if raw == nil {
		return ""
	}
	output.Reset()
	for _, s := range *raw {
		w.writer.Write([]byte(s + "\n"))
	}
	w.writer.Flush()
	return output.String()
}
