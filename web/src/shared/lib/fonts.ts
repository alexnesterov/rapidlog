const FONT_TIMEOUT_MS = 2000;

export function waitForFonts(): Promise<void> {
  if (!("fonts" in document)) return Promise.resolve();

  const ready = Promise.all([
    document.fonts.load("700 56px Caveat"),
    document.fonts.load("italic 400 16px Lora"),
  ]).then(() => undefined);

  const timeout = new Promise<void>((resolve) => setTimeout(resolve, FONT_TIMEOUT_MS));

  return Promise.race([ready, timeout]);
}
