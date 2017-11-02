package api

import (
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/mgo.v2"
)

// Session Make the mongoDB session a global variable
var Session *mgo.Session

// Storage interface to save/get webhooks and data from fixer.
type Storage interface {
	Init()
	Add(w Webhook) (string, error)
	AddCurrency(f Fixer) error
	Count() int
	Get(key string) (Webhook, bool)
	GetLatest(int) ([]Fixer, error)
	GetAll() []Webhook
	Remove(key string) bool
	GetDbURL() string
	GetDbName() string
}

// MongoDB struct for the mongoDB storage
type MongoDB struct {
	DatabaseURL string
	DatabaseName string
	WebhooksCollectionName string
	ExchangeCollectionName string
}

// GetDbURL get tha URL for database
func (db *MongoDB) GetDbURL() string {
	return db.DatabaseURL
}

// GetDbName get the name for database
func (db *MongoDB) GetDbName() string {
	return db.DatabaseName
}

// Init initializes the dabase
func (db *MongoDB) Init() {
	index := mgo.Index{
		Key: []string{"date"},
		Unique: true,
		DropDups: true,
		Background: true,
		Sparse: true,
	}

	err := Session.DB(db.DatabaseName).C(db.ExchangeCollectionName).EnsureIndex(index)
	if err != nil {
		panic(err)
	}
}

// Add a new webhook to database
func (db *MongoDB) Add(w Webhook) (string, error) {
	w.ID = bson.NewObjectId()

	err := w.Validate()
	if err != nil {
		return "", err
	}

	err = Session.DB(db.DatabaseName).C(db.WebhooksCollectionName).Insert(w)
	if err != nil {
		return "", err
	}

	return w.ID.Hex(), nil
}

// AddCurrency to the databse
func (db *MongoDB) AddCurrency(f Fixer) error {
	err := Session.DB(db.DatabaseName).C(db.ExchangeCollectionName).Insert(f)
	if err != nil {
		return err
	}

	return nil
}

// Count get number of webhooks in database
func (db *MongoDB) Count() int {
	return 0
}

// Get webhook by ID
func (db *MongoDB) Get(key string) (Webhook, bool) {
	webhook := Webhook{}
	ok := true

	if bson.IsObjectIdHex(key) == false {
		return webhook, false
	}

	id := bson.ObjectIdHex(key)

	err := Session.DB(db.DatabaseName).C(db.WebhooksCollectionName).
		Find(bson.M{"_id": id}).One(&webhook)
	if err != nil {
		ok = false
	}

	return webhook, ok
}

// GetLatest latest rates from db
func (db *MongoDB) GetLatest(days int) ([]Fixer, error) {
	fixers := []Fixer{}

	err := Session.DB(db.DatabaseName).C(db.ExchangeCollectionName).Find(bson.M{}).Sort("date", "1").Limit(days).All(&fixers)
	if err != nil {
		return fixers, err
	}

	return fixers, nil
}

// GetAll get all webhooks from db
func (db *MongoDB) GetAll() []Webhook {
	all := []Webhook{}

	err := Session.DB(db.DatabaseName).C(db.WebhooksCollectionName).Find(bson.M{}).All(&all)
	if err != nil {
		return []Webhook{}
	}

	return all
}

// Remove a webhook by ID from db
func (db *MongoDB) Remove(key string) bool {
	ok := true

	if bson.IsObjectIdHex(key) == false {
		return false
	}

	id := bson.ObjectIdHex(key)

	err := Session.DB(db.DatabaseName).C(db.WebhooksCollectionName).
		Remove(bson.M{"_id": id})
	if err != nil {
		ok = false
	}

	return ok
}
