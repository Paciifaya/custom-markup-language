package main

import (
	"fmt"
	"html"
	"os"
	"strings"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: mark <file.mark>")
		os.Exit(1)
	}

	src, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	out, err := Render(string(src))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Print(out)
}

func Render(src string) (string, error) {
	var out strings.Builder
	lines := strings.Split(strings.ReplaceAll(src, "\r\n", "\n"), "\n")

	for i := 0; i < len(lines); i++ {
		line := lines[i]

		switch {
		case line == "```":
			i = renderCode(&out, lines, i+1)

		case strings.HasPrefix(line, "::"):
			name := strings.TrimSpace(strings.TrimPrefix(line, "::"))
			i = renderDirective(&out, lines, i+1, name)

		case strings.HasPrefix(line, "#"):
			level := countPrefix(line, "#")
			text := strings.TrimSpace(line[level:])
			fmt.Fprintf(&out, "<h%d>%s</h%d>\n", level, escape(text), level)

		case strings.HasPrefix(line, "- "):
			items := []string{inline(line[2:])}
			i++

			for i < len(lines) && strings.HasPrefix(lines[i], "- ") {
				items = append(items, inline(lines[i][2:]))
				i++
			}
			i--

			out.WriteString("<ul>\n")
			for _, item := range items {
				fmt.Fprintf(&out, "  <li>%s</li>\n", item)
			}
			out.WriteString("</ul>\n")

		case strings.HasPrefix(line, "> "):
			quote := []string{inline(line[2:])}
			i++

			for i < len(lines) && strings.HasPrefix(lines[i], "> ") {
				quote = append(quote, inline(lines[i][2:]))
				i++
			}
			i--

			fmt.Fprintf(&out, "<blockquote>%s</blockquote>\n", strings.Join(quote, " "))

		case strings.TrimSpace(line) == "":
			// Ignore blank lines.

		default:
			fmt.Fprintf(&out, "<p>%s</p>\n", inline(line))
		}
	}

	return out.String(), nil
}

func renderCode(out *strings.Builder, lines []string, start int) int {
	var code []string

	for i := start; i < len(lines); i++ {
		if lines[i] == "```" {
			fmt.Fprintf(out, "<pre><code>%s</code></pre>\n", escape(strings.Join(code, "\n")))
			return i
		}
		code = append(code, lines[i])
	}

	return start - 1
}

func renderDirective(out *strings.Builder, lines []string, start int, name string) int {
	var body []string

	for i := start; i < len(lines); i++ {
		if strings.TrimSpace(lines[i]) == "::" {
			inner, _ := Render(strings.Join(body, "\n"))
			fmt.Fprintf(out, "<div class=\"callout %s\">\n%s</div>\n", escape(name), inner)
			return i
		}
		body = append(body, lines[i])
	}

	return start - 1
}

func inline(s string) string {
	var out strings.Builder

	for i := 0; i < len(s); i++ {
		switch {
		case strings.HasPrefix(s[i:], "*"):
			end := strings.Index(s[i+1:], "*")
			if end == -1 {
				out.WriteString("*")
				continue
			}

			text := s[i+1 : i+1+end]
			fmt.Fprintf(&out, "<em>%s</em>", escape(text))
			i += end + 1

		case strings.HasPrefix(s[i:], "`"):
			end := strings.Index(s[i+1:], "`")
			if end == -1 {
				out.WriteString("`")
				continue
			}

			text := s[i+1 : i+1+end]
			fmt.Fprintf(&out, "<code>%s</code>", escape(text))
			i += end + 1

		case strings.HasPrefix(s[i:], "["):
			end := strings.Index(s[i+1:], "](")
			if end == -1 {
				out.WriteString("[")
				continue
			}

			text := s[i+1 : i+1+end]
			rest := s[i+end+2:]
			urlEnd := strings.Index(rest, ")")
			if urlEnd == -1 {
				out.WriteString("[")
				continue
			}

			url := rest[:urlEnd]
			fmt.Fprintf(&out, "<a href=%q>%s</a>", escape(url), escape(text))
			i += end + urlEnd + 2

		default:
			out.WriteString(string(s[i]))
		}
	}

	return out.String()
}

func countPrefix(s, prefix string) int {
	count := 0
	for strings.HasPrefix(s[count:], prefix) {
		count++
	}
	return count
}

func escape(s string) string {
	return html.EscapeString(s)
}