const FONT_TIMEOUT_MS = 2000;

const REQUIRED_FONTS = ["700 56px Caveat", "400 30px Caveat", '400 16px "Golos Text"'];

export function fontsLoaded(): boolean {
  if (!("fonts" in document)) return true;
  return REQUIRED_FONTS.every((font) => document.fonts.check(font));
}

export function waitForFonts(): Promise<void> {
  if (!("fonts" in document)) return Promise.resolve();

  const ready = Promise.all(REQUIRED_FONTS.map((font) => document.fonts.load(font))).then(() => undefined);

  const timeout = new Promise<void>((resolve) => setTimeout(resolve, FONT_TIMEOUT_MS));

  return Promise.race([ready, timeout]);
}
