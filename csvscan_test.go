package csvscan

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"testing"
)

type CustomerID string

func (ci *CustomerID) UnmarshalText(text []byte) error {
	panic("not implemented")
}

func TestXXX(t *testing.T) {
	type Customer struct {
		CustomerID       *CustomerID `csv:"Customer Id"`
		FirstName        *string     `csv:"First Name"`
		LastName         string      `csv:"Last Name"`
		Company          *string     `csv:"Company"`
		City             string      `csv:"City"`
		Country          string      `csv:"Country"`
		Phone_1          string      `csv:"Phone 1"`
		Phone_2          string      `csv:"Phone 2"`
		Email            string      `csv:"Email"`
		SubscriptionDate string      `csv:"Subscription Date"`
		Website          string      `csv:"Website"`
		Index            int         `csv:"Index"`
	}

	file, err := os.OpenFile("customers-10.csv", os.O_RDONLY, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	scn := New()
	if err := scn.Init(file, (*Customer)(nil)); err != nil {
		t.Fatal(err)
	}

	for {
		var tmp Customer

		if err := scn.Scan(&tmp); err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			t.Fatal(err)
			return
		}

		buf, err := json.MarshalIndent(tmp, "", "\t")
		if err != nil {
			t.Fatal(err)
		}

		fmt.Fprint(os.Stderr, "\n", string(buf), "\n")
	}
}

type X string

func (x *X) UnmarshalJSON(_ []byte) error {
	panic("not implemented") // TODO: Implement
}

func TestJSON(t *testing.T) {
	type Y struct {
		X *X
	}

	var z Y
	err := json.Unmarshal([]byte(`{"x":""}`), &z)
	if err != nil {
		t.Fatal(err)
	}
}
