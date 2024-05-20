package cacheMall

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/kataras/iris/v12"
	"product/common"
	"product/datamodels"
	"product/fronted/middlerware"
	"product/rabbitmq"
	"product/repositories"
	"strconv"
	"sync"
	"time"
)

// PreHotStock 从MySQL中获取所有产品的库存信息，并将其存储到Redis中
// 修改：不要一次获取所有商品，而是输入一个商品名称列表
// 对于每个i
// 先访问缓存，再访问数据库，使用
func PreHotStock(ctx iris.Context) {
	// 连接Redis
	client, err := common.NewRedisClient()
	if err != nil {
		fmt.Println(err)
		return
	}

	// 连接MySQL
	db, err := common.NewMysqlConn()
	if err != nil {
		fmt.Println(err)
		return
	}

	// 获取所有产品
	productRepository := repositories.NewProductManager("product", db)
	productArray, err := productRepository.SelectAll()
	if err != nil {
		fmt.Println(err)
		return
	}

	// 将每个产品的库存信息存储到Redis中
	for _, product := range productArray {
		setOneStock(ctx, client, product.ID, product.ProductNum, 3*time.Hour)
	}
}

// SetStock 将产品库存信息存储到Redis中
// 需要加上分布式锁
func setOneStock(ctx iris.Context, redisClient *redis.Client, productID int64, stockNum int64, expiration time.Duration) error {
	redisKey := "stock_ProductID:" + strconv.FormatInt(productID, 10)
	// 库存信息转换为字符串
	stockStr := strconv.Itoa(int(stockNum))
	// 存储到Redis中，设置过期时间为3小时
	err := redisClient.Set(ctx, redisKey, stockStr, expiration).Err()
	if err != nil {
		return err
	}
	return nil
}

// GetStock 从Redis中获取产品库存信息
func GetStock(ctx iris.Context, redisClient *redis.Client, productID int64) (string, error) {

	return getStock(ctx, redisClient, productID)
}
func getStock(ctx iris.Context, redisClient *redis.Client, productID int64) (string, error) {
	redisKey := "stock_ProductID:" + strconv.FormatInt(productID, 10)
	// 从Redis中获取库存信息
	stockStr, err := redisClient.Get(ctx, redisKey).Result()
	if err == redis.Nil {
		// 如果缓存中不存在，则返回空字符串
		return "", nil
	} else if err != nil {
		return "", err
	}
	return stockStr, nil
}

// decrementStock 减少产品库存数量
func decrementStock(ctx iris.Context, redisClient *redis.Client, productID int64) (int64, error) {
	redisKey := "stock_ProductID:" + strconv.FormatInt(productID, 10)
	// 从Redis中获取库存信息
	stockNum, err := redisClient.Decr(ctx, redisKey).Result()
	if err != nil {
		return 0, err
	}
	return stockNum, nil
}

// isFirstPurchase 检查用户是否是首次购买某商品
func isFirstPurchase(ctx iris.Context, redisClient *redis.Client, userID, productID int64) (bool, error) {
	// 构建购买商品的 Redis 键名
	key := "buy_ProductID:" + strconv.FormatInt(productID, 10)
	// 检查 Redis 中是否存在该键
	exists, err := redisClient.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}

	// 如果不存在该键，说明没有人购买过该商品，返回 true
	if exists == 0 {
		return true, nil
	}

	// 检查用户是否购买过该商品
	isMember, err := redisClient.SIsMember(ctx, key, userID).Result()
	if err != nil {
		return false, err
	}

	// 如果用户已购买过该商品，返回 false，否则返回 true
	return !isMember, nil
}

// addUserToPurchaseRecord 将用户添加到购买商品记录中
func addUserToPurchaseRecord(ctx iris.Context, redisClient *redis.Client, userID, productID int64) error {
	// 构建购买商品的 Redis 键名
	key := "buy_ProductID:" + strconv.FormatInt(productID, 10)
	// 将用户ID添加到购买记录中
	_, err := redisClient.SAdd(ctx, key, userID).Result()
	if err != nil {
		return err
	}
	return nil
}

func BuyProductWithLock(ctx iris.Context, redisClient *redis.Client, userID, productID int64) (bool, error) {
	// 申请分布式锁
	lock := middlerware.NewRedisLock("buy_lock", "lock_value", 5*time.Second, productID, redisClient)
	defer lock.UnLock() // 确保在函数结束时释放锁

	// 尝试获取锁
	locked, err := lock.TryLock(100 * time.Millisecond)
	if err != nil {
		return false, err
	}
	if !locked {
		// 获取锁失败，可能有其他用户在购买该商品，返回false
		return false, nil
	}

	// 获取库存数量
	stockStr, err := getStock(ctx, redisClient, productID)
	if err != nil {
		return false, err
	}
	stockNum, err := strconv.Atoi(stockStr)
	if err != nil {
		return false, err
	}

	if stockNum <= 0 {
		// 库存为0，无法购买，返回false
		return false, nil
	}

	// 减少库存
	_, err = decrementStock(ctx, redisClient, productID)
	if err != nil {
		return false, err
	}

	// 将用户添加到购买记录中
	err = addUserToPurchaseRecord(ctx, redisClient, userID, productID)
	if err != nil {
		// 如果添加失败，需要恢复库存
		//_, _ = IncrementStock(ctx, redisClient, productID)
		return false, err
	}

	// 购买成功，返回true
	return true, nil
}

func BuyByProductID(ctx iris.Context, redisClient *redis.Client, userID, productID int64) (bool, error) {
	// 检查是否是首次购买
	firstPurchase, err := isFirstPurchase(ctx, redisClient, userID, productID)
	if err != nil {
		return false, err
	}
	if !firstPurchase {
		// 不是首次购买，直接返回false
		return false, nil
	}

	buyResult, err := BuyByProductID(ctx, redisClient, userID, productID)
	if err != nil {
		return false, err
	}
	if !buyResult {
		fmt.Println(userID, "已购买")
		return false, nil
	}

	// 获取库存数量,如果已经小于等于0，返回false
	//由于扣减库存 是原子操作，所以获取库存不需要分布式锁
	stockStr, err := getStock(ctx, redisClient, productID)
	if err != nil {
		return false, err
	}
	stockNum, err := strconv.Atoi(stockStr)
	if err != nil {
		return false, err
	}
	if stockNum <= 0 {
		fmt.Println("没有库存了")
		return false, err
	}

	// 将用户添加到购买记录中
	err = addUserToPurchaseRecord(ctx, redisClient, userID, productID)
	if err != nil {
		// 如果添加失败，需要恢复库存
		//_, _ = IncrementStock(ctx, redisClient, productID)
		return false, err
	}

	// 发送消息
	//3.创建消息体
	message := datamodels.NewMessage(userID, productID)
	//类型转化
	byteMessage, err := json.Marshal(message)
	if err != nil {
		fmt.Println(err)
		return false, err
	}

	//4.生产消息
	rabbitMqValidate := rabbitmq.NewRabbitMQSimple("imoocProduct")
	defer rabbitMqValidate.Destory()
	err = rabbitMqValidate.PublishSimple(string((byteMessage)))
	if err != nil {
		fmt.Println(err)
		return false, err
	}

	// 购买成功，返回true
	return true, nil
}

// 购买多件商品
func BuyMultipleProducts(ctx iris.Context, redisClient *redis.Client, userID int64, productIDs []int64) (int, error) {
	var wg sync.WaitGroup
	successCount := 0

	wg.Add(len(productIDs))
	for _, productID := range productIDs {

		go func(pid int64) {
			defer wg.Done()
			success, err := BuyByProductID(ctx, redisClient, userID, pid)
			if err != nil {
				fmt.Println("Error while buying product:", err)
				return
			}
			if success {
				successCount++
			}
		}(productID)
	}

	wg.Wait()

	return successCount, nil
}
