package grmath

import (
	"errors"
	"fmt"
)

func UnMaskWithPartialMd5(mark string, startPos int, endPos int, partialMd5 string, markType int) (unMark string, err error) {

	lengthMark := len(mark)
	if lengthMark <= endPos || endPos < startPos || startPos < 0 {
		err = errors.New("pos not good")
		return
	}

	//markStringRune := []rune(mark)

	if markType == 1 {
		// identity
		if lengthMark != 18 || startPos != 9 || endPos != 14 {
			err = fmt.Errorf("length is not good for identity. %d", lengthMark)
			return
		}

		for i := 0; i <= 9; i++ {
			// pos 9, a random digit

			//pos 10-13 4year 2month 2day

			for month := 1; month <= 12; month++ {
				for day := 1; day <= 31; day++ {
					for j := 0; j <= 9; j++ {

						genString := fmt.Sprintf("%d%02d%02d%d", i, month, day, j)

						//fmt.Println(genString)
						md5 := Md5(genString)
						if md5 == partialMd5 {
							unMark = fmt.Sprintf("%s%s%s", mark[0:startPos], genString, mark[endPos+1:lengthMark])
							return
						}
					}
				}

			}
		}

	}

	err = fmt.Errorf("can not unmark")
	return
}

func UnMaskWithWholeMd5(mark string, startPos int, endPos int, wholeMd5 string, markType int) (unMark string, err error) {

	lengthMark := len(mark)
	if lengthMark <= endPos || endPos < startPos || startPos < 0 {
		err = errors.New("pos not good")
		return
	}

	//markStringRune := []rune(mark)

	if markType == 1 {
		// digit
		max := Pow(10, endPos-startPos+1)

		format := fmt.Sprintf("%%s%%%dd%%s", endPos-startPos+1)

		for i := 0; i < max; i++ {
			// pos 5, a random digit

			genString := fmt.Sprintf(format, mark[0:startPos], i, mark[endPos+1:lengthMark])

			//fmt.Println(genString)
			md5 := Md5(genString)
			if md5 == wholeMd5 {
				unMark = genString
				return
			}

		}

	}

	err = fmt.Errorf("can not unmark")
	return
}
func MaskWithWholeMd5(unMark string, startPos int, endPos int) (mark string, wholeMd5 string, err error) {

	lengthUnMark := len(unMark)
	if lengthUnMark <= endPos || endPos < startPos || startPos < 0 {
		err = errors.New("pos not good")
		return
	}

	markStringRune := []rune(unMark)

	for tmpStartPos := startPos; tmpStartPos <= endPos; tmpStartPos = tmpStartPos + 1 {

		markStringRune[tmpStartPos] = rune('*')

	}

	mark = string(markStringRune)

	wholeMd5 = Md5(unMark)

	return
}

func MaskWithPartialMd5(unMark string, startPos int, endPos int) (mark string, partialMd5 string, err error) {

	lengthUnMark := len(unMark)
	if lengthUnMark <= endPos || endPos < startPos || startPos < 0 {
		err = errors.New("pos not good")
		return
	}

	markStringRune := []rune(unMark)

	partialString := ""

	for tmpStartPos := startPos; tmpStartPos <= endPos; tmpStartPos = tmpStartPos + 1 {
		partialString += string(markStringRune[tmpStartPos])
		markStringRune[tmpStartPos] = rune('*')

	}

	mark = string(markStringRune)

	partialMd5 = Md5(partialString)

	return
}
