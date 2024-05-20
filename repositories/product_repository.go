package repositories

import (
	"database/sql"
	"fmt"
	"github.com/kataras/iris/v12"
	"product/cacheMall"
	"product/common"
	"product/datamodels"
	"strconv"
)

// 第一步，先开发对应的接口
// 第二步，实现定义的接口
type IProduct interface {
	//连接数据
	Conn() error
	Insert(*datamodels.Product) (int64, error)
	Delete(int64) bool
	Update(*datamodels.Product) error
	SelectByKey(int64) (*datamodels.Product, error)
	SelectAll() ([]*datamodels.Product, error)
	SubProductNum(productID int64) error
	SelectByKey1(ctx iris.Context, productID int64) (*datamodels.Product, error)
}

type ProductManager struct {
	table     string
	mysqlConn *sql.DB
}

func NewProductManager(table string, db *sql.DB) IProduct {
	return &ProductManager{table: table, mysqlConn: db}
}

// 数据连接
func (p *ProductManager) Conn() (err error) {
	if p.mysqlConn == nil {
		mysql, err := common.NewMysqlConn()
		if err != nil {
			return err
		}
		p.mysqlConn = mysql
	}
	if p.table == "" {
		p.table = "product"
	}
	return
}

// 插入
func (p *ProductManager) Insert(product *datamodels.Product) (productId int64, err error) {
	//1.判断连接是否存在
	if err = p.Conn(); err != nil {
		return
	}
	//2.准备sql
	sql := "INSERT product SET productName=?,productNum=?,productImage=?,productUrl=?"
	stmt, errSql := p.mysqlConn.Prepare(sql)
	if errSql != nil {
		return 0, errSql
	}
	//3.传入参数
	result, errStmt := stmt.Exec(product.ProductName, product.ProductNum, product.ProductImage, product.ProductUrl)
	if errStmt != nil {
		return 0, errStmt
	}
	return result.LastInsertId()
}

// 商品的删除
// 注意这里的两个err
func (p *ProductManager) Delete(productID int64) bool {
	//1.判断连接是否存在
	if err := p.Conn(); err != nil {
		return false
	}
	sql := "delete from product where ID=?"
	stmt, err := p.mysqlConn.Prepare(sql)
	if err != nil {
		return false
	}
	//这个err可能是因为sql语法错误、数据库连接问题、权限问题、sql注入攻击等
	_, err = stmt.Exec(strconv.FormatInt(productID, 10))
	if err != nil {
		return false
	}
	//这个err代表sql语句没有正常执行，如参数错误
	return true
}

// 商品的更新
func (p *ProductManager) Update(product *datamodels.Product) error {
	//1.判断连接是否存在
	if err := p.Conn(); err != nil {
		return err
	}

	sql := "Update product set productName=?,productNum=?,productImage=?,productUrl=? where ID=" + strconv.FormatInt(product.ID, 10)

	stmt, err := p.mysqlConn.Prepare(sql)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(product.ProductName, product.ProductNum, product.ProductImage, product.ProductUrl)
	if err != nil {
		return err
	}
	return nil
}

// 根据商品ID查询商品
func (p *ProductManager) SelectByKey1(ctx iris.Context, productID int64) (productResult *datamodels.Product, err error) {
	// 先尝试从 Redis 中获取数据
	redisClient, _ := common.NewRedisClient()
	productCache, errCache := cacheMall.GetProductFromRedis(ctx, redisClient, productID)
	if errCache != nil {
		return &datamodels.Product{}, errCache
	}
	if productCache != nil {
		fmt.Println("从缓存读出")
		return productCache, nil
	}
	fmt.Println("缓存中没有有数据！")
	// 如果 Redis 中没有缓存数据，则从数据库中获取数据
	if err = p.Conn(); err != nil {
		return &datamodels.Product{}, err
	}
	sql := "Select * from " + p.table + " where ID =" + strconv.FormatInt(productID, 10)
	row, errRow := p.mysqlConn.Query(sql)
	defer row.Close()
	if errRow != nil {
		return &datamodels.Product{}, errRow
	}
	result := common.GetResultRow(row)
	if len(result) == 0 {
		return &datamodels.Product{}, nil
	}
	productResult = &datamodels.Product{}
	common.DataToStructByTagSql(result, productResult)
	if productCache == nil {
		err = cacheMall.SetProductToRedis(ctx, redisClient, productID, productCache)
		if err != nil {
			fmt.Println("写入缓存未实现")
		}
	}
	return productResult, nil
}

func (p *ProductManager) SelectByKey(productID int64) (productResult *datamodels.Product, err error) {
	//1.判断连接是否存在
	if err = p.Conn(); err != nil {
		return &datamodels.Product{}, err
	}
	sql := "Select * from " + p.table + " where ID =" + strconv.FormatInt(productID, 10)
	row, errRow := p.mysqlConn.Query(sql)
	defer row.Close()
	if errRow != nil {
		return &datamodels.Product{}, errRow
	}
	result := common.GetResultRow(row)
	if len(result) == 0 {
		return &datamodels.Product{}, nil
	}
	productResult = &datamodels.Product{}
	common.DataToStructByTagSql(result, productResult)
	return

}

// 获取所有商品
func (p *ProductManager) SelectAll() (productArray []*datamodels.Product, errProduct error) {
	//1.判断连接是否存在
	if err := p.Conn(); err != nil {
		return nil, err
	}
	sql := "Select * from " + p.table
	rows, err := p.mysqlConn.Query(sql)
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	result := common.GetResultRows(rows)
	if len(result) == 0 {
		return nil, nil
	}

	for _, v := range result {
		product := &datamodels.Product{}
		common.DataToStructByTagSql(v, product)
		productArray = append(productArray, product)
	}
	return
}

func (p *ProductManager) SubProductNum(productID int64) error {
	if err := p.Conn(); err != nil {
		return err
	}
	sql := "update " + p.table + " set " + " productNum=productNum-1 where ID=" + strconv.FormatInt(productID, 10)
	stmt, err := p.mysqlConn.Prepare(sql)
	if err != nil {
		return err
	}
	_, err = stmt.Exec()
	return err
}
