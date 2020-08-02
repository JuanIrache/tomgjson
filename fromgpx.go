package tomgjson

import (
	"encoding/xml"
	"fmt"
	"math"
	"time"
)

func degreesToRadians(degrees float64) float64 {
	return degrees * math.Pi / 180
}

func radiansToDegrees(radians float64) float64 {
	return radians * 180 / math.Pi
}

func distanceInMBetweenEarthCoordinates(lat1, lon1, lat2, lon2 float64) float64 {
	earthRadiusM := 6378137.0

	dLat := degreesToRadians(lat2 - lat1)
	dLon := degreesToRadians(lon2 - lon1)

	lat1 = degreesToRadians(lat1)
	lat2 = degreesToRadians(lat2)

	a := math.Sin(dLat/2)*math.Sin(dLat/2) + math.Sin(dLon/2)*math.Sin(dLon/2)*math.Cos(lat1)*math.Cos(lat2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return earthRadiusM * c
}

func angleFromCoordinate(lat1, lon1, lat2, lon2, prev float64) float64 {

	lat1 = degreesToRadians(lat1)
	lon1 = degreesToRadians(lon1)
	lat2 = degreesToRadians(lat2)
	lon2 = degreesToRadians(lon2)

	x := math.Cos(lat2) * math.Sin(lon2-lon1)

	y := math.Cos(lat1)*math.Sin(lat2) - math.Sin(lat1)*math.Cos(lat2)*math.Cos(lon2-lon1)

	course := radiansToDegrees(math.Atan2(x, y))

	for math.Abs(course-prev) > 180 {
		if math.Signbit(course - prev) {
			course += 360
		} else {
			course -= 360
		}
	}

	return course
}

var ids = []string{
	// Supported GPX streams
	"lat (°)",
	"lon (°)",
	"ele (m)",
	"magvar (°)",
	"geoidheight (m)",
	"fix",
	"sat",
	"hdop",
	"vdop",
	"pdop",
	"ageofdgpsdata (s)",
	"dgpsid",
	// Calculated streams
	"distance2d (m)",
	"distance3d (m)",
	"verticalSpeed (m/s)",
	"speed2d (m/s)",
	"speed3d (m/s)",
	"acceleration2d (m/s²)",
	"acceleration3d (m/s²)",
	"verticalAcceleration (m/s²)",
	"course (°)",
	"slope (°)",
	// Additional explicit date string
	"time",
}

// Return index of stream
func idx(item string) int {
	for i, v := range ids {
		if v == item {
			return i
		}
	}
	return -1
}

func appendToFloatStream(data FormattedData, v *float64, n string) Stream {
	st := data.Streams[idx(n)]
	if v != nil {
		st.Values = append(st.Values, *v)
		// Name confirmed streams
		st.Label = n
	} else {
		// Appending zeros for now, could potentially add interpolated values based on previous, next and time
		st.Values = append(st.Values, 0)
	}
	return st
}

func appendToStringStream(data FormattedData, s *string, n string) Stream {
	st := data.Streams[idx(n)]
	st.Strings = append(st.Strings, *s)
	// Name confirmed streams
	st.Label = n
	return st
}

// FromGPX formats a compatible GPX file as a struct ready for mgJSON and returns it. Or returns an error
// The optional extra bool will compute additional streams based on the existing data
func FromGPX(src []byte, extra bool) (FormattedData, error) {

	var data FormattedData

	utc, err := time.LoadLocation("UTC")
	if err != nil {
		return data, err
	}

	type Trkpt struct {
		XMLName       xml.Name `xml:"trkpt"`
		Time          *string  `xml:"time"`
		Lat           *float64 `xml:"lat,attr"`
		Lon           *float64 `xml:"lon,attr"`
		Ele           *float64 `xml:"ele"`
		Magvar        *float64 `xml:"magvar"`
		Geoidheight   *float64 `xml:"geoidheight"`
		Fix           *float64 `xml:"fix"`
		Sat           *float64 `xml:"sat"`
		Hdop          *float64 `xml:"hdop"`
		Vdop          *float64 `xml:"vdop"`
		Pdop          *float64 `xml:"pdop"`
		Ageofdgpsdata *float64 `xml:"ageofdgpsdata"`
		Dgpsid        *float64 `xml:"dgpsid"`
	}

	type Trkseg struct {
		XMLName xml.Name `xml:"trkseg"`
		Trkpt   []Trkpt  `xml:"trkpt"`
	}

	type Trk struct {
		XMLName xml.Name `xml:"trk"`
		Trkseg  []Trkseg `xml:"trkseg"`
	}

	type Gpx struct {
		XMLName xml.Name `xml:"gpx"`
		Trk     []Trk    `xml:"trk"`
	}

	gpx := Gpx{}

	err = xml.Unmarshal(src, &gpx)
	if err != nil {
		return data, err
	}

	if len(gpx.Trk) < 1 {
		return data, fmt.Errorf("Error: No GPX tracks")
	}

	// Just reading one track for now
	if len(gpx.Trk[0].Trkseg) < 1 {
		return data, fmt.Errorf("Error: No GPX trkseg")
	}

	trkpts := []Trkpt{}

	for _, trkseg := range gpx.Trk[0].Trkseg {
		trkpts = append(trkpts, trkseg.Trkpt...)
	}

	if len(trkpts) < 2 {
		return data, fmt.Errorf("Error: Not enough GPX trkpt")
	}

	// One Stream for each of the supported trkpt and custom fields
	data.Streams = make([]Stream, len(ids))
	data.Timing = make([]time.Time, len(trkpts))

	for _, st := range data.Streams {
		st.Values = make([]float64, len(trkpts))
	}

	for i, trkpt := range trkpts {
		if trkpt.Time == nil {
			return data, fmt.Errorf("Error: Missing timiing data in GPX")
		}
		t, err := time.Parse(time.RFC3339, *trkpt.Time)
		if err != nil {
			return data, err
		}

		t = t.In(utc)

		data.Timing[i] = t
		data.Streams[idx("lat (°)")] = appendToFloatStream(data, trkpt.Lat, "lat (°)")
		data.Streams[idx("lon (°)")] = appendToFloatStream(data, trkpt.Lon, "lon (°)")
		data.Streams[idx("ele (m)")] = appendToFloatStream(data, trkpt.Ele, "ele (m)")
		data.Streams[idx("magvar (°)")] = appendToFloatStream(data, trkpt.Magvar, "magvar (°)")
		data.Streams[idx("geoidheight (m)")] = appendToFloatStream(data, trkpt.Geoidheight, "geoidheight (m)")
		data.Streams[idx("fix")] = appendToFloatStream(data, trkpt.Fix, "fix")
		data.Streams[idx("sat")] = appendToFloatStream(data, trkpt.Sat, "sat")
		data.Streams[idx("hdop")] = appendToFloatStream(data, trkpt.Hdop, "hdop")
		data.Streams[idx("vdop")] = appendToFloatStream(data, trkpt.Vdop, "vdop")
		data.Streams[idx("pdop")] = appendToFloatStream(data, trkpt.Pdop, "pdop")
		data.Streams[idx("ageofdgpsdata (s)")] = appendToFloatStream(data, trkpt.Ageofdgpsdata, "ageofdgpsdata (s)")
		data.Streams[idx("dgpsid")] = appendToFloatStream(data, trkpt.Dgpsid, "dgpsid")

		data.Streams[idx("time")] = appendToStringStream(data, trkpt.Time, "time")

		// Computed streams
		if extra && trkpt.Lat != nil && trkpt.Lon != nil {
			var distance2d float64
			var speed2d float64
			var acceleration2d float64
			var course float64
			var slope float64
			var distance3d float64
			var speed3d float64
			var acceleration3d float64
			var verticalSpeed float64
			var verticalAcceleration float64
			if i > 0 {
				prevLat := data.Streams[idx("lat (°)")].Values[i-1]
				prevLon := data.Streams[idx("lon (°)")].Values[i-1]
				distance2d = distanceInMBetweenEarthCoordinates(*trkpt.Lat, *trkpt.Lon, prevLat, prevLon)
				duration := data.Timing[i].Sub(data.Timing[i-1]).Seconds()
				//Make sure duration is not zero
				duration = math.Max(duration, 1e-9)
				speed2d = distance2d / duration
				acceleration2d = speed2d
				prevCourse := data.Streams[idx("course (°)")].Values[i-1]
				course = angleFromCoordinate(*trkpt.Lat, *trkpt.Lon, prevLat, prevLon, prevCourse)
				if trkpt.Ele != nil {
					prevEle := data.Streams[idx("ele (m)")].Values[i-1]
					verticalDist := *trkpt.Ele - prevEle
					slope = math.Atan2(verticalDist, distance2d)
					slope = radiansToDegrees(slope)
					distance3d = math.Sqrt(math.Pow(verticalDist, 2) + math.Pow(distance2d, 2))
					speed3d = distance3d / duration
					acceleration3d = speed3d
					verticalSpeed = verticalDist / duration
					verticalAcceleration = verticalSpeed
				}
				if i > 1 {
					prevDistance := data.Streams[idx("distance2d (m)")].Values[i-1]
					distance2d += prevDistance
					prevSpeed2d := data.Streams[idx("speed2d (m/s)")].Values[i-1]
					speed2dChange := speed2d - prevSpeed2d
					acceleration2d = speed2dChange / duration
					if trkpt.Ele != nil {
						prevDistance3d := data.Streams[idx("distance3d (m)")].Values[i-1]
						distance3d += prevDistance3d
						prevSpeed3d := data.Streams[idx("speed3d (m/s)")].Values[i-1]
						speed3dChange := speed3d - prevSpeed3d
						acceleration3d = speed3dChange / duration
						prevVerticalSpeed := data.Streams[idx("verticalSpeed (m/s)")].Values[i-1]
						verticalSpeedChange := verticalSpeed - prevVerticalSpeed
						verticalAcceleration = verticalSpeedChange / duration
					}
				}
			}

			data.Streams[idx("distance2d (m)")] = appendToFloatStream(data, &distance2d, "distance2d (m)")
			data.Streams[idx("speed2d (m/s)")] = appendToFloatStream(data, &speed2d, "speed2d (m/s)")
			data.Streams[idx("acceleration2d (m/s²)")] = appendToFloatStream(data, &acceleration2d, "acceleration2d (m/s²)")
			data.Streams[idx("course (°)")] = appendToFloatStream(data, &course, "course (°)")
			data.Streams[idx("slope (°)")] = appendToFloatStream(data, &slope, "slope (°)")
			data.Streams[idx("distance3d (m)")] = appendToFloatStream(data, &distance3d, "distance3d (m)")
			data.Streams[idx("speed3d (m/s)")] = appendToFloatStream(data, &speed3d, "speed3d (m/s)")
			data.Streams[idx("acceleration3d (m/s²)")] = appendToFloatStream(data, &acceleration3d, "acceleration3d (m/s²)")
			data.Streams[idx("verticalSpeed (m/s)")] = appendToFloatStream(data, &verticalSpeed, "verticalSpeed (m/s)")
			data.Streams[idx("verticalAcceleration (m/s²)")] = appendToFloatStream(data, &verticalAcceleration, "verticalAcceleration (m/s²)")
		}
	}

	// Clean up unconfirmed streams
	for i := len(data.Streams) - 1; i >= 0; i-- {
		if len(data.Streams[i].Label) < 1 {
			copy(data.Streams[i:], data.Streams[i+1:])
			data.Streams[len(data.Streams)-1] = Stream{}
			data.Streams = data.Streams[:len(data.Streams)-1]
		}
	}

	return data, nil
}
