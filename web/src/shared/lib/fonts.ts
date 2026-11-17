const FONT_TIMEOUT_MS = 2000;

// всё, что видно на первом экране: логотип и записи (Caveat), текст, бейдж и кнопки (Golos 400/500/600)
const REQUIRED_FONTS = [
  "700 56px Caveat",
  "400 30px Caveat",
  '400 16px "Golos Text"',
  '500 16px "Golos Text"',
  '600 16px "Golos Text"',
];

// Google Fonts режет каждый шрифт на файлы по unicode-range (латиница, кириллица…).
// Без текста load смотрит только на пробел — грузится латинский файл, а кириллица
// догружается уже после показа страницы и подменяется на глазах. Пример текста — из обоих наборов.
const SAMPLE_TEXT = "RapidLog Журнал";

export function waitForFonts(): Promise<void> {
  if (!("fonts" in document)) return Promise.resolve();

  const ready = Promise.all(REQUIRED_FONTS.map((font) => document.fonts.load(font, SAMPLE_TEXT))).then(() => undefined);

  const timeout = new Promise<void>((resolve) => setTimeout(resolve, FONT_TIMEOUT_MS));

  return Promise.race([ready, timeout]);
}
