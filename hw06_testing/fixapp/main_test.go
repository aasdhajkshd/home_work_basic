package fixapp

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadJSON(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "testData.json")
	if err != nil {
		t.Fatal(err)
	} 
	// Мы писали, мы писали, наши пальчики устали...
	testData := `[{"userId": 3, "age": 44, "name": "Васиья", "departmentId": 7}]`
	_, err = tmpFile.Write([]byte(testData))
	if err != nil {
		t.Fatal(err)
	} 

	defer os.Remove(tmpFile.Name()) // какая прелесть

	// Тестируем функцию ReadJSON по пути "github.com"
	employee, err := ReadJSON(tmpFile.Name())
	assert.NoError(t, err, "На случай ошибки в JSON")
	assert.NotNil(t, employee, "Employees не должны быть nil")
	assert.Len(t, employee, 1, "Шикарная штука testify...")
	// как и поиск в Slice, чем перебирать...

	// Проверяем, что значения структуры совпадают с ожидаемыми
	expectedEmployee := Employee{UserID: 3, Age: 44, Name: "Васиья", DepartmentID: 7}
	assert.Equal(t, expectedEmployee, employee[0], "Жаль, мы так старались...")
}

func TestFixApp(t *testing.T) {
	err := FixApp()

	// Просто проверим, что нет не возникло ошибок :-(
	assert.NoError(t, err)
}
