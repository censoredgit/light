package session

import "time"

type Data struct {
	id         string
	data       map[string]string
	expireTime time.Time
	isModified bool
	isNew      bool
}

func NewData(id string, ttl time.Duration) *Data {
	return &Data{
		id:         id,
		data:       make(map[string]string),
		expireTime: time.Now().Add(ttl),
		isNew:      true,
		isModified: true,
	}
}

func NewExistsData(id string, ttl time.Duration) *Data {
	return &Data{
		id:         id,
		data:       make(map[string]string),
		expireTime: time.Now().Add(ttl),
		isNew:      false,
	}
}

func (d *Data) IsNew() bool {
	return d.isNew
}

func (d *Data) Id() string {
	return d.id
}

func (d *Data) IsModified() bool {
	return d.isModified
}

func (d *Data) Get(key string) string {
	return d.data[key]
}

func (d *Data) Has(key string) bool {
	_, ok := d.data[key]
	return ok
}

func (d *Data) Set(key, value string) {
	d.data[key] = value
	d.isModified = true
}

func (d *Data) Delete(key string) {
	delete(d.data, key)
	d.isModified = true
}

func (d *Data) Expire() time.Time {
	return d.expireTime
}

func (d *Data) isExpired() bool {
	return d.expireTime.Before(time.Now())
}

func (d *Data) All() map[string]string {
	return d.data
}

func (d *Data) IsEmpty() bool {
	return len(d.data) == 0
}
func (d *Data) Empty() {
	clear(d.data)
}

func (d *Data) OnSaved() {
	d.isModified = false
}
