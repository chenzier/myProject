package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"product/common"
	"product/datamodels"
	"strconv"
	"sync"
)

type ICartRepository interface {
	Conn() error
	AddItem(item *datamodels.CartItem) (int64, error)
	Delete(int64) bool
	Update(item *datamodels.CartItem) error
}
type CartMangerRepository struct {
	table     string
	mysqlConn *sql.DB
}

func NewCartMangerRepository(table string, sql *sql.DB) ICartRepository {
	return &CartMangerRepository{table: table, mysqlConn: sql}
}
func (c *CartMangerRepository) Conn() error {
	if c.mysqlConn == nil {
		mysql, err := common.NewMysqlConn()
		if err != nil {
			return err
		}
		c.mysqlConn = mysql
	}
	if c.table == "" {
		c.table = "ShoppingCart"
	}
	return nil
}

// 从mysql中查询所有的商品item
func (c *CartMangerRepository) getFromMysql(userID int64) (*sync.Map, error) {
	// 连接数据库
	if c.mysqlConn == nil {
		return nil, errors.New("MySQL connection is not initialized")
	}

	// 准备 SQL 查询语句
	query := "SELECT productID, productNum FROM " + c.table + " WHERE userID = ?"

	rows, err := c.mysqlConn.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// 创建 sync.Map 用于存储结果
	result := &sync.Map{}

	// 迭代查询结果，将数据存储到 sync.Map 中
	for rows.Next() {
		var productID, productNum int
		if err := rows.Scan(&productID, &productNum); err != nil {
			return nil, err
		}
		result.Store(productID, productNum)
	}

	// 检查是否有错误发生在迭代过程中
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

// 将数据存入redis中
func (c *CartMangerRepository) loadCartFromMysql(userID int64) error {

	redisClient, _ := common.NewRedisClient()

	//新建redis列
	ctx := context.Background()

	//从mysql中导入
	cartMap, err := c.getFromMysql(userID)
	if err != nil {
		fmt.Println(err)
	}

	redisMap := make(map[string]interface{})
	cartMap.Range(func(key, value interface{}) bool {
		redisMap[fmt.Sprintf("%v", key)] = value
		return true
	})

	// 将购物车数据存储到 Redis 中
	redisKey := "stock_ProductID:" + strconv.FormatInt(userID, 10)
	err = redisClient.HSet(ctx, redisKey, redisMap).Err()
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

// 根据uid返回购物车结构体
func (c *CartMangerRepository) GetUserCart(userID int64) (*datamodels.Cart, error) {
	// 连接 Redis
	redisClient, _ := common.NewRedisClient()
	ctx := context.Background()

	// 构建 Redis key
	key := "stock_ProductID:" + strconv.FormatInt(userID, 10)

	// 从 Redis 中读取缓存
	cartMap := &sync.Map{}
	val, err := redisClient.HGetAll(ctx, key).Result()
	if err != nil {
		// 如果读取不到缓存，则从 MySQL 中加载
		if err == redis.Nil {
			err = c.loadCartFromMysql(userID)
			if err != nil {
				fmt.Println("Failed to load cart from MySQL:", err)
				return nil, err
			}
		} else {
			fmt.Println("Failed to get cart from Redis:", err)
			return nil, err
		}
	}

	// 从 Redis 中读取到缓存(上一步以确保一定能)，则构建 sync.Map
	for k, v := range val {
		cartMap.Store(k, v)
	}

	// 返回构建的 Cart 结构体
	return &datamodels.Cart{
		UserID:  userID,
		CartMap: cartMap,
	}, nil
}

// 增加一件商品
func (c *CartMangerRepository) AddCartItem(userID int64, productID int64, productNum int) error {
	// 连接 Redis
	redisClient, _ := common.NewRedisClient()
	ctx := context.Background()

	// 构建 Redis key
	key := "stock_ProductID:" + strconv.FormatInt(userID, 10)

	// 访问 Redis，增加一个键值对
	err := redisClient.HSet(ctx, key, productID, productNum).Err()
	if err != nil {
		fmt.Println("Failed to add item to Redis cart:", err)
		return err
	}

	// 访问 MySQL，增加一条记录
	query := "INSERT INTO ShoppingCart (userID, productID, productNum) VALUES (?, ?, ?)"
	_, err = c.mysqlConn.Exec(query, userID, productID, productNum)
	if err != nil {
		fmt.Println("Failed to add item to MySQL cart:", err)
		return err
	}

	return nil
}

func (c *CartMangerRepository) RemoveCartItem(userID int64, productID int64) error {
	// 连接 Redis
	redisClient, _ := common.NewRedisClient()
	ctx := context.Background()

	// 构建 Redis key
	key := "stock_ProductID:" + strconv.FormatInt(userID, 10)

	// 删除 Redis 中的键值对
	strProductID := strconv.FormatInt(productID, 10)
	err := redisClient.HDel(ctx, key, strProductID).Err()
	if err != nil {
		fmt.Println("Failed to remove item from Redis cart:", err)
		return err
	}

	// 删除 MySQL 中的记录
	query := "DELETE FROM ShoppingCart WHERE userID = ? AND productID = ?"
	_, err = c.mysqlConn.Exec(query, userID, productID)
	if err != nil {
		fmt.Println("Failed to remove item from MySQL cart:", err)
		return err
	}

	return nil
}

func (c *CartMangerRepository) UpdateCartItem(userID int64, productID int64, newProductNum int) error {
	// 连接 Redis
	redisClient, _ := common.NewRedisClient()
	ctx := context.Background()

	// 构建 Redis key
	key := "stock_ProductID:" + strconv.FormatInt(userID, 10)

	// 更新 Redis 中的商品数量
	err := redisClient.HSet(ctx, key, productID, newProductNum).Err()
	if err != nil {
		fmt.Println("Failed to update item in Redis cart:", err)
		return err
	}

	// 更新 MySQL 中的商品数量
	query := "UPDATE ShoppingCart SET productNum = ? WHERE userID = ? AND productID = ?"
	_, err = c.mysqlConn.Exec(query, newProductNum, userID, productID)
	if err != nil {
		fmt.Println("Failed to update item in MySQL cart:", err)
		return err
	}

	return nil
}

// 统计购物车商品数量
func (c *CartMangerRepository) GetCartItemCount(userID int64) (int64, error) {
	// 连接 Redis
	redisClient, _ := common.NewRedisClient()
	ctx := context.Background()

	// 构建 Redis key
	key := "stock_ProductID:" + strconv.FormatInt(userID, 10)

	// 获取 Redis 哈希结构中的字段数量
	count, err := redisClient.HLen(ctx, key).Result()
	if err != nil {
		fmt.Println("Failed to get cart item count from Redis:", err)
		return 0, err
	}

	return count, nil
}
