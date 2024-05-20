package cacheMall

import (
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"github.com/kataras/iris/v12"
	"product/common"
	"product/datamodels"
	"strconv"
)

//改成hash
func GetProductFromRedis(ctx iris.Context, redisClient *redis.Client, productID int64) (*datamodels.Product, error) {

	redisKey := "product:" + strconv.FormatInt(productID, 10)
	cachedProductJSON, err := common.GetRedisString(ctx, redisClient, redisKey)
	if err != nil {
		return nil, err
	}
	if cachedProductJSON != "" {
		productResult := &datamodels.Product{}
		err = json.Unmarshal([]byte(cachedProductJSON), productResult) //json序列化和反序列化
		if err != nil {
			return nil, err
		}
		return productResult, nil
	}
	return nil, nil
}
func SetProductToRedis(ctx iris.Context, redisClient *redis.Client, productID int64, productResult *datamodels.Product) error {
	redisKey := "product:" + strconv.FormatInt(productID, 10)
	productJSON, err := json.Marshal(productResult)
	if err != nil {
		return err
	}
	err = common.SetRedisString(ctx, redisClient, redisKey, string(productJSON))
	if err != nil {
		return err
	}
	return nil
}
