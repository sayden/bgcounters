# Cli

```shell
Usage: counters <command>

Flags:
  -h, --help    Show context-sensitive help.

Commands:
  image --output-content=STRING --input-path=STRING

  json --input-path=STRING --output-content=STRING

  vassal --csv=STRING --vassal-output-folder=STRING
    Create a vassal module for testing. It searches for the 'template.xml' in the same folder

Run "counters <command> --help" for more information on a command.
```
## Commands
### `image`
```shell
Usage: counters image --output-content=STRING --input-path=STRING

Flags:
  -h, --help                     Show context-sensitive help.

      --output-content=STRING    InputContent to produce: counters, back-counters, cards, fow-counters or events
  -i, --input-path=STRING        Input path of the file to read. Be aware that some outputs requires specific inputs.
  -o, --output-path=STRING       Path to the folder to write the image(s)
      --tiled                    Write a sheet of 7x10 items per parge
      --individual               Write a file for each counter/card
```

### `json`
```shell
Usage: counters json --input-path=STRING --output-content=STRING

Flags:
  -h, --help                             Show context-sensitive help.

  -i, --input-path=STRING                Input path of the file to read. Be aware that some outputs requires specific inputs.
  -o, --output-path=STRING               Path to the folder to write the content
      --output-content=STRING            InputContent to produce: counters, back-counters, cards, fow-counters or events
      --events-pool-file=STRING          A file to take 'events' from
      --back-image=STRING                The image for the back of the cards
      --output-destination=STRING        When generating a JSON Template, this contains the destination folder for images inside the template
      --card-template-filepath=STRING    When writing cards, a template for those cards must be provided
```

### `vassal`
```shell
Usage: counters vassal --csv=STRING --vassal-output-folder=STRING

Create a vassal module for testing. It searches for the 'template.xml' in the same folder

Flags:
  -h, --help                           Show context-sensitive help.

      --csv=STRING                     Input path of the file to read. Be aware that some outputs requires specific inputs.
      --vassal-output-folder=STRING    Path to the folder to write the image(s).
```

#### Usage of the Vassal mode

Vassal mode requires 2 things:
1. The CSV input string must follow a strict format.
2. The path + name of the vassal file `*.vmod` to write

This is an example CSV file to use as input:

```csv
0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,Side,Bg
,,,Batallion,,,,,,2,,,,9,,,,Red,F00
,,,Company,,,,,,3,,,,3,,,,Red,F00
,,,Screen,,,,,,,,,,,,,,Markers,E0E0E0
```

It has a header, here it is `0-16, Side, Bg`. `0-16` refers to the position in a counter, being 0 the vertical and 
horizontal center, and 1 to 16 are positions in the edge of the counter, being 1 the top left, 5 the top right, 
9 the bottom right and 13 the bottom left. The in-between numbers are in-between positions.

Side refers to a Red/Blue or Good Guys / Bad guys type of side. The module will have markers for as many sides it finds
here. So you can create red side, blue side, markers and others, for example. Each will have its own tab.

Finally, Bg refers to the background color in hexadecimal without `#` in the prefix.

# Commands
## `assets`
### CSV
Cards can also be generated using a CSV input file. the format of the file is like this:
```csv
config,,,,,,
rows,7,,,,,
cols,10,,,,,
width,150,,,,,
height,200,,,,,
font_height,12,,,,,
multiplier,first,second,third,fourth,fifth,sixth
30,first,second,third,fourth,fifth,sixth
```
`first,second,third,fourth,fifth,sixth` are just placeholders, you can use as many as you want but probably no more than 6-8. The minimum is to have the `multiplier column` and one column more, which will create a card with a single text in the middle. 

## `json`
## `vassal`

# Flags
## `input-path`

## `output-path`

## `output-content`
