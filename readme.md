# tomgjson

Experimenting with converting other time based data sources to [Adobe's mgJSON format](https://github.com/JuanIrache/mgjson)

The first iteration will be very limited in terms of customization and supported formats, but the provided parser funcs can potentially be adapted to any source that associates time with values.

For instructions on how to use mgJSON files in After Effects, see [Work with Data-driven animation](https://helpx.adobe.com/after-effects/using/data-driven-animations.html)

## Sample project templates

You can find sample After Effects projects that use mgJSON files on the [GoPro Telemetry Extractor page](https://goprotelemetryextractor.com). Look for the **Lite templates**.

## Videos made with mgJSON

- [mgJSON playlist on YouTube](https://www.youtube.com/playlist?list=PLgoeWSWqXedI7FbZccAEudt2_t8qPX0Px)

If you create something with mgJSON, let me know and I'll add it to the list.

## Software that supports mgJSON

These apps can output mgJSON files:

- [GoPro Telemetry Extractor](https://goprotelemetryextractor.com)
- [DJI Telemetry Overlay](https://djitelemetryoverlay.com)

## To-Do

- Create documentation (tests/examples)
- Validate GPX results
- Create tutorial
- Use in production tool

## Might-Do

- Import from json and other formats
- Enable creating static fields? Non interpolable fields? Text fields? Multidimensional fields?
