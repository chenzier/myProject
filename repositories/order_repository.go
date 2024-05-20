package repositories

import (
	"database/sql"
	"product/common"
	"product/datamodels"
	"strconv"
)

type IOrderRepository interface {
	Conn() error
	Insert(*datamodels.Order) (int64, error)
	Delete(int64) bool
	Update(*datamodels.Order) error
	SelectByKey(int64) (*datamodels.Order, error)
	SelectAll() ([]*datamodels.Order, error)
	SelectAllWithInfo() (map[int]map[string]string, error)
}
type OrderMangerRepository struct {
	table     string
	mysqlConn *sql.DB
}

func NewOrderMangerRepository(table string, sql *sql.DB) IOrderRepository {
	return &OrderMangerRepository{table: table, mysqlConn: sql}
}
func (o *OrderMangerRepository) Conn() error {
	if o.mysqlConn == nil {
		mysql, err := common.NewMysqlConn()
		if err != nil {
			return err
		}
		o.mysqlConn = mysql
	}
	if o.table == "" {
		o.table = "order"
	}
	return nil
}

func (o *OrderMangerRepository) Insert(order *datamodels.Order) (productID int64, err error) {
	if err = o.Conn(); err != nil {
		return
	}

	sql := "INSERT `order` SET userID=?,productID=?,orderStatus=?"

	stmt, errStmt := o.mysqlConn.Prepare(sql)
	if errStmt != nil {
		return productID, errStmt
	}

	result, errResult := stmt.Exec(order.UserID, order.ProductID, order.OrderStatus)
	if errResult != nil {
		return productID, errResult
	}
	return result.LastInsertId()
}

func (o *OrderMangerRepository) Delete(productID int64) (isOk bool) {
	if err := o.Conn(); err != nil {
		return
	}
	sql := "delete from " + o.table + " where ID=?"
	stmt, errStmt := o.mysqlConn.Prepare(sql)
	if errStmt != nil {
		return
	}
	_, err := stmt.Exec(productID)
	if err != nil {
		return
	}
	return true
}

func (o *OrderMangerRepository) Update(order *datamodels.Order) (err error) {
	if errConn := o.Conn(); errConn != nil {
		return errConn
	}

	sql := "Update " + o.table + " ser userID=?,productID=?,orderStatus=? Where ID=" + strconv.FormatInt(order.ID, 10)
	// strconv.FormatInt(order.ID, 10) 是 Go 语言中用于将整数转换为字符串的函数调用
	//strconv.FormatInt() 是 Go 语言标准库 strconv 中的一个函数，用于将整数转换为指定进制的字符串表示形式。
	//在这个函数调用中，第一个参数是要转换的整数值，即 order.ID，第二个参数是要转换的进制，这里是 10，表示使用十进制。
	//函数会返回转换后的字符串表示形式，即订单的 ID 的十进制字符串表示。

	stmt, errStmt := o.mysqlConn.Prepare(sql)
	if errStmt != nil {
		return
	}
	_, err = stmt.Exec(order.UserID, order.ProductID, order.OrderStatus)
	if err != nil {
		return err
	}
	return
}

// 在 SelectAll()  ，最后的 return 语句没有跟随具体的返回值，这种语法称为裸返回语句 。
// 在 Go 语言中，如果一个函数的返回值已经在函数签名中明确声明了，那么在函数体中的 return 语句可以省略返回值列表，此时会将当前函数的所有返回值都作为返回值返回。
// 因此，在函数中最后的 return 语句没有指定具体的返回值，它会将 orderArray 和 errOrder 作为返回值返回，即使它们之前没有在 return 语句中显式地指定。这样做的好处是可以让代码更加简洁，减少重复。
func (o *OrderMangerRepository) SelectByKey(orderID int64) (order *datamodels.Order, err error) {
	if errConn := o.Conn(); errConn != nil {
		return &datamodels.Order{}, errConn
	}
	sql := "Select * From" + o.table + " Where ID=" + strconv.FormatInt(order.ID, 10)
	row, errRow := o.mysqlConn.Query(sql)
	//对于只涉及查询的操作，使用 Query 方法即可
	//这里就不用占位符了
	if errRow != nil {
		return &datamodels.Order{}, errRow
	}

	result := common.GetResultRow(row)
	if len(result) == 0 {
		return &datamodels.Order{}, err
	}
	order = &datamodels.Order{}
	common.DataToStructByTagSql(result, order)
	return
}

func (o *OrderMangerRepository) SelectAll() (orderArray []*datamodels.Order, errOrder error) {
	//1.判断连接是否存在
	if err := o.Conn(); err != nil {
		return nil, err
	}
	sql := "Select * from " + o.table
	rows, errRows := o.mysqlConn.Query(sql)
	defer rows.Close()
	if errRows != nil {
		return nil, errRows
	}

	result := common.GetResultRows(rows)
	if len(result) == 0 {
		return nil, nil
	}

	for _, v := range result {
		order := &datamodels.Order{}
		common.DataToStructByTagSql(v, order)
		orderArray = append(orderArray, order)
	}
	return
}

func (o *OrderMangerRepository) SelectAllWithInfo() (OrderMap map[int]map[string]string, err error) { //查询订单及其相关商品的信息
	if errConn := o.Conn(); errConn != nil {
		return nil, errConn
	}

	//sql := "Select o.ID,p.productName,o.orderStatus From imooc.order as o left join product as p on o.productID=p.ID"
	//构造 SQL 查询语句： 函数构造了一个 SQL 查询语句，该语句使用了 LEFT JOIN 将订单表（imooc.order）和产品表（product）连接起来，以获取订单的详细信息。
	//查询语句选择了订单的 ID、产品名称以及订单状态。
	//fmt.Println(sql)
	sql := "SELECT o.ID, p.productName, o.orderStatus, u.UserName FROM imooc.order AS o LEFT JOIN product AS p ON o.productID = p.ID LEFT JOIN user AS u ON o.userID = u.ID"
	//SELECT o.ID, p.productName, o.orderStatus, u.UserName
	//FROM imooc.order o
	//LEFT JOIN product p ON o.productID = p.ID
	//LEFT JOIN user u ON o.userID = u.ID;

	rows, errRows := o.mysqlConn.Query(sql)
	defer rows.Close()
	if errRows != nil {
		return nil, errRows
	}

	return common.GetResultRows(rows), err

}
