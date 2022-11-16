package bucketstore_test

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hiveot/hub.go/pkg/logging"
	"github.com/hiveot/hub.go/pkg/thing"
	"github.com/hiveot/hub.go/pkg/vocab"
	"github.com/hiveot/hub/pkg/bucketstore"
	"github.com/hiveot/hub/pkg/bucketstore/cmd"
)

var testBucketID = "default"

var testBackendType = bucketstore.BackendKVBTree

// var testBackendType = bucketstore.BackendBBolt
// var testBackendType = bucketstore.BackendPebble
var testBackendDirectory = "/tmp/test-bucketstore"
var testClientID = "buckettestclient"

const (
	doc1ID = "doc1"
	doc2ID = "doc2"
)

var doc1 = []byte(`{
  "id": "doc1",
  "title": "Title of doc 1",
  "@type": "sensor",
  "properties": 
     { "title": {
         "name": "title1" 
       }
     }
}`)
var doc2 = []byte(`{
  "id": "doc2",
  "title": "Title of doc 2",
  "properties": [
     { "title": "title2" }
  ]
}`)

// Create the bucket store using the backend
func openNewStore() (store bucketstore.IBucketStore, err error) {
	_ = os.RemoveAll(testBackendDirectory)
	store = cmd.NewBucketStore(testBackendDirectory, testClientID, testBackendType)
	err = store.Open()
	return store, err
}

// Create a TD document
func createTD(id string) *thing.ThingDescription {
	td := &thing.ThingDescription{
		ID:         id,
		Title:      fmt.Sprintf("test TD %s", id),
		AtType:     string(vocab.DeviceTypeSensor),
		Properties: make(map[string]*thing.PropertyAffordance),
		Events:     make(map[string]*thing.EventAffordance),
	}
	td.Properties[vocab.PropNameTitle] = &thing.PropertyAffordance{
		DataSchema: thing.DataSchema{
			Title:       "Sensor title",
			Description: "This is a smart sensor",
			Type:        vocab.WoTDataTypeString,
			Default:     "Default value",
		},
	}
	td.Properties[vocab.PropNameSoftwareVersion] = &thing.PropertyAffordance{
		DataSchema: thing.DataSchema{
			Title:       "Version",
			Description: "Embedded firmware",
			Type:        vocab.WoTDataTypeString,
			Default:     "Default value",
			Const:       "v1.0",
		},
	}
	td.Events[vocab.PropNameValue] = &thing.EventAffordance{
		Title:       "Event 1",
		Description: "Name of this event",
		Data: thing.DataSchema{
			Type:        vocab.WoTDataTypeString,
			Const:       "123",
			Title:       "Event name data",
			Description: "String with friendly name of the event"},
	}
	td.Events[vocab.PropNameBattery] = &thing.EventAffordance{
		Title: "Event 2",
		Data: thing.DataSchema{
			Type:        vocab.WoTDataTypeInteger,
			Title:       "Battery level",
			Unit:        vocab.UnitNamePercent,
			Description: "Battery level update in % of device"},
	}
	return td
}

// AddDocs adds documents doc1, doc2 and given nr additional docs
func addDocs(store bucketstore.IBucketStore, bucketID string, count int) error {
	const batchSize = 50000
	bucket := store.GetBucket(bucketID)

	// these docs have values used for testing
	err := bucket.Set(doc1ID, doc1)
	err = bucket.Set(doc2ID, doc2)
	if err != nil {
		return err
	}

	// breakup in batches to limit the transaction size
	// fill remainder with generated docs
	// don't sort order of id
	iBatch := 0
	docs := make(map[string][]byte)
	for i := count; i > 2; i-- {
		rn := rand.Intn(count * 33) // enough spread to avoid duplicates
		id := fmt.Sprintf("addDocs-%6d", rn)
		td := createTD(id)
		_ = td
		jsonDoc := []byte("hello world")
		jsonDoc, _ = json.Marshal(td) // 900msec
		docs[id] = jsonDoc
		// restart the batch
		iBatch++
		// close the bucket/transaction and reopen
		if iBatch == batchSize {
			err = bucket.SetMultiple(docs)
			err = bucket.Close()
			if err != nil {
				panic(fmt.Sprintf("SetMultiple failed: %s", err))
			}
			// next batch
			docs = make(map[string][]byte)
			bucket = store.GetBucket(bucketID)
			iBatch = 0
		}

	}
	// finish the remainder
	_ = bucket.SetMultiple(docs)
	err = bucket.Close()
	logrus.Infof("Added '%d' records to the store", count)
	return err
}

func TestMain(m *testing.M) {
	logging.SetLogging("info", "")

	res := m.Run()
	os.Exit(res)
}

// Generic directory store testcases
func TestStartStop(t *testing.T) {
	store, err := openNewStore()
	require.NoError(t, err)
	err = store.Close()
	assert.NoError(t, err)
}

func TestCreateStoreBadFolder(t *testing.T) {
	badDir := "/folder/does/not/exist/"
	store := cmd.NewBucketStore(badDir, testClientID, testBackendType)
	err := store.Open()
	assert.Error(t, err)
}

func TestCreateStoreReadOnlyFolder(t *testing.T) {
	badDir := "/var/"
	store := cmd.NewBucketStore(badDir, testClientID, testBackendType)
	err := store.Open()
	assert.Error(t, err)
}

func TestCreateStoreCantReadFile(t *testing.T) {
	badDir := "/bin"
	store := cmd.NewBucketStore(badDir, "yes", testBackendType)
	err := store.Open()
	assert.Error(t, err)
}

func TestWriteRead(t *testing.T) {
	const id1 = "id1"
	const id5 = "id5"
	const id22 = "id22"

	store, err := openNewStore()
	assert.NoError(t, err)
	err = addDocs(store, testBucketID, 3)

	bucket := store.GetBucket(testBucketID)
	assert.NotNil(t, bucket)

	require.NoError(t, err)

	// write docs
	td1 := createTD(id1)
	td1json, _ := json.Marshal(td1)
	err = bucket.Set(id1, td1json)
	assert.NoError(t, err)
	td22 := createTD(id22)
	td22json, _ := json.Marshal(td22)
	err = bucket.Set(id22, td22json)
	assert.NoError(t, err)
	td5 := createTD(id5)
	td5json, _ := json.Marshal(td5)
	err = bucket.Set(id5, td5json)
	assert.NoError(t, err)

	// kvstore flushes to file in autosave loop every 3 seconds
	// needs to be tested
	time.Sleep(time.Second * 4)

	err = bucket.Close()
	assert.NoError(t, err)
	err = store.Close()
	assert.NoError(t, err)
	time.Sleep(time.Second)

	// --- reopen ---
	err = store.Open() // reopen
	require.NoError(t, err)
	bucket = store.GetBucket(testBucketID)
	assert.NotNil(t, bucket)

	// Read and compare
	resp, err := bucket.Get(id22)
	if assert.NotNil(t, resp) {
		assert.Equal(t, td22json, resp)
	}
	resp, err = bucket.Get(id1)
	if assert.NotNil(t, resp) {
		assert.Equal(t, td1json, resp)
	}
	resp, err = bucket.Get(id5)
	if assert.NotNil(t, resp) {
		assert.Equal(t, td5json, resp)
	}
	// Delete
	err = bucket.Delete(id1)
	assert.NoError(t, err)
	resp, err = bucket.Get(id1)
	assert.NoError(t, err)
	assert.Nil(t, resp)

	err = bucket.Close()
	assert.NoError(t, err)
	err = store.Close()
	assert.NoError(t, err)

	// Read again should fail
	// (pebble throws a panic :(
	//_, err = store.Get(testBucketID, doc1ID)
	//assert.Error(t, err)
}

func TestWriteBadData(t *testing.T) {
	store, err := openNewStore()
	require.NoError(t, err)
	defer store.Close()
	bucket := store.GetBucket(testBucketID)
	defer bucket.Close()
	// not json
	err = bucket.Set(doc1ID, []byte("not-json"))
	assert.NoError(t, err)
	// missing key
	err = bucket.Set("", []byte("{}"))
	assert.Error(t, err)

}

func TestWriteReadMultiple(t *testing.T) {
	const id1 = "id1"
	const id5 = "id5"
	const id22 = "id22"
	docs := make(map[string][]byte)

	store, err := openNewStore()
	require.NoError(t, err)
	err = addDocs(store, testBucketID, 3)
	require.NoError(t, err)

	bucket := store.GetBucket(testBucketID)
	assert.NotNil(t, bucket)
	defer store.Close()
	defer bucket.Close() // last defer completes first

	// write docs
	docs[id1], _ = json.Marshal(createTD(id1))
	docs[id22], _ = json.Marshal(createTD(id22))
	docs[id5], _ = json.Marshal(createTD(id5))
	err = bucket.SetMultiple(docs)
	assert.NoError(t, err)

	// Read and compare

	resp, err := bucket.GetMultiple([]string{id22, id1, id5})
	assert.NoError(t, err)
	assert.Equal(t, docs[id1], resp[id1])
	assert.Equal(t, docs[id5], resp[id5])
	assert.Equal(t, docs[id22], resp[id22])

	// Delete
	err = bucket.Delete(id1)
	assert.NoError(t, err)
	resp2, err := bucket.GetMultiple([]string{id22, id1, id5})
	assert.NoError(t, err)
	assert.Equal(t, 2, len(resp2))
}

func TestSeek(t *testing.T) {
	const count = 1000
	const seekCount = 200
	const base = 500

	store, err := openNewStore()
	require.NoError(t, err)
	err = addDocs(store, testBucketID, count)
	require.NoError(t, err)
	bucket := store.GetBucket(testBucketID)
	require.NotNil(t, bucket)
	defer store.Close()
	defer bucket.Close()

	// set cursor 'base' records forward
	cursor := bucket.Cursor()
	k1, v1, valid := cursor.First()
	assert.True(t, valid)
	for i := 0; i < base; i++ {
		k1, v1, valid = cursor.Next()
		assert.True(t, valid)
		assert.NotEmpty(t, k1)
		assert.NotEmpty(t, v1)
	}
	// k1 now holds the key at the base N'th record

	// seek of the current key should bring us back here, at the base Nth record
	k2, v2, valid2 := cursor.Seek(k1)
	assert.True(t, valid2)
	assert.Equal(t, k1, k2)
	assert.Equal(t, v1, v2)

	// test that keys increment
	for i := 0; i < seekCount; i++ {
		k, v, valid := cursor.Next()
		require.GreaterOrEqual(t, k, k2)
		assert.True(t, valid)
		if !assert.NotEmpty(t, v) {
			logrus.Infof("unexpected")
		}
		k2 = k
	}

	// step seekCount nr backwards should lead us right back to k1
	k2, v2, valid2 = cursor.Prev()
	assert.True(t, valid2)
	for i := 1; i < seekCount; i++ {
		k, v, valid := cursor.Prev()
		assert.True(t, valid)
		require.LessOrEqual(t, k, k2)
		if !assert.NotEmpty(t, v) {
			logrus.Infof("unexpected")
		}
		k2 = k
	}
	// how to test Last?
	_, _, valid = cursor.Last()
	assert.True(t, valid)
	assert.Equal(t, k1, k2)
	cursor.Release()
}

func TestPrevNextN(t *testing.T) {
	const count = 1000
	const seekCount = 200
	const base = 500

	// setup
	store, err := openNewStore()
	require.NoError(t, err)
	err = addDocs(store, testBucketID, count)
	require.NoError(t, err)
	bucket := store.GetBucket(testBucketID)
	require.NotNil(t, bucket)
	defer store.Close()
	defer bucket.Close()

	// test NextN
	cursor := bucket.Cursor()
	k1, v1, valid := cursor.First()
	assert.True(t, valid)
	docs, itemsRemaining := cursor.NextN(seekCount)
	assert.True(t, itemsRemaining)
	assert.Equal(t, seekCount, len(docs))

	docs2, itemsRemaining := cursor.PrevN(seekCount - 1)
	assert.True(t, itemsRemaining)
	assert.Equal(t, seekCount-1, len(docs2))

	// one step further we're at the begin again
	k2, v2, valid2 := cursor.Prev()
	assert.True(t, valid2)
	assert.Equal(t, k1, k2)
	assert.Equal(t, v1, v2)

	cursor.Release()
}

//func TestQuery(t *testing.T) {
//	store, err := createNewStore()
//	require.NoError(t, err)
//	err = addDocs(store, 20)
//	require.NoError(t, err)
//
//	// filter on key 'id' == doc1
//	//args := &svc.Query_Args{JsonPathQuery: `$[?(@.id=="doc1")]`}
//	jsonPath := `$[?(@.id=="doc1")]`
//	resp, err := store.Query(jsonPath, 0, 0, nil)
//	require.NoError(t, err)
//	assert.NotEmpty(t, resp)
//
//	// regular nested filter comparison. note that a TD does not hold values
//	jsonPath = `$[?(@.properties.title.name=="title1")]`
//	resp, err = store.Query(jsonPath, 0, 0, nil)
//	require.NoError(t, err)
//	assert.NotEmpty(t, resp)
//
//	// filter with nested notation. some examples that return a list of TDs matching the filter
//	//res, err = fileStore.Query(`$[?(@.properties.title.value=="title1")]`, 0, 0)
//	// res, err = fileStore.Query(`$[?(@.*.title.value=="title1")]`, 0, 0)
//	// res, err = fileStore.Query(`$[?(@['properties']['title']['value']=="title1")]`, 0, 0)
//	jsonPath = `$[?(@..title.name=="title1")]`
//	resp, err = store.Query(jsonPath, 0, 0, nil)
//	assert.NoError(t, err)
//
//	// these only return the properties - not good
//	// res, err = fileStore.Query(`$.*.properties[?(@.value=="title1")]`, 0, 0) // returns list of props, not tds
//	//res, err = fileStore.Query(`$.*.*[?(@.value=="title1")]`, 0, 0) // returns list of props, not tds
//	// res, err = fileStore.Query(`$[?(@...value=="title1")]`, 0, 0)
//	assert.NotEmpty(t, resp)
//
//	// filter with bracket notation
//	jsonPath = `$[?(@["id"]=="doc1")]`
//	resp, err = store.Query(jsonPath, 0, 0, nil)
//	require.NoError(t, err)
//	assert.NotEmpty(t, resp)
//
//	// filter with bracket notation and current object literal (for search @type)
//	// only supported by: ohler55/ojg
//	jsonPath = `$[?(@['@type']=="sensor")]`
//	resp, err = store.Query(jsonPath, 0, 0, nil)
//	assert.NoError(t, err)
//	assert.Greater(t, len(resp), 1)
//
//	// bad query expression
//	jsonPath = `$[?(.id=="doc1")]`
//	resp, err = store.Query(jsonPath, 0, 0, nil)
//	assert.Error(t, err)
//}

// tests to figure out how to use jp parse with bracket notation
//func TestQueryBracketNotationA(t *testing.T) {
//	store := make(map[string]interface{})
//	query1 := `$[?(@['type']=="type1")]`
//	query2 := `$[?(@['@type']=="sensor")]`
//
//	jsonDoc := `{
//		"thing1": {
//			"id": "thing1",
//			"type": "type1",
//			"@type": "sensor",
//			"properties": {
//				"title": "title1"
//			}
//		},
//		"thing2": {
//			"id": "thing2",
//			"type": "type2",
//			"@type": "sensor",
//			"properties": {
//				"title": "title1"
//			}
//		}
//	}`
//
//	err := json.Unmarshal([]byte(jsonDoc), &store)
//	assert.NoError(t, err)
//
//	jpExpr, err := jp.ParseString(query1)
//	assert.NoError(t, err)
//	result := jpExpr.Get(store)
//	assert.NotEmpty(t, result)
//
//	jpExpr, err = jp.ParseString(query2)
//	assert.NoError(t, err)
//	result = jpExpr.Get(store)
//	assert.NotEmpty(t, result)
//}

// tests to figure out how to use jp parse with bracket notation
//func TestQueryBracketNotationB(t *testing.T) {
//	//store := make(map[string]interface{})
//	queryString := "$[?(@['@type']==\"sensor\")]"
//	id1 := "thing1"
//	id2 := "thing2"
//	td1 := thing.ThingDescription{
//		ID:         id1,
//		Title:      "test TD 1",
//		AtType:     string(vocab.DeviceTypeSensor),
//		Properties: make(map[string]*thing.PropertyAffordance),
//	}
//	//td1 := thing.CreateTD(id1, "test TD", vocab.DeviceTypeSensor)
//	td1.Properties[vocab.PropNameTitle] = &thing.PropertyAffordance{
//		DataSchema: thing.DataSchema{
//			Title: "Sensor title",
//			Type:  vocab.WoTDataTypeString,
//		},
//	}
//	td1.Properties[vocab.PropNameValue] = &thing.PropertyAffordance{
//		DataSchema: thing.DataSchema{
//			Title: "Sensor value",
//			Type:  vocab.WoTDataTypeNumber,
//		},
//	}
//
//	td2 := thing.ThingDescription{
//		ID:         id2,
//		Title:      "test TD 2",
//		AtType:     string(vocab.DeviceTypeSensor),
//		Properties: make(map[string]*thing.PropertyAffordance),
//	}
//	td2.Properties[vocab.PropNameTitle] = &thing.PropertyAffordance{
//		DataSchema: thing.DataSchema{
//			Title: "The switch",
//			Type:  vocab.WoTDataTypeBool,
//		},
//	}
//
//	store, err := createNewStore()
//	require.NoError(t, err)
//
//	//td1json, err := json.MarshalIndent(td1, "", "")
//	td1json, err := json.Marshal(&td1)
//	td2json, err := json.Marshal(&td2)
//	_ = store.Write(id1, string(td1json))
//	err = store.Write(id2, string(td2json))
//	assert.NoError(t, err)
//
//	// query returns 2 sensors.
//	resp, err := store.Query(queryString, 0, 0, nil)
//	require.NoError(t, err)
//	require.Equal(t, 2, len(resp))
//
//	var readTD1 thing.ThingDescription
//	err = json.Unmarshal([]byte(resp[0]), &readTD1)
//	require.NoError(t, err)
//	read1type := readTD1.AtType
//	assert.Equal(t, string(vocab.DeviceTypeSensor), read1type)
//}

// test query with reduced list of IDs
//func TestQueryFiltered(t *testing.T) {
//	queryString := "$..id"
//
//	store, err := createNewStore()
//	require.NoError(t, err)
//	_ = addDocs(store, 10)
//
//	// result of a normal query
//	resp, err := store.Query(queryString, 0, 0, nil)
//	require.NoError(t, err)
//	assert.Equal(t, 10, len(resp))
//}
