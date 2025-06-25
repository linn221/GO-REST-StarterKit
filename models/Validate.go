package models

import (
	"linn221/shop/utils"

	"gorm.io/gorm"
)

type Rule interface {
	Init() bool
	CountResults(*gorm.DB, *int64) error
}

func Validate(db *gorm.DB, rules ...Rule) error {
	var count int64
	for _, rule := range rules {
		if ok := rule.Init(); !ok {
			continue
		}
		err := rule.CountResults(db, &count)
		if err != nil {
			return err
		}
	}
	return nil
}
func ValidateInBatch(db *gorm.DB, rules ...Rule) []error {
	var count int64
	errors := make([]error, 0)
	for _, rule := range rules {
		if ok := rule.Init(); !ok {
			continue
		}
		err := rule.CountResults(db, &count)
		if err != nil {
			errors = append(errors, err)
		}
	}
	return errors
}

type HasFilter struct {
	Cond         string
	FilterValues []interface{}
}

func (f *HasFilter) ApplyFilter(dbCtx *gorm.DB) {
	if f != nil {
		dbCtx.Where(f.Cond, f.FilterValues...)
	}
}

func NewFilter(cond string, values ...interface{}) *HasFilter {
	return &HasFilter{
		Cond:         cond,
		FilterValues: values,
	}
}

func NewShopFilter(shopId string) *HasFilter {
	return &HasFilter{
		Cond:         "shop_id = ?",
		FilterValues: []interface{}{shopId},
	}
}

// check if resource exists (where business_id = ?)
type ruleExists struct {
	statusCode *int
	table      string
	id         interface{}
	err        error
	do         *bool
	*HasFilter
}

// specifies When to validate
// if When is not specified, will validate by default
func (rule ruleExists) When(when bool) ruleExists {
	rule.do = &when
	return rule
}

func (vr ruleExists) Init() bool {
	// skip validation if user specifies when
	return vr.do == nil || *vr.do
}

func (vr ruleExists) CountResults(dbCtx *gorm.DB, count *int64) error {
	dbCtx = dbCtx.Table(vr.table).Where("id = ?", vr.id)
	vr.ApplyFilter(dbCtx)
	if err := dbCtx.Count(count).Error; err != nil {
		return dbError(err)
	}
	if *count <= 0 {
		return vr.err
	}

	return nil
}

func NewExistsRule(table string, id interface{}, err error, filter *HasFilter) ruleExists {
	return ruleExists{
		table:     table,
		id:        id,
		HasFilter: filter,
		err:       err,
	}
}

// check if slice of resource id exists (where business_id IN ?)
type RuleMassExists[ID comparable] struct {
	Table         string
	Ids           []ID
	Err           error
	NoDuplicateID bool
	*HasFilter
}

func (r RuleMassExists[ID]) Init() bool {
	return len(r.Ids) > 0
}

func (r RuleMassExists[ID]) CountResults(dbCtx *gorm.DB, count *int64) error {

	uniqIds := utils.UniqueSlice(r.Ids)
	dbCtx = dbCtx.Table(r.Table).Where("id IN ?", uniqIds)
	err := dbCtx.Count(count).Error
	if err != nil {
		return dbError(err)
	}
	if *count != int64(len(uniqIds)) {
		return r.Err
	}

	return nil
}

type ruleUnique struct {
	table    string
	err      error
	column   string
	value    interface{}
	exceptId int
	do       *bool

	*HasFilter
}

func (rule ruleUnique) When(cond bool) ruleUnique {
	rule.do = &cond
	return rule
}

func (rule ruleUnique) Filter(cond string, values ...interface{}) ruleUnique {
	rule.HasFilter = &HasFilter{
		Cond:         cond,
		FilterValues: values,
	}
	return rule
}

func (rule ruleUnique) Init() bool {

	if rule.do != nil && !*rule.do {
		return false
	}
	return true
}

func (r ruleUnique) CountResults(dbCtx *gorm.DB, count *int64) error {
	dbCtx = dbCtx.Table(r.table).Where("`"+r.column+"`"+" = ?", r.value)
	if r.exceptId > 0 {
		dbCtx.Where("id != ?", r.exceptId)
	}
	r.ApplyFilter(dbCtx)
	err := dbCtx.Count(count).Error
	if err != nil {
		return dbError(err)
	}

	if *count > 0 {
		return r.err
	}

	return nil
}

func NewUniqueRule(table string, column string, value interface{}, exceptId int, err error, filter *HasFilter) ruleUnique {
	// var v T
	return ruleUnique{
		table:     table,
		column:    column,
		value:     value,
		exceptId:  exceptId,
		err:       err,
		HasFilter: filter,
	}
}

type noResultRule struct {
	statusCode *int
	table      string
	err        error
	do         *bool
	*HasFilter
}

// specifies When to validate
// if When is not specified, will validate by default
func (rule noResultRule) When(when bool) noResultRule {
	rule.do = &when
	return rule
}

func (rule noResultRule) OverrideStatusCode(i int) noResultRule {
	rule.statusCode = &i
	return rule
}

func (vr noResultRule) Init() bool {
	// skip validation if user specifies when
	return vr.do == nil || *vr.do
}

func (vr noResultRule) CountResults(dbCtx *gorm.DB, count *int64) error {
	dbCtx = dbCtx.Table(vr.table)
	vr.ApplyFilter(dbCtx)
	if err := dbCtx.Count(count).Error; err != nil {
		return dbError(err)
	}
	if *count > 0 {
		return vr.err
	}

	return nil
}

func NewNoResultRule(table string, err error, filter *HasFilter) noResultRule {
	return noResultRule{
		table:     table,
		err:       err,
		HasFilter: filter,
	}
}
