package csvscan

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"os"
	"testing"
)

type CustomerID string

func (ci *CustomerID) UnmarshalText(text []byte) error {
	*ci = CustomerID(text)
	return nil
}

const expected_output = `[{"CustomerID":"dE014d010c7ab0c","FirstName":"Andrew","LastName":"Goodman","Company":"Stewart-Flynn","City":"Rowlandberg","Country":"Macao","Phone_1":"846-790-4623x4715","Phone_2":"(422)787-2331x71127","Email":"marieyates@gomez-spencer.info","SubscriptionDate":"2021-07-26","Website":"http://www.shea.biz/","Index":1},{"CustomerID":"2B54172c8b65eC3","FirstName":"Alvin","LastName":"Lane","Company":"Terry, Proctor and Lawrence","City":"Bethside","Country":"Papua New Guinea","Phone_1":"124-597-8652x05682","Phone_2":"321.441.0588x6218","Email":"alexandra86@mccoy.com","SubscriptionDate":"2021-06-24","Website":"http://www.pena-cole.com/","Index":2},{"CustomerID":"d794Dd48988d2ac","FirstName":"Jenna","LastName":"Harding","Company":"Bailey Group","City":"Moniquemouth","Country":"China","Phone_1":"(335)987-3085x3780","Phone_2":"001-680-204-8312","Email":"justincurtis@pierce.org","SubscriptionDate":"2020-04-05","Website":"http://www.booth-reese.biz/","Index":3},{"CustomerID":"3b3Aa4aCc68f3Be","FirstName":"Fernando","LastName":"Ford","Company":"Moss-Maxwell","City":"Leeborough","Country":"Macao","Phone_1":"(047)752-3122","Phone_2":"048.779.5035x9122","Email":"adeleon@hubbard.org","SubscriptionDate":"2020-11-29","Website":"http://www.hebert.com/","Index":4},{"CustomerID":"D60df62ad2ae41E","FirstName":"Kara","LastName":"Woods","Company":"Mccarthy-Kelley","City":"Port Jacksonland","Country":"Nepal","Phone_1":"+1-360-693-4419x19272","Phone_2":"163-627-2565","Email":"jesus90@roberson.info","SubscriptionDate":"2022-04-22","Website":"http://merritt.com/","Index":5},{"CustomerID":"8aaa5d0CE9ee311","FirstName":"Marissa","LastName":"Gamble","Company":"Cherry and Sons","City":"Webertown","Country":"Sudan","Phone_1":"001-645-334-5514x0786","Phone_2":"(751)980-3163","Email":"katieallison@leonard.com","SubscriptionDate":"2021-11-17","Website":"http://www.kaufman.org/","Index":6},{"CustomerID":"73B22Ac8A43DD1A","FirstName":"Julie","LastName":"Cooley","Company":"Yu, Norman and Sharp","City":"West Sandra","Country":"Japan","Phone_1":"+1-675-243-7422x9177","Phone_2":"(703)337-5903","Email":"priscilla88@stephens.info","SubscriptionDate":"2022-03-26","Website":"http://www.sexton-chang.com/","Index":7},{"CustomerID":"DC94CCd993D311b","FirstName":"Lauren","LastName":"Villa","Company":"French, Travis and Hensley","City":"New Yolanda","Country":"Fiji","Phone_1":"081.226.1797x647","Phone_2":"186.540.9690x605","Email":"colehumphrey@austin-caldwell.com","SubscriptionDate":"2020-08-14","Website":"https://www.kerr.com/","Index":8},{"CustomerID":"9Ba746Cb790FED9","FirstName":"Emily","LastName":"Bryant","Company":"Moon, Strickland and Combs","City":"East Normanchester","Country":"Seychelles","Phone_1":"430-401-5228x35091","Phone_2":"115-835-3840","Email":"buckyvonne@church-lutz.com","SubscriptionDate":"2020-12-30","Website":"http://grimes.com/","Index":9},{"CustomerID":"aAa1EDfaA70DA0c","FirstName":"Marie","LastName":"Estrada","Company":"May Inc","City":"Welchton","Country":"United Arab Emirates","Phone_1":"001-648-790-9244","Phone_2":"973-767-3611","Email":"christie44@mckenzie.biz","SubscriptionDate":"2020-09-03","Website":"https://www.salinas.net/","Index":10}]`

func TestScaner(t *testing.T) {
	type Customer struct {
		CustomerID       CustomerID `csv:"Customer Id"`
		FirstName        *string    `csv:"First Name"`
		LastName         string     `csv:"Last Name"`
		Company          *string    `csv:"Company"`
		City             string     `csv:"City"`
		Country          string     `csv:"Country"`
		Phone_1          string     `csv:"Phone 1"`
		Phone_2          string     `csv:"Phone 2"`
		Email            string     `csv:"Email"`
		SubscriptionDate string     `csv:"Subscription Date"`
		Website          string     `csv:"Website"`
		Index            int        `csv:"Index"`
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

	var customs []Customer

	for {
		var tmp Customer
		if err := scn.Scan(&tmp); err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			t.Fatal(err)
			return
		}
		customs = append(customs, tmp)
	}

	var bb bytes.Buffer
	if err := json.NewEncoder(&bb).Encode(customs); err != nil {
		t.Fatal(err)
	}

	if bb.String() != expected_output {
		t.Fatal("unexpected result")
	}
}
