package option

type RangeOption struct {
	DelFlag   bool   //遍历完是否删除该key
	SleepMs   int    //处理2个Key的间隔，单位ms；为 0 表示不 sleep
	MaxCount  int    //最多处理的Key的个数，为 0 表示全部
	LoadValue bool   //是否加载Value, true 加载
	Prefix    string //前缀遍历
}

func DefaultOption() *RangeOption {
	return &RangeOption{
		DelFlag:   false,
		SleepMs:   0,
		MaxCount:  0,
		LoadValue: true,
		Prefix:    "",
	}
}

func NewRangeOption(del bool, sleep int, max int, loadV bool, prefix string) *RangeOption {
	return &RangeOption{
		DelFlag:   del,
		SleepMs:   sleep,
		MaxCount:  max,
		LoadValue: loadV,
		Prefix:    prefix,
	}
}
