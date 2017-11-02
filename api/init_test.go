package api

import (
	"testing"
	"io/ioutil"
	"os"
	"gopkg.in/mgo.v2/dbtest"
)

// This was taken from https://medium.com/@mvmaasakkers/writing-integration-tests-with-mongodb-support-231580a566cd
// and modified for my use.

// This file is part of https://github.com/mvmaasakkers/gohttptestmongodb/
// Server holds the dbtest DBServer
var Server dbtest.DBServer

// TestMain wraps all tests with the needed initialized mock DB and fixtures
func TestMain(m *testing.M) {
	Db = &MongoDB{
		DatabaseName: "exchange_test",
		WebhooksCollectionName: "webhooks",
		ExchangeCollectionName: "currencyrates",
	}

	// The tempdir is created so MongoDB has a location to store its files.
	// Contents are wiped once the server stops
	tempDir, _ := ioutil.TempDir("", "testing")
	Server.SetPath(tempDir)

	// My main session var is now set to the temporary MongoDB instance
	Session = Server.Session()
	Db.Init()

	// Run the test suite
	retCode := m.Run()

	// Make sure we DropDatabase so we make absolutely sure nothing is left or locked while wiping the data and
	// close session
	Session.DB(Db.GetDbName()).DropDatabase()
	Session.Close()

	// Stop shuts down the temporary server and removes data on disk.
	defer Server.Stop()

	// call with result of m.Run()
	os.Exit(retCode)
}
