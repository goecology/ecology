package dao

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/i2eco/ecology/appgo/model/mysql"
	"github.com/i2eco/ecology/appgo/model/trans"
	"github.com/i2eco/muses/pkg/logger"
	"github.com/jinzhu/gorm"
	"go.uber.org/zap"
)

type attachment struct {
	logger *logger.Client
	db     *gorm.DB
}

func InitAttachment(logger *logger.Client, db *gorm.DB) *attachment {
	return &attachment{
		logger: logger,
		db:     db,
	}
}

// Create 新增一条记
func (g *attachment) Create(c *gin.Context, db *gorm.DB, data *mysql.Attachment) (err error) {
	data.CreateTime = time.Now()
	if err = db.Create(data).Error; err != nil {
		g.logger.Error("create attachment create error", zap.Error(err))
		return
	}
	return nil
}

// Update 根据主键更新一条记录
func (g *attachment) Update(c *gin.Context, db *gorm.DB, paramAttachmentId int, ups mysql.Ups) (err error) {
	var sql = "`attachment_id`=?"
	var binds = []interface{}{paramAttachmentId}

	if err = db.Table("attachment").Where(sql, binds...).Updates(ups).Error; err != nil {
		g.logger.Error("attachment update error", zap.Error(err))
		return
	}
	return
}

// UpdateX Update的扩展方法，根据Cond更新一条或多条记录
func (g *attachment) UpdateX(c *gin.Context, db *gorm.DB, conds mysql.Conds, ups mysql.Ups) (err error) {

	sql, binds := mysql.BuildQuery(conds)
	if err = db.Table("attachment").Where(sql, binds...).Updates(ups).Error; err != nil {
		g.logger.Error("attachment update error", zap.Error(err))
		return
	}
	return
}

// Delete 根据主键删除一条记录。如果有delete_time则软删除，否则硬删除。
func (g *attachment) Delete(c *gin.Context, db *gorm.DB, paramAttachmentId int) (err error) {
	var sql = "`attachment_id`=?"
	var binds = []interface{}{paramAttachmentId}

	if err = db.Table("attachment").Where(sql, binds...).Delete(&mysql.Attachment{}).Error; err != nil {
		g.logger.Error("attachment delete error", zap.Error(err))
		return
	}

	return
}

// DeleteX Delete的扩展方法，根据Cond删除一条或多条记录。如果有delete_time则软删除，否则硬删除。
func (g *attachment) DeleteX(c *gin.Context, db *gorm.DB, conds mysql.Conds) (err error) {
	sql, binds := mysql.BuildQuery(conds)

	if err = db.Table("attachment").Where(sql, binds...).Delete(&mysql.Attachment{}).Error; err != nil {
		g.logger.Error("attachment delete error", zap.Error(err))
		return
	}

	return
}

// Info 根据PRI查询单条记录
func (g *attachment) Info(c *gin.Context, paramAttachmentId int) (resp mysql.Attachment, err error) {
	var sql = "`attachment_id`=?"
	var binds = []interface{}{paramAttachmentId}

	if err = g.db.Table("attachment").Where(sql, binds...).First(&resp).Error; err != nil {
		g.logger.Error("attachment info error", zap.Error(err))
		return
	}
	return
}

// InfoX Info的扩展方法，根据Cond查询单条记录
func (g *attachment) InfoX(c *gin.Context, conds mysql.Conds) (resp mysql.Attachment, err error) {
	sql, binds := mysql.BuildQuery(conds)

	if err = g.db.Table("attachment").Where(sql, binds...).First(&resp).Error; err != nil {
		g.logger.Error("attachment info error", zap.Error(err))
		return
	}
	return
}

// List 查询list，extra[0]为sorts
func (g *attachment) List(c *gin.Context, conds mysql.Conds, extra ...string) (resp []mysql.Attachment, err error) {
	sql, binds := mysql.BuildQuery(conds)

	sorts := ""
	if len(extra) >= 1 {
		sorts = extra[0]
	}
	if err = g.db.Table("attachment").Where(sql, binds...).Order(sorts).Find(&resp).Error; err != nil {
		g.logger.Error("attachment info error", zap.Error(err))
		return
	}
	return
}

// ListMap 查询map，map遍历的时候是无序的，所以指定sorts参数没有意义
func (g *attachment) ListMap(c *gin.Context, conds mysql.Conds) (resp map[int]mysql.Attachment, err error) {
	sql, binds := mysql.BuildQuery(conds)

	mysqlSlice := make([]mysql.Attachment, 0)
	resp = make(map[int]mysql.Attachment, 0)
	if err = g.db.Table("attachment").Where(sql, binds...).Find(&mysqlSlice).Error; err != nil {
		g.logger.Error("attachment info error", zap.Error(err))
		return
	}
	for _, value := range mysqlSlice {
		resp[value.AttachmentId] = value
	}
	return
}

// ListPage 根据分页条件查询list
func (g *attachment) ListPage(c *gin.Context, conds mysql.Conds, reqList *trans.ReqPage) (total int, respList []mysql.Attachment) {
	if reqList.PageSize == 0 {
		reqList.PageSize = 10
	}
	if reqList.Current == 0 {
		reqList.Current = 1
	}
	sql, binds := mysql.BuildQuery(conds)

	db := g.db.Table("attachment").Where(sql, binds...)
	respList = make([]mysql.Attachment, 0)
	db.Count(&total)
	db.Order(reqList.Sort).Offset((reqList.Current - 1) * reqList.PageSize).Limit(reqList.PageSize).Find(&respList)
	return
}
