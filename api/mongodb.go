package api

import (
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/mgo.v2"
)

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
}

// MongoDB struct for the mongoDB storage
type MongoDB struct {
	DatabaseURL string
	DatabaseName string
	WebhooksCollectionName string
	ExchangeCollectionName string
}

// Init initializes the dabase
func (db *MongoDB) Init() {
	session, err := mgo.Dial(db.DatabaseURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	index := mgo.Index{
		Key: []string{"date"},
		Unique: true,
		DropDups: true,
		Background: true,
		Sparse: true,
	}

	err = session.DB(db.DatabaseName).C(db.ExchangeCollectionName).EnsureIndex(index)
	if err != nil {
		panic(err)
	}
}

// Add a new webhook to database
func (db *MongoDB) Add(w Webhook) (string, error) {
	session, err := mgo.Dial(db.DatabaseURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	w.ID = bson.NewObjectId()

	err = w.Validate()
	if err != nil {
		return "", err
	}

	err = session.DB(db.DatabaseName).C(db.WebhooksCollectionName).Insert(w)
	if err != nil {
		return "", err
	}

	return w.ID.Hex(), nil
}

// AddCurrency to the databse
func (db *MongoDB) AddCurrency(f Fixer) error {
	session, err := mgo.Dial(db.DatabaseURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	err = session.DB(db.DatabaseName).C(db.ExchangeCollectionName).Insert(f)
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
	session, err := mgo.Dial(db.DatabaseURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	webhook := Webhook{}
	ok := true

	if bson.IsObjectIdHex(key) == false {
		return webhook, false
	}

	id := bson.ObjectIdHex(key)

	err = session.DB(db.DatabaseName).C(db.WebhooksCollectionName).
		Find(bson.M{"_id": id}).One(&webhook)
	if err != nil {
		ok = false
	}

	return webhook, ok
}

// GetLatest latest rates from db
func (db *MongoDB) GetLatest(days int) ([]Fixer, error) {
	session, err := mgo.Dial(db.DatabaseURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	fixers := []Fixer{}

	err = session.DB(db.DatabaseName).C(db.ExchangeCollectionName).Find(bson.M{}).Sort("date", "1").Limit(days).All(&fixers)
	if err != nil {
		return fixers, err
	}

	return fixers, nil
}

// GetAll get all webhooks from db
func (db *MongoDB) GetAll() []Webhook {
	session, err := mgo.Dial(db.DatabaseURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	all := []Webhook{}

	err = session.DB(db.DatabaseName).C(db.WebhooksCollectionName).Find(bson.M{}).All(&all)
	if err != nil {
		return []Webhook{}
	}

	return all
}

// Remove a webhook by ID from db
func (db *MongoDB) Remove(key string) bool {
	session, err := mgo.Dial(db.DatabaseURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	ok := true

	if bson.IsObjectIdHex(key) == false {
		return false
	}

	id := bson.ObjectIdHex(key)

	err = session.DB(db.DatabaseName).C(db.WebhooksCollectionName).
		Remove(bson.M{"_id": id})
	if err != nil {
		ok = false
	}

	return ok
}
