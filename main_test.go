// go test -gcflags=all=-l
package main

import (
	"fmt"
	"linkshortener/lib/shorten"
	"linkshortener/lib/tool"
	"linkshortener/model"
	"strconv"
	"testing"
	"time"

	"bou.ke/monkey"
	"github.com/go-playground/assert/v2"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

func TestGetToken(t *testing.T) {
	testLength := 10000
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	mt.Run("Token Length Zero", func(mt *mtest.T) {
		token, err := tool.GetToken(0)
		assert.Equal(t, 0, len(token))
		assert.NotEqual(t, err, nil)
		fmt.Println(tool.ConcatStrings("Token Length Zero Success len:", strconv.Itoa(len(token))))
	})

	mt.Run("Token Length Large", func(mt *mtest.T) {
		token, err := tool.GetToken(4096)
		assert.Equal(t, 4096, len(token))
		assert.Equal(t, err, nil)
		fmt.Println(tool.ConcatStrings("Token Length Large Success len:", strconv.Itoa(len(token))))
	})

	mt.Run("Token length Equal", func(mt *mtest.T) {
		for i := 1; i < testLength; i++ {
			token, err := tool.GetToken(i)
			assert.Equal(t, i, len(token))
			assert.Equal(t, err, nil)
		}
		fmt.Println(tool.ConcatStrings("TestGetToken length Equal Success loop:", strconv.Itoa(testLength)))
	})

	testLength = 10000000

	mt.Run("Token  Length 16 NotEqual", func(mt *mtest.T) {
		token := make(map[string]bool)

		for i := 0; i < testLength; i++ {
			token01, err01 := tool.GetToken(16)
			token[token01] = true
			assert.Equal(t, err01, nil)
		}
		assert.Equal(t, len(token), testLength)
		fmt.Println(tool.ConcatStrings("TestGetToken Token Length 16 NotEqual Success loop:", strconv.Itoa(testLength)))
	})

	mt.Run("Token Length 32 NotEqual", func(mt *mtest.T) {
		token := make(map[string]bool)

		for i := 0; i < testLength; i++ {
			token01, err01 := tool.GetToken(32)
			token[token01] = true
			assert.Equal(t, err01, nil)
		}
		assert.Equal(t, len(token), testLength)
		fmt.Println(tool.ConcatStrings("TestGetToken Token Length 32 NotEqual Success loop:", strconv.Itoa(testLength)))
	})

}

func TestGenerateShortenLink(t *testing.T) {
	testLength := 10000000
	var request model.InsertLinkReq
	token, _ := tool.GetToken(32)

	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("ShortenLinkHash Equal", func(mt *mtest.T) {
		mockTime := time.Date(2023, 10, 11, 11, 10, 11, 11, time.UTC)
		patchTimeTime := monkey.Patch(time.Now, func() time.Time { return mockTime })
		defer patchTimeTime.Unpatch()

		patchToolGlobalCounterSafeAdd := monkey.Patch(tool.GlobalCounterSafeAdd, func(uint64) uint64 { return uint64(100000) })
		defer patchToolGlobalCounterSafeAdd.Unpatch()

		patchToolGetToken := monkey.Patch(tool.GetToken, func(int2 int) (string, error) { return token, nil })
		defer patchToolGetToken.Unpatch()

		request.URL = "https://www.lioat.cn/"
		mt.AddMockResponses(mtest.CreateSuccessResponse())

		t01 := shorten.GenerateShortenLink(request)
		t02 := shorten.GenerateShortenLink(request)
		assert.Equal(t, t01.ShortHash, t02.ShortHash)

		fmt.Println("TestGenerateShortenLink Equal Success")
	})

	mt.Run("ShortenLinkHash NotEqual (Monkey Time)", func(mt *mtest.T) {
		mockTime := time.Date(2023, 10, 11, 11, 10, 11, 11, time.UTC)
		patchTimeTime := monkey.Patch(time.Now, func() time.Time { return mockTime })
		defer patchTimeTime.Unpatch()

		patchToolGetToken := monkey.Patch(tool.GetToken, func(int2 int) (string, error) { return "TestToken0000", nil })
		defer patchToolGetToken.Unpatch()

		shortHashes := make(map[string]int)
		for i := 0; i < testLength; i++ {
			request.URL = "https://www.lioat.cn/"
			mt.AddMockResponses(mtest.CreateSuccessResponse())
			linkInfo := shorten.GenerateShortenLink(request)
			if index, ok := shortHashes[linkInfo.ShortHash]; ok {
				fmt.Printf("Test Length [%d-%d] ShortHash Equal (Monkey Time) ': %s\n", index, i, linkInfo.ShortHash)
			} else {
				shortHashes[linkInfo.ShortHash] = i
			}
		}
		fmt.Printf("TestGenerateShortenLink NotEqual (Monkey Time) loop:%d collision:%d collisionRate:%.2f%%\n", testLength, testLength-len(shortHashes), float64(testLength-len(shortHashes))/float64(testLength)*100)
		assert.Equal(t, len(shortHashes) >= int(float64(testLength)*0.998), true)
		fmt.Println("TestGenerateShortenLink NotEqual (Monkey Time) Success")
	})

}
