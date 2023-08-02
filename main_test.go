// go test -gcflags=all=-l
package main

import (
	"bou.ke/monkey"
	"fmt"
	"github.com/go-playground/assert/v2"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
	"linkshortener/lib/shorten"
	"linkshortener/lib/tool"
	"linkshortener/model"
	"strconv"
	"testing"
	"time"
)

func TestGetToken(t *testing.T) {
	testLength := 10000
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()
	mt.Run("Token length Equal", func(mt *mtest.T) {
		for i := 0; i < testLength; i++ {
			token, _ := tool.GetToken(i)
			assert.Equal(t, i, len(token))
		}
		fmt.Println(tool.ConcatStrings("TestGetToken length Equal Success loop:", strconv.Itoa(testLength)))
	})

	mt.Run("Token NotEqual", func(mt *mtest.T) {
		for i := 6; i < testLength/2+6; i++ {
			token01, _ := tool.GetToken(i)
			token02, _ := tool.GetToken(i)
			assert.NotEqual(t, token01, token02)
		}
		fmt.Println(tool.ConcatStrings("TestGetToken Token NotEqual Success loop:", strconv.Itoa(testLength)))
	})
}

func TestGenerateShortenLink(t *testing.T) {
	testLength := 10000
	var request model.InsertLinkReq

	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("ShortenLinkHash NotEqual", func(mt *mtest.T) {
		for i := 0; i < testLength; i++ {
			token, _ := tool.GetToken(i)
			request.URL = tool.ConcatStrings("https://www.lioat.cn/", token)
			mt.AddMockResponses(mtest.CreateSuccessResponse())

			t01 := shorten.GenerateShortenLink(request)
			t02 := shorten.GenerateShortenLink(request)
			assert.NotEqual(t, t01.ShortHash, t02.ShortHash)
		}
	})
	fmt.Println(tool.ConcatStrings("TestGenerateShortenLink NotEqual Success loop:", strconv.Itoa(testLength)))

	mt.Run("ShortenLinkHash Equal", func(mt *mtest.T) {
		mockTime := time.Date(2023, 10, 11, 11, 10, 11, 11, time.UTC)
		patchTimeTime := monkey.Patch(time.Now, func() time.Time { return mockTime })
		defer patchTimeTime.Unpatch()

		patchToolGlobalCounterSafeAdd := monkey.Patch(tool.GlobalCounterSafeAdd, func(uint64) uint64 { return uint64(100000) })
		defer patchToolGlobalCounterSafeAdd.Unpatch()

		patchToolGetToken := monkey.Patch(tool.GetToken, func(int2 int) (string, error) { return "TestToken0000", nil })
		defer patchToolGetToken.Unpatch()

		for i := 0; i < testLength*10; i++ {
			request.URL = "https://www.lioat.cn/"
			mt.AddMockResponses(mtest.CreateSuccessResponse())

			t01 := shorten.GenerateShortenLink(request)
			t02 := shorten.GenerateShortenLink(request)
			assert.Equal(t, t01.ShortHash, t02.ShortHash)
		}
	})

	fmt.Println(tool.ConcatStrings("TestGenerateShortenLink Equal Success loop:", strconv.Itoa(testLength*10)))
}
