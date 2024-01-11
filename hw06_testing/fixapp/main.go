package fixapp

import (
	"fmt"

	"github.com/aasdhajkshd/home_work_basic/hw02_fix_app/reader"
	"github.com/aasdhajkshd/home_work_basic/hw02_fix_app/types"
)

func PrintStaff(staff []types.Employee) {
	for i := 0; i < len(staff); i++ {
		fmt.Println(staff[i])
	}
}

func FixApp() error {
	var path string

	fmt.Printf("Enter data file path: ")
	fmt.Scanln(&path)

	var err error
	var staff []types.Employee

	if len(path) == 0 {
		path = "../../hw02_fix_app/data.json"
	} else {
		fmt.Println("Successfully read data.json")
	}

	staff, err = reader.ReadJSON(path)

	if err == nil {
		PrintStaff(staff)
	} else {
		return err
	}
	return nil
}
