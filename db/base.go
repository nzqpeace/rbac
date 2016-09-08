package db

import (
	"gopkg.in/mgo.v2"

	"gopkg.in/mgo.v2/bson"
)

type (
	Base struct {
		db         *DataBase
		collection string
	}
)

func NewBase(db *DataBase, collection string) *Base {
	return &Base{
		db:         db,
		collection: collection,
	}
}

func (m *Base) Invoke(fn func(col *mgo.Collection) error) error {

	s := m.db.Copy()
	defer s.Close()
	err := fn(s.C(m.collection))
	return err
}

func (m *Base) Insert(model interface{}) error {
	return m.Invoke(func(col *mgo.Collection) error {
		return col.Insert(model)
	})
}

func (m *Base) Upsert(query, change interface{}) error {
	return m.Invoke(func(col *mgo.Collection) (err error) {
		_, err = col.Upsert(query, change)
		return
	})
}

func (m *Base) Update(query, change interface{}) error {
	return m.Invoke(func(col *mgo.Collection) error {
		return col.Update(query, change)
	})
}

func (m *Base) UpdateAll(query, change interface{}) error {
	return m.Invoke(func(col *mgo.Collection) error {
		_, err := col.UpdateAll(query, change)
		return err
	})
}

func (m *Base) Find(query bson.M, model interface{}) error {
	return m.Invoke(func(col *mgo.Collection) error {
		return col.Find(query).One(model)
	})
}

func (m *Base) FindAll(query bson.M, models interface{}, skip, limit int, sorts ...string) error {
	return m.Invoke(func(col *mgo.Collection) error {
		return col.Find(query).Skip(skip).Limit(limit).Sort(sorts...).All(models)
	})
}

func (m *Base) Distinct(query bson.M, models interface{}, key string) error {
	return m.Invoke(func(col *mgo.Collection) error {
		return col.Find(query).Distinct(key, models)
	})
}

// atomic update object and return old object
func (m *Base) FindAndModify(query, change bson.M, model interface{}) error {
	return m.Invoke(func(col *mgo.Collection) error {
		_, err := col.Find(query).Apply(mgo.Change{
			Update: change,
		}, model)
		return err
	})
}

func (m *Base) Remove(query bson.M) error {
	return m.Invoke(func(col *mgo.Collection) error {
		return col.Remove(query)
	})
}

func (m *Base) RemoveAll(query bson.M) error {
	return m.Invoke(func(col *mgo.Collection) error {
		_, err := col.RemoveAll(query)
		return err
	})
}

func (m *Base) Count(query bson.M) (n int, err error) {
	err = m.Invoke(func(col *mgo.Collection) error {
		n, err = col.Find(query).Count()
		return err
	})

	return
}
