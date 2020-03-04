# tomgjson

Experimenting with converting other time based data sources to [Adobe's mgJSON format](https://github.com/JuanIrache/mgjson)

The first iteration will be very limited in terms of customization.

## Sample projects

You can find sample projects that use mgJSON files on the [GoPro Telemetry Extractor page](https://goprotelemetryextractor.com#ae). Look for the Lite templates.

## Videos made with mgJSON

- [mgJSON playlist on YouTube](https://youtu.be/TAdxsTv4hPU?list=PLgoeWSWqXedI7FbZccAEudt2_t8qPX0Px)

If you create something with mgJSON, let me know and I'll add it to the list.

## Software that supports mgJSON

These apps can output mgJSON files:

- [GoPro Telemetry Extractor](https://goprotelemetryextractor.com)
- [DJI Telemetry Overlay](https://djitelemetryoverlay.com)

## To-Do

- Read from CSV
  - Test with other csv standards or enconding (\r\n?)
- Read from GPX
  - Read standard track
  - Interpret additional data (speeds, acceleration, course direction...)
- Add func to read formatted DataSource directly
- Structure module properly
- Create documentation (tests/examples)
- Create tutorial
- Use in production tool

## Might-Do

- Enable creating static fields? Non interpolable fields? Text fields?
