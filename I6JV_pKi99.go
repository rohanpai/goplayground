         
func TestInfluxdb(t *testing.T) {
	machineName := &#34;mymachine&#34;
	tablename := &#34;table&#34;
	database := &#34;db&#34;
	username := &#34;root&#34;
	password := &#34;root&#34;
	hostname := &#34;localhost:8086&#34;
	percentilesDuration := 10 * time.Minute
	config := &amp;influxdb.ClientConfig{
		Host:     hostname,
		Username: username,
		Password: password,
		Database: database,
	}
	client, err := influxdb.NewClient(config)
	if err != nil {
		t.Fatal(err)
	}
	deleteAll := fmt.Sprintf(&#34;drop series %v&#34;, tablename)
	_, err = client.Query(deleteAll)
	if err != nil {
		t.Fatal(err)
	}

	defer client.Query(deleteAll)

	series := &amp;influxdb.Series{
		Name:    tablename,
		Columns: []string{&#34;col1&#34;},
		Points: [][]interface{}{
			[]interface{}{1.0},
		},
	}
	err = client.WriteSeries([]*influxdb.Series{series})
	if err != nil {
		t.Fatal(err)
	}
}   