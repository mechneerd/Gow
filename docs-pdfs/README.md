# GoW Documentation - PDF Versions

This folder contains clean, PDF-friendly Markdown versions of the most important beginner guides for the GoW Framework.

> **Note**: These files now include the **"What Can You Build With GoW?"** section, which explains all supported tech stacks (API only, React + Inertia, Vue + Inertia, Goblade, HTMX, etc.).

## Available PDF Sources

| File | Description | Best For |
|------|-------------|----------|
| `GoW_Getting_Started_Detailed.md` | Very detailed step-by-step guide (zero experience → React/Vue apps) | New users who want everything explained |
| `GoW_Quick_Start.md` | Fast 10-20 minute guide with tech stack options | People who want to get running quickly |
| `GoW_Quick_Start_One_Page.md` | Ultra-short printable version (fits on 1 page) | Quick reference / printing while coding |

---

## How to Convert These to Real PDFs

### Easiest Method (No Installation)

1. Open any of the `.md` files in **VS Code**.
2. Press `Ctrl + P` (or `Cmd + P` on Mac).
3. Type `Markdown Preview Enhanced: Print to PDF` (you may need to install the "Markdown Preview Enhanced" extension first).
4. Or simply open the file in browser and use **Print → Save as PDF**.

### Recommended Tools

**Option 1: VS Code + Extension (Best for most people)**
- Install extension: **Markdown Preview Enhanced**
- Right-click the Markdown file → "Markdown Preview Enhanced: Print to PDF"

**Option 2: Pandoc (Most Professional)**
Install from: https://pandoc.org

Then run:

```bash
# Detailed guide
pandoc GoW_Getting_Started_Detailed.md -o "GoW_Getting_Started_Detailed.pdf" --pdf-engine=xelatex -V geometry:margin=1in

# Quick Start
pandoc GoW_Quick_Start.md -o "GoW_Quick_Start.pdf" --pdf-engine=xelatex -V geometry:margin=0.8in
```

**Option 3: Online Converters**
- https://md-to-pdf.fly.dev
- https://www.markdowntopdf.com

Just drag and drop the `.md` file.

---

## Suggested PDF Naming for Distribution

When converting, consider renaming the output files like this:

- `GoW_Getting_Started_Detailed.pdf`
- `GoW_Quick_Start_Guide.pdf`
- `GoW_Quick_Start_One_Page_Cheat_Sheet.pdf`

---

## When to Use Each Version

- **New user with no experience** → Give them `GoW_Getting_Started_Detailed.pdf`
- **Someone who wants to try quickly** → Give them `GoW_Quick_Start.pdf`
- **Print-and-follow while coding** → Give them `GoW_Quick_Start_One_Page.pdf`

---

Last updated: 2026-05-23 (Feature Parity Release)
