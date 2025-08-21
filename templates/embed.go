package templates

import "embed"

// Embed the default template directory recursively.
// Paths are relative to this directory; naming a directory embeds its full tree.
//go:embed default default/.*
var FS embed.FS
