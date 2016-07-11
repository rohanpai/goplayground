package main

import (
	&#34;fmt&#34;
	&#34;strconv&#34;
)

func main() {
	fmt.Println(&#34;Float32ToString:&#34;, Float32ToString(float32(55.55)))
	fmt.Println(&#34;Float64ToString:&#34;, Float64ToString(float64(55.55)))

	fmt.Println(&#34;Int32ToString:&#34;, Int32ToString(int32(-44)))
	fmt.Println(&#34;Int64ToString:&#34;, Int64ToString(int64(-44)))

	fmt.Println(&#34;Uint32ToString:&#34;, Uint32ToString(uint32(44)))
	fmt.Println(&#34;Uint64ToString:&#34;, Uint64ToString(uint64(44)))

	string1, _ := strconv.Atoi(&#34;66&#34;)
	fmt.Println(&#34;StringToInt with strconv.Atoi:&#34;, string1)
	fmt.Println(&#34;IntToString with strconv.Itoa (Actually it is int64):&#34;, strconv.Itoa(11))
	fmt.Println(&#34;BoolToString:&#34;, BoolToString(true))

	float1, _ := StringToFloat32(&#34;22.22&#34;)
	float2, _ := StringToFloat64(&#34;22.22&#34;)
	fmt.Println(&#34;StringToFloat32:&#34;, float1)
	fmt.Println(&#34;StringToFloat64:&#34;, float2)

	int1, _ := StringToInt32(&#34;-33&#34;)
	int2, _ := StringToInt64(&#34;-33&#34;)
	fmt.Println(&#34;StringToInt32:&#34;, int1)
	fmt.Println(&#34;StringToInt64:&#34;, int2)

	uint1, _ := StringToUint32(&#34;33&#34;)
	uint2, _ := StringToUint64(&#34;33&#34;)
	fmt.Println(&#34;StringToUint32:&#34;, uint1)
	fmt.Println(&#34;StringToUint64:&#34;, uint2)
}

func Float32ToString(value float32) string {
	return Float64ToString(float64(value))
}

func Float64ToString(value float64) string {
	return strconv.FormatFloat(value, &#39;f&#39;, -1, 64)
}

func Int32ToString(value int32) string {
	return Int64ToString(int64(value))
}

func Int64ToString(value int64) string {
	return strconv.FormatInt(value, 10)
}

func Uint32ToString(value uint32) string {
	return Uint64ToString(uint64(value))
}

func Uint64ToString(value uint64) string {
	return strconv.FormatUint(value, 10)
}

func BoolToString(value bool) string {
	return strconv.FormatBool(value)
}

func StringToFloat32(value string) (float32, error) {
	result, err := StringToFloat64(value)
	if err != nil {
		return 0, err
	}
	return float32(result), nil
}

func StringToFloat64(value string) (float64, error) {
	result, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0, err
	}
	return result, nil
}

func StringToInt32(value string) (int32, error) {
	result, err := StringToInt64(value)
	if err != nil {
		return 0, err
	}
	return int32(result), nil
}

func StringToInt64(value string) (int64, error) {
	result, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0, err
	}
	return result, nil
}

func StringToUint32(value string) (uint32, error) {
	result, err := StringToUint64(value)
	if err != nil {
		return 0, err
	}
	return uint32(result), nil
}

func StringToUint64(value string) (uint64, error) {
	result, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return 0, err
	}
	return result, nil
}