package confluence2md

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"strings"

	"github.com/beevik/etree"
)

func notImplemented(t etree.Token, w io.Writer) {
	fmt.Fprintf(w, "<!-- not implemented: %#v -->", t)
}

func walk(t etree.Token, w io.Writer, level int) {
	switch t := t.(type) {
	case *etree.CharData:
		fmt.Fprint(w, t.Data)

	case *etree.Document:
		walk(&t.Element, w, level)

	case *etree.Element:
		for _, c := range t.Child {
			switch c := c.(type) {
			case *etree.CharData:
				fmt.Fprint(w, c.Data)

			case *etree.Element:
				switch c.Space {
				case "":
					switch c.Tag {
					case "b", "strong":
						fmt.Fprint(w, "**")
						walk(c, w, level)
						fmt.Fprint(w, "**")

					case "i", "em":
						fmt.Fprint(w, "_")
						walk(c, w, level)
						fmt.Fprint(w, "_")

					case "code":
						fmt.Fprint(w, "`")
						walk(c, w, level)
						fmt.Fprint(w, "`")

					case "br":
						fmt.Fprint(w, "\n")

					case "p":
						walk(c, w, level)
						fmt.Fprint(w, "\n\n")

					case "h1", "h2", "h3", "h4", "h5", "h6":
						level := int(rune(c.Tag[1]) - rune('0'))
						fmt.Fprint(w, strings.Repeat("#", level)+" ")
						walk(c, w, level)
						fmt.Fprint(w, "\n\n")

					case "a":
						href := c.SelectAttrValue("href", "")
						if href == "" {
							walk(c, w, level)
						} else {
							fmt.Fprint(w, "[")
							walk(c, w, level)
							fmt.Fprint(w, "]("+href+")")
						}

					case "span":
						walk(c, w, level)

					case "ol", "ul":
						if level > 0 {
							fmt.Fprint(w, "\n")
							var b bytes.Buffer
							walk(c, &b, level)
							lines := strings.Split(b.String(), "\n")
							for _, line := range lines {
								fmt.Fprint(w, strings.Repeat("    ", level)+line+"\n")
							}
						} else {
							walk(c, w, level)
						}
						fmt.Fprint(w, "\n")

					case "li":
						switch c.Parent().Tag {
						case "ol":
							fmt.Fprint(w, "1. ")
						default:
							fmt.Fprint(w, "* ")
						}

						walk(c, w, level+1)
						fmt.Fprint(w, "\n")

					default:
						notImplemented(c, w)
						walk(c, w, level)
					}

				case "ac":
					switch c.Tag {
					case "toc":
						// Unsupported

					case "structured-macro":
						name := c.SelectAttrValue("ac:name", "")
						switch name {
						case "code":
							var lang string
							if el := c.FindElement("ac:parameter"); el != nil && el.SelectAttrValue("ac:name", "") == "language" {
								lang = el.Text()
							}

							fmt.Fprint(w, "```"+lang+"\n")
							walk(c.FindElement("ac:plain-text-body"), w, level)
							fmt.Fprint(w, "\n```\n")

						default:
							notImplemented(c, w)
							walk(c, w, level)
						}

					default:
						notImplemented(c, w)
						walk(c, w, level)
					}

				default:
					notImplemented(c, w)
					walk(c, w, level)
				}

			default:
				notImplemented(c, w)
				walk(c, w, level)
			}
		}
	}
}

func Convert(r io.Reader, w io.Writer) error {
	doc := etree.NewDocument()
	doc.ReadSettings.Entity = xml.HTMLEntity
	doc.ReadSettings.Permissive = true

	if _, err := doc.ReadFrom(r); err != nil {
		return err
	}

	walk(doc, w, 0)
	return nil
}
