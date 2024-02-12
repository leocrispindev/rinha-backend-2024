package util

import "strconv"

func StringToInt(str string) (int, error) {

	int64Val, err := strconv.ParseInt(str, 10, 32)

	if err != nil {
		return 0, err
	}

	intValue := int(int64Val)
	if int64Val != int64(intValue) {
		return 0, strconv.ErrRange
	}

	return intValue, nil
}
