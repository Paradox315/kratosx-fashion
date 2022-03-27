// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package query

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"

	"kratosx-fashion/app/system/internal/data/model"

	"gorm.io/gen"
	"gorm.io/gen/field"
)

func newLoginLog(db *gorm.DB) loginLog {
	_loginLog := loginLog{}

	_loginLog.loginLogDo.UseDB(db)
	_loginLog.loginLogDo.UseModel(&model.LoginLog{})

	tableName := _loginLog.loginLogDo.TableName()
	_loginLog.ALL = field.NewField(tableName, "*")
	_loginLog.ID = field.NewUint(tableName, "id")
	_loginLog.CreatedAt = field.NewTime(tableName, "created_at")
	_loginLog.UpdatedAt = field.NewTime(tableName, "updated_at")
	_loginLog.DeletedAt = field.NewField(tableName, "deleted_at")
	_loginLog.UserID = field.NewUint64(tableName, "user_id")
	_loginLog.Ip = field.NewString(tableName, "ip")
	_loginLog.Location = field.NewString(tableName, "location")
	_loginLog.LoginType = field.NewUint8(tableName, "login_type")
	_loginLog.Agent = field.NewString(tableName, "agent")
	_loginLog.OS = field.NewString(tableName, "os")
	_loginLog.Device = field.NewString(tableName, "device")
	_loginLog.DeviceType = field.NewUint8(tableName, "device_type")

	_loginLog.fillFieldMap()

	return _loginLog
}

type loginLog struct {
	loginLogDo loginLogDo

	ALL        field.Field
	ID         field.Uint
	CreatedAt  field.Time
	UpdatedAt  field.Time
	DeletedAt  field.Field
	UserID     field.Uint64
	Ip         field.String
	Location   field.String
	LoginType  field.Uint8
	Agent      field.String
	OS         field.String
	Device     field.String
	DeviceType field.Uint8

	fieldMap map[string]field.Expr
}

func (l loginLog) Table(newTableName string) *loginLog {
	l.loginLogDo.UseTable(newTableName)
	return l.updateTableName(newTableName)
}

func (l loginLog) As(alias string) *loginLog {
	l.loginLogDo.DO = *(l.loginLogDo.As(alias).(*gen.DO))
	return l.updateTableName(alias)
}

func (l *loginLog) updateTableName(table string) *loginLog {
	l.ALL = field.NewField(table, "*")
	l.ID = field.NewUint(table, "id")
	l.CreatedAt = field.NewTime(table, "created_at")
	l.UpdatedAt = field.NewTime(table, "updated_at")
	l.DeletedAt = field.NewField(table, "deleted_at")
	l.UserID = field.NewUint64(table, "user_id")
	l.Ip = field.NewString(table, "ip")
	l.Location = field.NewString(table, "location")
	l.LoginType = field.NewUint8(table, "login_type")
	l.Agent = field.NewString(table, "agent")
	l.OS = field.NewString(table, "os")
	l.Device = field.NewString(table, "device")
	l.DeviceType = field.NewUint8(table, "device_type")

	l.fillFieldMap()

	return l
}

func (l *loginLog) WithContext(ctx context.Context) *loginLogDo { return l.loginLogDo.WithContext(ctx) }

func (l loginLog) TableName() string { return l.loginLogDo.TableName() }

func (l *loginLog) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := l.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (l *loginLog) fillFieldMap() {
	l.fieldMap = make(map[string]field.Expr, 12)
	l.fieldMap["id"] = l.ID
	l.fieldMap["created_at"] = l.CreatedAt
	l.fieldMap["updated_at"] = l.UpdatedAt
	l.fieldMap["deleted_at"] = l.DeletedAt
	l.fieldMap["user_id"] = l.UserID
	l.fieldMap["ip"] = l.Ip
	l.fieldMap["location"] = l.Location
	l.fieldMap["login_type"] = l.LoginType
	l.fieldMap["agent"] = l.Agent
	l.fieldMap["os"] = l.OS
	l.fieldMap["device"] = l.Device
	l.fieldMap["device_type"] = l.DeviceType
}

func (l loginLog) clone(db *gorm.DB) loginLog {
	l.loginLogDo.ReplaceDB(db)
	return l
}

type loginLogDo struct{ gen.DO }

func (l loginLogDo) Debug() *loginLogDo {
	return l.withDO(l.DO.Debug())
}

func (l loginLogDo) WithContext(ctx context.Context) *loginLogDo {
	return l.withDO(l.DO.WithContext(ctx))
}

func (l loginLogDo) Clauses(conds ...clause.Expression) *loginLogDo {
	return l.withDO(l.DO.Clauses(conds...))
}

func (l loginLogDo) Returning(value interface{}, columns ...string) *loginLogDo {
	return l.withDO(l.DO.Returning(value, columns...))
}

func (l loginLogDo) Not(conds ...gen.Condition) *loginLogDo {
	return l.withDO(l.DO.Not(conds...))
}

func (l loginLogDo) Or(conds ...gen.Condition) *loginLogDo {
	return l.withDO(l.DO.Or(conds...))
}

func (l loginLogDo) Select(conds ...field.Expr) *loginLogDo {
	return l.withDO(l.DO.Select(conds...))
}

func (l loginLogDo) Where(conds ...gen.Condition) *loginLogDo {
	return l.withDO(l.DO.Where(conds...))
}

func (l loginLogDo) Exists(subquery interface{ UnderlyingDB() *gorm.DB }) *loginLogDo {
	return l.Where(field.CompareSubQuery(field.ExistsOp, nil, subquery.UnderlyingDB()))
}

func (l loginLogDo) Order(conds ...field.Expr) *loginLogDo {
	return l.withDO(l.DO.Order(conds...))
}

func (l loginLogDo) Distinct(cols ...field.Expr) *loginLogDo {
	return l.withDO(l.DO.Distinct(cols...))
}

func (l loginLogDo) Omit(cols ...field.Expr) *loginLogDo {
	return l.withDO(l.DO.Omit(cols...))
}

func (l loginLogDo) Join(table schema.Tabler, on ...field.Expr) *loginLogDo {
	return l.withDO(l.DO.Join(table, on...))
}

func (l loginLogDo) LeftJoin(table schema.Tabler, on ...field.Expr) *loginLogDo {
	return l.withDO(l.DO.LeftJoin(table, on...))
}

func (l loginLogDo) RightJoin(table schema.Tabler, on ...field.Expr) *loginLogDo {
	return l.withDO(l.DO.RightJoin(table, on...))
}

func (l loginLogDo) Group(cols ...field.Expr) *loginLogDo {
	return l.withDO(l.DO.Group(cols...))
}

func (l loginLogDo) Having(conds ...gen.Condition) *loginLogDo {
	return l.withDO(l.DO.Having(conds...))
}

func (l loginLogDo) Limit(limit int) *loginLogDo {
	return l.withDO(l.DO.Limit(limit))
}

func (l loginLogDo) Offset(offset int) *loginLogDo {
	return l.withDO(l.DO.Offset(offset))
}

func (l loginLogDo) Scopes(funcs ...func(gen.Dao) gen.Dao) *loginLogDo {
	return l.withDO(l.DO.Scopes(funcs...))
}

func (l loginLogDo) Unscoped() *loginLogDo {
	return l.withDO(l.DO.Unscoped())
}

func (l loginLogDo) Create(values ...*model.LoginLog) error {
	if len(values) == 0 {
		return nil
	}
	return l.DO.Create(values)
}

func (l loginLogDo) CreateInBatches(values []*model.LoginLog, batchSize int) error {
	return l.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (l loginLogDo) Save(values ...*model.LoginLog) error {
	if len(values) == 0 {
		return nil
	}
	return l.DO.Save(values)
}

func (l loginLogDo) First() (*model.LoginLog, error) {
	if result, err := l.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*model.LoginLog), nil
	}
}

func (l loginLogDo) Take() (*model.LoginLog, error) {
	if result, err := l.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*model.LoginLog), nil
	}
}

func (l loginLogDo) Last() (*model.LoginLog, error) {
	if result, err := l.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*model.LoginLog), nil
	}
}

func (l loginLogDo) Find() ([]*model.LoginLog, error) {
	result, err := l.DO.Find()
	return result.([]*model.LoginLog), err
}

func (l loginLogDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.LoginLog, err error) {
	buf := make([]*model.LoginLog, 0, batchSize)
	err = l.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (l loginLogDo) FindInBatches(result *[]*model.LoginLog, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return l.DO.FindInBatches(result, batchSize, fc)
}

func (l loginLogDo) Attrs(attrs ...field.AssignExpr) *loginLogDo {
	return l.withDO(l.DO.Attrs(attrs...))
}

func (l loginLogDo) Assign(attrs ...field.AssignExpr) *loginLogDo {
	return l.withDO(l.DO.Assign(attrs...))
}

func (l loginLogDo) Joins(field field.RelationField) *loginLogDo {
	return l.withDO(l.DO.Joins(field))
}

func (l loginLogDo) Preload(field field.RelationField) *loginLogDo {
	return l.withDO(l.DO.Preload(field))
}

func (l loginLogDo) FirstOrInit() (*model.LoginLog, error) {
	if result, err := l.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*model.LoginLog), nil
	}
}

func (l loginLogDo) FirstOrCreate() (*model.LoginLog, error) {
	if result, err := l.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*model.LoginLog), nil
	}
}

func (l loginLogDo) FindByPage(offset int, limit int) (result []*model.LoginLog, count int64, err error) {
	if limit <= 0 {
		count, err = l.Count()
		return
	}

	result, err = l.Offset(offset).Limit(limit).Find()
	if err != nil {
		return
	}

	if size := len(result); 0 < size && size < limit {
		count = int64(size + offset)
		return
	}

	count, err = l.Offset(-1).Limit(-1).Count()
	return
}

func (l loginLogDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = l.Count()
	if err != nil {
		return
	}

	err = l.Offset(offset).Limit(limit).Scan(result)
	return
}

func (l *loginLogDo) withDO(do gen.Dao) *loginLogDo {
	l.DO = *do.(*gen.DO)
	return l
}
