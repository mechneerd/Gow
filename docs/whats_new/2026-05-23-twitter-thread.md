# GoW Feature Parity Release — Twitter Thread (Ready to Post)

Copy and paste each tweet in order.

---

**Tweet 1/8**

🚀 GoW just shipped its biggest release ever.

The **May 23, 2026 Feature Parity Release** is here.

After an intense implementation wave, GoW has closed almost every remaining critical gap.

The framework is now genuinely production-ready.

Thread 👇

---

**Tweet 2/8**

**Authentication & Security**

- Socialite (OAuth) — Google + GitHub support out of the box
- Two-Factor Authentication (pure Go TOTP)
- Account Lockout after failed attempts
- Remember Me tokens

You can now build modern, secure auth flows without writing everything yourself.

---

**Tweet 3/8**

**ORM Got Extremely Powerful**

New in this release:
- Polymorphic relations (MorphOne/Many/To)
- HasOneThrough / HasManyThrough
- Full Attribute Casting (datetime, json, bool, etc.)
- Accessors & Mutators
- Relation Touching, Upsert, Pessimistic Locking, Chunking, and more

Goquent is now one of the best ORMs in Go.

---

**Tweet 4/8**

**Developer Experience Upgrades**

- Dozens of new Artisan commands (`make:mail`, `make:policy`, `make:resource`, `cache:*`, `queue:retry`, `schedule:list`…)
- Signed URLs + middleware
- Named + per-user rate limiting
- Complete FormRequest validation with all hooks
- `actingAs()` + file upload testing helpers

Building apps just got a lot more pleasant.

---

**Tweet 5/8**

**Ecosystem Foundations**

We also added solid starting points for:
- Inertia.js
- Livewire-style components
- Scout (search)
- Cashier (billing)
- Telescope / Pulse
- Octane
- Prometheus metrics

The foundation is now there for advanced applications.

---

**Tweet 6/8**

**Documentation Overhaul**

New and heavily expanded guides:
- Socialite (OAuth)
- Testing
- Major updates to ORM, Authentication, and Routing

Plus a full **Upgrade Guide** for existing projects.

We made it easy to adopt everything.

---

**Tweet 7/8**

**Bottom Line**

Before this release, GoW was “promising”.

After today, GoW is **production-viable**.

You can now build real SaaS products, admin panels, and full-stack apps with a Laravel-like experience in Go — without constantly hitting missing pieces.

---

**Tweet 8/8**

Full release notes + upgrade guide:

https://github.com/.../gow/blob/main/docs/whats_new/2026-05-23.md

Upgrade guide: https://github.com/.../gow/blob/main/docs/UPGRADE.md

GitHub Release coming shortly.

We’re incredibly proud of this one.

Thank you to everyone following and supporting the project ❤️

#GoLang #WebDev #Laravel

---

**Optional extra tweet (if you want to go longer):**

We’ll continue polishing (better mail drivers, more notification channels, richer Inertia/Livewire experience, official starter kits), but the core is now extremely solid.

Go build something great with GoW.
