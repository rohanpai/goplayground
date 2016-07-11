package main

import "fmt"

// GeoPoint maps against Postgis geographical point
type GeoPoint struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

func (p *GeoPoint) String() string {
	return fmt.Sprintf("POINT(%v %v)", p.Lat, p.Lng)
}

// Scan implements the Scanner interface and will scan the Postgis POINT(x y) into the GeoPoint struct
func (p *GeoPoint) Scan(val interface{}) error {
	b, err := hex.DecodeString(string(val.([]uint8)))
	if err != nil {
		return err
	}

	r := bytes.NewReader(b)
	var wkbByteOrder uint8
	if err := binary.Read(r, binary.LittleEndian, &wkbByteOrder); err != nil {
		return err
	}

	var byteOrder binary.ByteOrder
	switch wkbByteOrder {
	case 0:
		byteOrder = binary.BigEndian
	case 1:
		byteOrder = binary.LittleEndian
	default:
		return fmt.Errorf("invalid byte order %u", wkbByteOrder)
	}

	var wkbGeometryType uint64
	if err := binary.Read(r, byteOrder, &wkbGeometryType); err != nil {
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
	fmt.Println("Hello, playground")
}
