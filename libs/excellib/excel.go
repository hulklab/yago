package excellib

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/hulklab/cast"
)

const sheetName = "Sheet1"

type E struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

// D{{"foo","bar","hello","world"}}
type D []E

type List = []interface{}
type Hash = map[string]interface{}

type ImportReq struct {
	Filename       string
	Reader         io.Reader
	DivideFirstRow bool
	SheetName      string
}

func (r *ImportReq) GetSheetName() string {
	if r.SheetName == "" {
		return sheetName
	}

	return r.SheetName
}

type ImportResp struct {
	Rows    [][]string
	Headers []string
}

type ExportReq struct {
	Rows      []List
	Headers   []string
	SheetName string
	AutoAlign bool
}

func (r *ExportReq) GetSheetName() string {
	if r.SheetName == "" {
		return sheetName
	}

	return r.SheetName
}

type ExportAssocReq struct {
	Rows      []Hash
	Headers   D
	SheetName string
	AutoAlign bool
}

func ImportFile(req *ImportReq) (resp *ImportResp, err error) {
	if len(req.Filename) == 0 {
		return nil, errors.New("filename is required")
	}

	f, err := excelize.OpenFile(req.Filename)
	if err != nil {
		return
	}
	return importExcel(f, req)
}

func ImportReader(req *ImportReq) (resp *ImportResp, err error) {
	if req.Reader == nil {
		return nil, errors.New("reader is required")
	}

	f, err := excelize.OpenReader(req.Reader)
	if err != nil {
		return nil, err
	}
	return importExcel(f, req)
}

func importExcel(f *excelize.File, req *ImportReq) (resp *ImportResp, err error) {
	contents := make([][]string, 0)
	headers := make([]string, 0)
	resp = &ImportResp{}

	rows, err := f.GetRows(req.GetSheetName())

	if err != nil {
		return nil, err
	}

	if len(rows) == 0 {
		return resp, nil
	}

	if req.DivideFirstRow {
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
		if i == 0 && req.DivideFirstRow {
			max = len(headers)
			break
		}

		n := len(row)
		if n > max {
			max = n
		}
	}

	for i, row := range rows {
		if i == 0 && req.DivideFirstRow {
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

	resp.Headers = headers
	resp.Rows = contents

	return resp, nil
}

// rows [{id:xx,name:xx,username},{}] headers [{id:"ID"},{name,姓名}]
func ExportAssoc(req *ExportAssocReq) (*excelize.File, error) {
	headerList := make([]string, 0)

	// 先处理表头
	for _, kv := range req.Headers {

		v := cast.ToString(kv.Value)

		headerList = append(headerList, v)
	}

	rowList := make([][]interface{}, 0)

	// 再处理内容
	for _, row := range req.Rows {

		rowData := List{}

		for _, kv := range req.Headers {
			val, ok := row[kv.Key]
			if !ok {
				val = ""
			}

			rowData = append(rowData, val)
		}

		rowList = append(rowList, rowData)
	}

	return Export(&ExportReq{
		Rows:      rowList,
		Headers:   headerList,
		SheetName: req.SheetName,
		AutoAlign: req.AutoAlign,
	})

}

func Export(req *ExportReq) (*excelize.File, error) {

	handle := excelize.NewFile()

	err := handle.SetSheetRow(req.GetSheetName(), "A1", &req.Headers)
	if err != nil {
		return nil, err
	}

	if req.AutoAlign {
		var rate = 1.2

		// 获取每一列的最大宽度
		for i := 0; i < len(req.Headers); i++ {

			colName, err := excelize.ColumnNumberToName(i + 1)
			if err != nil {
				return nil, err
			}

			maxCnt := searchCount(req.Headers[i])

			for j, row := range req.Rows {
				if len(row) < len(req.Headers) {
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

			// fmt.Println("col:",colName,"width:",fitColWidth)

			_ = handle.SetColWidth(req.GetSheetName(), colName, colName, fitColWidth)
		}

	}

	// 再处理内容
	for i, row := range req.Rows {
		axis := fmt.Sprintf("A%d", i+2)

		err := handle.SetSheetRow(req.GetSheetName(), axis, &row)
		if err != nil {
			return nil, err
		}
	}

	return handle, nil
}

// Excel 所有出现的值进行匹配给出对应宽度值
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
