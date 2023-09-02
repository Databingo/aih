package main

var (
	hooks_buffer map[string][]func(*Buffer)
)

func hook_buffer(name string, f func(*Buffer)) {
	if _, ok := hooks_buffer[name]; !ok {
		hooks_buffer[name] = []func(*Buffer){}
	}
	hooks_buffer[name] = append(hooks_buffer[name], f)
}

func hook_trigger_buffer(name string, b *Buffer) {
	if hooks, ok := hooks_buffer[name]; ok {
		for _, f := range hooks {
			f(b)
		}
	}
}

func init_hooks() {
	hooks_buffer = map[string][]func(*Buffer){}

	hook_buffer("moved", func(b *Buffer) {
		if currentViewTree.Leaf.Buf == b {
			currentViewTree.Leaf.AdjustScroll(b.LastRenderWidth, b.LastRenderHeight)
		}
	})
}
