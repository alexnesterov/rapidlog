import type { ReactNode } from "react";

const URL_RE = /https?:\/\/[^\s<>"]+/g;
const TRAILING_PUNCT_RE = /[.,!?;:'")\]}]+$/;

function trimTrailingPunctuation(url: string): string {
  const match = TRAILING_PUNCT_RE.exec(url);
  if (!match) return url;

  let trimmed = match[0];
  let core = url.slice(0, url.length - trimmed.length);

  // Keep a single closing paren that balances an opening one inside the URL
  // itself, e.g. https://en.wikipedia.org/wiki/Foo_(disambiguation).
  if (trimmed.startsWith(")")) {
    const opens = (core.match(/\(/g) ?? []).length;
    const closes = (core.match(/\)/g) ?? []).length;
    if (opens > closes) {
      core += ")";
      trimmed = trimmed.slice(1);
    }
  }

  return core;
}

export function linkifyText(text: string): ReactNode[] {
  const parts: ReactNode[] = [];
  let lastIndex = 0;
  let key = 0;

  URL_RE.lastIndex = 0;
  let match: RegExpExecArray | null;
  while ((match = URL_RE.exec(text)) !== null) {
    const url = trimTrailingPunctuation(match[0]);
    if (url.length === 0) continue;

    const start = match.index;
    const end = start + url.length;

    if (start > lastIndex) {
      parts.push(text.slice(lastIndex, start));
    }

    parts.push(
      <a key={`url-${key++}`} href={url} target="_blank" rel="noopener noreferrer">
        {url}
      </a>,
    );

    lastIndex = end;
    URL_RE.lastIndex = end;
  }

  if (lastIndex < text.length) {
    parts.push(text.slice(lastIndex));
  }

  return parts.length > 0 ? parts : [text];
}
