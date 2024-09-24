package main

import (
	"github.com/gocarina/gocsv"
	"log"
	"os"
	"time"
)

type Row struct {
	StrValue   string    `csv:"str_value"`
	IntValue   int       `csv:"int_value"`
	TimeValue  time.Time `csv:"time_value"`
	FloatValue float64   `csv:"float_value"`
	BoolValue  bool      `csv:"bool_value"`

	StrArrayValue   []string    `csv:"str_array_value"`
	IntArrayValue   []int       `csv:"int_array_value"`
	TimeArrayValue  []time.Time `csv:"time_array_value"`
	FloatArrayValue []float64   `csv:"float_array_value"`
	BoolArrayValue  []bool      `csv:"bool_array_value"`

	NotUsed string `csv:"-"`
}

type DateTime struct {
	T time.Time
}

func (dt *DateTime) String() string {
	return dt.T.Format(time.DateTime)
}

func (dt *DateTime) UnmarshalCSV(csv string) (err error) {
	dt.T, err = time.Parse(time.DateTime, csv)
	return err
}

func (dt *DateTime) MarshalCSV() (string, error) {
	return dt.String(), nil
}

func (dt *DateTime) MarshalJSON() ([]byte, error) {
	return []byte(`"` + dt.String() + `"`), nil
}

func (dt *DateTime) UnmarshalJSON(data []byte) (err error) {
	dt.T, err = time.Parse(`"`+time.DateTime+`"`, string(data))
	return err
}

type TimeTestRow struct {
	Time    DateTime    `csv:"time"`
	TimeArr []*DateTime `csv:"time_arr"`
}

func main() {
	now := time.Unix(time.Now().Unix(), 0).UTC()
	now = now.Add(-(time.Duration(now.Second()) * time.Second))

	writeCsv("test.csv", []*Row{
		{
			StrValue:   "hello",
			IntValue:   1,
			TimeValue:  now,
			FloatValue: 1.1,
			BoolValue:  true,

			StrArrayValue:   []string{"a", "b", "c"},
			IntArrayValue:   []int{1, 2, 3},
			TimeArrayValue:  []time.Time{now, now, now},
			FloatArrayValue: []float64{1.1, 2.2, 3.3},
			BoolArrayValue:  []bool{true, false, true},

			NotUsed: "not used",
		},
		{
			StrValue:   "world",
			IntValue:   2,
			TimeValue:  time.Now(),
			FloatValue: 2.2,
			BoolValue:  false,

			StrArrayValue:   []string{"d", "e", "f"},
			IntArrayValue:   []int{4, 5, 6},
			TimeArrayValue:  []time.Time{now, now, now},
			FloatArrayValue: []float64{4.4, 5.5, 6.6},
			BoolArrayValue:  []bool{false, true, false},
		},
	})

	rows := readCsv[Row]("test.csv")
	for _, row := range rows {
		log.Printf("%+v", row)
	}

	writeCsv("test2.csv", []*TimeTestRow{
		{
			Time:    DateTime{T: now},
			TimeArr: []*DateTime{{T: now}, {T: now}},
		},
	})
	timeTest := readCsv[TimeTestRow]("test2.csv")
	log.Printf("%+v", timeTest[0])
}

func writeCsv[T any](csvFile string, rows []*T) {
	f, err := os.Create(csvFile)
	if err != nil {
		log.Fatalf("failed to create file: %v", err)
	}
	defer f.Close()

	if err := gocsv.MarshalFile(&rows, f); err != nil {
		log.Fatalf("failed to write csv: %v", err)
	}

}

func readCsv[T any](csvFile string) []*T {
	f, err := os.Open(csvFile)
	if err != nil {
		log.Fatalf("failed to open file: %v", err)
	}
	defer f.Close()

	var rows []*T
	if err := gocsv.UnmarshalFile(f, &rows); err != nil {
		log.Fatalf("failed to read csv: %v", err)
	}

	return rows
}
