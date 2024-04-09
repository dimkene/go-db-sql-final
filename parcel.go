package main

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

type ParcelStore struct {
	db *sql.DB
}

func NewParcelStore(db *sql.DB) ParcelStore {
	return ParcelStore{db: db}
}

func (s ParcelStore) Add(p Parcel) (int, error) {
	query := "INSERT INTO parcel (client, status, address, created_at) VALUES (?, ?, ?, ?)"
	output, err := s.db.Exec(query, p.Client, p.Status, p.Address, p.CreatedAt)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	id, err := output.LastInsertId()
	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	return int(id), nil
}

func (s ParcelStore) Get(number int) (Parcel, error) {
	query := "SELECT * FROM parcel WHERE number = ?"
	output := s.db.QueryRow(query, number)

	p := Parcel{}
	err := output.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
	if err != nil {
		fmt.Println(err)
		return Parcel{}, err
	}

	return p, nil
}

func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {
	query := "SELECT * FROM parcel WHERE client = ?"
	output, err := s.db.Query(query, client)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer output.Close()

	var result []Parcel
	for output.Next() {
		var p Parcel
		err := output.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		result = append(result, p)
	}

	err = output.Err()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return result, nil
}

func (s ParcelStore) SetStatus(number int, status string) error {

	query := "UPDATE parcel SET status = ? WHERE number = ?"
	_, err := s.db.Exec(query, status, number)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func (s ParcelStore) SetAddress(number int, address string) error {

	query := "UPDATE parcel SET address = ? WHERE number = ? AND status = ?"
	_, err := s.db.Exec(query, address, number, ParcelStatusRegistered)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func (s ParcelStore) Delete(number int) error {

	query := "DELETE FROM parcel WHERE number = ? AND status = ?"
	_, err := s.db.Exec(query, number, ParcelStatusRegistered)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}
