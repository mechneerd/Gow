package livewire

import "fmt"

// Counter is a simple Livewire-style reactive component example.
// Place this in app/Livewire/Counter.go in your project.
type Counter struct {
	BaseComponent
	Count int
}

func (c *Counter) Mount() {
	c.Count = 0
}

func (c *Counter) Increment() {
	c.Count++
	c.MarkDirty("Count")
}

func (c *Counter) Decrement() {
	c.Count--
	c.MarkDirty("Count")
}

func (c *Counter) Reset() {
	c.Count = 0
	c.MarkDirty("Count")
}

func (c *Counter) Render() string {
	return fmt.Sprintf(`
<div wire:id="%s" class="p-6 bg-zinc-900 border border-zinc-700 rounded-3xl max-w-sm mx-auto">
    <div class="text-center">
        <div class="text-sm text-zinc-400 mb-1">Livewire Counter</div>
        <div class="text-6xl font-semibold tabular-nums tracking-tighter mb-6">%d</div>

        <div class="flex items-center justify-center gap-x-3">
            <button wire:click="Decrement" 
                    class="px-5 py-2.5 bg-zinc-800 hover:bg-zinc-700 rounded-2xl text-sm font-medium transition">
                −
            </button>
            <button wire:click="Increment" 
                    class="px-5 py-2.5 bg-white text-zinc-950 hover:bg-zinc-100 rounded-2xl text-sm font-semibold transition">
                +
            </button>
            <button wire:click="Reset" 
                    class="px-4 py-2.5 text-sm text-zinc-400 hover:text-white transition">
                Reset
            </button>
        </div>
    </div>
</div>`, c.GetID(), c.Count)
}
