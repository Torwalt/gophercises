// ideally this would be in its own package but for brevity its together with the urlshortener
package urlshortener

import (
	"encoding/binary"
	"encoding/json"
	"fmt"

	"github.com/boltdb/bolt"
)

func PopulateBoltDB(db *bolt.DB, paths []PathURL, bName string) error {
	for _, p := range paths {
		err := db.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(bName))

			bId, _ := b.NextSequence()
			id := int(bId)

			buf, err := json.Marshal(p)
			if err != nil {
				return err
			}
			err = b.Put(itob(id), buf)
			if err != nil {
				return fmt.Errorf("put: %s", err)
			}
			return nil
		})
		if err != nil {
			return fmt.Errorf("populate boltdb: %s", err)
		}
	}
	return nil
}

func InitBoltDB(db *bolt.DB, bName string) error {
	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bName))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})
	return nil
}

// itob returns an 8-byte big endian representation of v.
func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}
