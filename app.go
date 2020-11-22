package main

import (
	"flag"
	"fmt"
	"github.com/just1689/image-migrate/disk"
	"github.com/just1689/image-migrate/docker"
	"github.com/just1689/image-migrate/term"
	"github.com/just1689/image-migrate/util"
	"log"
	"os"
	"strconv"
	"strings"
)

var recursive = flag.Bool("r", false, "recursively check directories")
var update = flag.Bool("u", false, "update the YAML file after caching locally")
var registry = flag.String("registry", "", "registry URL (no trailing /)")

var skipPush = flag.Bool("skipPush", false, "skip pushing the image")
var skipPull = flag.Bool("skipPull", false, "skip pulling the image")

func main() {
	flag.Parse()
	p := os.Args[len(os.Args)-1]
	width := 19 + len(p)
	writer := term.Writer

	term.PrintWithColor(term.Repeat("-", width), term.MagentaText)
	term.PrintWithColor(fmt.Sprintf("%s%s%s", term.Repeat(" ", (width-24)/2), "starting image-migration", term.Repeat(" ", (width-24)/2)), term.YellowText)
	term.PrintWithColor(term.Repeat("-", width), term.MagentaText)

	// Header
	header := []string{"path", "recursive", "update"}
	term.PrintRow(writer, term.PaintRowUniformly(term.DefaultText, header))
	term.PrintRow(writer, term.PaintRowUniformly(term.DefaultText, term.AnonymizeRow(header))) // header separator

	colors := []term.Color{term.BrightYellowText, term.BrightGreenText, term.BrightGreenText, term.BrightGreenText}
	term.PrintRow(writer, term.PaintRow(colors, []string{p, strconv.FormatBool(*recursive), strconv.FormatBool(*update)}))

	writer.Flush()
	fmt.Println("")
	files := disk.ReadAllFiles(p, *recursive)
	colors = []term.Color{term.MagentaText}
	term.PrintRow(writer, term.PaintRow(colors, []string{"------------------------------------------"}))
	colors = []term.Color{term.YellowText}
	term.PrintRow(writer, term.PaintRow(colors, []string{" >            Getting files           "}))
	colors = []term.Color{term.MagentaText}
	term.PrintRow(writer, term.PaintRow(colors, []string{"------------------------------------------"}))
	for file := range files {
		if !strings.Contains(file, ".yaml") && !strings.Contains(file, ".yml") {
			colors = []term.Color{term.BrightYellowText}
			term.PrintRow(writer, term.PaintRow(colors, []string{"   ... skipping"}))
			continue
		}
		changeSet := make(map[string]string)
		colors = []term.Color{term.BrightGreenText}
		term.PrintRow(writer, term.PaintRow(colors, []string{fmt.Sprintf("   ::: %s", file)}))
		colors = []term.Color{term.Reset}
		lines := disk.ReadFile(file)
		for line := range lines {
			if strings.Contains(line, *registry) {
				colors = []term.Color{term.BrightYellowText}
				term.PrintRow(writer, term.PaintRow(colors, []string{"   ... skipping"}))
				continue
			}
			tabs := util.SplitStringChan(line)
			for tab := range tabs {
				if docker.IsDockerImage(tab) || (strings.Contains(line, "image:") && docker.IsDockerImageSquishy(tab)) {
					tab = strings.ReplaceAll(tab, "\"", "")
					tab = strings.ReplaceAll(tab, "'", "")
					tab = strings.ReplaceAll(tab, ",", "")
					colors = []term.Color{term.BrightYellowText}
					term.PrintRow(writer, term.PaintRow(colors, []string{fmt.Sprintf("   ... %s", tab)}))
					newTag := fmt.Sprintf("%s/%s", *registry, tab)
					if !*skipPull {
						err := docker.Pull(tab)
						if err != nil {
							log.Println("Failed", tab)
							continue
						}
						docker.Tag(tab, newTag)
					}
					if !*skipPush {
						err := docker.Push(newTag)
						if err != nil {
							fmt.Println(err)
							continue
						}
					}
					changeSet[tab] = newTag
				}
			}
		}
		if *update && len(changeSet) != 0 {
			newFile := fmt.Sprintf("%s.new", file)
			out := make(chan string)
			disk.NewWriter(newFile, out)
			in := disk.ReadFile(file)
			for next := range in {
				for prev, now := range changeSet {
					next = strings.ReplaceAll(next, prev, now)
				}
				out <- next
			}
			close(out)
			os.Remove(file)
			os.Rename(newFile, file)

		}
	}

}
