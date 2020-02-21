# tomgjson

Experimenting with converting other time based data sources to [Adobe's mgJSON format](https://github.com/JuanIrache/mgjson)

## To-Do

- Read from CSV
  - Support files with multiple streams (sharing time)
  - If specified time (milliseconds?), use it
  - If specified date, use it (maybe epoch, same column as milliseconds?)
  - Otherwise, if specified frame rate, use it, or use PAL standard
  - Test with other csv standards or enconding (\r\n?)
- General
  - Pad number strings
  - Recreate mgJSON structure
- Convert to mgJSON
- Create documentation (tests/examples)
- Create tutorial
- Use in production tool

## Might-Do

- Enable creating static fields? Non interpolable fields? Text fields?
