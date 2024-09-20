package embed

import _ "embed"

//go:embed code.mp3
var CodeBytes []byte

const SeparatorExpectedString = "crimson falcon, silent thunder, midnight echo"
