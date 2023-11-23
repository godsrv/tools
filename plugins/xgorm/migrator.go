package xgorm

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/gookit/color"
	"gorm.io/gorm"
)

// Migrator 数据迁移操作类
type Migrator struct {
	Folder   string
	DB       *gorm.DB
	Migrator gorm.Migrator
}

// Migration 对应数据的 migrations 表里的一条数据
type Migration struct {
	ID        uint64 `gorm:"primaryKey;autoIncrement;"`
	Migration string `gorm:"type:varchar(191);not null;unique;"`
	Batch     int
}

// NewMigrator 创建 Migrator 实例，用以执行迁移操作
func NewMigrator(db *gorm.DB, folder string) *Migrator {
	migrator := &Migrator{
		Folder:   folder,
		DB:       db,
		Migrator: db.Migrator(),
	}
	// migrations 不存在的话就创建它
	migrator.createMigrationsTable()
	return migrator
}

// 创建 migrations 表
func (m *Migrator) createMigrationsTable() {
	migration := Migration{}
	// 不存在才创建
	if !m.Migrator.HasTable(&migration) {
		m.Migrator.CreateTable(&migration)
	}
}

// Up 执行所有未迁移过的文件
func (m *Migrator) Up() {
	migrateFiles := m.readAllMigrationFiles()
	// 获取当前批次的值
	batch := m.getBatch()
	// 获取所有迁移数据
	migrations := []Migration{}
	m.DB.Find(&migrations)
	// 可以通过此值来判断数据库是否已是最新
	runed := false
	// 对迁移文件进行遍历，如果没有执行过，就执行 up 回调
	for _, mfile := range migrateFiles {
		// 对比文件名称，看是否已经运行过
		if mfile.isNotMigrated(migrations) {
			m.runUpMigration(mfile, batch)
			runed = true
		}
	}
	if !runed {
		color.Green.Println("[migrations] database is up to date.")
	}
}

// Rollback 回滚上一个操作
func (m *Migrator) Rollback() {
	// 获取最后一批次的迁移数据
	lastMigration := Migration{}
	m.DB.Order("id DESC").First(&lastMigration)
	migrations := []Migration{}
	m.DB.Where("batch = ?", lastMigration.Batch).Order("id DESC").Find(&migrations)
	if !m.rollbackMigrations(migrations) {
		color.Green.Println("[migrations] table is empty, nothing to rollback.")
	}
}

// 回退迁移，按照倒序执行迁移的 down 方法
func (m *Migrator) rollbackMigrations(migrations []Migration) bool {
	// 标记是否真的有执行了迁移回退的操作
	runed := false
	for _, _migration := range migrations {
		// 执行迁移文件的 down 方法
		mfile := getMigrationFile(_migration.Migration)
		if mfile.Down != nil {
			sqlDb, err := m.DB.DB()
			if err != nil {
				color.Errorln("[migrations] rollback " + "failed " + err.Error())
				os.Exit(0)
			}
			err = mfile.Down(m.DB.Migrator(), sqlDb)
			if err != nil {
				color.Errorln("[migrations] rollback " + mfile.FileName + "failed " + err.Error())
				os.Exit(0)
			}
		}
		runed = true
		// 回退成功了就删除掉这条记录
		m.DB.Delete(&_migration)
		// 打印运行状态
		color.Green.Println("[migrations] finished " + mfile.FileName)
	}
	return runed
}

// 获取当前这个批次的值
func (m *Migrator) getBatch() int {
	// 默认为 1
	batch := 1
	// 取最后执行的一条迁移数据
	lastMigration := Migration{}
	m.DB.Order("id DESC").First(&lastMigration)
	// 如果有值的话，加一
	if lastMigration.ID > 0 {
		batch = lastMigration.Batch + 1
	}
	return batch
}

// 从文件目录读取文件，依据时间正确排序
func (m *Migrator) readAllMigrationFiles() []MigrationFile {
	files, err := os.ReadDir(m.Folder)
	if err != nil {
		color.Errorln("[migrations] read  dir failed, err:", err.Error())
		os.Exit(0)
	}
	var migrateFiles []MigrationFile
	for _, f := range files {
		fName := f.Name()
		fileName := strings.TrimSuffix(fName, filepath.Ext(fName))
		mfile := getMigrationFile(fileName)
		if len(mfile.FileName) > 0 {
			migrateFiles = append(migrateFiles, mfile)
		}
	}
	return migrateFiles
}

// 执行迁移，执行迁移的 up 方法
func (m *Migrator) runUpMigration(mfile MigrationFile, batch int) {
	if mfile.Up != nil {
		sqlDb, err := m.DB.DB()
		if err != nil {
			color.Errorln("[migrations] up " + "failed " + err.Error())
			os.Exit(0)
		}
		err = mfile.Up(m.DB.Migrator(), sqlDb)
		if err != nil {
			color.Errorln("[migrations] failed " + err.Error())
			os.Exit(0)
		}
		color.Green.Println("migrated " + mfile.FileName)
	}
	// 入库执行迁移
	err := m.DB.Create(&Migration{Migration: mfile.FileName, Batch: batch}).Error
	if err != nil {
		color.Errorln("[migrations] failed " + err.Error())
		os.Exit(0)
	}
}

// Reset 回滚所有迁移
func (m *Migrator) Reset() {

	migrations := []Migration{}

	// 按照倒序读取所有迁移文件
	m.DB.Order("id DESC").Find(&migrations)

	// 回滚所有迁移
	if !m.rollbackMigrations(migrations) {
		color.Green.Println("[migrations] table is empty, nothing to reset.")
	}
}

// Refresh 回滚所有迁移，并运行所有迁移
func (m *Migrator) Refresh() {

	// 回滚所有迁移
	m.Reset()

	// 再次执行所有迁移
	m.Up()
}
