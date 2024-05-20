package others

//package main
//
//import (
//	"fmt"
//	"sync"
//	"time"
//)
//
//// 任务结构体
//type Task struct {
//	ID  int
//	Job func()
//}
//
//// 协程池结构体
//type Pool struct {
//	taskQueue chan Task      //一个任务队列 taskQueue
//	wg        sync.WaitGroup //一个 WaitGroup wg
//}
//
//// 创建协程池
//// NewPool 函数用于创建一个协程池，参数 numWorkers 指定了协程池中的工作协程数量。
//// 在 NewPool 函数中，会初始化 taskQueue 通道，并启动指定数量的工作协程。
//func NewPool(numWorkers int) *Pool {
//	p := &Pool{
//		taskQueue: make(chan Task),
//	}
//
//	p.wg.Add(numWorkers)
//	for i := 0; i < numWorkers; i++ {
//		//通过 go p.worker() 语句启动了指定数量的工作协程，这些工作协程会立即执行 worker 函数。
//		//在 worker 函数中，通过 for task := range p.taskQueue 循环
//		//工作协程会不断地从任务队列 taskQueue 中取出任务并执行。
//		//如果任务队列为空，则工作协程会阻塞在取任务的操作上，直到有新的任务到来或者任务队列被关闭
//		//因此初始时，这些协程会被阻塞
//		go p.worker(i)
//	}
//
//	return p
//}
//
//// 添加任务到协程池
//func (p *Pool) AddTask(task Task) {
//	p.taskQueue <- task
//}
//
//// 工作协程
//func (p *Pool) worker(workerID int) {
//	//对于每个协程，不断从taskQueue取得任务
//	for task := range p.taskQueue {
//		fmt.Printf("%d new start task %d\n", workerID, task.ID)
//		task.Job()
//		fmt.Printf("finished task %d\n", task.ID)
//	}
//	p.wg.Done()
//}
//
//// 等待所有任务完成
//func (p *Pool) Wait() {
//	close(p.taskQueue)
//	p.wg.Wait()
//}
//
//func main() {
//	// 创建一个协程池，设置工作协程数为3
//	pool := NewPool(10)
//
//	startTime := time.Now()
//	// 添加任务到协程池
//	for i := 0; i < 100; i++ {
//		taskID := i
//		task := Task{
//			ID: taskID,
//			Job: func() {
//				time.Sleep(time.Second)
//				fmt.Printf("Task %d is running\n", taskID)
//			},
//		}
//		pool.AddTask(task)
//	}
//
//	// 等待所有任务完成
//	pool.Wait()
//
//	endTime := time.Now()
//
//	// 计算处理时间并打印
//	processingTime := endTime.Sub(startTime)
//	fmt.Printf("Processing time: %v\n", processingTime)
//}
//
////package main
////
////import "fmt"
////
////type TreeNode struct {
////	Val   int
////	Left  *TreeNode
////	Right *TreeNode
////}
////
////func makeBinaryTree(rootVal []int) *TreeNode {
////	if len(rootVal) == 0 {
////		return nil
////	}
////	var makeNode func(index int) *TreeNode
////	makeNode = func(index int) *TreeNode {
////		if index >= len(rootVal) {
////			return nil
////		}
////		if rootVal[index] == -1 {
////			return nil
////		}
////		left := makeNode(index*2 + 1)
////		right := makeNode(index*2 + 2)
////		return &TreeNode{rootVal[index], left, right}
////	}
////	root := makeNode(0)
////	return root
////}
////func printTree(root *TreeNode) {
////	if root == nil {
////		fmt.Println()
////		return
////	}
////	res := []int{}
////	queue := []*TreeNode{}
////	queue = append(queue, root)
////	for len(queue) > 0 {
////		l := len(queue)
////		for i := 0; i < l; i++ {
////			node := queue[0]
////			queue = queue[1:]
////			res = append(res, node.Val)
////			if node.Left != nil {
////				queue = append(queue, node.Left)
////			}
////			if node.Right != nil {
////				queue = append(queue, node.Right)
////			}
////		}
////	}
////	fmt.Println(res)
////}
////func zigzagLevelOrder(root *TreeNode) [][]int {
////	res := [][]int{}
////	var dfs func(node *TreeNode, depth, flag int)
////	dfs = func(node *TreeNode, depth, flag int) {
////		if node == nil {
////			return
////		}
////		if depth == len(res) {
////			res = append(res, []int{})
////		}
////		if flag == -1 {
////			res[depth] = append(res[depth], node.Val)
////
////		} else {
////			res[depth] = append([]int{node.Val}, res[depth]...)
////		}
////		dfs(node.Left, depth+1, flag*-1)
////		dfs(node.Right, depth+1, flag*-1)
////	}
////	dfs(root, 0, -1)
////	return res
////}
////func main() {
////	rootVal := []int{3, 9, 20, -1, -1, 15, 7}
////	root := makeBinaryTree(rootVal)
////	printTree(root)
////	fmt.Println(zigzagLevelOrder(root))
////
////}
////
//////package main
//////
//////import (
//////	"fmt"
//////	"math/rand"
//////	//"fmt"
//////	"product/fronted/middlerware"
//////	"time"
//////)
//////
//////func main() {
//////	// 创建一个布隆过滤器
//////	n := 1000000
//////	p := 0.03125
//////	bf := middlerware.NewBloomFilter(n, p)
//////	bf.Print()
//////	// 生成 500 个随机产品ID
//////	rand.Seed(time.Now().UnixNano())
//////	productIDs := make([]int64, n)
//////	lenHalf := len(productIDs) / 2
//////	for i := 0; i < lenHalf*2; i++ {
//////		productIDs[i] = int64(i) // 假设产品ID在这个范围内
//////	}
//////
//////	for i := 0; i < lenHalf; i++ {
//////		bf.Add(productIDs[i])
//////	}
//////
//////	//// 测试命中率
//////	//hitCount := 0
//////	//for i := 0; i < lenHalf; i++ {
//////	//	if bf.Contains(productIDs[i]) {
//////	//		hitCount++
//////	//
//////	//	}
//////	//}
//////	//for i := 0; i < lenHalf; i++ {
//////	//	if !bf.Contains(productIDs[i+lenHalf]) {
//////	//		hitCount++
//////	//	}
//////	//}
//////
//////	// 测试命中率
//////	hitCount := 0
//////	for i := 0; i < 2*lenHalf; i++ {
//////		if !bf.Contains(productIDs[i] + int64(n)) {
//////			hitCount++
//////
//////		}
//////	}
//////
//////	hitRate := float64(hitCount) / float64(lenHalf*2) * 100.0
//////	fmt.Printf("命中率: %.2f%%\n", hitRate)
//////	fmt.Println(hitCount)
//////}
//////
////////package main
////////
////////import (
////////	"context"
////////	"fmt"
////////	"github.com/go-redis/redis/v8"
////////)
////////
////////func main() {
////////	// 创建 Redis 客户端连接
////////	client := redis.NewClient(&redis.Options{
////////		Addr:     "localhost:6379",
////////		Password: "123456",
////////		DB:       0,
////////	})
////////
////////	// 创建一个上下文对象
////////	ctx := context.Background()
////////
////////	// 设置一个键值对
////////	err := client.Set(ctx, "test_key", "test_value", 0).Err()
////////	if err != nil {
////////		fmt.Println("Failed to set key:", err)
////////		return
////////	}
////////
////////	// 查找键对应的值
////////	val, err := client.Get(ctx, "test_key").Result()
////////	if err != nil {
////////		fmt.Println("Failed to get key:", err)
////////		return
////////	}
////////
////////	fmt.Println("Value of 'test_key':", val)
////////}
////////
//////////package main
//////////
//////////import (
//////////	"database/sql"
//////////	"errors"
//////////	"fmt"
//////////	"product/datamodels"
//////////	"reflect"
//////////	"strconv"
//////////	"time"
//////////)
//////////import _ "github.com/go-sql-driver/mysql"
//////////
//////////// 创建mysql 连接
//////////func NewMysqlConn() (db *sql.DB, err error) {
//////////	db, err = sql.Open("mysql", "root:imooc@tcp(127.0.0.1:3306)/imooc?charset=utf8")
//////////	return
//////////}
//////////
//////////// 第一步，先开发对应的接口
//////////// 第二步，实现定义的接口
//////////type IProduct interface {
//////////	//连接数据
//////////	Conn() error
//////////	Insert(user *datamodels.Product) (int64, error)
//////////	Delete(int64) bool
//////////	Update(*datamodels.Product) error
//////////	SelectByKey(int64) (*datamodels.Product, error)
//////////	SelectAll() ([]*datamodels.Product, error)
//////////}
//////////
//////////type ProductManager struct {
//////////	table     string
//////////	mysqlConn *sql.DB
//////////}
//////////
//////////func NewProductManager(table string, db *sql.DB) IProduct {
//////////	return &ProductManager{table: table, mysqlConn: db}
//////////}
//////////
//////////// 数据连接
//////////func (p *ProductManager) Conn() (err error) {
//////////	if p.mysqlConn == nil {
//////////		mysql, err := NewMysqlConn()
//////////		if err != nil {
//////////			return err
//////////		}
//////////		p.mysqlConn = mysql
//////////	}
//////////	if p.table == "" {
//////////		p.table = "product"
//////////	}
//////////	return
//////////}
//////////
//////////// 插入
//////////func (p *ProductManager) Insert(product *datamodels.Product) (productId int64, err error) {
//////////	//1.判断连接是否存在
//////////	if err = p.Conn(); err != nil {
//////////		return
//////////	}
//////////	//2.准备sql
//////////	sql := "INSERT product SET productName=?,productNum=?,productImage=?,productUrl=?"
//////////	stmt, errSql := p.mysqlConn.Prepare(sql)
//////////	if errSql != nil {
//////////		return 0, errSql
//////////	}
//////////	//3.传入参数
//////////	result, errStmt := stmt.Exec(product.ProductName, product.ProductNum, product.ProductImage, product.ProductUrl)
//////////	if errStmt != nil {
//////////		return 0, errStmt
//////////	}
//////////	return result.LastInsertId()
//////////}
//////////
//////////// 商品的删除
//////////// 注意这里的两个err
//////////func (p *ProductManager) Delete(productID int64) bool {
//////////	//1.判断连接是否存在
//////////	if err := p.Conn(); err != nil {
//////////		return false
//////////	}
//////////	sql := "delete from product where ID=?"
//////////	stmt, err := p.mysqlConn.Prepare(sql)
//////////	if err != nil {
//////////		return false
//////////	}
//////////	//这个err可能是因为sql语法错误、数据库连接问题、权限问题、sql注入攻击等
//////////	_, err = stmt.Exec(strconv.FormatInt(productID, 10))
//////////	if err != nil {
//////////		return false
//////////	}
//////////	//这个err代表sql语句没有正常执行，如参数错误
//////////	return true
//////////}
//////////
//////////// 商品的更新
//////////func (p *ProductManager) Update(product *datamodels.Product) error {
//////////	//1.判断连接是否存在
//////////	if err := p.Conn(); err != nil {
//////////		return err
//////////	}
//////////
//////////	sql := "Update product set productName=?,productNum=?,productImage=?,productUrl=? where ID=" + strconv.FormatInt(product.ID, 10)
//////////
//////////	stmt, err := p.mysqlConn.Prepare(sql)
//////////	if err != nil {
//////////		return err
//////////	}
//////////
//////////	_, err = stmt.Exec(product.ProductName, product.ProductNum, product.ProductImage, product.ProductUrl)
//////////	if err != nil {
//////////		return err
//////////	}
//////////	return nil
//////////}
//////////
//////////// 根据商品ID查询商品
//////////func (p *ProductManager) SelectByKey(productID int64) (productResult *datamodels.Product, err error) {
//////////	//1.判断连接是否存在
//////////	if err = p.Conn(); err != nil {
//////////		return &datamodels.Product{}, err
//////////	}
//////////	sql := "Select * from " + p.table + " where ID =" + strconv.FormatInt(productID, 10)
//////////	row, errRow := p.mysqlConn.Query(sql)
//////////	fmt.Println("Type :", reflect.TypeOf(row))
//////////	fmt.Println("row", row)
//////////	defer row.Close()
//////////	if errRow != nil {
//////////		return &datamodels.Product{}, errRow
//////////	}
//////////	result := GetResultRow(row)
//////////	fmt.Println("result", result)
//////////	if len(result) == 0 {
//////////		return &datamodels.Product{}, nil
//////////	}
//////////	productResult = &datamodels.Product{}
//////////	DataToStructByTagSql(result, productResult)
//////////	return
//////////
//////////}
//////////
//////////// 获取所有商品
//////////func (p *ProductManager) SelectAll() (productArray []*datamodels.Product, errProduct error) {
//////////	//1.判断连接是否存在
//////////	if err := p.Conn(); err != nil {
//////////		return nil, err
//////////	}
//////////	sql := "Select * from " + p.table
//////////	rows, err := p.mysqlConn.Query(sql)
//////////	defer rows.Close()
//////////	if err != nil {
//////////		return nil, err
//////////	}
//////////
//////////	result := GetResultRows(rows)
//////////	if len(result) == 0 {
//////////		return nil, nil
//////////	}
//////////
//////////	for _, v := range result {
//////////		product := &datamodels.Product{}
//////////		DataToStructByTagSql(v, product)
//////////		productArray = append(productArray, product)
//////////	}
//////////	return
//////////}
//////////
//////////func PrintProduct(product *datamodels.Product) {
//////////	fmt.Println(product.ID, product.ProductName, product.ProductNum)
//////////}
//////////
//////////func PrintProductArray(productArray []*datamodels.Product) {
//////////	for i := 0; i < len(productArray); i++ {
//////////		PrintProduct(productArray[i])
//////////	}
//////////}
//////////
//////////func main() {
//////////	db, err := NewMysqlConn()
//////////	if err != nil {
//////////		fmt.Println(err)
//////////	}
//////////	//创建product数据库操作实例
//////////	product := NewProductManager("product", db)
//////////	aProduct, err := product.SelectByKey(1)
//////////	PrintProduct(aProduct)
//////////
//////////}
//////////
//////////// 获取返回值，获取一条
//////////func GetResultRow(rows *sql.Rows) map[string]string {
//////////	fmt.Println(456)
//////////	columns, _ := rows.Columns()
//////////	//fmt.Println(reflect.TypeOf(columns), columns)//[]string [ID productName productNum productImage productUrl]
//////////
//////////	scanArgs := make([]interface{}, len(columns))
//////////	values := make([][]byte, len(columns))
//////////	for j := range values {
//////////		scanArgs[j] = &values[j]
//////////	}
//////////	//fmt.Println(scanArgs, values)//scanArgs[j]存的是每个values[j]的地址，换句话说，scanArgs[j]是指针
//////////	record := make(map[string]string)
//////////	j := 233
//////////	for rows.Next() {
//////////		//rows.Next() 方法在调用时会检查是否还有未处理的行，如果有，它会将游标移动到下一行，并返回 true。
//////////		//如果已经没有更多的行了，它会返回 false，表示查询结果已经遍历完成。
//////////		j += 1
//////////		//将行数据保存到record字典
//////////		rows.Scan(scanArgs...) //将当前行的数据填充到 scanArgs 切片中指向的位置
//////////		fmt.Println(j, scanArgs, values)
//////////		for i, v := range values {
//////////			if v != nil {
//////////				//fmt.Println(reflect.TypeOf(col))
//////////				record[columns[i]] = string(v)
//////////			}
//////////		}
//////////	}
//////////	fmt.Println(789)
//////////	return record
//////////}
//////////
//////////// 获取所有
//////////func GetResultRows(rows *sql.Rows) map[int]map[string]string {
//////////	fmt.Println(321)
//////////	//返回所有列
//////////	columns, _ := rows.Columns() //调用 rows.Columns() 方法获取查询结果的列名，存储在 columns 变量中。
//////////	// 这个操作会返回一个字符串切片，包含了查询结果的列名。
//////////	vals := make([][]byte, len(columns)) //这里表示一行所有列的值，用[]byte表示
//////////	fmt.Println(vals)
//////////	scans := make([]interface{}, len(columns)) //这里表示一行填充数据
//////////	//这里scans引用vals，把数据填充到[]byte里
//////////	for k, _ := range vals {
//////////		scans[k] = &vals[k]
//////////	}
//////////	i := 0
//////////	fmt.Println(scans)
//////////	result := make(map[int]map[string]string)
//////////
//////////	for rows.Next() {
//////////		//rows.Next() 方法在调用时会检查是否还有未处理的行，如果有，它会将游标移动到下一行，并返回 true。
//////////		//如果已经没有更多的行了，它会返回 false，表示查询结果已经遍历完成。
//////////
//////////		//填充数据
//////////		rows.Scan(scans...)
//////////		//每行数据
//////////		row := make(map[string]string)
//////////		//把vals中的数据复制到row中
//////////		for k, v := range vals {
//////////			key := columns[k]
//////////			//这里把[]byte数据转成string
//////////			row[key] = string(v)
//////////		}
//////////		//放入结果集
//////////		result[i] = row
//////////		i++
//////////	}
//////////	fmt.Println(123)
//////////	return result
//////////}
//////////
//////////// 根据结构体中sql标签映射数据到结构体中并且转换类型
//////////func DataToStructByTagSql(data map[string]string, obj interface{}) {
//////////	objValue := reflect.ValueOf(obj).Elem()
//////////	for i := 0; i < objValue.NumField(); i++ {
//////////		//获取sql对应的值
//////////		value := data[objValue.Type().Field(i).Tag.Get("sql")]
//////////		//获取对应字段的名称
//////////		name := objValue.Type().Field(i).Name
//////////		//获取对应字段类型
//////////		structFieldType := objValue.Field(i).Type()
//////////		//获取变量类型，也可以直接写"string类型"
//////////		val := reflect.ValueOf(value)
//////////		var err error
//////////		if structFieldType != val.Type() {
//////////			//类型转换
//////////			val, err = TypeConversion(value, structFieldType.Name()) //类型转换
//////////			if err != nil {
//////////
//////////			}
//////////		}
//////////		//设置类型值
//////////		objValue.FieldByName(name).Set(val)
//////////	}
//////////}
//////////
//////////// 类型转换
//////////func TypeConversion(value string, ntype string) (reflect.Value, error) {
//////////	if ntype == "string" {
//////////		return reflect.ValueOf(value), nil
//////////	} else if ntype == "time.Time" {
//////////		t, err := time.ParseInLocation("2006-01-02 15:04:05", value, time.Local)
//////////		return reflect.ValueOf(t), err
//////////	} else if ntype == "Time" {
//////////		t, err := time.ParseInLocation("2006-01-02 15:04:05", value, time.Local)
//////////		return reflect.ValueOf(t), err
//////////	} else if ntype == "int" {
//////////		i, err := strconv.Atoi(value)
//////////		return reflect.ValueOf(i), err
//////////	} else if ntype == "int8" {
//////////		i, err := strconv.ParseInt(value, 10, 64)
//////////		return reflect.ValueOf(int8(i)), err
//////////	} else if ntype == "int32" {
//////////		i, err := strconv.ParseInt(value, 10, 64)
//////////		return reflect.ValueOf(int64(i)), err
//////////	} else if ntype == "int64" {
//////////		i, err := strconv.ParseInt(value, 10, 64)
//////////		return reflect.ValueOf(i), err
//////////	} else if ntype == "float32" {
//////////		i, err := strconv.ParseFloat(value, 64)
//////////		return reflect.ValueOf(float32(i)), err
//////////	} else if ntype == "float64" {
//////////		i, err := strconv.ParseFloat(value, 64)
//////////		return reflect.ValueOf(i), err
//////////	}
//////////
//////////	//else if .......增加其他一些类型的转换
//////////
//////////	return reflect.ValueOf(value), errors.New("未知的类型：" + ntype)
//////////}
//////////
////////////package main
////////////
////////////import "fmt"
////////////import "context"
////////////
////////////func main() {
////////////
////////////	a := context.Background()              // 创建上下文
////////////	b := context.WithValue(a, "k1", "v1")  // 塞入一个kv
////////////	c := context.WithValue(b, "k2", "v2")  // 塞入另外一个kv
////////////	d := context.WithValue(c, "k1", "vo1") // 覆盖一个kv
////////////
////////////	fmt.Printf("k1 of b: %s\n", b.Value("k1"))
////////////	fmt.Printf("k1 of d: %s\n", d.Value("k1"))
////////////	fmt.Printf("k2 of d: %s\n", d.Value("k2"))
////////////}
////////////
//////////////package main
//////////////
//////////////import "fmt"
//////////////
//////////////func main() {
//////////////	lru := Constructor(2)
//////////////	lru.Put(1, 1)
//////////////	lru.Put(2, 2)
//////////////	lru.Put(3, 3)
//////////////	fmt.Println(lru.Get(1))
//////////////	fmt.Println(lru.Get(2))
//////////////	lru.Put(4, 4)
//////////////	fmt.Println(lru.Get(3))
//////////////}
//////////////
//////////////type DoubleLinkedList struct {
//////////////	key, value int
//////////////	pre, next  *DoubleLinkedList
//////////////}
//////////////type LRUCache struct {
//////////////	cap        int
//////////////	hashmap    map[int]*DoubleLinkedList
//////////////	head, tail *DoubleLinkedList
//////////////}
//////////////
//////////////func Constructor(capacity int) LRUCache {
//////////////	head, tail := &DoubleLinkedList{}, &DoubleLinkedList{}
//////////////	head.next = tail
//////////////	tail.pre = head
//////////////	return LRUCache{
//////////////		cap:     capacity,
//////////////		hashmap: map[int]*DoubleLinkedList{},
//////////////		head:    head,
//////////////		tail:    tail,
//////////////	}
//////////////}
//////////////
//////////////func (this *LRUCache) Get(key int) int {
//////////////	if node, ok := this.hashmap[key]; ok {
//////////////		remove(node)
//////////////		insert2Head(this.head, node)
//////////////		return node.value
//////////////	}
//////////////	return -1
//////////////}
//////////////
//////////////func (this *LRUCache) Put(key int, value int) {
//////////////	if node, ok := this.hashmap[key]; ok {
//////////////		remove(node)
//////////////		insert2Head(this.head, node)
//////////////		node.value = value
//////////////	} else {
//////////////		node := &DoubleLinkedList{key: key, value: value}
//////////////		insert2Head(this.head, node)
//////////////		this.hashmap[key] = node
//////////////		if len(this.hashmap) > this.cap {
//////////////			tail := this.tail.pre//为什么？
//////////////			remove(tail)
//////////////			delete(this.hashmap, tail.key)
//////////////		}
//////////////	}
//////////////
//////////////}
//////////////
//////////////func remove(node *DoubleLinkedList) {
//////////////	node.pre.next, node.next.pre = node.next, node.pre
//////////////}
//////////////func insert2Head(head, node *DoubleLinkedList) {
//////////////	next := head.next
//////////////	head.next, node.pre = node, head
//////////////	node.next, next.pre = next, node
//////////////}
