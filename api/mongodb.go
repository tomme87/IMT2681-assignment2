package api

import (
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/mgo.v2"
)

type Storage interface {
	Init()
	Add(w Webhook) (string, error)
	AddCurrency(f Fixer) error
	Count() int
	Get(key string) (Webhook, bool)
	Remove(key string) bool
}

type MongoDB struct {
	DatabaseURL string
	DatabaseName string
	WebhooksCollectionName string
	ExchangeCollectionName string
}

func (db *MongoDB) Init() {
	session, err := mgo.Dial(db.DatabaseURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()
}

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

func (db *MongoDB) Count() int {
	return 0
}

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
