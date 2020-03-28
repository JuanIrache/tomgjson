# tomgjson

Converts time based data sources to [Adobe's mgJSON format](https://github.com/JuanIrache/mgjson) (Motion Graphics JSON).

A live version of this app can be found here: [To mgJSON](https://goprotelemetryextractor.com/csv-gpx-to-mgjson/)

The first iterations will be very limited in terms of customization and supported formats, but the provided parser funcs can potentially be adapted to any source that associates time with values.

For instructions on how to use mgJSON files in After Effects, see [Work with Data-driven animation](https://helpx.adobe.com/after-effects/using/data-driven-animations.html)

## Supported input formats

See the **sample_sources** directory for compatible samples.

## Usage

```go
src, _ := ioutil.ReadFile("./sample_sources/multiple-data.csv")
converted := FromCSV(src, 0)
doc := tomgjson.ToMgjson(converted, "Author Name")
f, _ := os.Create("./out.mgjson")
f.Write(doc)
f.Close()
```

See **all_test.go** for implementation examples.

### CSV

The simplest CSV file supported is a column with numbers. When a frame rate is specified, every value will be assigned a time based on the frame rate. Optionally, a header can be included in order to label the data. If the desired times do not correspond to the frame rate, a left-aligned "milliseconds" column can be used to specify the times relative to the beginning of the video. Additional columns with different labels can be appended to the right-hand side of the document to create new streams.

### GPX

GPS tracks with time fields can be parsed. For now, only the first track and its first track segment will be read. Based on the parsed data, additional data streams can be computed (speed, acceleration, course direction, distance...).

## Sample project templates

You can find sample After Effects projects that use mgJSON files on the [GoPro Telemetry Extractor page](https://goprotelemetryextractor.com). Look for the **Lite templates**.

## Videos made with mgJSON

- [mgJSON playlist on YouTube](https://www.youtube.com/playlist?list=PLgoeWSWqXedI7FbZccAEudt2_t8qPX0Px)

If you create something with mgJSON, let me know and I'll add it to the list.

## Software that supports mgJSON

These apps can output mgJSON files:

- [GoPro Telemetry Extractor](https://goprotelemetryextractor.com)
- [DJI Telemetry Overlay](https://djitelemetryoverlay.com)

## Contribution

Please make your changes to the **dev** branch, so that automated tests can be run before merging to **master**. Also, if possible, provide tests for new functionality.

## To-Do

- Support text columns in csv
- Support parsing column of strings in csv

## Might-Do

- Import from json and other formats
- Enable creating static fields? Non interpolable fields? Multidimensional fields?
