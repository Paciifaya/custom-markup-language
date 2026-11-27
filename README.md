# MarkLang

MarkLang is a small custom markup language written in Go that converts `.mark` files into HTML.

## Features

- Headings, paragraphs, lists, quotes, and code blocks
- Inline emphasis, code, and links
- Custom callout blocks
- HTML escaping

## Requirements

- Go 1.21 or newer

## Getting started

```bash
go mod init custom-markup-language
go run . example.mark
```

## Usage

```bash
go run . <file.mark>
```

## Example

Create `example.mark`:

```mark
# My Notes

This is *important* and `useful`.

- Learn Go
- Build a parser

> Keep it simple.

[facebook](https://www.facebook.com)

::warning
Be careful with HTML escaping.
::
```

Run:

```bash
go run . example.mark
```

Output:

```html
<h1>My Notes</h1>
<p>This is <em>important</em> and <code>useful</code>.</p>
<ul>
  <li>Learn Go</li>
  <li>Build a parser</li>
</ul>
<blockquote>Keep it simple.</blockquote>
<p><a href="[https://www.facebook.com](https://www.facebook.com)">Facebook</a></p>
<div class="callout warning">
<p>Be careful with HTML escaping.</p>
</div>
```

## Syntax

| Syntax | HTML |
|---|---|
| `# Text` | `<h1>` |
| `## Text` | `<h2>` |
| Plain line | `<p>` |
| `- Text` | `<li>` |
| `> Text` | `<blockquote>` |
| ` ``` ` | `<pre><code>` |
| `*text*` | `<em>` |
| `` `code` `` | `<code>` |
| `[text](url)` | `<a href="url">` |
| `::name` … `::` | `<div class="callout name">` |

