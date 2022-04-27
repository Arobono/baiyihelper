package dbw

import (
	"bytes"
	"database/sql"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
	"unsafe"

	"github.com/Arobono/baiyihelper/helper/cal"
	"github.com/go-sql-driver/mysql"
	"github.com/gogf/gf/v2/util/guid"
)

type BaseKLine struct {
	Cdatetime time.Time //日期"`
	Copen     float64   //开盘价"`
	Chigh     float64   //最高价"`
	Clow      float64   //最低价"`
	Cclose    float64   //收盘价"`
	Cvolume   int64     //成交量"`
	Cturnover float64   //成交额"`
	Cinterval string    //时间间隔"`

	CAverage float64 //均价
}

type DbWorker struct {
	Dsn string
	Db  *sql.DB
}

func InitDb(link string) (dbw DbWorker) {

	dbw = DbWorker{Dsn: link}
	dbtemp, err := sql.Open("mysql", dbw.Dsn)
	dbw.Db = dbtemp
	if err != nil {
		panic(err)
	}
	return
}

//获取数字货币K线
func (dbw DbWorker) GetKLine(Symbol, interval, beginTime, endTime string) (baseKLineList []BaseKLine, err error) {
	defer dbw.Db.Close()
	var sqlBuffer bytes.Buffer
	sqlBuffer.WriteString(" WHERE Cinterval = '")
	sqlBuffer.WriteString(interval)
	sqlBuffer.WriteString("' AND Cdatetime >='")
	sqlBuffer.WriteString(beginTime)
	sqlBuffer.WriteString("' AND Cdatetime <='")
	sqlBuffer.WriteString(endTime)
	sqlBuffer.WriteString("';")
	sqlwhere := sqlBuffer.String()
	var rows *sql.Rows

	switch Symbol {

	case "Futures":
		{
			if interval == "1m" {
				rows, err = dbw.Db.Query("SELECT * FROM `ths_hisonefutures` WHERE  Cdatetime >='" + beginTime + "' AND AND Cdatetime <='" + endTime + "';")
			} else {
				rows, err = dbw.Db.Query("SELECT * FROM `ths_hisfutures`" + sqlwhere)
			}

		}
	case "BTC":
		{
			rows, err = dbw.Db.Query("SELECT * FROM `tb_dcp_btcperpetualfuture`" + sqlwhere)
		}
	case "ETH":
		{
			rows, err = dbw.Db.Query("SELECT * FROM `tb_dcp_ethperpetualfuture`" + sqlwhere)
		}
	default:
		{
			err = errors.New("未支持" + Symbol)
			return
		}
	}
	defer rows.Close()
	if err != nil {
		return
	}
	columns, _ := rows.Columns()
	values := make([]sql.RawBytes, len(columns))
	scans := make([]interface{}, len(columns))
	for i := range values {
		scans[i] = &values[i]
	}
	i := 0
	baseKLineList = make([]BaseKLine, 0)
	for rows.Next() {
		err = rows.Scan(scans...)
		var u = new(BaseKLine)
		value := reflect.ValueOf(u)
		if value.Kind() == reflect.Ptr {
			for i, col := range values {
				item := string(col)
				elem := value.Elem()
				name := elem.FieldByName(cal.Capitalize(columns[i]))
				switch name.Kind() {
				case reflect.String:
					*(*string)(unsafe.Pointer(name.Addr().Pointer())) = item
				case reflect.Int64:
					intitem, _ := strconv.ParseInt(item, 10, 64)
					*(*int64)(unsafe.Pointer(name.Addr().Pointer())) = intitem
				case reflect.Float64:
					floattitem, _ := strconv.ParseFloat(item, 64)
					*(*float64)(unsafe.Pointer(name.Addr().Pointer())) = floattitem
				case reflect.Struct:
					loc, _ := time.LoadLocation("Asia/Shanghai") //获取当地时区
					location, _ := time.ParseInLocation("2006-01-02 15:04:05", item, loc)
					*(*time.Time)(unsafe.Pointer(name.Addr().Pointer())) = location
				default:
				}
			}
		}
		u.CAverage = cal.Divide(cal.Add(u.Copen, u.Clow, u.Chigh, u.Cclose), 4)
		baseKLineList = append(baseKLineList, *u)
		i++
	}
	return
}

//统计收益
func (dbw DbWorker) StatProfit(ployid string) (err error) {

	return
}

//批量插入
func (dbw DbWorker) BulkLoaderInsert(data []interface{}, tableName string, iscover bool) (issuccess bool, err error) {
	//创建Csv文件
	issuccess = false
	filename, failedAr, err := toCsv(data)
	if err != nil {
		fmt.Println(err)
		return
	}
	if filename == "" {
		msg := "文件名空了， " + filename
		fmt.Println(msg)
		err = errors.New(msg)
		return
	}

	//文件地址
	path := "./" + filename

	//导出后删除文件
	defer func() {
		err := os.Remove(path)
		if err != nil {
			fmt.Println("删除文件失败", err)
		}
	}()

	//数据库链接
	dbtemp, err := sql.Open("mysql", dbw.Dsn)
	if err != nil {
		fmt.Println("数据库链接失败", err)
		return
	}
	dbw.Db = dbtemp
	defer dbw.Db.Close()

	var handName = guid.S()
	mysql.RegisterReaderHandler(handName, func() io.Reader {
		file, err := os.Open(path)
		if err != nil {
			return nil
		}
		return file
	})
	var sqlBuild strings.Builder
	sqlBuild.WriteString(" LOAD DATA LOCAL INFILE 'Reader::")
	sqlBuild.WriteString(handName)
	sqlBuild.WriteString("' ")
	if iscover {
		sqlBuild.WriteString(" REPLACE ")
	}
	sqlBuild.WriteString(" INTO TABLE ")
	sqlBuild.WriteString(tableName)
	sqlBuild.WriteString(" CHARACTER SET UTF8 ")
	sqlBuild.WriteString(" FIELDS TERMINATED BY ',' ")
	sqlBuild.WriteString(" IGNORE 1 LINES ")
	sqlBuild.WriteString(" (`")
	sqlBuild.WriteString(strings.Join(failedAr, "`,`"))
	sqlBuild.WriteString("`) ")
	sqlBuild.WriteString(";")

	_, err = dbw.Db.Exec(sqlBuild.String())
	if err != nil {
		fmt.Println("数据库插入失败", err)
		return
	}
	return true, nil
}

// 实现csv文件写入的过程
func toCsv(data []interface{}) (filename string, failedAr []string, err error) {
	strTime := time.Now().Format("20060102150405")
	filename = fmt.Sprintf("临时%s.csv", strTime)
	xlsFile, err := os.OpenFile("./"+filename, os.O_RDWR|os.O_CREATE, 0766)
	if err != nil {
		fmt.Println("创建失败", err)
		return
	}
	defer xlsFile.Close()

	xlsFile.WriteString("\xEF\xBB\xBF")
	wStr := csv.NewWriter(xlsFile)

	//var failedAr = GetField(data[0])
	failedAr = cal.GetTagField(data[0], "json")
	wStr.Write(failedAr)
	for _, item := range data {
		//var valueAr = GetValue(item)
		var valueAr = cal.GetTagValue(item, "json")
		wStr.Write(valueAr)
	}
	wStr.Flush()
	return
}
