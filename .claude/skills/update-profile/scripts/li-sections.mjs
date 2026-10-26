// Usage: node li-sections.mjs <state.json> <out-dir>
// Headless: reads the profile headline/About and each details section with a saved
// session, writing li-<section>.txt files to <out-dir>. Exits 3 if the session expired.
import { chromium } from 'playwright';
import { mkdirSync, writeFileSync } from 'node:fs';
import { join } from 'node:path';

const [state, out] = process.argv.slice(2);
if (!state || !out) { console.error('usage: node li-sections.mjs <state.json> <out-dir>'); process.exit(2); }
mkdirSync(out, { recursive: true });

const BASE = 'https://www.linkedin.com/in/evanmcdowell31';
const browser = await chromium.launch();
const context = await browser.newContext({ storageState: state, viewport: { width: 1280, height: 900 } });
const page = await context.newPage();

const grab = async (url, name) => {
  await page.goto(url, { waitUntil: 'load' });
  await page.waitForTimeout(3000);
  if (/\/(login|authwall|checkpoint)/.test(page.url())) {
    console.log('Session expired; run li-login.mjs again');
    await browser.close();
    process.exit(3);
  }
  for (let i = 0; i < 6; i++) { await page.mouse.wheel(0, 1500); await page.waitForTimeout(700); }
  const text = await page.locator('main').innerText().catch(() => '');
  writeFileSync(join(out, `li-${name}.txt`), `${page.url()}\n\n${text}`);
  console.log(name, text.length);
};

await grab(`${BASE}/`, 'profile');
for (const section of ['experience', 'education', 'skills', 'certifications', 'projects']) {
  await grab(`${BASE}/details/${section}/`, section);
}
await browser.close();
