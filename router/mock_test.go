package router

type mockTrier struct {
	get route
}

func (m mockTrier) Get(string) interface{} {
	return m.get
}

func (m mockTrier) Put(string, interface{}) {
}
