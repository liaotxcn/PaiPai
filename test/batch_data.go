package main

// 批量数据插入

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"math/rand"
	"time"
)

// User对应数据库表结构
type User struct {
	ID        string    `gorm:"column:id;type:varchar(24);primaryKey"`
	Avatar    string    `gorm:"column:avatar;type:varchar(191);default:''"`
	Nickname  string    `gorm:"column:nickname;type:varchar(24);not null"`
	Phone     string    `gorm:"column:phone;type:varchar(20);not null"`
	Password  string    `gorm:"column:password;type:varchar(191)"`
	Status    int8      `gorm:"column:status;type:tinyint"`
	Sex       int8      `gorm:"column:sex;type:tinyint"`
	CreatedAt time.Time `gorm:"column:created_at;type:timestamp"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:timestamp"`
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}

func main() {
	// 初始化GORM连接
	db, err := initGormDB()
	if err != nil {
		log.Fatal("数据库连接失败:", err)
	}

	// 生成测试数据
	userCount := 100000 // 10万条数据
	users := generateUsers(userCount)

	start := time.Now()

	// 方法1: 使用CreateInBatches (最简单)
	if err := batchInsertWithCreateInBatches(db, users, 2000); err != nil {
		log.Fatal("批量插入失败:", err)
	}

	// 方法2: 使用原生SQL批量插入 (最快)
	// if err := batchInsertWithRawSQL(db, users, 3000); err != nil {
	// 	log.Fatal("原生SQL插入失败:", err)
	// }

	// 方法3：并发批量插入 (推荐)
	//if err := concurrentBatchInsert(db, users, 2500, 6); err != nil {
	//	log.Fatal("并发插入失败:", err)
	//}

	elapsed := time.Since(start)
	fmt.Printf("成功插入 %d 条用户数据\n", userCount)
	fmt.Printf("总耗时: %v\n", elapsed)
	fmt.Printf("平均速度: %.0f 条/秒\n", float64(userCount)/elapsed.Seconds())
}

// initGormDB 初始化GORM数据库连接
func initGormDB() (*gorm.DB, error) {
	dsn := "root:123456@tcp(localhost:3306)/paipai?charset=utf8mb4&parseTime=True&loc=Local"

	// 配置GORM日志级别，生产环境可以关闭或减少日志
	newLogger := logger.New(
		log.New(log.Writer(), "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Warn, // 只显示警告和错误
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
		},
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger:                                   newLogger,
		PrepareStmt:                              true, // 开启预处理语句
		SkipDefaultTransaction:                   true, // 跳过默认事务
		DisableForeignKeyConstraintWhenMigrating: true, // 迁移时禁用外键
		CreateBatchSize:                          1000, // 批量创建大小
	})
	if err != nil {
		return nil, err
	}

	// 获取底层sql.DB并配置连接池
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxOpenConns(30)
	sqlDB.SetMaxIdleConns(20)
	sqlDB.SetConnMaxLifetime(10 * time.Minute)
	sqlDB.SetConnMaxIdleTime(5 * time.Minute)

	return db, nil
}

// generateUsers 生成测试用户数据
func generateUsers(count int) []User {
	users := make([]User, count)
	now := time.Now()
	rand.New(rand.NewSource(time.Now().UnixNano())) // 初始化随机种子

	avatars := []string{
		"https://example.com/avatar1.jpg",
		"https://example.com/avatar2.png",
		"https://example.com/avatar3.webp",
		"https://example.com/default.png",
	}

	nicknamePrefixes := []string{"用户", "会员", "客户", "访客", "测试"}

	for i := 0; i < count; i++ {
		users[i] = User{
			ID:        generateShortUUID(),
			Avatar:    avatars[rand.Intn(len(avatars))],
			Nickname:  fmt.Sprintf("%s%d", nicknamePrefixes[rand.Intn(len(nicknamePrefixes))], 10000+i),
			Phone:     generatePhoneNumber(13000000000 + i),
			Password:  generatePasswordHash("123456"),
			Status:    int8(rand.Intn(2)),
			Sex:       int8(rand.Intn(3)),
			CreatedAt: now.Add(time.Duration(i) * time.Millisecond),
			UpdatedAt: now.Add(time.Duration(i) * time.Millisecond),
		}
	}
	return users
}

// generateShortUUID 生成24字符的短UUID
func generateShortUUID() string {
	uuidObj := uuid.New()
	return hex.EncodeToString(uuidObj[:16])
}

// generatePhoneNumber 生成手机号
func generatePhoneNumber(base int) string {
	phone := fmt.Sprintf("1%010d", base%10000000000)
	if len(phone) > 11 {
		return phone[:11]
	}
	return phone
}

// generatePasswordHash 生成密码哈希
func generatePasswordHash(password string) string {
	// 使用MD5生成固定32字符哈希
	hasher := md5.New()
	hasher.Write([]byte(password))
	return hex.EncodeToString(hasher.Sum(nil))
}

// batchInsertWithCreateInBatches 使用GORM的CreateInBatches方法
func batchInsertWithCreateInBatches(db *gorm.DB, users []User, batchSize int) error {
	// 先测试小批量数据
	if len(users) == 0 {
		return fmt.Errorf("用户数据为空")
	}

	fmt.Printf("开始插入，总共 %d 条数据，每批 %d 条\n", len(users), batchSize)

	// 临时调整MySQL设置优化性能
	if err := db.Exec("SET FOREIGN_KEY_CHECKS = 0").Error; err != nil {
		log.Printf("警告: 禁用外键检查失败: %v", err)
	}

	// 使用事务确保数据一致性
	tx := db.Begin()
	if tx.Error != nil {
		return fmt.Errorf("开始事务失败: %v", tx.Error)
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Printf("事务回滚 due to panic: %v", r)
		}
	}()

	// 使用CreateInBatches批量插入
	if err := tx.CreateInBatches(users, batchSize).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("批量插入失败: %v", err)
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("提交事务失败: %v", err)
	}

	// 恢复MySQL设置
	if err := db.Exec("SET FOREIGN_KEY_CHECKS = 1").Error; err != nil {
		log.Printf("警告: 启用外键检查失败: %v", err)
	}

	return nil
}

//batchInsertWithRawSQL 使用原生SQL进行批量插入 (最快)
//func batchInsertWithRawSQL(db *gorm.DB, users []User, batchSize int) error {
//	// 优化MySQL设置
//	if err := db.Exec("SET FOREIGN_KEY_CHECKS = 0").Error; err != nil {
//		return err
//	}
//	if err := db.Exec("SET UNIQUE_CHECKS = 0").Error; err != nil {
//		return err
//	}
//	if err := db.Exec("SET AUTOCOMMIT = 0").Error; err != nil {
//		return err
//	}
//
//	defer func() {
//		db.Exec("SET FOREIGN_KEY_CHECKS = 1")
//		db.Exec("SET UNIQUE_CHECKS = 1")
//		db.Exec("SET AUTOCOMMIT = 1")
//	}()
//
//	tx := db.Begin()
//	defer func() {
//		if r := recover(); r != nil {
//			tx.Rollback()
//		}
//	}()
//
//	total := len(users)
//	for start := 0; start < total; start += batchSize {
//		end := start + batchSize
//		if end > total {
//			end = total
//		}
//
//		batch := users[start:end]
//		if err := insertBatchWithRawSQL(tx, batch); err != nil {
//			tx.Rollback()
//			return fmt.Errorf("批次 %d-%d 插入失败: %v", start, end, err)
//		}
//
//		// 显示进度
//		if start%(batchSize*10) == 0 && start > 0 {
//			fmt.Printf("已处理: %d/%d (%.1f%%)\n", end, total, float64(end)*100/float64(total))
//		}
//	}
//
//	return tx.Commit().Error
//}
//
//// insertBatchWithRawSQL 使用原生SQL插入单个批次
//func insertBatchWithRawSQL(tx *gorm.DB, users []User) error {
//	if len(users) == 0 {
//		return nil
//	}
//
//	// 手动构建SQL语句
//	sql := "INSERT INTO users (id, avatar, nickname, phone, password, status, sex, created_at, updated_at) VALUES "
//	values := []interface{}{}
//
//	for i, user := range users {
//		if i > 0 {
//			sql += ","
//		}
//		sql += "(?, ?, ?, ?, ?, ?, ?, ?, ?)"
//		values = append(values,
//			user.ID, user.Avatar, user.Nickname, user.Phone, user.Password,
//			user.Status, user.Sex, user.CreatedAt, user.UpdatedAt,
//		)
//	}
//
//	sql += " ON DUPLICATE KEY UPDATE updated_at = VALUES(updated_at)"
//
//	return tx.Exec(sql, values...).Error
//}

// concurrentBatchInsert 并发批量插入
//func concurrentBatchInsert(db *gorm.DB, users []User, batchSize, numWorkers int) error {
//	// 优化MySQL设置
//	if err := db.Exec("SET FOREIGN_KEY_CHECKS = 0").Error; err != nil {
//		return err
//	}
//	if err := db.Exec("SET UNIQUE_CHECKS = 0").Error; err != nil {
//		return err
//	}
//
//	defer func() {
//		db.Exec("SET FOREIGN_KEY_CHECKS = 1")
//		db.Exec("SET UNIQUE_CHECKS = 1")
//	}()
//
//	var wg sync.WaitGroup
//	errCh := make(chan error, numWorkers)
//	userCh := make(chan []User, numWorkers*2)
//
//	// 启动worker
//	for i := 0; i < numWorkers; i++ {
//		wg.Add(1)
//		go func(workerID int) {
//			defer wg.Done()
//
//			// 每个worker使用独立的数据库会话
//			workerDB := db.Session(&gorm.Session{})
//			for batch := range userCh {
//				if err := insertBatchWithRawSQL(workerDB, batch); err != nil {
//					errCh <- fmt.Errorf("worker %d 插入失败: %v", workerID, err)
//					return
//				}
//			}
//		}(i)
//	}
//
//	// 分发数据
//	total := len(users)
//	for start := 0; start < total; start += batchSize {
//		end := start + batchSize
//		if end > total {
//			end = total
//		}
//		userCh <- users[start:end]
//
//		if start%(batchSize*10) == 0 {
//			fmt.Printf("已分发: %d/%d (%.1f%%)\n", end, total, float64(end)*100/float64(total))
//		}
//	}
//	close(userCh)
//
//	wg.Wait()
//	close(errCh)
//
//	// 检查错误
//	for err := range errCh {
//		if err != nil {
//			return err
//		}
//	}
//
//	return nil
//}

// 可选：使用GORM钩子自动生成ID和时间戳
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == "" {
		u.ID = generateShortUUID()
	}
	now := time.Now()
	if u.CreatedAt.IsZero() {
		u.CreatedAt = now
	}
	if u.UpdatedAt.IsZero() {
		u.UpdatedAt = now
	}
	return nil
}
