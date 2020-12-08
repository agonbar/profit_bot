package main

import (
	"fmt"
	"strconv"
	"time"

	bolt "go.etcd.io/bbolt"
)

func initDB() *bolt.DB {
	// INIT DATABASE
	db, _ := bolt.Open("bot.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	// Users bucket
	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("users"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})
	// Nodes bucket
	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("nodes"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})

	return db
}

func saveUser(db *bolt.DB, ID int, value string) {
	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("users"))
		err := b.Put([]byte(strconv.Itoa(ID)), []byte(value))
		return err
	})
}

func getUsers(db *bolt.DB) [][2]string {
	var users [][2]string
	db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte("users"))

		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			users = append(users, [2]string{string(k), string(v)})
			fmt.Printf("key=%s, value=%s\n", k, v)
		}

		return nil
	})
	return users
}

func saveNode(db *bolt.DB, name string, url string) {
	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("nodes"))
		err := b.Put([]byte(name), []byte(url))
		return err
	})
}

func deleteNode(db *bolt.DB, name string) {
	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("nodes"))
		err := b.Delete([]byte(name))
		return err
	})
}

func getNodes(db *bolt.DB) [][2]string {
	var nodes [][2]string
	db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte("nodes"))

		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			nodes = append(nodes, [2]string{string(k), string(v)})
			fmt.Printf("key=%s, value=%s\n", k, v)
		}

		return nil
	})
	return nodes
}
