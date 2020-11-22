package main

import (
	"flag"
	"fmt"
	"github.com/just1689/image-migrate/disk"
	"github.com/just1689/image-migrate/docker"
	"github.com/just1689/image-migrate/util"
	"io"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"
)

var interactive = flag.Bool("i", false, "interactive")
var recursive = flag.Bool("r", false, "recursive")
var update = flag.Bool("u", false, "update")
var registry = flag.String("registry", "", "registry URL")

type Color string

// Color codes interpretted by the terminal
// NOTE: all codes must be of the same length or they will throw off the field alignment of tabwriter
const (
	Reset                   Color = "\x1b[0000m"
	Bright                        = "\x1b[0001m"
	BlackText                     = "\x1b[0030m"
	RedText                       = "\x1b[0031m"
	GreenText                     = "\x1b[0032m"
	YellowText                    = "\x1b[0033m"
	BlueText                      = "\x1b[0034m"
	MagentaText                   = "\x1b[0035m"
	CyanText                      = "\x1b[0036m"
	WhiteText                     = "\x1b[0037m"
	DefaultText                   = "\x1b[0039m"
	BrightRedText                 = "\x1b[1;31m"
	BrightGreenText               = "\x1b[1;32m"
	BrightYellowText              = "\x1b[1;33m"
	BrightBlueText                = "\x1b[1;34m"
	BrightMagentaText             = "\x1b[1;35m"
	BrightCyanText                = "\x1b[1;36m"
	BrightWhiteText               = "\x1b[1;37m"
	BlackBackground               = "\x1b[0040m"
	RedBackground                 = "\x1b[0041m"
	GreenBackground               = "\x1b[0042m"
	YellowBackground              = "\x1b[0043m"
	BlueBackground                = "\x1b[0044m"
	MagentaBackground             = "\x1b[0045m"
	CyanBackground                = "\x1b[0046m"
	WhiteBackground               = "\x1b[0047m"
	BrightBlackBackground         = "\x1b[0100m"
	BrightRedBackground           = "\x1b[0101m"
	BrightGreenBackground         = "\x1b[0102m"
	BrightYellowBackground        = "\x1b[0103m"
	BrightBlueBackground          = "\x1b[0104m"
	BrightMagentaBackground       = "\x1b[0105m"
	BrightCyanBackground          = "\x1b[0106m"
	BrightWhiteBackground         = "\x1b[0107m"
)

// Color implements the Stringer interface for interoperability with string
func (c *Color) String() string {
	return fmt.Sprintf("%v", c)
}

func Paint(color Color, value string) string {
	return fmt.Sprintf("%v%v%v", color, value, Reset)
}

func PaintRow(colors []Color, row []string) []string {
	paintedRow := make([]string, len(row))
	for i, v := range row {
		paintedRow[i] = Paint(colors[i], v)
	}
	return paintedRow
}

func PaintRowUniformly(color Color, row []string) []string {
	colors := make([]Color, len(row))
	for i, _ := range colors {
		colors[i] = color
	}
	return PaintRow(colors, row)
}

func AnonymizeRow(row []string) []string {
	anonRow := make([]string, len(row))
	for i, v := range row {
		anonRow[i] = strings.Repeat("-", len(v))
	}
	return anonRow
}

func PrintRow(writer io.Writer, line []string) {
	fmt.Fprintln(writer, strings.Join(line, "\t"))
}

func main() {
	flag.Parse()
	writer := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

	colors := []Color{MagentaText}
	PrintRow(writer, PaintRow(colors, []string{"------------------------------------"}))
	colors = []Color{YellowText}
	PrintRow(writer, PaintRow(colors, []string{"      starting image-migration      "}))
	colors = []Color{MagentaText}
	PrintRow(writer, PaintRow(colors, []string{"------------------------------------"}))

	// Header
	header := []string{"path", "recursive", "Update", "interactive"}
	PrintRow(writer, PaintRowUniformly(DefaultText, header))
	PrintRow(writer, PaintRowUniformly(DefaultText, AnonymizeRow(header))) // header separator

	p := os.Args[len(os.Args)-1]
	colors = []Color{BrightYellowText, BrightGreenText, BrightGreenText, BrightGreenText}
	PrintRow(writer, PaintRow(colors, []string{p, strconv.FormatBool(*recursive), strconv.FormatBool(*update), strconv.FormatBool(*interactive)}))

	writer.Flush()
	fmt.Println("")
	files := disk.ReadAllFiles(p)
	colors = []Color{MagentaText}
	PrintRow(writer, PaintRow(colors, []string{"------------------------------------------"}))
	colors = []Color{YellowText}
	PrintRow(writer, PaintRow(colors, []string{" >            Getting files           "}))
	colors = []Color{MagentaText}
	PrintRow(writer, PaintRow(colors, []string{"------------------------------------------"}))
	for file := range files {
		colors = []Color{BrightGreenText}
		PrintRow(writer, PaintRow(colors, []string{fmt.Sprintf("   ::: %s", file)}))
		colors = []Color{Reset}
		lines := disk.ReadFile(file)
		for line := range lines {
			tabs := util.SplitStringChan(line)
			for tab := range tabs {
				if docker.IsDockerImage(tab) || (strings.Contains(line, "image:") && docker.IsDockerImageSquishy(tab)) {
					tab = strings.ReplaceAll(tab, "\"", "")
					colors = []Color{BrightYellowText}
					PrintRow(writer, PaintRow(colors, []string{fmt.Sprintf("   ... %s", tab)}))
					docker.Pull(tab)
					newTag := fmt.Sprintf("%s/%s", *registry, tab)
					docker.Tag(tab, newTag)
					docker.Push(newTag)
				}
			}
		}

	}

}
