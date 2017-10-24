package main

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type WebhooksStorage interface {
	Init()
	Add(w Webhook) (string, error)
	Count() int
	Get(key string) (Webhook, bool)
	Remove(key string) bool
}

type Webhook struct {
	ID bson.ObjectId `json:"_id" bson:"_id"`
	WebhookURL string
	BaseCurrency string
	TargetCurrency string
	MinTriggerValue int
	MaxTriggerValue int
}

type WebhooksMongoDB struct {
	DatabaseURL string
	DatabaseName string
	WebhooksCollectionName string
}

func (db *WebhooksMongoDB) Init() {
	session, err := mgo.Dial(db.DatabaseURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()
}

func (db *WebhooksMongoDB) Add(w Webhook) (string, error) {
	session, err := mgo.Dial(db.DatabaseURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	w.ID = bson.NewObjectId()
	err = session.DB(db.DatabaseName).C(db.WebhooksCollectionName).Insert(w)
	if err != nil {
		return "", err
	}

	return w.ID.Hex(), nil
}

func (db *WebhooksMongoDB) Count() int {
	return 0
}

func (db *WebhooksMongoDB) Get(key string) (Webhook, bool) {
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

func (db *WebhooksMongoDB) Remove(key string) bool {
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