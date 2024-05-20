package others

//
//import "fmt"
//
//func main() {
//	var a []int
//	a = new([]int)
//	a = append(a, 1)
//	fmt.Println(a)
//}
//
////package main
////
////import (
////	"database/sql"
////	"errors"
////	"fmt"
////	_ "github.com/go-sql-driver/mysql"
////	"product/datamodels"
////	"reflect"
////	"strconv"
////	"time"
////)
////
////// 创建mysql 连接
////func NewMysqlConn() (db *sql.DB, err error) {
////	db, err = sql.Open("mysql", "root:imooc@tcp(127.0.0.1:3306)/imooc?charset=utf8")
////	return
////}
////
////// 获取返回值，获取一条
////func GetResultRow(rows *sql.Rows) map[string]string {
////	columns, _ := rows.Columns()
////	scanArgs := make([]interface{}, len(columns))
////	values := make([][]byte, len(columns))
////	for j := range values {
////		scanArgs[j] = &values[j]
////	}
////	record := make(map[string]string)
////	for rows.Next() {
////		//将行数据保存到record字典
////		rows.Scan(scanArgs...)
////		for i, v := range values {
////			if v != nil {
////				//fmt.Println(reflect.TypeOf(col))
////				record[columns[i]] = string(v)
////			}
////		}
////	}
////	return record
////}
////
////// 获取所有
////func GetResultRows(rows *sql.Rows) map[int]map[string]string {
////	//返回所有列
////	columns, _ := rows.Columns()
////	//这里表示一行所有列的值，用[]byte表示
////	vals := make([][]byte, len(columns))
////	//这里表示一行填充数据
////	scans := make([]interface{}, len(columns))
////	//这里scans引用vals，把数据填充到[]byte里
////	for k, _ := range vals {
////		scans[k] = &vals[k]
////	}
////	i := 0
////	result := make(map[int]map[string]string)
////	for rows.Next() {
////		//填充数据
////		rows.Scan(scans...)
////		//每行数据
////		row := make(map[string]string)
////		//把vals中的数据复制到row中
////		for k, v := range vals {
////			key := columns[k]
////			//这里把[]byte数据转成string
////			row[key] = string(v)
////		}
////		//放入结果集
////		result[i] = row
////		i++
////	}
////	return result
////}
////
////// 第一步，先开发对应的接口
////// 第二步，实现定义的接口
////type IProduct interface {
////	//连接数据
////	Conn() error
////	Insert(user *datamodels.Product) (int64, error)
////	Delete(int64) bool
////	Update(*datamodels.Product) error
////	SelectByKey(int64) (*datamodels.Product, error)
////	SelectAll() ([]*datamodels.Product, error)
////}
////
////type ProductManager struct {
////	table     string
////	mysqlConn *sql.DB
////}
////
////func NewProductManager(table string, db *sql.DB) IProduct {
////	return &ProductManager{table: table, mysqlConn: db}
////}
////
////// 数据连接
////func (p *ProductManager) Conn() (err error) {
////	if p.mysqlConn == nil {
////		mysql, err := NewMysqlConn()
////		if err != nil {
////			return err
////		}
////		p.mysqlConn = mysql
////	}
////	if p.table == "" {
////		p.table = "product"
////	}
////	return
////}
////
////// 插入
////func (p *ProductManager) Insert(product *datamodels.Product) (productId int64, err error) {
////	//1.判断连接是否存在
////	if err = p.Conn(); err != nil {
////		return
////	}
////	//2.准备sql
////	sql := "INSERT product SET productName=?,productNum=?,productImage=?,productUrl=?"
////	stmt, errSql := p.mysqlConn.Prepare(sql)
////	if errSql != nil {
////		return 0, errSql
////	}
////	//3.传入参数
////	result, errStmt := stmt.Exec(product.ProductName, product.ProductNum, product.ProductImage, product.ProductUrl)
////	if errStmt != nil {
////		return 0, errStmt
////	}
////	return result.LastInsertId()
////}
////
////// 商品的删除
////// 注意这里的两个err
////func (p *ProductManager) Delete(productID int64) bool {
////	//1.判断连接是否存在
////	if err := p.Conn(); err != nil {
////		return false
////	}
////	sql := "delete from product where ID=?"
////	stmt, err := p.mysqlConn.Prepare(sql)
////	if err != nil {
////		return false
////	}
////	//这个err可能是因为sql语法错误、数据库连接问题、权限问题、sql注入攻击等
////	_, err = stmt.Exec(strconv.FormatInt(productID, 10))
////	if err != nil {
////		return false
////	}
////	//这个err代表sql语句没有正常执行，如参数错误
////	return true
////}
////
////// 商品的更新
////func (p *ProductManager) Update(product *datamodels.Product) error {
////	//1.判断连接是否存在
////	if err := p.Conn(); err != nil {
////		return err
////	}
////
////	sql := "Update product set productName=?,productNum=?,productImage=?,productUrl=? where ID=" + strconv.FormatInt(product.ID, 10)
////
////	stmt, err := p.mysqlConn.Prepare(sql)
////	if err != nil {
////		return err
////	}
////
////	_, err = stmt.Exec(product.ProductName, product.ProductNum, product.ProductImage, product.ProductUrl)
////	if err != nil {
////		return err
////	}
////	return nil
////}
////
////// 根据商品ID查询商品
////func (p *ProductManager) SelectByKey(productID int64) (productResult *datamodels.Product, err error) {
////	//1.判断连接是否存在
////	if err = p.Conn(); err != nil {
////		return &datamodels.Product{}, err
////	}
////	sql := "Select * from " + p.table + " where ID =" + strconv.FormatInt(productID, 10)
////	row, errRow := p.mysqlConn.Query(sql)
////	defer row.Close()
////	if errRow != nil {
////		return &datamodels.Product{}, errRow
////	}
////	result := GetResultRow(row)
////	if len(result) == 0 {
////		return &datamodels.Product{}, nil
////	}
////	productResult = &datamodels.Product{}
////	DataToStructByTagSql(result, productResult)
////	return
////
////}
////
////// 获取所有商品
////func (p *ProductManager) SelectAll() (productArray []*datamodels.Product, errProduct error) {
////	//1.判断连接是否存在
////	if err := p.Conn(); err != nil {
////		return nil, err
////	}
////	sql := "Select * from " + p.table
////	rows, err := p.mysqlConn.Query(sql)
////	defer rows.Close()
////	if err != nil {
////		return nil, err
////	}
////
////	result := GetResultRows(rows)
////	if len(result) == 0 {
////		return nil, nil
////	}
////
////	for _, v := range result {
////		product := &datamodels.Product{}
////		DataToStructByTagSql(v, product)
////		productArray = append(productArray, product)
////	}
////	return
////}
////
////func PrintProduct(product *datamodels.Product) {
////	fmt.Println(product.ID, product.ProductName, product.ProductNum)
////}
////
////func PrintProductArray(productArray []*datamodels.Product) {
////	for i := 0; i < len(productArray); i++ {
////		PrintProduct(productArray[i])
////	}
////}
////
////func main() {
////	db, err := NewMysqlConn()
////	if err != nil {
////		fmt.Println(err)
////	}
////	//创建product数据库操作实例
////	product := NewProductManager("product", db)
////	res, err := product.SelectAll()
////	if err != nil {
////		fmt.Println(err)
////	}
////	PrintProductArray(res)
////
////}
////
////// 根据结构体中sql标签映射数据到结构体中并且转换类型
////func DataToStructByTagSql(data map[string]string, obj interface{}) {
////	objValue := reflect.ValueOf(obj).Elem()
////	fmt.Println("objValue", objValue)
////	for i := 0; i < objValue.NumField(); i++ {
////		//获取sql对应的值
////		value := data[objValue.Type().Field(i).Tag.Get("sql")]
////		//获取对应字段的名称
////		name := objValue.Type().Field(i).Name
////		//获取对应字段类型
////		structFieldType := objValue.Field(i).Type()
////		//获取变量类型，也可以直接写"string类型"
////		val := reflect.ValueOf(value)
////		fmt.Println(i, value, name, structFieldType, val)
////		var err error
////		if structFieldType != val.Type() {
////			//类型转换
////			val, err = TypeConversion(value, structFieldType.Name()) //类型转换
////			if err != nil {
////
////			}
////		}
////		//设置类型值
////		objValue.FieldByName(name).Set(val)
////	}
////}
////
////// 类型转换
////func TypeConversion(value string, ntype string) (reflect.Value, error) {
////	if ntype == "string" {
////		return reflect.ValueOf(value), nil
////	} else if ntype == "time.Time" {
////		t, err := time.ParseInLocation("2006-01-02 15:04:05", value, time.Local)
////		return reflect.ValueOf(t), err
////	} else if ntype == "Time" {
////		t, err := time.ParseInLocation("2006-01-02 15:04:05", value, time.Local)
////		return reflect.ValueOf(t), err
////	} else if ntype == "int" {
////		i, err := strconv.Atoi(value)
////		return reflect.ValueOf(i), err
////	} else if ntype == "int8" {
////		i, err := strconv.ParseInt(value, 10, 64)
////		return reflect.ValueOf(int8(i)), err
////	} else if ntype == "int32" {
////		i, err := strconv.ParseInt(value, 10, 64)
////		return reflect.ValueOf(int64(i)), err
////	} else if ntype == "int64" {
////		i, err := strconv.ParseInt(value, 10, 64)
////		return reflect.ValueOf(i), err
////	} else if ntype == "float32" {
////		i, err := strconv.ParseFloat(value, 64)
////		return reflect.ValueOf(float32(i)), err
////	} else if ntype == "float64" {
////		i, err := strconv.ParseFloat(value, 64)
////		return reflect.ValueOf(i), err
////	}
////
////	//else if .......增加其他一些类型的转换
////
////	return reflect.ValueOf(value), errors.New("未知的类型：" + ntype)}
