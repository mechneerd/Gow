// GoW Livewire JavaScript Client
// This provides basic Livewire-like behavior: wire:click and wire:model

(function () {
    'use strict';

    window.Livewire = {
        components: {},

        // Initialize Livewire on the page
        start: function () {
            document.addEventListener('DOMContentLoaded', () => {
                this.scan();
            });
        },

        // Scan DOM for livewire components and attach listeners
        scan: function () {
            document.querySelectorAll('[wire\\:id]').forEach(el => {
                const id = el.getAttribute('wire:id');
                if (!id) return;

                this.components[id] = { el: el };

                // Attach wire:click handlers
                el.querySelectorAll('[wire\\:click]').forEach(btn => {
                    const method = btn.getAttribute('wire:click');
                    btn.addEventListener('click', (e) => {
                        e.preventDefault();
                        this.call(id, method);
                    });
                });

                // Attach wire:model handlers (two-way binding)
                el.querySelectorAll('[wire\\:model]').forEach(input => {
                    const property = input.getAttribute('wire:model');

                    input.addEventListener('input', () => {
                        this.updateProperty(id, property, input.value);
                    });

                    input.addEventListener('change', () => {
                        this.updateProperty(id, property, input.value);
                    });
                });

                // Attach wire:submit handlers (for forms)
                el.querySelectorAll('form[wire\\:submit]').forEach(form => {
                    const method = form.getAttribute('wire:submit');
                    form.addEventListener('submit', (e) => {
                        e.preventDefault();
                        const formData = new FormData(form);
                        const data = {};
                        formData.forEach((value, key) => data[key] = value);
                        this.call(id, method, [data]);
                    });
                });
            });
        },

        // Call a method on the component
        call: function (id, method, params = []) {
            this.sendUpdate(id, method, params);
        },

        // Update a property and sync to server
        updateProperty: function (id, property, value) {
            // Optimistic update in DOM
            const component = this.components[id];
            if (!component) return;

            // Send update to server
            this.sendUpdate(id, null, [], { [property]: value });
        },

        // Core function that talks to the backend
        sendUpdate: function (id, method, params, stateUpdate = {}) {
            const componentEl = document.querySelector(`[wire\\:id="${id}"]`);
            if (!componentEl) return;

            // Get current state from data attributes or hidden inputs (simplified)
            const payload = {
                id: id,
                method: method || '',
                params: params || [],
                state: stateUpdate
            };

            const componentEl = document.querySelector(`[wire\\:id="${id}"]`);

            // Show loading states
            if (componentEl) {
                componentEl.querySelectorAll('[wire\\:loading]').forEach(el => {
                    el.style.display = '';
                });
            }

            fetch('/livewire/update', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Accept': 'application/json'
                },
                body: JSON.stringify(payload)
            })
            .then(res => res.json())
            .then(data => {
                if (data.html && data.id) {
                    this.replaceComponent(data.id, data.html);
                }

                // Hide loading states after update
                const newEl = document.querySelector(`[wire\\:id="${id}"]`);
                if (newEl) {
                    newEl.querySelectorAll('[wire\\:loading]').forEach(el => {
                        el.style.display = 'none';
                    });
                }
            })
            .catch(err => {
                console.error('Livewire update failed:', err);
                if (componentEl) {
                    componentEl.querySelectorAll('[wire\\:loading]').forEach(el => {
                        el.style.display = 'none';
                    });
                }
            });
        },

        // Replace the component HTML with new content from server
        replaceComponent: function (id, newHTML) {
            const oldEl = document.querySelector(`[wire\\:id="${id}"]`);
            if (!oldEl) return;

            const wrapper = document.createElement('div');
            wrapper.innerHTML = newHTML.trim();
            const newEl = wrapper.firstElementChild;

            if (newEl) {
                // Preserve elements with wire:ignore
                const ignored = oldEl.querySelectorAll('[wire\\:ignore]');
                ignored.forEach(el => {
                    const selector = this.getUniqueSelector(el);
                    const newIgnored = newEl.querySelector(selector);
                    if (newIgnored) {
                        newIgnored.replaceWith(el.cloneNode(true));
                    }
                });

                oldEl.replaceWith(newEl);
                this.scanNewElement(newEl);
            }
        },

        // Helper to generate a simple selector for ignored elements
        getUniqueSelector: function (el) {
            if (el.id) return '#' + el.id;
            if (el.className) return '.' + el.className.split(' ').join('.');
            return el.tagName.toLowerCase();
        },

        scanNewElement: function (el) {
            // Re-attach listeners to the newly rendered component
            const id = el.getAttribute('wire:id');

            // wire:click
            el.querySelectorAll('[wire\\:click]').forEach(btn => {
                const method = btn.getAttribute('wire:click');
                btn.addEventListener('click', (e) => {
                    e.preventDefault();
                    this.call(id, method);
                });
            });

            // wire:model
            el.querySelectorAll('[wire\\:model]').forEach(input => {
                const property = input.getAttribute('wire:model');
                input.addEventListener('input', () => {
                    this.updateProperty(id, property, input.value);
                });
                input.addEventListener('change', () => {
                    this.updateProperty(id, property, input.value);
                });
            });
        }
    };

    // Auto-start Livewire
    if (document.readyState === 'loading') {
        document.addEventListener('DOMContentLoaded', () => window.Livewire.start());
    } else {
        window.Livewire.start();
    }
})();
