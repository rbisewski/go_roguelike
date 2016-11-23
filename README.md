# go-roguelike - straightforward roguelike written in golang

go-roguelike is a plain ol' roguelike game implemented in golang. It
resembles other roguelikes, such as nethack.

The code itself is based off of code found from other online sources and
uses the gocurses library in order to interface with the ever-popular
ncurses library. 

Perhaps one day it will be completed, but for now it is still missing
key features like inventory or proper monster generation or multiple
levels.


# Requirements

Specifically, the following packages are required:

* golang
* ncurses

Nothing else should be needed, but in my experience some exotic (i.e.
non-x86) platforms tend to have gimpy or partly functional ncurses
libraries, so it ought to go without saying that certain bugs may be
present. 

I tend to write software on the x86-64 and arm7 architectures, so feel
free to shoot me an email if you discover bugs on other platforms.


# Installation

Enter the following command to build the executable (if necessary as root):

    go build

Afterwards run the binary from the commandline, as you would any typical
golang program.

# Running go-roguelike

Simply run the compiled file from the commandline and it should work as
intended.


# Authors

This software utilizes a golang wrapper for the ncurses library. For more
information, consider contacting the original author:

* GitHub Repo -> https://github.com/tncardoso/gocurses  

Some code was gleamed from the golang code of this developer, so naturally
a tip-of-the-hat is in order:

* GitHub Repo -> https://github.com/GGalizzi 

The remainder was created by Robert Bisewski at Ibis Cybernetics. For
more information, contact:

* Website -> www.ibiscybernetics.com

* Email -> contact@ibiscybernetics.com
