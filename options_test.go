package skiplist

import "testing"

func TestInitOption_Apply(t *testing.T) {
	_ = New[string, string](MaxLevels(-1), Concurrent(false))
	_ = New[string, string](MaxLevels(99), Concurrent(false))
	_ = New[string, string](MaxLevels(24), Concurrent(false))

}

func TestSetOption_Apply(t *testing.T) {
	list := New[string, string]()
	_ = list.Set("user-001", "XXX")
	if err := list.Set("user-001", "XXX", OnNotExist()); err != ErrDuplicatedKey {
		t.Error("expected ErrDuplicate")
	}
}

func TestRangeOption_Apply(t *testing.T) {
	list := New[string, string]()

	list.Range(func(key string, value string) bool {
		return true
	}, From("user-003", true), To("user-034", true))

	_ = list.Set("user-002", "XXX")
	_ = list.Set("user-003", "XXX")
	_ = list.Set("user-021", "XXX")
	_ = list.Set("user-034", "XXX")
	_ = list.Set("user-062", "XXX")

	list.Range(func(key string, value string) bool {
		return true
	}, From("user-003", true), To("user-034", true))

	list.Range(func(key string, value string) bool {
		return true
	}, From("user-003", false), To("user-034", false))

	list.Range(func(key string, value string) bool {
		return false
	}, From("user-003", false), To("user-034", false))
}

func TestGetOption_Apply(t *testing.T) {
	list := New[string, string]()
	_ = list.Set("1", "1")
	_, err := list.Get("2", GetOrDefault())
	if err != nil {
		t.Errorf("expected err=nil, got err=%v", err)
	}
}
