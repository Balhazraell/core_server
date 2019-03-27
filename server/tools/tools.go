package tools

//DeleElementFromArraByIndex - Метод удаляет элемент массива по индексу.
func DeleElementFromArraByIndex(array []int, index int) []int {
	array[index] = array[len(array)-1]
	return array[:len(array)-1]
}

//FindElementInArray - Метод находит первое вхождение заданного элемента в массиве.
func FindElementInArray(array []int, val int) int {
	var result = -1

	for i, v := range array {
		if val == v {
			return i
		}
	}

	return result
}
