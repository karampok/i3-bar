## My i3-bar using Barista

This is my bar optimized for Lenovo ThinkPad X1 Carbon (6th) running arch linux and for my needs.
It is based on [barista](https://github.com/soumya92/barista).

## Installation

With a working GO env:

    go get github.com/karampok/i3-bar

Fonts/icons from dir of `~/.{fonts,icons}`

# Usage

in ~/.config/i3/config:

```
bar {
  position top
  status_command exec ${HOME}/bin/i3-bar
  font xft:Sans 11
}
```
