## My i3-bar using Barista

[barista](https://github.com/soumya92/barista) is i3status replacement in go.


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
