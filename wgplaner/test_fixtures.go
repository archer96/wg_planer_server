package wgplaner

import (
	"gopkg.in/testfixtures.v2"
)

var fixtures *testfixtures.Context

// InitFixtures initialize test fixtures for a test database
func InitFixtures(helper testfixtures.Helper, dir string) (err error) {
	testfixtures.SkipDatabaseNameCheck(true) // TODO
	fixtures, err = testfixtures.NewFolder(OrmEngine.DB().DB, helper, dir)
	return err
}

// LoadFixtures load fixtures for a test database
func LoadFixtures() error {
	return fixtures.Load()
}
