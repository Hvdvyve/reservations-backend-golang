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
	"github.com/gitpod/mycli/cmd"
	ravendb "github.com/ravendb/ravendb-go-client"
)

func main() {
	cmd.Execute()
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
