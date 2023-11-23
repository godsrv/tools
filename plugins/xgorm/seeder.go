package xgorm

import (
	"context"

	"gorm.io/gorm"
)

// Seeder 对应每一个 database/seeders 目录下的 Seeder 文件
type Seeder struct {
	Ctx  context.Context
	Name string
	Func SeederFunc
}

// 存放所有 Seeder
var seeders []Seeder

// 按顺序执行的 Seeder 数组
// 支持一些必须按顺序执行的 seeder，例如 topic 创建的
// 时必须依赖于 user, 所以 TopicSeeder 应该在 UserSeeder 后执行
var orderedSeederNames []string

type SeederFunc func(*gorm.DB)

// Add 注册到 seeders 数组中
func AddSeeder(name string, fn SeederFunc) {
	seeders = append(seeders, Seeder{
		Name: name,
		Func: fn,
	})
}

// SetRunOrder 设置『按顺序执行的 Seeder 数组』
func SetRunOrder(names []string) {
	orderedSeederNames = names
}

// GetSeeder 通过名称来获取 Seeder 对象
func GetSeeder(name string) Seeder {
	for _, sdr := range seeders {
		if name == sdr.Name {
			return sdr
		}
	}
	return Seeder{}
}

// RunAll 运行所有 Seeder
func RunAll(db *gorm.DB) {
	// 先运行 ordered 的
	executed := make(map[string]string)
	for _, name := range orderedSeederNames {
		sdr := GetSeeder(name)
		if len(sdr.Name) > 0 {
			sdr.Func(db)
			executed[name] = name

		}
	}
	// 再运行剩下的
	for _, sdr := range seeders {
		// 过滤已运行
		if _, ok := executed[sdr.Name]; !ok {
			sdr.Func(db)
		}
	}
}

// RunSeeder 运行单个 Seeder
func RunSeeder(db *gorm.DB, name string) {
	for _, sdr := range seeders {
		if name == sdr.Name {
			sdr.Func(db)
			break
		}
	}
}
