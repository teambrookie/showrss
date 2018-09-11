package db

import (
	"os"
	"reflect"
	"testing"
)

var media = Media{
	ID:     "1234",
	Name:   "Wicked 3",
	Magnet: "Wicked 3 Magnetized",
	Seeds:  10,
	Leechs: 10,
}

func before(t *testing.T) *BoltMediaStore {
	db, err := Open("test.db")
	if err != nil {
		t.Fatalf("Unable to open the test db : %s", err)
	}

	return db
}

func after(db *BoltMediaStore) {
	db.Close()
	os.Remove("test.db")
}

func TestGetCollection(t *testing.T) {

	db := before(t)
	defer after(db)
	err := db.AddMedia(media, FOUND)
	if err != nil {
		t.Fatalf("Unable to add media: %s", err)
	}
	result, err := db.GetCollection(FOUND)
	if err != nil {
		t.Fatalf("Unable to get collection: %s", err)
	}
	expected := []Media{media}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected : %+v \n Result : %+v", expected, result)
	}
	result, err = db.GetCollection(NOTFOUND)
	if result != nil {
		t.Errorf("Expected : %+v \n Result : %+v", nil, result)
	}
}

func TestGetMedia(t *testing.T) {

	db := before(t)
	defer after(db)
	err := db.AddMedia(media, NOTFOUND)
	if err != nil {
		t.Fatalf("Unable to add media: %s", err)
	}
	// We try to get the media in the right collection
	result, err := db.GetMedia(media.ID, NOTFOUND)
	if err != nil {
		t.Fatalf("Unable to get media: %s", err)
	}
	if result != media {
		t.Errorf("Expected: %+v \n Result: %+v", media, result)
	}

	//We try to get the media in the wrong collection
	result, err = db.GetMedia(media.ID, FOUND)
	if err != nil {
		t.Fatalf("Unable to get media: %s", err)
	}
	if result != (Media{}) {
		t.Errorf("Expected: %+v \n Result: %+v", Media{}, result)
	}

	//we try to get the media with the wrong id
	//We try to get the media in the wrong collection
	result, err = db.GetMedia("4321", NOTFOUND)
	if err != nil {
		t.Fatalf("Unable to get media: %s", err)
	}
	if result != (Media{}) {
		t.Errorf("Expected: %+v \n Result: %+v", Media{}, result)
	}

}

func TestAddMedia(t *testing.T) {

	db := before(t)
	defer after(db)
	err := db.AddMedia(media, FOUND)
	if err != nil {
		t.Fatalf("Unable to add media: %s", err)
	}
	newMedia := media
	newMedia.Name = "YOLO"
	err = db.AddMedia(newMedia, FOUND)
	if err == nil {
		t.Fatalf("Two media with same ID can't be added: %s", err)
	}
	result, err := db.GetMedia(media.ID, FOUND)
	if err != nil {
		t.Fatalf("Unable to get media: %s", err)
	}
	if result != media {
		t.Errorf("Expected: %+v \n Result: %+v", media, result)
	}

}

func TestUpdateMedia(t *testing.T) {

	db := before(t)
	defer after(db)
	err := db.AddMedia(media, FOUND)
	if err != nil {
		t.Fatalf("Unable to add media: %s", err)
	}
	newMedia := media
	newMedia.Name = "NEW MEDIA BITCH"
	err = db.UpdateMedia(newMedia, FOUND)
	if err != nil {
		t.Fatalf("Unable to add media: %s", err)
	}
	result, err := db.GetMedia(media.ID, FOUND)
	if err != nil {
		t.Fatalf("Unable to get media: %s", err)
	}
	if result != newMedia {
		t.Errorf("Expected: %+v \n Result: %+v", media, result)
	}

}

func TestDeleteMedia(t *testing.T) {

	db := before(t)
	defer after(db)
	err := db.AddMedia(media, FOUND)
	if err != nil {
		t.Fatalf("Unable to add media: %s", err)
	}
	err = db.DeleteMedia(media.ID, FOUND)
	if err != nil {
		t.Fatalf("Unable to delete media: %s", err)
	}
	result, err := db.GetMedia(media.ID, FOUND)
	if err != nil {
		t.Fatalf("Unable to get media: %s", err)
	}
	if result != (Media{}) {
		t.Errorf("Expected: %+v \n Result: %+v", Media{}, result)
	}

}
