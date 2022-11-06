package pebble

import (
	"fmt"
	"strings"

	"github.com/cockroachdb/pebble"
)

type PebbleCursor struct {
	//db       *pebble.DB
	//bucket   *PebbleBucket
	bucketPrefix string // prefix to remove from keys returned by get/set/seek/first/lasst
	bucketID     string
	clientID     string
	iterator     *pebble.Iterator
}

// Close the cursor
// This does not release the transaction that created the bucket
func (cursor *PebbleCursor) Close() error {
	return nil
}

// First moves the cursor to the first item
func (cursor *PebbleCursor) First() (key string, value []byte) {
	isValid := cursor.iterator.First()
	_ = isValid
	return cursor.getKV()
}

// Return the iterator current key and value
// This removes the bucket prefix
func (cursor *PebbleCursor) getKV() (key string, value []byte) {
	k := string(cursor.iterator.Key())
	v, err := cursor.iterator.ValueAndErr()
	if strings.HasPrefix(k, cursor.bucketPrefix) {
		key = k[len(cursor.bucketPrefix):]
	} else {
		err = fmt.Errorf("bucket key '%s' has no prefix '%s'", k, cursor.bucketPrefix)
	}
	// what to do in case of error?
	_ = err
	return key, v
}

// Last moves the cursor to the last item
func (cursor *PebbleCursor) Last() (key string, value []byte) {
	cursor.iterator.Last()
	return cursor.getKV()
}

// Next iterates to the next key from the current cursor
func (cursor *PebbleCursor) Next() (key string, value []byte) {
	cursor.iterator.Next()
	return cursor.getKV()
}

// Prev iterations to the previous key from the current cursor
func (cursor *PebbleCursor) Prev() (key string, value []byte) {
	cursor.iterator.Prev()
	return cursor.getKV()
}

// Seek returns a cursor with Next() and Prev() iterators
func (cursor *PebbleCursor) Seek(searchKey string) (key string, value []byte) {
	bucketKey := cursor.bucketPrefix + searchKey
	cursor.iterator.SeekGE([]byte(bucketKey))
	return cursor.getKV()
}

func NewPebbleCursor(clientID, bucketID string, bucketPrefix string, iterator *pebble.Iterator) *PebbleCursor {
	cursor := &PebbleCursor{
		bucketPrefix: bucketPrefix,
		bucketID:     bucketID,
		clientID:     clientID,
		iterator:     iterator,
	}
	return cursor
}
