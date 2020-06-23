package excellib

import (
	"fmt"
	"io"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/hulklab/cast"
)

type E struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

// g.D{{"foo","bar","hello","world"}}
type D []E

type Arg struct {
	SheetName      string
	DivideFirstRow bool // 分离第一行
	AutoAlign      bool // 自动对齐
}

type Option func(arg *Arg)

func WithSheetName(name string) Option {
	return func(arg *Arg) {
		arg.SheetName = name
	}
}

func WithDivideFirstRow(b bool) Option {
	return func(arg *Arg) {
		arg.DivideFirstRow = b
	}
}

func WithAutoAlign() Option {
	return func(arg *Arg) {
		arg.AutoAlign = true
	}
}

func extractArg(opts ...Option) Arg {
	sheetName := "Sheet1"

	arg := Arg{
		SheetName: sheetName,
	}

	for _, opt := range opts {
		opt(&arg)
	}

	return arg
}

func ImportFile(filename string, divideFirst bool, opts ...Option) ([][]string, []string, error) {

	f, err := excelize.OpenFile(filename)
	if err != nil {
		return nil, nil, err
	}
	opts = append(opts, WithDivideFirstRow(divideFirst))
	return importExcel(f, opts...)
}

func ImportReader(r io.Reader, divideFirst bool, opts ...Option) ([][]string, []string, error) {
	f, err := excelize.OpenReader(r)
	if err != nil {
		return nil, nil, err
	}
	opts = append(opts, WithDivideFirstRow(divideFirst))
	return importExcel(f, opts...)
}

func importExcel(f *excelize.File, opts ...Option) ([][]string, []string, error) {
	arg := extractArg(opts...)

	contents := make([][]string, 0)
	headers := make([]string, 0)

	rows, err := f.GetRows(arg.SheetName)

	if err != nil {
		return contents, nil, err
	}

	if len(rows) == 0 {
		return contents, headers, nil
	}

	if arg.DivideFirstRow {
		headers = rows[0]
		// 矫整列数据
		for i := len(headers) - 1; i >= 0; i-- {
			// 去掉末尾空串
			if len(headers[i]) == 0 {
				headers = headers[0:i]
			} else {
				break
			}
		}
	}

	max := 0
	for i, row := range rows {
		// 如果有标题，以标题的长度为准
		if i == 0 && arg.DivideFirstRow {
			max = len(headers)
			break
		}

		n := len(row)
		if n > max {
			max = n
		}
	}

	for i, row := range rows {
		if i == 0 && arg.DivideFirstRow {
			continue
		}

		if len(row) == 0 {
			continue
		}

		// 补齐
		if len(row) < max {
			for i := len(row); i < max; i++ {
				row = append(row, "")
			}
		}

		// 每一列都空直接过滤掉
		flag := true
		for _, col := range row {
			if len(col) > 0 {
				flag = false
				break
			}
		}

		if flag {
			continue
		}

		contents = append(contents, row)
	}

	return contents, headers, nil
}

// rows [{id:xx,name:xx,username},{}] headers [{id:"ID"},{name,姓名}]
func ExportAssoc(rows []map[string]interface{}, headers D, opts ...Option) (*excelize.File, error) {
	headerList := make([]string, 0)

	// 先处理表头
	for _, kv := range headers {

		v := cast.ToString(kv.Value)

		headerList = append(headerList, v)
	}

	rowList := make([][]interface{}, 0)

	// 再处理内容
	for _, row := range rows {

		rowData := []interface{}{}

		for _, kv := range headers {
			val, ok := row[kv.Key]
			if !ok {
				val = ""
			}

			rowData = append(rowData, val)
		}

		rowList = append(rowList, rowData)
	}

	return Export(rowList, headerList, opts...)

}

// rows [[id_value,name_value,username_value]] headers  [ID,姓名]
func Export(rows [][]interface{}, headers []string, opts ...Option) (*excelize.File, error) {
	arg := extractArg(opts...)

	handle := excelize.NewFile()

	err := handle.SetSheetRow(arg.SheetName, "A1", &headers)
	if err != nil {
		return nil, err
	}

	if arg.AutoAlign {
		var rate = 1.2

		// 获取每一列的最大宽度
		for i := 0; i < len(headers); i++ {

			colName, err := excelize.ColumnNumberToName(i + 1)
			if err != nil {
				return nil, err
			}

			maxCnt := searchCount(headers[i])

			for j, row := range rows {
				if len(row) < len(headers) {
					return nil, fmt.Errorf("第 %d 行数据少于 header 个数", j+1)
				}

				str, err := cast.ToStringE(row[i])
				if err != nil {
					return nil, fmt.Errorf("第 %d 行 %d 列不是有效的字符串", j+1, i+1)
				}

				cnt := searchCount(str)
				if cnt > maxCnt {
					maxCnt = cnt
				}
			}

			fitColWidth := float64(maxCnt) * rate
			if fitColWidth > 255 {
				fitColWidth = 255
			} else if fitColWidth < 9 {
				fitColWidth = 9
			}

			//fmt.Println("col:",colName,"width:",fitColWidth)

			_ = handle.SetColWidth(arg.SheetName, colName, colName, fitColWidth)
		}

	}

	// 再处理内容
	for i, row := range rows {
		axis := fmt.Sprintf("A%d", i+2)

		err := handle.SetSheetRow(arg.SheetName, axis, &row)
		if err != nil {
			return nil, err
		}
	}

	return handle, nil
}

// Excel所有出现的值进行匹配给出对应宽度值
func searchCount(src string) int {
	letters := "abcdefghijklmnopqrstuvwxyz"
	letters = letters + strings.ToUpper(letters)
	nums := "0123456789"
	chars := "()/#"

	numCount := 0
	letterCount := 0
	othersCount := 0
	charsCount := 0

	for _, i := range src {
		switch {
		case strings.ContainsRune(letters, i) == true:
			letterCount += 1
		case strings.ContainsRune(nums, i) == true:
			numCount += 1
		case strings.ContainsRune(chars, i) == true:
			charsCount += 1
		default:
			othersCount += 1
		}
	}

	return numCount*1 + letterCount*1 + charsCount*1 + othersCount*2
}
