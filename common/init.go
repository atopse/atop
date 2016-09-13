package common

//hook function to run
type hookfunc func() error

var (
	hooks      = make(map[int][]hookfunc, 0) //hook function slice to store the hookfunc
	maxPriotiy = 10
)

// AddStartHook is used to register the hookfunc
// The hookfuncs will run in RunStartHook()
// priority 范围 0-10；执行顺序为：0->10.
func AddStartHook(hf hookfunc, priority ...int) {
	p := maxPriotiy
	if len(priority) > 0 {
		p = priority[0]
		if p < 0 {
			p = 0
		} else if p > maxPriotiy {
			p = maxPriotiy
		}
	}
	hooks[p] = append(hooks[p], hf)
}

// RunStartHook is
func RunStartHook() {
	if len(hooks) == 0 {
		return
	}
	for i := 0; i <= maxPriotiy; i++ {
		fs := hooks[i]
		for _, f := range fs {
			if err := f(); err != nil {
				panic(err)
			}
		}
	}

}
