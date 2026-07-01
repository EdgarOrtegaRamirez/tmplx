# Security

## Template Injection

tmplx processes user-provided template strings. If templates come from untrusted sources, ensure:

1. **Sanitize input**: Validate template syntax before rendering
2. **Limit complexity**: Set maximum template size and nesting depth
3. **Sandbox partials**: Only register trusted partials, never load from user paths

## Variable Data

Variables passed to templates should be validated and sanitized. Template filters do not escape HTML by default - use `html_escape` filter when rendering user content to HTML.

## Dependencies

tmplx has zero external dependencies, reducing the attack surface significantly. All functionality is implemented in pure Go using the standard library only.

## Reporting

If you discover a security vulnerability, please report it responsibly. Do not open a public issue.
