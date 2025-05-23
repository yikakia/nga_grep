// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package gen

import (
	"context"

	"github.com/yikakia/nga_grep/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"

	"gorm.io/gen"
	"gorm.io/gen/field"

	"gorm.io/plugin/dbresolver"
)

func newThreadLatestData(db *gorm.DB, opts ...gen.DOOption) threadLatestData {
	_threadLatestData := threadLatestData{}

	_threadLatestData.threadLatestDataDo.UseDB(db, opts...)
	_threadLatestData.threadLatestDataDo.UseModel(&model.ThreadLatestData{})

	tableName := _threadLatestData.threadLatestDataDo.TableName()
	_threadLatestData.ALL = field.NewAsterisk(tableName)
	_threadLatestData.TID = field.NewInt(tableName, "tid")
	_threadLatestData.LastTime = field.NewTime(tableName, "last_time")
	_threadLatestData.LastReplyCount = field.NewInt(tableName, "last_reply_count")

	_threadLatestData.fillFieldMap()

	return _threadLatestData
}

type threadLatestData struct {
	threadLatestDataDo

	ALL            field.Asterisk
	TID            field.Int
	LastTime       field.Time
	LastReplyCount field.Int

	fieldMap map[string]field.Expr
}

func (t threadLatestData) Table(newTableName string) *threadLatestData {
	t.threadLatestDataDo.UseTable(newTableName)
	return t.updateTableName(newTableName)
}

func (t threadLatestData) As(alias string) *threadLatestData {
	t.threadLatestDataDo.DO = *(t.threadLatestDataDo.As(alias).(*gen.DO))
	return t.updateTableName(alias)
}

func (t *threadLatestData) updateTableName(table string) *threadLatestData {
	t.ALL = field.NewAsterisk(table)
	t.TID = field.NewInt(table, "tid")
	t.LastTime = field.NewTime(table, "last_time")
	t.LastReplyCount = field.NewInt(table, "last_reply_count")

	t.fillFieldMap()

	return t
}

func (t *threadLatestData) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := t.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (t *threadLatestData) fillFieldMap() {
	t.fieldMap = make(map[string]field.Expr, 3)
	t.fieldMap["tid"] = t.TID
	t.fieldMap["last_time"] = t.LastTime
	t.fieldMap["last_reply_count"] = t.LastReplyCount
}

func (t threadLatestData) clone(db *gorm.DB) threadLatestData {
	t.threadLatestDataDo.ReplaceConnPool(db.Statement.ConnPool)
	return t
}

func (t threadLatestData) replaceDB(db *gorm.DB) threadLatestData {
	t.threadLatestDataDo.ReplaceDB(db)
	return t
}

type threadLatestDataDo struct{ gen.DO }

type IThreadLatestDataDo interface {
	gen.SubQuery
	Debug() IThreadLatestDataDo
	WithContext(ctx context.Context) IThreadLatestDataDo
	WithResult(fc func(tx gen.Dao)) gen.ResultInfo
	ReplaceDB(db *gorm.DB)
	ReadDB() IThreadLatestDataDo
	WriteDB() IThreadLatestDataDo
	As(alias string) gen.Dao
	Session(config *gorm.Session) IThreadLatestDataDo
	Columns(cols ...field.Expr) gen.Columns
	Clauses(conds ...clause.Expression) IThreadLatestDataDo
	Not(conds ...gen.Condition) IThreadLatestDataDo
	Or(conds ...gen.Condition) IThreadLatestDataDo
	Select(conds ...field.Expr) IThreadLatestDataDo
	Where(conds ...gen.Condition) IThreadLatestDataDo
	Order(conds ...field.Expr) IThreadLatestDataDo
	Distinct(cols ...field.Expr) IThreadLatestDataDo
	Omit(cols ...field.Expr) IThreadLatestDataDo
	Join(table schema.Tabler, on ...field.Expr) IThreadLatestDataDo
	LeftJoin(table schema.Tabler, on ...field.Expr) IThreadLatestDataDo
	RightJoin(table schema.Tabler, on ...field.Expr) IThreadLatestDataDo
	Group(cols ...field.Expr) IThreadLatestDataDo
	Having(conds ...gen.Condition) IThreadLatestDataDo
	Limit(limit int) IThreadLatestDataDo
	Offset(offset int) IThreadLatestDataDo
	Count() (count int64, err error)
	Scopes(funcs ...func(gen.Dao) gen.Dao) IThreadLatestDataDo
	Unscoped() IThreadLatestDataDo
	Create(values ...*model.ThreadLatestData) error
	CreateInBatches(values []*model.ThreadLatestData, batchSize int) error
	Save(values ...*model.ThreadLatestData) error
	First() (*model.ThreadLatestData, error)
	Take() (*model.ThreadLatestData, error)
	Last() (*model.ThreadLatestData, error)
	Find() ([]*model.ThreadLatestData, error)
	FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.ThreadLatestData, err error)
	FindInBatches(result *[]*model.ThreadLatestData, batchSize int, fc func(tx gen.Dao, batch int) error) error
	Pluck(column field.Expr, dest interface{}) error
	Delete(...*model.ThreadLatestData) (info gen.ResultInfo, err error)
	Update(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	Updates(value interface{}) (info gen.ResultInfo, err error)
	UpdateColumn(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateColumnSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	UpdateColumns(value interface{}) (info gen.ResultInfo, err error)
	UpdateFrom(q gen.SubQuery) gen.Dao
	Attrs(attrs ...field.AssignExpr) IThreadLatestDataDo
	Assign(attrs ...field.AssignExpr) IThreadLatestDataDo
	Joins(fields ...field.RelationField) IThreadLatestDataDo
	Preload(fields ...field.RelationField) IThreadLatestDataDo
	FirstOrInit() (*model.ThreadLatestData, error)
	FirstOrCreate() (*model.ThreadLatestData, error)
	FindByPage(offset int, limit int) (result []*model.ThreadLatestData, count int64, err error)
	ScanByPage(result interface{}, offset int, limit int) (count int64, err error)
	Scan(result interface{}) (err error)
	Returning(value interface{}, columns ...string) IThreadLatestDataDo
	UnderlyingDB() *gorm.DB
	schema.Tabler
}

func (t threadLatestDataDo) Debug() IThreadLatestDataDo {
	return t.withDO(t.DO.Debug())
}

func (t threadLatestDataDo) WithContext(ctx context.Context) IThreadLatestDataDo {
	return t.withDO(t.DO.WithContext(ctx))
}

func (t threadLatestDataDo) ReadDB() IThreadLatestDataDo {
	return t.Clauses(dbresolver.Read)
}

func (t threadLatestDataDo) WriteDB() IThreadLatestDataDo {
	return t.Clauses(dbresolver.Write)
}

func (t threadLatestDataDo) Session(config *gorm.Session) IThreadLatestDataDo {
	return t.withDO(t.DO.Session(config))
}

func (t threadLatestDataDo) Clauses(conds ...clause.Expression) IThreadLatestDataDo {
	return t.withDO(t.DO.Clauses(conds...))
}

func (t threadLatestDataDo) Returning(value interface{}, columns ...string) IThreadLatestDataDo {
	return t.withDO(t.DO.Returning(value, columns...))
}

func (t threadLatestDataDo) Not(conds ...gen.Condition) IThreadLatestDataDo {
	return t.withDO(t.DO.Not(conds...))
}

func (t threadLatestDataDo) Or(conds ...gen.Condition) IThreadLatestDataDo {
	return t.withDO(t.DO.Or(conds...))
}

func (t threadLatestDataDo) Select(conds ...field.Expr) IThreadLatestDataDo {
	return t.withDO(t.DO.Select(conds...))
}

func (t threadLatestDataDo) Where(conds ...gen.Condition) IThreadLatestDataDo {
	return t.withDO(t.DO.Where(conds...))
}

func (t threadLatestDataDo) Order(conds ...field.Expr) IThreadLatestDataDo {
	return t.withDO(t.DO.Order(conds...))
}

func (t threadLatestDataDo) Distinct(cols ...field.Expr) IThreadLatestDataDo {
	return t.withDO(t.DO.Distinct(cols...))
}

func (t threadLatestDataDo) Omit(cols ...field.Expr) IThreadLatestDataDo {
	return t.withDO(t.DO.Omit(cols...))
}

func (t threadLatestDataDo) Join(table schema.Tabler, on ...field.Expr) IThreadLatestDataDo {
	return t.withDO(t.DO.Join(table, on...))
}

func (t threadLatestDataDo) LeftJoin(table schema.Tabler, on ...field.Expr) IThreadLatestDataDo {
	return t.withDO(t.DO.LeftJoin(table, on...))
}

func (t threadLatestDataDo) RightJoin(table schema.Tabler, on ...field.Expr) IThreadLatestDataDo {
	return t.withDO(t.DO.RightJoin(table, on...))
}

func (t threadLatestDataDo) Group(cols ...field.Expr) IThreadLatestDataDo {
	return t.withDO(t.DO.Group(cols...))
}

func (t threadLatestDataDo) Having(conds ...gen.Condition) IThreadLatestDataDo {
	return t.withDO(t.DO.Having(conds...))
}

func (t threadLatestDataDo) Limit(limit int) IThreadLatestDataDo {
	return t.withDO(t.DO.Limit(limit))
}

func (t threadLatestDataDo) Offset(offset int) IThreadLatestDataDo {
	return t.withDO(t.DO.Offset(offset))
}

func (t threadLatestDataDo) Scopes(funcs ...func(gen.Dao) gen.Dao) IThreadLatestDataDo {
	return t.withDO(t.DO.Scopes(funcs...))
}

func (t threadLatestDataDo) Unscoped() IThreadLatestDataDo {
	return t.withDO(t.DO.Unscoped())
}

func (t threadLatestDataDo) Create(values ...*model.ThreadLatestData) error {
	if len(values) == 0 {
		return nil
	}
	return t.DO.Create(values)
}

func (t threadLatestDataDo) CreateInBatches(values []*model.ThreadLatestData, batchSize int) error {
	return t.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (t threadLatestDataDo) Save(values ...*model.ThreadLatestData) error {
	if len(values) == 0 {
		return nil
	}
	return t.DO.Save(values)
}

func (t threadLatestDataDo) First() (*model.ThreadLatestData, error) {
	if result, err := t.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*model.ThreadLatestData), nil
	}
}

func (t threadLatestDataDo) Take() (*model.ThreadLatestData, error) {
	if result, err := t.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*model.ThreadLatestData), nil
	}
}

func (t threadLatestDataDo) Last() (*model.ThreadLatestData, error) {
	if result, err := t.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*model.ThreadLatestData), nil
	}
}

func (t threadLatestDataDo) Find() ([]*model.ThreadLatestData, error) {
	result, err := t.DO.Find()
	return result.([]*model.ThreadLatestData), err
}

func (t threadLatestDataDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.ThreadLatestData, err error) {
	buf := make([]*model.ThreadLatestData, 0, batchSize)
	err = t.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (t threadLatestDataDo) FindInBatches(result *[]*model.ThreadLatestData, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return t.DO.FindInBatches(result, batchSize, fc)
}

func (t threadLatestDataDo) Attrs(attrs ...field.AssignExpr) IThreadLatestDataDo {
	return t.withDO(t.DO.Attrs(attrs...))
}

func (t threadLatestDataDo) Assign(attrs ...field.AssignExpr) IThreadLatestDataDo {
	return t.withDO(t.DO.Assign(attrs...))
}

func (t threadLatestDataDo) Joins(fields ...field.RelationField) IThreadLatestDataDo {
	for _, _f := range fields {
		t = *t.withDO(t.DO.Joins(_f))
	}
	return &t
}

func (t threadLatestDataDo) Preload(fields ...field.RelationField) IThreadLatestDataDo {
	for _, _f := range fields {
		t = *t.withDO(t.DO.Preload(_f))
	}
	return &t
}

func (t threadLatestDataDo) FirstOrInit() (*model.ThreadLatestData, error) {
	if result, err := t.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*model.ThreadLatestData), nil
	}
}

func (t threadLatestDataDo) FirstOrCreate() (*model.ThreadLatestData, error) {
	if result, err := t.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*model.ThreadLatestData), nil
	}
}

func (t threadLatestDataDo) FindByPage(offset int, limit int) (result []*model.ThreadLatestData, count int64, err error) {
	result, err = t.Offset(offset).Limit(limit).Find()
	if err != nil {
		return
	}

	if size := len(result); 0 < limit && 0 < size && size < limit {
		count = int64(size + offset)
		return
	}

	count, err = t.Offset(-1).Limit(-1).Count()
	return
}

func (t threadLatestDataDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = t.Count()
	if err != nil {
		return
	}

	err = t.Offset(offset).Limit(limit).Scan(result)
	return
}

func (t threadLatestDataDo) Scan(result interface{}) (err error) {
	return t.DO.Scan(result)
}

func (t threadLatestDataDo) Delete(models ...*model.ThreadLatestData) (result gen.ResultInfo, err error) {
	return t.DO.Delete(models)
}

func (t *threadLatestDataDo) withDO(do gen.Dao) *threadLatestDataDo {
	t.DO = *do.(*gen.DO)
	return t
}
