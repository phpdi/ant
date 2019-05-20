package validation

import (
	"ant"
	"testing"
)

func TestMax(t *testing.T) {

	a:=ant.StockHsas{}
	v:=NewValidation()

	//零值不验证
	err:=v.Valid(a)
	if err != nil {
		t.Error(err)
	}
	t.Log("零值不验证,通过")

	//非零值,小于最大值验证
	a.Id= 9
	err=v.Clear().Valid(a)
	if err != nil {
		t.Error(err)
	}
	t.Log("非零值,小于最大值验证,通过")


	//非零值,小于最大值验证
	a.Id= 20
	err=v.Clear().Valid(a)
	if err != nil {
		t.Error(err)
	}
	t.Log("非零值,等于最大值验证,通过")

	//非零值,大于最大值
	a.Id= 21
	err=v.Clear().Valid(a)
	if err != nil {
		t.Log("非零值,大于最大值验证,通过")
	}

}