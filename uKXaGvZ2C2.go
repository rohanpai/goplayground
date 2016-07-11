package main

import &#34;fmt&#34;

// GeoPoint maps against Postgis geographical point
type GeoPoint struct {
	Lat float64 `json:&#34;lat&#34;`
	Lng float64 `json:&#34;lng&#34;`
}

func (p *GeoPoint) String() string {
	return fmt.Sprintf(&#34;POINT(%v %v)&#34;, p.Lat, p.Lng)
}

// Scan implements the Scanner interface and will scan the Postgis POINT(x y) into the GeoPoint struct
func (p *GeoPoint) Scan(val interface{}) error {
	b, err := hex.DecodeString(string(val.([]uint8)))
	if err != nil {
		return err
	}

	r := bytes.NewReader(b)
	var wkbByteOrder uint8
	if err := binary.Read(r, binary.LittleEndian, &amp;wkbByteOrder); err != nil {
		return err
	}

	var byteOrder binary.ByteOrder
	switch wkbByteOrder {
	case 0:
		byteOrder = binary.BigEndian
	case 1:
		byteOrder = binary.LittleEndian
	default:
		return fmt.Errorf(&#34;invalid byte order %u&#34;, wkbByteOrder)
	}

	var wkbGeometryType uint64
	if err := binary.Read(r, byteOrder, &amp;wkbGeometryType); err != nil {
		return err
	}

	if err := binary.Read(r, byteOrder, p); err != nil {
		return err
	}

	return nil
}

// Value implements the driver Valuer interface and will return the string representation of the GeoPoint struct by calling the String() method
func (p *GeoPoint) Value() (driver.Value, error) {
	return p.String(), nil
}

func main() {
	fmt.Println(&#34;Hello, playground&#34;)
}
