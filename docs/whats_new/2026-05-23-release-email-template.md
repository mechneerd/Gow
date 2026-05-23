# GoW Framework — Release Email / Newsletter Template

**Subject Line Options** (pick one):

1. GoW Feature Parity Release – Now Production Ready
2. The Biggest GoW Update Ever Just Landed
3. GoW v0.9.0 – Socialite, 2FA, Advanced ORM & More
4. GoW is Ready for Real Production Apps

---

**Email Body (HTML + Plain Text friendly)**

---

Hi there,

We’re thrilled to announce the **May 23, 2026 Feature Parity Release** of the GoW Framework.

This is the largest and most important release in GoW’s history. After an intensive implementation wave, we’ve closed almost every remaining critical gap that was holding the framework back from real-world production use.

### What’s New

**Authentication & Security**
- Official Socialite (OAuth) support – Google & GitHub ready
- Two-Factor Authentication (pure Go TOTP)
- Account Lockout protection
- Remember Me tokens

**ORM – Major Upgrades**
- Polymorphic relations (MorphOne/Many/To)
- HasOneThrough / HasManyThrough
- Full Attribute Casting system
- Accessors & Mutators
- Relation Touching, Upsert, Pessimistic Locking, Chunking, and more

**Developer Experience**
- Dozens of new Artisan commands (make:mail, make:policy, make:resource, cache:*, queue:*, schedule:list, etc.)
- Signed URLs + middleware
- Advanced named rate limiting
- Complete FormRequest validation
- Much better testing tools (`actingAs()`, file uploads, JSON assertions)

**Ecosystem Foundations**
- Inertia.js support
- Livewire-style components
- Scout, Cashier, Telescope, Octane, and Prometheus metrics starters

We’ve also published new and heavily expanded documentation, including a full **Upgrade Guide** for existing projects.

### Why This Matters

Before this release, many teams hit hard limits when trying to build serious applications with GoW.

After today, GoW is genuinely production-viable for SaaS products, admin panels, APIs, and full-stack applications — with a Laravel-like experience in Go.

### Get Started

- Full Release Notes: https://github.com/your-org/gow/blob/main/docs/whats_new/2026-05-23.md
- Upgrade Guide: https://github.com/your-org/gow/blob/main/docs/UPGRADE.md
- GitHub Release: [Coming soon]

We’re incredibly proud of how far GoW has come with this release.

Thank you for being part of the journey.

Happy building,

The GoW Team

---

**Plain Text Version** (for email clients that strip HTML)

```
Hi there,

We’re thrilled to announce the May 23, 2026 Feature Parity Release of the GoW Framework.

This is the largest and most important release in GoW’s history. We’ve closed almost every remaining critical gap.

Highlights:
- Socialite (OAuth) – Google & GitHub
- Two-Factor Authentication (TOTP)
- Account Lockout
- Advanced ORM (Polymorphic, Casting, Accessors, Through relations, Locking, etc.)
- Dozens of new Artisan commands
- Signed URLs, better rate limiting, full FormRequest
- Testing improvements (actingAs, file uploads, assertions)
- New docs + Upgrade Guide

GoW is now genuinely production-ready for real applications.

Full details + upgrade guide:
https://github.com/your-org/gow/blob/main/docs/whats_new/2026-05-23.md

Thank you for supporting the project.

The GoW Team
```

---

**Recommended Email Settings**

- From: GoW Team <team@gow.dev> (or your official address)
- Reply-To: team@gow.dev
- Preheader: The biggest GoW release yet – Socialite, 2FA, advanced ORM, and more

---

**Optional Add-ons**

- Link to the Twitter thread
- Link to the GitHub Release once published
- Call-to-action button: “Read the Release Notes” pointing to the full notes

---

This template is ready to copy into your email tool (Mailchimp, ConvertKit, Brevo, etc.).
