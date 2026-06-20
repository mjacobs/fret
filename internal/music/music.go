// Package music models pitch classes, scales, chord-tone arpeggios, tunings,
// and the fretboard grid that the UI renders.
package music

import "strings"

// Note is a pitch class in [0, 12), with C = 0.
type Note int

const (
	C Note = iota
	Db
	D
	Eb
	E
	F
	Gb
	G
	Ab
	A
	Bb
	B
)

var (
	sharpNames = [...]string{"C", "C#", "D", "D#", "E", "F", "F#", "G", "G#", "A", "A#", "B"}
	flatNames  = [...]string{"C", "Db", "D", "Eb", "E", "F", "Gb", "G", "Ab", "A", "Bb", "B"}
)

func (n Note) Normalize() Note {
	return Note((int(n)%12 + 12) % 12)
}

// Add returns the pitch class some semitones away.
func (n Note) Add(semitones int) Note {
	return Note(int(n) + semitones).Normalize()
}

// Name spells the note with sharps or flats.
func (n Note) Name(flats bool) string {
	if flats {
		return flatNames[n.Normalize()]
	}
	return sharpNames[n.Normalize()]
}

// BothNames is the "C#/Db" spelling for pickers; natural notes return "C".
func (n Note) BothNames() string {
	s, f := n.Name(false), n.Name(true)
	if s == f {
		return s
	}
	return s + "/" + f
}

// PrefersFlats reports whether the key is conventionally spelled with flats.
func (n Note) PrefersFlats() bool {
	switch n.Normalize() {
	case F, Bb, Eb, Ab, Db:
		return true
	}
	return false
}

// Roots lists the twelve possible keys in chromatic order from C.
func Roots() []Note {
	return []Note{C, Db, D, Eb, E, F, Gb, G, Ab, A, Bb, B}
}

var intervalNames = [...]string{
	"root", "minor 2nd", "major 2nd", "minor 3rd", "major 3rd", "perfect 4th",
	"tritone", "perfect 5th", "minor 6th", "major 6th", "minor 7th", "major 7th",
}

// IntervalName names the interval of a semitone distance from the root.
func IntervalName(semitones int) string {
	return intervalNames[((semitones%12)+12)%12]
}

// Kind separates scales from chord-tone arpeggios.
type Kind int

const (
	KindScale Kind = iota
	KindArpeggio
)

func (k Kind) String() string {
	if k == KindArpeggio {
		return "arpeggio"
	}
	return "scale"
}

// NoteSet is a named collection of intervals: a scale or an arpeggio.
type NoteSet struct {
	Name    string
	Kind    Kind
	Formula []int    // semitone offsets from the root, ascending, starting at 0
	Degrees []string // labels aligned with Formula; "R" marks the root
	Vibe    string   // one-line personality blurb
}

var Scales = []NoteSet{
	{
		Name:    "Major",
		Kind:    KindScale,
		Formula: []int{0, 2, 4, 5, 7, 9, 11},
		Degrees: []string{"R", "2", "3", "4", "5", "6", "7"},
		Vibe:    "bright, settled, home-base friendly",
	},
	{
		Name:    "Natural minor",
		Kind:    KindScale,
		Formula: []int{0, 2, 3, 5, 7, 8, 10},
		Degrees: []string{"R", "2", "b3", "4", "5", "b6", "b7"},
		Vibe:    "darker, melodic, very song-shaped",
	},
	{
		Name:    "Harmonic minor",
		Kind:    KindScale,
		Formula: []int{0, 2, 3, 5, 7, 8, 11},
		Degrees: []string{"R", "2", "b3", "4", "5", "b6", "7"},
		Vibe:    "dramatic, candle-lit, faintly haunted",
	},
	{
		Name:    "Melodic minor",
		Kind:    KindScale,
		Formula: []int{0, 2, 3, 5, 7, 9, 11},
		Degrees: []string{"R", "2", "b3", "4", "5", "6", "7"},
		Vibe:    "minor with jazz-school posture",
	},
	{
		Name:    "Major pentatonic",
		Kind:    KindScale,
		Formula: []int{0, 2, 4, 7, 9},
		Degrees: []string{"R", "2", "3", "5", "6"},
		Vibe:    "open, singable, hard to make sound bad",
	},
	{
		Name:    "Minor pentatonic",
		Kind:    KindScale,
		Formula: []int{0, 3, 5, 7, 10},
		Degrees: []string{"R", "b3", "4", "5", "b7"},
		Vibe:    "the people's scale: immediate and sturdy",
	},
	{
		Name:    "Blues",
		Kind:    KindScale,
		Formula: []int{0, 3, 5, 6, 7, 10},
		Degrees: []string{"R", "b3", "4", "b5", "5", "b7"},
		Vibe:    "minor pentatonic plus the spicy passing tone",
	},
	{
		Name:    "Dorian",
		Kind:    KindScale,
		Formula: []int{0, 2, 3, 5, 7, 9, 10},
		Degrees: []string{"R", "2", "b3", "4", "5", "6", "b7"},
		Vibe:    "minor with a wink; funk and folk approved",
	},
	{
		Name:    "Phrygian",
		Kind:    KindScale,
		Formula: []int{0, 1, 3, 5, 7, 8, 10},
		Degrees: []string{"R", "b2", "b3", "4", "5", "b6", "b7"},
		Vibe:    "flamenco shadows and metal thunder",
	},
	{
		Name:    "Lydian",
		Kind:    KindScale,
		Formula: []int{0, 2, 4, 6, 7, 9, 11},
		Degrees: []string{"R", "2", "3", "#4", "5", "6", "7"},
		Vibe:    "major that floats two inches off the ground",
	},
	{
		Name:    "Mixolydian",
		Kind:    KindScale,
		Formula: []int{0, 2, 4, 5, 7, 9, 10},
		Degrees: []string{"R", "2", "3", "4", "5", "6", "b7"},
		Vibe:    "major in a leather jacket; lives over dom7",
	},
	{
		Name:    "Locrian",
		Kind:    KindScale,
		Formula: []int{0, 1, 3, 5, 6, 8, 10},
		Degrees: []string{"R", "b2", "b3", "4", "b5", "b6", "b7"},
		Vibe:    "the unstable one; rarely home, always tense",
	},
}

var Arpeggios = []NoteSet{
	{
		Name:    "Major triad",
		Kind:    KindArpeggio,
		Formula: []int{0, 4, 7},
		Degrees: []string{"R", "3", "5"},
		Vibe:    "the big friendly yes",
	},
	{
		Name:    "Minor triad",
		Kind:    KindArpeggio,
		Formula: []int{0, 3, 7},
		Degrees: []string{"R", "b3", "5"},
		Vibe:    "the beautiful sigh",
	},
	{
		Name:    "Diminished triad",
		Kind:    KindArpeggio,
		Formula: []int{0, 3, 6},
		Degrees: []string{"R", "b3", "b5"},
		Vibe:    "suspense in three notes",
	},
	{
		Name:    "Augmented triad",
		Kind:    KindArpeggio,
		Formula: []int{0, 4, 8},
		Degrees: []string{"R", "3", "#5"},
		Vibe:    "dream-sequence shimmer",
	},
	{
		Name:    "Sus2",
		Kind:    KindArpeggio,
		Formula: []int{0, 2, 7},
		Degrees: []string{"R", "2", "5"},
		Vibe:    "neither major nor minor; pleasantly undecided",
	},
	{
		Name:    "Sus4",
		Kind:    KindArpeggio,
		Formula: []int{0, 5, 7},
		Degrees: []string{"R", "4", "5"},
		Vibe:    "the held breath before the chord resolves",
	},
	{
		Name:    "Major 6",
		Kind:    KindArpeggio,
		Formula: []int{0, 4, 7, 9},
		Degrees: []string{"R", "3", "5", "6"},
		Vibe:    "vintage sweetness, western swing approved",
	},
	{
		Name:    "Minor 6",
		Kind:    KindArpeggio,
		Formula: []int{0, 3, 7, 9},
		Degrees: []string{"R", "b3", "5", "6"},
		Vibe:    "noir minor with a raised eyebrow",
	},
	{
		Name:    "Major 7",
		Kind:    KindArpeggio,
		Formula: []int{0, 4, 7, 11},
		Degrees: []string{"R", "3", "5", "7"},
		Vibe:    "sunset on a rooftop",
	},
	{
		Name:    "Dominant 7",
		Kind:    KindArpeggio,
		Formula: []int{0, 4, 7, 10},
		Degrees: []string{"R", "3", "5", "b7"},
		Vibe:    "the engine of the blues; wants to resolve",
	},
	{
		Name:    "Minor 7",
		Kind:    KindArpeggio,
		Formula: []int{0, 3, 7, 10},
		Degrees: []string{"R", "b3", "5", "b7"},
		Vibe:    "smooth, warm, endlessly loopable",
	},
	{
		Name:    "Minor 7 flat 5",
		Kind:    KindArpeggio,
		Formula: []int{0, 3, 6, 10},
		Degrees: []string{"R", "b3", "b5", "b7"},
		Vibe:    "jazz's favorite question mark",
	},
	{
		Name:    "Diminished 7",
		Kind:    KindArpeggio,
		Formula: []int{0, 3, 6, 9},
		Degrees: []string{"R", "b3", "b5", "bb7"},
		Vibe:    "perfectly symmetrical mischief",
	},
}

// DegreeOf locates a note within the set. It always returns the semitone
// distance from the root; ok reports whether the note belongs to the set.
func (s NoteSet) DegreeOf(root, note Note) (label string, semitones int, ok bool) {
	d := int(note.Normalize()-root.Normalize()+12) % 12
	for i, iv := range s.Formula {
		if iv == d {
			return s.Degrees[i], d, true
		}
	}
	return "", d, false
}

// Tuning is a set of open-string pitches, index 0 being the highest-pitched
// string (the top row of the fretboard display).
type Tuning struct {
	Name string
	Desc string
	Open []Note
}

var Tunings = []Tuning{
	{
		Name: "Standard",
		Desc: "the home you know",
		Open: []Note{E, B, G, D, A, E},
	},
	{
		Name: "Drop D",
		Desc: "low E drops a step; one-finger power chords",
		Open: []Note{E, B, G, D, A, D},
	},
	{
		Name: "DADGAD",
		Desc: "modal, droney, instant film score",
		Open: []Note{D, A, G, D, A, D},
	},
	{
		Name: "Open G",
		Desc: "strum it open, hear a G chord",
		Open: []Note{D, B, G, D, G, D},
	},
	{
		Name: "Open D",
		Desc: "slide-friendly, big and ringing",
		Open: []Note{D, A, Gb, D, A, D},
	},
	{
		Name: "Eb Standard",
		Desc: "everything a half step down, stadium style",
		Open: []Note{Eb, Bb, Gb, Db, Ab, Eb},
	},
}

// Label names a string for the left gutter; the highest string is lowercase
// by guitar convention.
func (t Tuning) Label(i int) string {
	name := t.Open[i].Name(t.Open[i].PrefersFlats())
	if i == 0 {
		return strings.ToLower(name)
	}
	return name
}

// Wound reports whether the string is drawn as a wound (thicker) string.
func (t Tuning) Wound(i int) bool {
	return i >= len(t.Open)-3
}

// LowToHigh spells the tuning the way guitarists say it: low string first.
func (t Tuning) LowToHigh(flats bool) string {
	parts := make([]string, len(t.Open))
	for i, n := range t.Open {
		parts[len(t.Open)-1-i] = n.Name(flats)
	}
	return strings.Join(parts, " ")
}

// Position is one fret on one string, annotated against the active set.
type Position struct {
	String   int
	Fret     int
	Note     Note
	Degree   string
	Interval int // semitones above the root
	InSet    bool
	IsRoot   bool
}

// Fretboard generates the annotated grid: one row per string, frets 0..frets.
func Fretboard(root Note, set NoteSet, tuning Tuning, frets int) [][]Position {
	board := make([][]Position, len(tuning.Open))
	for s, open := range tuning.Open {
		row := make([]Position, frets+1)
		for f := 0; f <= frets; f++ {
			note := open.Add(f)
			degree, interval, in := set.DegreeOf(root, note)
			row[f] = Position{
				String:   s,
				Fret:     f,
				Note:     note,
				Degree:   degree,
				Interval: interval,
				InSet:    in,
				IsRoot:   in && interval == 0,
			}
		}
		board[s] = row
	}
	return board
}

// Lens is a practice window: a named fret range to focus on.
type Lens struct {
	Name  string
	Min   int
	Max   int
	Blurb string
}

var Lenses = []Lens{
	{"Whole neck", 0, 24, "every note from nut to horizon"},
	{"Open position", 0, 4, "first frets plus open strings; cowboy-chord country"},
	{"Low lane", 2, 5, "a snug four-fret window near the nut"},
	{"Middle lane", 5, 9, "mid-neck territory where shapes start linking up"},
	{"High lane", 9, 12, "upper-neck real estate for melodic mischief"},
	{"Octave lane", 12, 15, "the whole map again, one octave up"},
}

func (l Lens) Includes(fret int) bool {
	return fret >= l.Min && fret <= l.Max
}
