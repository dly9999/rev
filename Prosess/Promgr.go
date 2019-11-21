package Prosess

type MgrquePro struct {
	Mgr map[string]ProInterface
}

var mgr = MgrquePro{
	make(map[string]ProInterface),
}

func (p MgrquePro) registerProc(fucname string, pro ProInterface) {
	p.Mgr[fucname] = pro
}
func RegisterProc(fucname string, pro ProInterface) {
	mgr.registerProc(fucname, pro)
}
func DoProsess(fucname string) ProInterface {
	return mgr.Mgr[fucname]
}
