package readline

// Character codes
const (
	charCtrlA     = iota + 1 // "^A"
	charCtrlB                // "^B"
	charCtrlC                // "^C"
	charEOF                  // "^D"
	charCtrlE                // "^E"
	charCtrlF                // "^F"
	charCtrlG                // "^G"
	charBackspace            // "^?" // ISO 646
	charTab                  // "^C"
	// charCtrlH                      // "^H" // Added by me
	charCtrlJ                 // "^J"
	charCtrlK                 // "^K"
	charCtrlL                 // "^L"
	charCtrlM                 // "^M"
	charCtrlN                 // "^N"
	charCtrlO                 // "^O"
	charCtrlP                 // "^P"
	charCtrlQ                 // "^Q"
	charCtrlR                 // "^R"
	charCtrlS                 // "^S"
	charCtrlT                 // "^T"
	charCtrlU                 // "^U"
	charCtrlV                 // "^V"
	charCtrlW                 // "^W"
	charCtrlX                 // "^X"
	charCtrlY                 // "^Y"
	charCtrlZ                 // "^Z"
	charEscape                // ^[
	charCtrlSlash             // ^\
	charCtrlCloseSquare       // ^]
	charCtrlHat               // ^^
	charCtrlUnderscore        // ^_
	charBackspace2      = 127 // ASCII 1963
)

// This block maps all NON-nil-character special keys,
// including their combinations with key modifiers.
// Equivalents of []\^_ with Alt modifier
var ()

// This block maps all nil-character special keys, including
// their combinations with key modifiers.
var (
	// Sequences with modifiers
	seqInsert     = string([]byte{27, 91, 50, 126}) // ^[[2~
	seqDelete     = string([]byte{27, 91, 51, 126}) // ^[[3~
	seqHome       = string([]byte{27, 91, 72})      // ^[[H
	seqPageUp     = string([]byte{27, 91, 53, 126}) // ^[[5~
	seqPageDown   = string([]byte{27, 91, 54, 126}) // ^[[6~
	seqEnd        = string([]byte{27, 91, 70})      // ^[[F
	seqArrowUp    = string([]byte{27, 91, 65})      // ^[[A
	seqArrowDown  = string([]byte{27, 91, 66})      // ^[[B
	seqArrowRight = string([]byte{27, 91, 67})      // ^[[C
	seqArrowLeft  = string([]byte{27, 91, 68})      // ^[[D

	// Modifier (Ctrl)
	seqCtrlInsert     = string([]byte{27, 91, 50, 59, 53, 126}) // ^[[2;5~
	seqCtrlDelete     = string([]byte{27, 91, 51, 59, 53, 126}) // ^[[3;5~
	seqCtrlHome       = string([]byte{27, 91, 49, 59, 53, 72})  // ^[[1;5H
	seqCtrlPageUp     = string([]byte{27, 91, 53, 59, 53, 126}) // ^[[5;5~
	seqCtrlPageDown   = string([]byte{27, 91, 54, 59, 53, 126}) // ^[[6;5~
	seqCtrlEnd        = string([]byte{27, 91, 49, 59, 53, 70})  // ^[[1;5F
	seqCtrlArrowUp    = string([]byte{27, 91, 49, 59, 53, 65})  // ^[[1;5A
	seqCtrlArrowDown  = string([]byte{27, 91, 49, 59, 53, 66})  // ^[[1;5B
	seqCtrlArrowRight = string([]byte{27, 91, 49, 59, 53, 67})  // ^[[1;5C
	seqCtrlArrowLeft  = string([]byte{27, 91, 49, 59, 53, 68})  // ^[[1;5D

	// Modifier (Alt)
	seqAltInsert     = string([]byte{27, 91, 50, 59, 51, 126}) // ^[[2;3~
	seqAltDelete     = string([]byte{27, 91, 51, 59, 51, 126}) // ^[[3;3~
	seqAltHome       = string([]byte{27, 91, 49, 59, 51, 72})  // ^[[1;3H
	seqAltPageUp     = string([]byte{27, 91, 53, 59, 51, 126}) // ^[[5;3~
	seqAltPageDown   = string([]byte{27, 91, 54, 59, 51, 126}) // ^[[6;3~
	seqAltEnd        = string([]byte{27, 91, 49, 59, 51, 70})  // ^[[1;3F
	seqAltArrowUp    = string([]byte{27, 91, 49, 59, 51, 65})  // ^[[1;3A
	seqAltArrowDown  = string([]byte{27, 91, 49, 59, 51, 66})  // ^[[1;3B
	seqAltArrowRight = string([]byte{27, 91, 49, 59, 51, 67})  // ^[[1;3C
	seqAltArrowLeft  = string([]byte{27, 91, 49, 59, 51, 68})  // ^[[1;3D
)

// Other escape sequences
var (
	seqHomeSc   = string([]byte{27, 91, 49, 126})
	seqEndSc    = string([]byte{27, 91, 52, 126})
	seqShiftTab = string([]byte{27, 91, 90})
	seqAltQuote = string([]byte{27, 34})  // Added for showing registers ^["
	seqAltR     = string([]byte{27, 114}) // Used for alternative history
)

const (
	seqPosSave    = "\x1b[s"
	seqPosRestore = "\x1b[u"

	seqClearLineAfer    = "\x1b[0k"
	seqClearLineBefore  = "\x1b[1k"
	seqClearLine        = "\x1b[2k"
	seqClearScreenBelow = "\x1b[0J"
	seqClearScreen      = "\x1b[2J" // Clears screen fully
	seqCursorTopLeft    = "\x1b[H"  // Clears screen and places cursor on top-left

	seqGetCursorPos = "\x1b6n" // response: "\x1b{Line};{Column}R"

	seqCtrlLeftArrow  = "\x1b[1;5D"
	seqCtrlRightArrow = "\x1b[1;5C"
)

// Text effects
const (
	seqReset      = "\x1b[0m"
	seqBold       = "\x1b[1m"
	seqUnderscore = "\x1b[4m"
	seqBlink      = "\x1b[5m"
)

// Text colours
const (
	seqFgBlack   = "\x1b[30m"
	seqFgRed     = "\x1b[31m"
	seqFgGreen   = "\x1b[32m"
	seqFgYellow  = "\x1b[33m"
	seqFgBlue    = "\x1b[34m"
	seqFgMagenta = "\x1b[35m"
	seqFgCyan    = "\x1b[36m"
	seqFgWhite   = "\x1b[37m"

	seqFgBlackBright   = "\x1b[1;30m"
	seqFgRedBright     = "\x1b[1;31m"
	seqFgGreenBright   = "\x1b[1;32m"
	seqFgYellowBright  = "\x1b[1;33m"
	seqFgBlueBright    = "\x1b[1;34m"
	seqFgMagentaBright = "\x1b[1;35m"
	seqFgCyanBright    = "\x1b[1;36m"
	seqFgWhiteBright   = "\x1b[1;37m"
)

// Background colours
const (
	seqBgBlack   = "\x1b[40m"
	seqBgRed     = "\x1b[41m"
	seqBgGreen   = "\x1b[42m"
	seqBgYellow  = "\x1b[43m"
	seqBgBlue    = "\x1b[44m"
	seqBgMagenta = "\x1b[45m"
	seqBgCyan    = "\x1b[46m"
	seqBgWhite   = "\x1b[47m"

	seqBgBlackBright   = "\x1b[1;40m"
	seqBgRedBright     = "\x1b[1;41m"
	seqBgGreenBright   = "\x1b[1;42m"
	seqBgYellowBright  = "\x1b[1;43m"
	seqBgBlueBright    = "\x1b[1;44m"
	seqBgMagentaBright = "\x1b[1;45m"
	seqBgCyanBright    = "\x1b[1;46m"
	seqBgWhiteBright   = "\x1b[1;47m"
)

// Xterm 256 colors
const (
	seqCtermFg255 = "\033[48;5;255m"
)
