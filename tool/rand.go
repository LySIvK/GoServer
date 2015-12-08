package tool

import (
	"math/rand"
	"time"
)

const MaxRandNum = 10000

var (
	randValueList [MaxRandNum]int16
	nCurIndex     = 0
)

func disOrder() {
	var nIndex int
	var nTemp int16
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < MaxRandNum; i++ {
		nIndex = rand.Int() % (i + 1)
		if nIndex != i {
			nTemp = randValueList[i]
			randValueList[i] = randValueList[nIndex]
			randValueList[nIndex] = nTemp
		}
	}
}

func initRandom() {
	for i := int16(0); i < MaxRandNum; i++ {
		randValueList[i] = i
	}

	disOrder()
}

func GetRandValue16() int16 {

	nCurIndex = (nCurIndex + 1) % MaxRandNum

	return randValueList[nCurIndex]
}

func GetRandValueInt() int {

	nCurIndex = (nCurIndex + 1) % MaxRandNum

	return int(randValueList[nCurIndex])
}

func HitRandTest(value int) bool {
	if value > GetRandValueInt() {
		return true
	}

	return false
}
