package reader

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/aasdhajkshd/home_work_basic/hw02_fix_app/types"
)

func ReadJSON(filePath string) ([]types.Employee, error) {
	f, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}

	bytes, err := io.ReadAll(f)
	if err != nil {
		fmt.Println("Error reading JSON content:", err)
		return nil, err
	}

	var data []types.Employee

	err = json.Unmarshal(bytes, &data)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return nil, err
	}

	res := data
	return res, nil
}
