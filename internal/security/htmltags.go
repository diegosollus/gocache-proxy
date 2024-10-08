package security

import (
	"regexp"
	"strings"
)

// Lista de tags HTML.
var tags = []string{
	"a", "abbr", "acronym", "address", "applet", "area", "audioscope",
	"b", "base", "basefront", "bdo", "bgsound", "big", "blackface",
	"blink", "blockquote", "body", "bq", "br", "button", "caption",
	"center", "cite", "code", "col", "colgroup", "comment", "dd",
	"del", "dfn", "dir", "div", "dl", "dt", "em", "embed", "fieldset",
	"fn", "font", "form", "frame", "frameset", "h1", "head", "hr",
	"html", "i", "iframe", "ilayer", "img", "input", "ins", "isindex",
	"kdb", "keygen", "label", "layer", "legend", "li", "limittext",
	"link", "listing", "map", "marquee", "menu", "meta", "multicol",
	"nobr", "noembed", "noframes", "noscript", "nosmartquotes",
	"object", "ol", "optgroup", "option", "p", "param", "plaintext",
	"pre", "q", "rt", "ruby", "s", "samp", "script", "select",
	"server", "shadow", "sidebar", "small", "spacer", "span",
	"strike", "strong", "style", "sub", "sup", "table", "tbody",
	"td", "textarea", "tfoot", "th", "thead", "title", "tr", "tt",
	"u", "ul", "var", "wbr", "xml", "xmp",
}

// TagRegex é a expressão regular que detecta tags HTML maliciosas.
var TagRegex = regexp.MustCompile(`<(?:` + strings.Join(tags, "|") + `)\W`)
