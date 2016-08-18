package common

//hook function to run
type hookfunc func() error

var (
	hooks      = make(map[int][]hookfunc, 0) //hook function slice to store the hookfunc
	maxPriotiy = 10
)

// AddAPPStartHook is used to register the hookfunc
// The hookfuncs will run in RunStartHook()
// such as load Config
func AddAPPStartHook(priority int, hf hookfunc) {
	if priority < 0 {
		priority = 0
	} else if priority > maxPriotiy {
		priority = maxPriotiy
	}
	hooks[priority] = append(hooks[priority], hf)
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
