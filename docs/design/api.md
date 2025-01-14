# API

## Initial proposed API

### POST /api/v1/game

Create a game.

Returns: `201 Created` with a `Location` header pointing to the new game

### GET /api/v1/game/{game_id}

Get the current state of the game

Returns: `200 OK` with a JSON representation of the game state

```json
{
  "game_id": "abcd",
  "state":"awaiting_announcement/awaiting_placement/game_over",
  "playerCount": 2,
  "squaresFilled": 0,
  "currentAnnouncingPlayerId": 0
}
```

### GET /api/v1/game/{game_id}/players/{player_id}

Get the given player's current state and grid

Returns: `200 OK` with a JSON representation of the player's state

```json
{
  "player_id": 0,
  "grid": [
    ["", "", "", "", ""],
    ["", "", "", "", ""],
    ["", "", "", "", ""],
    ["", "", "", "", ""],
    ["", "", "", "", ""]
  ]
}
```

### POST /api/v1/game/{game_id}/players/{player_id}/announce

Announce a letter for the game

Request body: JSON object with a single key `letter` containing a single letter

Returns: `200 OK` if the player can announce a letter, `400 Bad Request`
if not

### POST /api/v1/game/{game_id}/players/{player_id}/place

Place a letter in the player's grid

Request body: JSON object with keys `row` and `column` containing the row and
column to place the letter in, and a key `letter` containing the letter to place

Returns: `200 OK` if the letter can be placed, `400 Bad Request` if not

### GET /api/v1/game/{game_id}/players/{player_id}/score

Get the given player's score.

Returns: `200 OK` with a JSON object containing the player's score if the game
is over, otherwise `400 Bad Request`.

```json
{
  "player_id": 0,
  "score": 0
}
```