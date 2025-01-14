# 'Crossword game' overview

The idea here is an implementation of a simple word game I played a bunch with
my family as a kid. Never had a name for it other than as a 'crossword game',
although it is not a _crossword_, clearly. I should probably come up with
something better.

## Short overview of the game idea

The crossword game is played by 2-4 players on individual 5x5 grids, which start 
empty.

Each player takes turns in announcing a letter, which all players must then 
place in an empty square of their grid.

Once the grid is full, the players score based on words they have created either
vertically top-to-bottom or horizontally left-to-right in the grid. The scoring
rules are:

* Each letter may only be used once in each direction
  * i.e. the same letter may be used for a horizontal and a vertical word, but
    not in two overlapping words in the same direction
* Words of length 1-4 score their length
* Words of length 5 score double their length (i.e. 10 points)

Being that the grid is 5x5, the players will generally not get an equal number
of turns. For now at least this unfairness is accepted.

## Some vague design thoughts

* The intention here is to start with the engine of the game made accessible via
  an HTTP API – at first without even e.g. authentication for the players
* The game's data model is simple and so can be stored just about anywhere (to 
  start with, in memory)
* Once a basic game can be played, extensions might be:
  * Some kind of lobby support and identification/auth of players
  * An htmx-based web UI for the game
