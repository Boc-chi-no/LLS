package main

import (
	"errors"
	"fmt"
	"linkshortener/db"
	"linkshortener/lib/tool"
	"linkshortener/log"
	"linkshortener/model"
	"linkshortener/setting"
	"strings"
	"testing"
	"time"

	"bou.ke/monkey"
	"github.com/dgraph-io/badger/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// mockConfig is a helper function to set up a mock configuration
func mockConfig(dbType string) {
	setting.Cfg.DB = model.DBConfig{
		Type:     dbType,
		Database: "TEST",
	}
}

func TestDB(t *testing.T) {
	testingT := t
	patchPanicPrint := monkey.Patch(log.PanicPrint, func(format string, values ...interface{}) {
		msg := fmt.Sprintf(format, values...)
		testingT.Logf("PanicPrint called: %s", msg)
	})
	defer patchPanicPrint.Unpatch()

	patchErrorPrint := monkey.Patch(log.ErrorPrint, func(format string, values ...interface{}) {
		msg := fmt.Sprintf(format, values...)
		testingT.Logf("ErrorPrint called: %s", msg)
	})
	defer patchErrorPrint.Unpatch()

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("Panic occurred: %v", r)
		}
	}()

	setting.InitSetting()
	log.InitLog()

	tests := []struct {
		name    string
		dbType  string
		wantErr bool
	}{
		{"BadgerDB", "BADGERDB", false},
		{"MongoDB", "MONGODB", false},
		{"Invalid", "INVALID", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testingT = t
			mockConfig(tt.dbType)
			db.InitDB()

			if tt.dbType == "BADGERDB" && db.BadgerDB == nil {
				t.Errorf("InitDB() BadgerDB is nil")
			}
			if tt.dbType == "MONGODB" && db.MongoDB == nil {
				t.Errorf("InitDB() MongoDB is nil")
			}

			t.Run("Tabler.NewModel", func(t *testing.T) {
				testingT = t
				defer func() {
					if r := recover(); r != nil {
						t.Fatalf("Panic occurred: %v", r)
					}
				}()
				testNewModel(t)
			})

			t.Run("Tabler.SetModel", func(t *testing.T) {
				testingT = t
				defer func() {
					if r := recover(); r != nil {
						t.Fatalf("Panic occurred: %v", r)
					}
				}()
				got := testSetModel(t, tt.dbType)
				testToken, _ := tool.GetToken(16)

				t.Run("Tabler.FindByID", func(t *testing.T) {
					wantType, wantNil := checkGot(t, got)
					if !wantNil {
						var result model.Link
						err := got.FindByID(testToken, &result)
						if err == nil {
							t.Errorf("%s.FindByID() expected an error, but got none", wantType)
							return
						} else if !errors.Is(err, badger.ErrKeyNotFound) && !errors.Is(err, mongo.ErrNoDocuments) { //TODO: 统一错误
							t.Errorf("%s.FindByID() error = %v", wantType, err)
							return
						}
						if result.ShortHash != "" {
							t.Errorf("%s.FindByID() expected ShortHash to be empty, but got %v", wantType, result.ShortHash)
							return
						}
					}
				})

				t.Run("Tabler.InsertOne", func(t *testing.T) {
					wantType, wantNil := checkGot(t, got)
					if !wantNil {
						var link = model.Link{
							ShortHash: testToken,
							URL:       "https://www.example.com/",
							Created:   time.Now().Unix(),
							Memo:      "Test Hash",
							Delete:    false,
						}

						link.Token, _ = tool.GetToken(16)
						insertedID, err := got.InsertOne(link, false)
						if insertedID != link.ShortHash {
							t.Errorf("%s().InsertOne() Inserted document mismatch. Got %v, want %v", wantType, insertedID, link.Token)
							return
						}
						if err != nil {
							t.Errorf("%s().InsertOne() error = %v", wantType, err)
							return
						} else {
							t.Run("Tabler.FindByID Without AutoKey", func(t *testing.T) {
								var result model.Link
								err = got.FindByID(link.ShortHash, &result)
								if err != nil {
									t.Errorf("%s().FindByID() Failed to find inserted document: %v", wantType, err)
									return
								}
								if result.ShortHash != link.ShortHash {
									t.Errorf("%s().FindByID() ShortHash mismatch. Got %v, want %v", wantType, result.ShortHash, link.ShortHash)
									return
								}
								if result.Token != link.Token {
									t.Errorf("%s().FindByID() Find Inserted document mismatch. Got %v, want %v", wantType, result.Token, link.Token)
									return
								}

							})
						}

						link.Token, _ = tool.GetToken(16)
						link.ShortHash = ""
						insertedID, err = got.InsertOne(link, true)
						if insertedID == link.ShortHash {
							t.Errorf("%s().InsertOne() Inserted document match. Got %v, want %v", wantType, insertedID, link.Token)
							return
						}
						if err != nil {
							t.Errorf("%s().InsertOne() error = %v", wantType, err)
							return
						} else {
							t.Run("Tabler.FindByID With AutoKey", func(t *testing.T) {
								var result model.Link
								err = got.FindByID(insertedID, &result)
								if err != nil {
									t.Errorf("%s().FindByID() Failed to find inserted document: %v", wantType, err)
									return
								}

								switch typedID := insertedID.(type) {
								case string:
									if result.ShortHash != typedID {
										t.Errorf("%s().FindByID() ShortHash mismatch. Got %v, want %v", wantType, result.ShortHash, typedID)
										return
									}
								case primitive.ObjectID:
									if result.ShortHash != typedID.Hex() {
										t.Errorf("%s().FindByID ShortHash mismatch. Got %v, want %v", wantType, result.ShortHash, typedID.Hex())
										return
									}
								default:
									t.Errorf("%s().FindByID() Unexpected type for insertedID: %T", wantType, insertedID)
									return
								}

								if result.Token != link.Token {
									t.Errorf("%s().FindByID() Inserted document mismatch. Got %v, want %v", wantType, result.Token, link.Token)
									return
								}

							})
						}

					}
				})

				t.Run("Tabler.UpdateByID", func(t *testing.T) {
					wantType, wantNil := checkGot(t, got)
					if !wantNil {
						var updatedLink model.Link
						updatePassword, _ := tool.GetToken(16)

						err := got.FindByID(testToken, &updatedLink)
						if err != nil {
							t.Errorf("%s().FindByID() Failed to find updated document: %v", wantType, err)
							return
						}
						if updatedLink.Password == updatePassword {
							t.Errorf("%s().FindByID() expected Password to differ from update value, but got %v", wantType, updatedLink.Password)
							return
						}

						// Update the document
						update := bson.M{"$set": bson.M{"password": updatePassword}}
						err = got.UpdateByID(testToken, update)
						if err != nil {
							t.Errorf("%s().UpdateByID() error = %v, wantErr %v", wantType, err, tt.wantErr)
							return
						}

						// Verify the update
						err = got.FindByID(testToken, &updatedLink)
						if err != nil {
							t.Errorf("%s().FindByID() Failed to find updated document: %v", wantType, err)
							return
						}
						if updatedLink.Password != updatePassword {
							t.Errorf("%s().FindByID() Password mismatch. Got %v, want %v", wantType, updatedLink.Password, updatePassword)
							return
						}
						if updatedLink.Token != updatedLink.Token {
							t.Errorf("%s().FindByID() Token mismatch. Got %v, want %v", wantType, updatedLink.Token, updatedLink.Token)
							return
						}
					}
				})

				t.Run("Tabler.Find", func(t *testing.T) {
					wantType, wantNil := checkGot(t, got)
					if !wantNil {
						// First, insert multiple test documents
						testLinks := []model.Link{}
						countToken, _ := tool.GetToken(16)
						for i := 0; i < 5; i++ {
							findToken, _ := tool.GetToken(16)

							link := model.Link{
								URL:      fmt.Sprintf("https://www.example%s.com/", findToken),
								Created:  time.Now().Unix(),
								Memo:     fmt.Sprintf("Test Link %s", findToken),
								Delete:   false,
								Password: countToken,
							}
							link.Token, _ = tool.GetToken(16)
							link.ShortHash, _ = tool.GetToken(8)
							if i == 0 {
								link.Password, _ = tool.GetToken(16)
							}
							_, err := got.InsertOne(link, false)
							if err != nil {
								t.Errorf("%s().InsertOne() error = %v", wantType, err)
								return
							}
							testLinks = append(testLinks, link)
						}

						// Test count documents
						t.Run("Tabler.CountDocuments", func(t *testing.T) {
							filter := bson.M{"password": countToken}
							opts := db.Find()

							count, err := got.CountDocuments(filter, opts)
							if err != nil {
								t.Errorf("%s().CountDocuments() error = %v", wantType, err)
								return
							}

							// Should count at least the number of documents we inserted
							if count != 4 {
								t.Errorf("%s().CountDocuments() count mismatch. Got %d, want %d", wantType, count, 4)
								return
							}

							filter = bson.M{"password": testLinks[0].Password}
							opts = db.Find()

							count, err = got.CountDocuments(filter, opts)
							if err != nil {
								t.Errorf("%s().CountDocuments() error = %v", wantType, err)
								return
							}

							// Should count at least the number of documents we inserted
							if count != 1 {
								t.Errorf("%s().CountDocuments() count mismatch. Got %d, want %d", wantType, count, 1)
								return
							}

						})

						// Test finding all documents
						t.Run("Tabler.Find All Documents", func(t *testing.T) {
							var results []model.Link
							filter := bson.M{"delete": false}
							opts := db.Find().SetLimit(10).SetSkip(0)

							err := got.Find(filter, &results, opts)
							if err != nil {
								t.Errorf("%s().Find() error = %v", wantType, err)
								return
							}

							// Should find at least the number of documents we inserted
							if len(results) < len(testLinks) {
								t.Errorf("%s().Find() found fewer documents than expected. Got %d, want at least %d",
									wantType, len(results), len(testLinks))
								return
							}
						})

						// Test with filter
						t.Run("Tabler.Find With Filter", func(t *testing.T) {
							var results []model.Link
							// Use the memo of the first test link as filter
							filter := bson.M{"memo": testLinks[0].Memo}
							opts := db.Find().SetLimit(10).SetSkip(0)

							err := got.Find(filter, &results, opts)
							if err != nil {
								t.Errorf("%s().Find() error = %v", wantType, err)
								return
							}

							if len(results) != 1 {
								t.Errorf("%s().Find() with specific memo filter found %d documents, want 1",
									wantType, len(results))
								return
							}

							if results[0].Memo != testLinks[0].Memo {
								t.Errorf("%s().Find() memo mismatch. Got %v, want %v",
									wantType, results[0].Memo, testLinks[0].Memo)
								return
							}
						})

						// Test with limit and skip
						t.Run("Tabler.Find With Limit and Skip", func(t *testing.T) {
							var results []model.Link
							filter := bson.M{"delete": false}
							opts := db.Find().SetLimit(2).SetSkip(1)

							err := got.Find(filter, &results, opts)
							if err != nil {
								t.Errorf("%s().Find() error = %v", wantType, err)
								return
							}

							if len(results) != 2 {
								t.Errorf("%s().Find() with limit 2 and skip 1 found %d documents, want 2",
									wantType, len(results))
								return
							}
						})

						// Test with key for BadgerDB
						if wantType == "*db.BadgerDBTable" {
							t.Run("Tabler.Find with Key for BadgerDB", func(t *testing.T) {
								var results []model.Link
								filter := bson.M{}
								opts := db.Find().SetKey(testLinks[0].ShortHash)

								err := got.Find(filter, &results, opts)
								if err != nil {
									t.Errorf("%s().Find() error = %v", wantType, err)
									return
								}

								found := false
								for _, result := range results {
									if result.ShortHash == testLinks[0].ShortHash {
										found = true
										break
									}
								}

								if !found {
									t.Errorf("%s().Find() with key did not find the expected document", wantType)
									return
								}
							})

							t.Run("Tabler.Find with PrefixScans for BadgerDB", func(t *testing.T) {
								var results []model.Link
								filter := bson.M{}
								// Use a common prefix from the first test link's hash
								prefix := testLinks[0].ShortHash[:2]
								opts := db.Find().SetKey(prefix).SetPrefixScans(true).SetLimit(10)

								err := got.Find(filter, &results, opts)
								// We don't check for errors here as not all hashes may have the same prefix
								if err == nil && len(results) > 0 {
									for _, result := range results {
										if !strings.HasPrefix(result.ShortHash, prefix) {
											t.Errorf("%s().Find() with prefix scan returned document with non-matching prefix: %s",
												wantType, result.ShortHash)
											return
										}
									}
								}
							})
						}
					}
				})
			})
		})
	}
}

func testNewModel(t *testing.T) {
	t.Run("Tabler.SetDB", func(t *testing.T) {
		got := db.NewModel(setting.Cfg.DB.Database, "testTable")
		wantType, wantNil := checkGot(t, got)
		if !wantNil {
			t.Run("Tabler.CreateOneIndex", func(t *testing.T) {
				statsIndex := mongo.IndexModel{
					Keys: bson.M{
						"token": 1,
					},
					Options: options.Index().SetName("IDX_TOKEN"),
				}
				err := got.CreateOneIndex(statsIndex)
				if err != nil {
					t.Errorf("%s.CreateOneIndex() error = %v", wantType, err)
				}
			})
		}
	})
}

func checkGot(t *testing.T, got db.Tabler) (string, bool) {
	wantNil := setting.Cfg.DB.Type == "INVALID"
	wantType := ""
	switch setting.Cfg.DB.Type {
	case "BADGERDB":
		wantType = "*db.BadgerDBTable"
	case "MONGODB":
		wantType = "*db.MongoDBTable"
	}

	if (got == nil) != wantNil {
		t.Errorf("got = %v, want nil = %v", got, wantNil)
	}
	if !wantNil && wantType != "" {
		if gotType := fmt.Sprintf("%T", got); gotType != wantType {
			t.Errorf("got type = %v, want type %v", gotType, wantType)
		}
	}
	return wantType, wantNil
}

func testSetModel(t *testing.T, dbType string) db.Tabler {
	got := db.SetModel(setting.Cfg.DB.Database, "testTable")
	wantNil := dbType == "INVALID"
	wantType := ""
	switch dbType {
	case "BADGERDB":
		wantType = "*db.BadgerDBTable"
	case "MONGODB":
		wantType = "*db.MongoDBTable"
	}

	if (got == nil) != wantNil {
		t.Errorf("got = %v, want nil = %v", got, wantNil)
	}
	if !wantNil && wantType != "" {
		if gotType := fmt.Sprintf("%T", got); gotType != wantType {
			t.Errorf("got type = %v, want type %v", gotType, wantType)
		}
	}
	if !wantNil {
		return got
	}
	return nil
}
