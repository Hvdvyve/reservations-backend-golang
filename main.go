/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"sort"

	reservation "github.com/gitpod/mycli/Reservations"
	"github.com/gitpod/mycli/cmd"
	"github.com/kr/pretty"
	ravendb "github.com/ravendb/ravendb-go-client"
)

var (
	dbName        = "Reservations"
	serverNodeURL = "https://a.reservationdata.ravendb.community/studio/index.html"

	// if true, we'll show summary of HTTP requests made to the server
	// and dump full info about failed HTTP requests
	verboseLogging = true

	// if true, logs all http requests/responses to a file for further inspection
	// this is for use in tests so the file has a fixed location:
	// logs/trace_${test_name}_go.txt
	logAllRequests = false

	// if logAllRequests is true, this is a path of a file where we log
	// info about all HTTP requests
	logAllRequestsPath = "http_requests_log.txt"
)

func main() {
	cmd.Execute()
	getDocumentStore(dbName)
	openSession(dbName)
	queryGetConsumables()
	queryGetRoomsAndDesks()
	crudDeleteReservation("1")
	crudAddConsumable("Chocolate", "Food", "20", "True")
}

func getDocumentStore(databaseName string) (*ravendb.DocumentStore, error) {
	cerPath := "admin.client.certificate.reservationdata.crt"
	keyPath := "admin.client.certificate.reservationdata.key"
	serverNodes := []string{"https://a.reservationdata.ravendb.community/studio/index.html"}

	cer, err := tls.LoadX509KeyPair(cerPath, keyPath)
	if err != nil {
		return nil, err
	}
	store := ravendb.NewDocumentStore(serverNodes, databaseName)
	store.Certificate = &cer
	x509cert, err := x509.ParseCertificate(cer.Certificate[0])
	if err != nil {
		return nil, err
	}
	store.TrustStore = x509cert
	if err := store.Initialize(); err != nil {
		return nil, err
	}
	return store, nil
}

func openSession(databaseName string) (*ravendb.DocumentStore, *ravendb.DocumentSession, error) {
	store, err := getDocumentStore(dbName)
	if err != nil {
		return nil, nil, fmt.Errorf("getDocumentStore() failed with %s", err)
	}

	session, err := store.OpenSession("")
	if err != nil {
		return nil, nil, fmt.Errorf("store.OpenSession() failed with %s", err)
	}
	return store, session, nil
}

func printRQL(q *ravendb.DocumentQuery) {
	iq, err := q.GetIndexQuery()
	if err != nil {
		log.Fatalf("q.GetIndexQuery() returned '%s'\n", err)
	}
	fmt.Printf("RQL: %s\n", iq.GetQuery())
	params := iq.GetQueryParameters()
	if len(params) == 0 {
		return
	}
	fmt.Printf("Parameters:\n")
	var keys []string
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, key := range keys {
		fmt.Printf("  $%s: %#v\n", key, params[key])
	}
	fmt.Print("\n")
}

func queryGetConsumables() {
	store, session, err := openSession(dbName)
	if err != nil {
		log.Fatalf("openSession() failed with %s\n", err)
	}
	defer store.Close()
	defer session.Close()

	q := session.QueryCollection("Consumables")
	printRQL(q)

	var results []*reservation.Consumable
	err = q.GetResults(&results)
	if err != nil {
		log.Fatalf("q.GetResults() failed with '%s'\n", err)
	}
	if len(results) > 0 {
		fmt.Print("First result:\n")
		pretty.Print(results[0])
	}
}

// query collection of a given name
func queryGetRoomsAndDesks() {
	store, session, err := openSession(dbName)
	if err != nil {
		log.Fatalf("openSession() failed with %s\n", err)
	}
	defer store.Close()
	defer session.Close()

	q := session.QueryCollection("Consumables")
	printRQL(q)

	var results []*reservation.Consumable
	err = q.GetResults(&results)
	if err != nil {
		log.Fatalf("q.GetResults() failed with '%s'\n", err)
	}
	if len(results) > 0 {
		fmt.Print("First result:\n")
		pretty.Print(results[0])
	}
}

func crudDeleteReservation(reservationID string) {
	store, err := getDocumentStore(dbName)
	if err != nil {
		log.Fatalf("getDocumentStore() failed with %s\n", err)
	}
	defer store.Close()

	// delete the entity
	{
		session, err := store.OpenSession("")
		if err != nil {
			log.Fatalf("store.OpenSession() failed with %s\n", err)
		}

		var p *reservation.Reservations
		err = session.Load(&p, reservationID)
		if err != nil {
			log.Fatalf("session.Load() failed with %s\n", err)
		}

		err = session.Delete(p)
		if err != nil {
			log.Fatalf("session.Delete() failed with %s\n", err)
		}

		err = session.SaveChanges()
		if err != nil {
			log.Fatalf("session.SaveChanges() failed with %s\n", err)
		}

		session.Close()
	}

	// verify entity was deleted
	{
		session, err := store.OpenSession("")
		if err != nil {
			log.Fatalf("store.OpenSession() failed with %s\n", err)
		}

		// try to load deleted entity from the server
		var p *reservation.Reservations
		err = session.Load(&p, reservationID)
		if err != nil {
			log.Fatalf("session.Load() failed with %s\n", err)
		}

		if p == nil {
			fmt.Printf("Success: we deleted Reservation with this id so we get nil when we try to load it\n")
		} else {
			fmt.Printf("Error: this entity was deleted so we shouldn't be able to load it\n")
		}

		session.Close()
	}
}

func crudDeleteConsumable(consumableID string) {
	store, err := getDocumentStore(dbName)
	if err != nil {
		log.Fatalf("getDocumentStore() failed with %s\n", err)
	}
	defer store.Close()

	// delete the entity
	{
		session, err := store.OpenSession("")
		if err != nil {
			log.Fatalf("store.OpenSession() failed with %s\n", err)
		}

		var p *reservation.Consumable
		err = session.Load(&p, consumableID)
		if err != nil {
			log.Fatalf("session.Load() failed with %s\n", err)
		}

		err = session.Delete(p)
		if err != nil {
			log.Fatalf("session.Delete() failed with %s\n", err)
		}

		err = session.SaveChanges()
		if err != nil {
			log.Fatalf("session.SaveChanges() failed with %s\n", err)
		}

		session.Close()
	}

	// verify entity was deleted
	{
		session, err := store.OpenSession("")
		if err != nil {
			log.Fatalf("store.OpenSession() failed with %s\n", err)
		}

		// try to load deleted entity from the server
		var p *reservation.Consumable
		err = session.Load(&p, consumableID)
		if err != nil {
			log.Fatalf("session.Load() failed with %s\n", err)
		}

		if p == nil {
			fmt.Printf("Success: we deleted Consumable with this id so we get nil when we try to load it\n")
		} else {
			fmt.Printf("Error: this entity was deleted so we shouldn't be able to load it\n")
		}

		session.Close()
	}
}

func crudAddConsumable(Name string, Type string, Stock string, Available string) {
	store, session, err := openSession(dbName)
	if err != nil {
		log.Fatalf("openSession() failed with %s\n", err)
	}
	defer store.Close()
	defer session.Close()

	consumable := &reservation.Consumable{
		Name:      Name,
		Type:      Type,
		Stock:     Stock,
		Available: Available,
	}
	err = session.Store(consumable)
	if err != nil {
		log.Fatalf("session.Store() failed with %s\n", err)
	}
	fmt.Printf("Consumable ID: %s\n", consumable.ID)
	err = session.SaveChanges()
	if err != nil {
		log.Fatalf("session.SaveChanges() failed with %s\n", err)
	}
}
