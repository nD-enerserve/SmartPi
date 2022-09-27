package smartpi

import (
	// "database/sql"
	// "fmt"
	// "os"

	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/FransTheekrans/SmartPi/src/smartpi/network"
	client "github.com/influxdata/influxdb1-client/v2"
	log "github.com/sirupsen/logrus"
)

// type MinuteValues struct {
// 	Date                                                                                                                                                                                                                                                            time.Time
// 	Current_1, Current_2, Current_3, Current_4, Voltage_1, Voltage_2, Voltage_3, Power_1, Power_2, Power_3, Cosphi_1, Cosphi_2, Cosphi_3, Frequency_1, Frequency_2, Frequency_3, Energy_pos_1, Energy_pos_2, Energy_pos_3, Energy_neg_1, Energy_neg_2, Energy_neg_3 float64
// }

func InsertInfluxData(c *Config, t time.Time, v ReadoutAccumulator, consumedWattHourBalanced float64, producedWattHourBalanced float64) {

	dbc, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     c.Influxdatabase,
		Username: c.Influxuser,
		Password: c.Influxpassword,
	})
	if err != nil {
		log.Printf("Error creating InfluxDB Client: ", err.Error())
	}
	defer dbc.Close()

	// Create a new point batch
	bp, _ := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  "MeteringData",
		Precision: "s",
	})

	// Create a point and add to batch
	macaddress := network.GetMacAddr()
	tags := map[string]string{"serial": macaddress, "type": "electric"}
	fields := map[string]interface{}{
		"I1":      float64(v.Current[PhaseA]),
		"I2":      float64(v.Current[PhaseB]),
		"I3":      float64(v.Current[PhaseC]),
		"I4":      float64(v.Current[PhaseN]),
		"U1":      float64(v.Voltage[PhaseA]),
		"U2":      float64(v.Voltage[PhaseB]),
		"U3":      float64(v.Voltage[PhaseC]),
		"P1":      float64(v.ActiveWatts[PhaseA]),
		"P2":      float64(v.ActiveWatts[PhaseB]),
		"P3":      float64(v.ActiveWatts[PhaseC]),
		"CosPhi1": float64(v.CosPhi[PhaseA]),
		"CosPhi2": float64(v.CosPhi[PhaseB]),
		"CosPhi3": float64(v.CosPhi[PhaseC]),
		"F1":      float64(v.Frequency[PhaseA]),
		"F2":      float64(v.Frequency[PhaseB]),
		"F3":      float64(v.Frequency[PhaseC]),
		"Ec1":     float64(v.WattHoursConsumed[PhaseA]),
		"Ec2":     float64(v.WattHoursConsumed[PhaseB]),
		"Ec3":     float64(v.WattHoursConsumed[PhaseC]),
		"Ep1":     float64(v.WattHoursProduced[PhaseA]),
		"Ep2":     float64(v.WattHoursProduced[PhaseB]),
		"Ep3":     float64(v.WattHoursProduced[PhaseC]),
		"Eb1":     float64(v.WattHoursBalanced[PhaseA]),
		"Eb2":     float64(v.WattHoursBalanced[PhaseB]),
		"Eb3":     float64(v.WattHoursBalanced[PhaseC]),
		"bEc":     float64(consumedWattHourBalanced),
		"bEp":     float64(producedWattHourBalanced),
	}
	pt, err := client.NewPoint("data", tags, fields, time.Now())
	if err != nil {
		log.Printf("Error: ", err.Error())
	}
	bp.AddPoint(pt)

	// Write the batch
	err = dbc.Write(bp)
	if err != nil {
		log.Printf("Error: ", err.Error())
	}
}

func ReadCSVData(c *Config, starttime time.Time, endtime time.Time) string {

	loc, _ := time.LoadLocation("UTC")
	// endtime := time.Now()
	// starttime := endtime.Add(time.Minute * -10)

	querystring := url.QueryEscape("SELECT mean(\"P1\") AS \"P1\", mean(\"P2\") AS \"P2\", mean(\"P3\") AS \"P3\", mean(\"bEc\") AS \"bEc\", mean(\"bEp\") AS \"bEp\", mean(\"I1\") AS \"i1\", mean(\"I2\") AS \"i2\", mean(\"I3\") AS \"i3\", mean(\"I4\") AS \"i4\", mean(\"U1\") AS \"u1\", mean(\"U2\") AS \"u2\", mean(\"U3\") AS \"u3\", mean(\"Ec1\") AS \"Ec1\", mean(\"Ec2\") AS \"Ec2\", mean(\"Ec3\") AS \"Ec3\", mean(\"Ep1\") AS \"Ep1\", mean(\"Ep2\") AS \"Ep2\", mean(\"Ep3\") AS \"Ep3\", mean(\"CosPhi1\") AS \"CosPhi1\", mean(\"CosPhi2\") AS \"CosPhi2\", mean(\"CosPhi3\") AS \"CosPhi3\", mean(\"F1\") AS \"f1\", mean(\"F2\") AS \"f2\", mean(\"F3\") AS \"f3\" FROM \"data\" WHERE time >= '" + starttime.In(loc).Format("2006-01-02 15:04:05") + "' AND time <= '" + endtime.In(loc).Format("2006-01-02 15:04:05") + "' GROUP BY time(1m) fill(null)")
	req, err := http.NewRequest("POST", c.Influxdatabase+"/query?db=MeteringData&u="+c.Influxuser+"&p="+c.Influxpassword+"&q="+querystring, nil)
	if err != nil {
		// handle err
	}
	req.Header.Set("Accept", "application/csv")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		// handle err
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	stringbody := strings.Replace(string(body), ",", ";", -1)
	stringbody = strings.Replace(string(stringbody), ".", c.CSVdecimalpoint, -1)

	return stringbody
}

func ExampleClient_query(c *Config) {

	dbc, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     c.Influxdatabase,
		Username: c.Influxuser,
		Password: c.Influxpassword,
	})
	if err != nil {
		log.Printf("Error creating InfluxDB Client: ", err.Error())
	}
	defer dbc.Close()

	q := client.NewQuery("SELECT * FROM MeteringData", "MeteringData", "ns")
	if response, err := dbc.Query(q); err == nil && response.Error() == nil {
		fmt.Println(response.Results)
	}
}
