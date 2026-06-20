# fret

`fret` is a Bubble Tea TUI for learning the guitar fretboard as a connected
map — scales, chords, and arpeggios drawn across a neck that actually looks
like a neck: string lines, fret wires, a nut, and inlay dots.

```
      0    1     2     3     4     5     6     7     8     9    10    11    12
  e ─ 5 ─┃─────│─────│─ b7 ─│─────│─ R ──│─────│─────│─ b3 ─│─────│── 4 ─│─────│── 5 ─│
  B ─────┃─ b3 │─────│── 4 ─│─────│── 5 ─│─────│─────│─ b7 ─│─────│─ R ──│─────│─────│
  G ─ b7 ┃─────│─ R ──│─────│─────│─ b3 ─│─────│── 4 ─│─────│── 5 ─│─────│─────│─ b7 ─│
  D ─ 4 ─┃═════│══ 5 ═│═════│═════│═ b7 ═│═════│═ R ══│═════│═════│═ b3 ═│═════│══ 4 ═│
  A ─ R ─┃═════│═════│═ b3 ═│═════│══ 4 ═│═════│══ 5 ═│═════│═════│═ b7 ═│═════│═ R ══│
  E ─ 5 ─┃═════│═════│═ b7 ═│═════│═ R ══│═════│═════│═ b3 ═│═════│══ 4 ═│═════│══ 5 ═│
                       ∙           ∙           ∙           ∙                 ∙∙
```

Every interval gets its own color (root is peach, thirds are pink, the fifth
is blue, ...) so shapes become visible before the theory sinks in. Hover a
note and all its octave twins light up across the neck.

## Run

```sh
make run
```

or

```sh
make build
./bin/fret
```

## Controls

| Key | Action |
| --- | --- |
| `←↓↑→` / `hjkl` | move the cursor around the neck |
| `g` / `G` | jump to the nut / the last fret |
| `space` / `enter` | pin the hovered position |
| `x` | clear every pin |
| `r` | pick a key (root note) |
| `,` / `.` | transpose down / up a half step |
| `s` | pick a scale |
| `c` (or `a`) | pick a chord to see as an arpeggio |
| `t` | pick a tuning |
| `p` | pick a practice lens (a focused fret window) |
| `tab` / `n` | cycle labels: degrees → notes → both |
| `b` | cycle spelling: auto → sharps → flats |
| `[` / `]` | shrink / grow the neck (`0` resets to 12 frets) |
| `?` | help overlay |
| `q` / `ctrl+c` | quit |

## What's inside

- **12 scales**: major, natural/harmonic/melodic minor, both pentatonics,
  blues, and the modes (dorian, phrygian, lydian, mixolydian, locrian)
- **13 chord-tone arpeggios**: triads, sus chords, sixths, and sevenths
- **6 tunings**: Standard, Drop D, DADGAD, Open G, Open D, Eb Standard
- **Practice lenses**: plain-named fret windows (open position, low lane,
  middle lane, high lane, octave lane) for picking a small region to drill
- **Smart spelling**: flat keys spell flat, sharp keys spell sharp, and you
  can override either way
- The board auto-fits your terminal width and tells you how many frets are
  parked off-screen (`⋯+N`)

## Design notes

The tone is a playful music trainer, not a solemn theory encyclopedia.
Pattern names appear when they teach something transferable; otherwise the UI
uses plain practice labels like "middle lane" instead of inventing canon.

Layout: `internal/music` is the theory (pure, tested), `internal/ui` is the
Bubble Tea model and Lip Gloss rendering. The fretboard is the instrument
panel; the inspector and legend stay compact beneath it.
