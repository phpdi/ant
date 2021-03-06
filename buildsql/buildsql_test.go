package buildsql

import (
	"fmt"
	"testing"
)

//插入
func TestBuildSql_Insert(t *testing.T) {
	sql, err := NewModel(StockHsas{Code: "131", OpenToday: 111}).Insert()
	if err != nil {
		t.Error(err)
	}
	fmt.Println(sql)
}

//更新
func TestBuildSql_Update(t *testing.T) {
	sql, err := NewModel(StockHsas{Code: "131"}).Where("id", 1).Update()
	if err != nil {
		t.Error(err)
	}
	fmt.Println(sql)
}

//删除
func TestBuildSql_Delete(t *testing.T) {
	sql, err := NewModel(StockHsas{}).Where("id", 1).Delete()
	if err != nil {
		t.Error(err)
	}
	fmt.Println(sql)
}

//查询
func TestBuildSql_Select(t *testing.T) {
	sql, err := NewModel(StockHsas{}).Where("id", 1).Select()
	if err != nil {
		t.Error(err)
	}
	fmt.Println(sql)
}

func TestBuildSql_LeftJoin(t *testing.T) {

	sql, err := NewModel(StockHsas{}, "A").
		LeftJoin(StockHsas{}, "B", "A.code=B.code").
		Where("A.id", 1).
		Select()
	if err != nil {
		t.Error(err)
	}
	fmt.Println(sql)

}

func TestBuildSql_Field(t *testing.T) {

	sql, err := NewModel(StockHsas{}, "A").
		LeftJoin(StockHsas{}, "B", "A.code=B.code").
		Where("A.id", 1).
		Field("A.id", "A.code", "B.id").
		Select()
	if err != nil {
		t.Error(err)
	}
	fmt.Println(sql)

}

//func TestIsSatisfied(t *testing.T)  {
//
//	res:=IsSatisfied(float64(0))
//	fmt.Println(res)
//}

func TestGetColumnName(t *testing.T) {
	str := getColumnName("column(id);table(rms_hsas)", "table")

	fmt.Println(str)
}

func TestBuildSql_GetTableNameFromModel(t *testing.T) {

	tableName, err := GetTableNameFromModel(StockHsas{})
	if err != nil {
		t.Error(err)
	}

	fmt.Println(tableName)

}
