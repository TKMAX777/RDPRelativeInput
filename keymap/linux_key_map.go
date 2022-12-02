package keymap

type linuxKeyMap map[uint32]LinuxKey

type LinuxKey struct {
	Constant    string
	Value       uint32
	Description string
	WindowsKey  uint32
	EventInput  string
}

type linuxKeyMapFromEventInput map[string]LinuxKey

var linuxKeysFromEventInput = linuxKeyMapFromEventInput{}

func init() {
	for _, k := range linuxKeys {
		linuxKeysFromEventInput[k.EventInput] = k
	}
}

// Linux keycode to key details
var linuxKeys = linuxKeyMap{
	53:  {Constant: "Esc", Value: 53, Description: "Esc", EventInput: "Escape", WindowsKey: 0x1B},
	122: {Constant: "F1", Value: 122, Description: "F1", EventInput: "F1", WindowsKey: 0x70},
	120: {Constant: "F2", Value: 120, Description: "F2", EventInput: "F2", WindowsKey: 0x71},
	99:  {Constant: "F3", Value: 99, Description: "F3", EventInput: "F3", WindowsKey: 0x72},
	118: {Constant: "F4", Value: 118, Description: "F4", EventInput: "F4", WindowsKey: 0x73},
	96:  {Constant: "F5", Value: 96, Description: "F5", EventInput: "F5", WindowsKey: 0x74},
	97:  {Constant: "F6", Value: 97, Description: "F6", EventInput: "F6", WindowsKey: 0x75},
	98:  {Constant: "F7", Value: 98, Description: "F7", EventInput: "F7", WindowsKey: 0x76},
	100: {Constant: "F8", Value: 100, Description: "F8", EventInput: "F8", WindowsKey: 0x77},
	101: {Constant: "F9", Value: 101, Description: "F9", EventInput: "F9", WindowsKey: 0x78},
	109: {Constant: "F10", Value: 109, Description: "F10", EventInput: "F10", WindowsKey: 0x79},
	103: {Constant: "F11", Value: 103, Description: "F11", EventInput: "F11", WindowsKey: 0x7A},
	127: {Constant: "F12", Value: 127, Description: "F12", EventInput: "F12", WindowsKey: 0x7B},
	105: {Constant: "PrintScrn", Value: 105, Description: "PrintScrn", EventInput: "Print", WindowsKey: 0x2A},
	107: {Constant: "Scroll Lock", Value: 107, Description: "Scroll Lock", EventInput: "Scroll_Lock", WindowsKey: 0x91},
	113: {Constant: "Pause", Value: 113, Description: "Pause", EventInput: "Pause", WindowsKey: 0x13},
	50:  {Constant: "`", Value: 50, Description: "`"},
	18:  {Constant: "1", Value: 18, Description: "1", EventInput: "1", WindowsKey: 0x31},
	19:  {Constant: "2", Value: 19, Description: "2", EventInput: "2", WindowsKey: 0x32},
	20:  {Constant: "3", Value: 20, Description: "3", EventInput: "3", WindowsKey: 0x33},
	21:  {Constant: "4", Value: 21, Description: "4", EventInput: "4", WindowsKey: 0x34},
	23:  {Constant: "5", Value: 23, Description: "5", EventInput: "5", WindowsKey: 0x35},
	22:  {Constant: "6", Value: 22, Description: "6", EventInput: "6", WindowsKey: 0x36},
	26:  {Constant: "7", Value: 26, Description: "7", EventInput: "7", WindowsKey: 0x37},
	28:  {Constant: "8", Value: 28, Description: "8", EventInput: "8", WindowsKey: 0x38},
	25:  {Constant: "9", Value: 25, Description: "9", EventInput: "9", WindowsKey: 0x39},
	29:  {Constant: "0", Value: 29, Description: "0", EventInput: "0", WindowsKey: 0x30},
	27:  {Constant: "-", Value: 27, Description: "-", EventInput: "minus", WindowsKey: 0x6D},
	24:  {Constant: "=", Value: 24, Description: "="},
	51:  {Constant: "Backspace", Value: 51, Description: "Backspace", EventInput: "BackSpace", WindowsKey: 0x08},
	114: {Constant: "Insert", Value: 114, Description: "Insert", EventInput: "Insert", WindowsKey: 0x2D},
	115: {Constant: "Home", Value: 115, Description: "Home", EventInput: "Home", WindowsKey: 0x24},
	116: {Constant: "Page Up", Value: 116, Description: "Page Up", EventInput: "Page_Up", WindowsKey: 0x21},
	71:  {Constant: "Num Lock", Value: 71, Description: "Num Lock", EventInput: "Num_Lock", WindowsKey: 0x90},
	75:  {Constant: "KP /", Value: 75, Description: "KP /"},
	67:  {Constant: "KP *", Value: 67, Description: "KP *"},
	78:  {Constant: "KP -", Value: 78, Description: "KP -"},
	48:  {Constant: "Tab", Value: 48, Description: "Tab", EventInput: "Tab", WindowsKey: 0x09},
	12:  {Constant: "Q", Value: 12, Description: "Q", EventInput: "q", WindowsKey: 0x51},
	13:  {Constant: "W", Value: 13, Description: "W", EventInput: "w", WindowsKey: 0x57},
	14:  {Constant: "E", Value: 14, Description: "E", EventInput: "e", WindowsKey: 0x45},
	15:  {Constant: "R", Value: 15, Description: "R", EventInput: "r", WindowsKey: 0x52},
	17:  {Constant: "T", Value: 17, Description: "T", EventInput: "t", WindowsKey: 0x54},
	16:  {Constant: "Y", Value: 16, Description: "Y", EventInput: "y", WindowsKey: 0x59},
	32:  {Constant: "U", Value: 32, Description: "U", EventInput: "u", WindowsKey: 0x55},
	34:  {Constant: "I", Value: 34, Description: "I", EventInput: "i", WindowsKey: 0x49},
	31:  {Constant: "O", Value: 31, Description: "O", EventInput: "o", WindowsKey: 0x4F},
	35:  {Constant: "P", Value: 35, Description: "P", EventInput: "p", WindowsKey: 0x50},
	33:  {Constant: "[", Value: 33, Description: "[", EventInput: "["},
	30:  {Constant: "]", Value: 30, Description: "]", EventInput: "]"},
	36:  {Constant: "Return", Value: 36, Description: "Return", EventInput: "Return", WindowsKey: 0x0D},
	117: {Constant: "Delete", Value: 117, Description: "Delete", EventInput: "Delete", WindowsKey: 0x2E},
	119: {Constant: "End", Value: 119, Description: "End", EventInput: "End", WindowsKey: 0x23},
	121: {Constant: "Page Down", Value: 121, Description: "Page Down", EventInput: "Page_Down", WindowsKey: 0x22},
	89:  {Constant: "KP 7", Value: 89, Description: "KP 7", EventInput: ""},
	91:  {Constant: "KP 8", Value: 91, Description: "KP 8", EventInput: ""},
	92:  {Constant: "KP 9", Value: 92, Description: "KP 9", EventInput: ""},
	69:  {Constant: "KP +", Value: 69, Description: "KP +", EventInput: ""},
	57:  {Constant: "Caps Lock", Value: 57, Description: "Caps Lock", EventInput: "Caps_Lock", WindowsKey: 0x14},
	0:   {Constant: "A", Value: 0, Description: "A", EventInput: "a", WindowsKey: 0x41},
	1:   {Constant: "S", Value: 1, Description: "S", EventInput: "s", WindowsKey: 0x53},
	2:   {Constant: "D", Value: 2, Description: "D", EventInput: "d", WindowsKey: 0x44},
	3:   {Constant: "F", Value: 3, Description: "F", EventInput: "f", WindowsKey: 0x46},
	5:   {Constant: "G", Value: 5, Description: "G", EventInput: "g", WindowsKey: 0x47},
	4:   {Constant: "H", Value: 4, Description: "H", EventInput: "h", WindowsKey: 0x48},
	38:  {Constant: "J", Value: 38, Description: "J", EventInput: "j", WindowsKey: 0x4A},
	40:  {Constant: "K", Value: 40, Description: "K", EventInput: "k", WindowsKey: 0x4B},
	37:  {Constant: "L", Value: 37, Description: "L", EventInput: "l", WindowsKey: 0x4C},
	41:  {Constant: ";", Value: 41, Description: ";", EventInput: ";"},
	39:  {Constant: "'", Value: 39, Description: "'", EventInput: "'"},
	86:  {Constant: "KP 4", Value: 86, Description: "KP 4", EventInput: ""},
	87:  {Constant: "KP 5", Value: 87, Description: "KP 5", EventInput: ""},
	88:  {Constant: "KP 6", Value: 88, Description: "KP 6", EventInput: ""},
	56:  {Constant: "Shift Left", Value: 56, Description: "Shift Left", EventInput: "shift", WindowsKey: 0xA0},
	6:   {Constant: "Z", Value: 6, Description: "Z", EventInput: "z", WindowsKey: 0x5A},
	7:   {Constant: "X", Value: 7, Description: "X", EventInput: "x", WindowsKey: 0x58},
	8:   {Constant: "C", Value: 8, Description: "C", EventInput: "c", WindowsKey: 0x43},
	9:   {Constant: "V", Value: 9, Description: "V", EventInput: "v", WindowsKey: 0x56},
	11:  {Constant: "B", Value: 11, Description: "B", EventInput: "b", WindowsKey: 0x42},
	45:  {Constant: "N", Value: 45, Description: "N", EventInput: "n", WindowsKey: 0x4D},
	46:  {Constant: "M", Value: 46, Description: "M", EventInput: "m", WindowsKey: 0x53},
	43:  {Constant: ",", Value: 43, Description: ",", EventInput: "comma"},
	47:  {Constant: ".", Value: 47, Description: ".", EventInput: "period", WindowsKey: 0x6E},
	44:  {Constant: "/", Value: 44, Description: "/", EventInput: "slash", WindowsKey: 0x6F},
	42:  {Constant: "\\", Value: 42, Description: "\\", EventInput: "backslash"},
	62:  {Constant: "Cursor Up", Value: 62, Description: "Cursor Up", EventInput: "Up", WindowsKey: 0x26},
	83:  {Constant: "KP 1", Value: 83, Description: "KP 1", EventInput: ""},
	84:  {Constant: "KP 2", Value: 84, Description: "KP 2", EventInput: ""},
	85:  {Constant: "KP 3", Value: 85, Description: "KP 3", EventInput: ""},
	76:  {Constant: "KP Enter", Value: 76, Description: "KP Enter", EventInput: ""},
	54:  {Constant: "Ctrl Left", Value: 54, Description: "Ctrl Left", EventInput: "ctrl", WindowsKey: 0xA2},
	58:  {Constant: "Logo Left (-> Option)", Value: 58, Description: "Logo Left (-> Option)", EventInput: "", WindowsKey: 0x5B},
	55:  {Constant: "Alt Left (-> Command)", Value: 55, Description: "Alt Left (-> Command)", EventInput: "alt", WindowsKey: 0x12},
	49:  {Constant: "Space", Value: 49, Description: "Space", EventInput: "space", WindowsKey: 0x20},
	59:  {Constant: "Cursor Left", Value: 59, Description: "Cursor Left", EventInput: "Left", WindowsKey: 0x25},
	61:  {Constant: "Cursor Down", Value: 61, Description: "Cursor Down", EventInput: "Down", WindowsKey: 0x28},
	60:  {Constant: "Cursor Right", Value: 60, Description: "Cursor Right", EventInput: "Right", WindowsKey: 0x27},
	82:  {Constant: "KP 0", Value: 82, Description: "KP 0", EventInput: ""},
	65:  {Constant: "KP .", Value: 65, Description: "KP .", EventInput: ""},
}
