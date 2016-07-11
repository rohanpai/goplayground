         
func TestInfluxdb(t *testing.T) {
	machineName := "mymachine"
	tablename := "table"
	database := "db"
	username := "root"
	password := "root"
	hostname := "localhost:8086"
	percentilesDuration := 10 * time.Minute
	config := &influxdb.ClientConfig{
		Host:     hostname,
		Username: username,
		Password: password,
		Database: database,
	}
	client, err := influxdb.NewClient(config)
	if err != nil {
		t.Fatal(err)
	}
	deleteAll := fmt.Sprintf("drop series %v", tablename)
	_, err = client.Query(deleteAll)
	if err != nil {
		t.Fatal(err)
	}

	defer client.Query(deleteAll)

	series := &influxdb.Series{
		Name:    tablename,
		Columns: []string{"col1"},
		Points: [][]interface{}{
			[]interface{}{1.0},
		},
	}
	err = client.WriteSeries([]*influxdb.Series{series})
	if err != nil {
		t.Fatal(err)
	}
}   