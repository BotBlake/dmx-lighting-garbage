package device

type Registry struct {
	devices map[string]*Instance
}

func NewRegistry() *Registry {
	return &Registry{devices: make(map[string]*Instance)}
}

func (r *Registry) Add(d *Instance) {
	r.devices[d.Name] = d
}

func (r *Registry) Get(name string) (*Instance, bool) {
	d, ok := r.devices[name]
	return d, ok
}
