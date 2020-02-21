# tomgjson

Experimenting with converting other time based data sources to [Adobe's mgJSON format](https://github.com/JuanIrache/mgjson)

## To-Do

- Read from CSV
  - If headers, use for labelling
  - If specified time (milliseconds?), use it
  - If specified date, use it (maybe epoch, same column as milliseconds?)
  - Otherwise, if specified frame rate, use it
- General
  - Pad number strings
  - Recreate mgJSON structure
- Convert to mgJSON
- Create documentation
- Create tutorial
- Use in production tool

## Might-Do

- Enable creating static fields? Non interpolable fields? Text fields?
- Support files with multiple streams (sharing time)
