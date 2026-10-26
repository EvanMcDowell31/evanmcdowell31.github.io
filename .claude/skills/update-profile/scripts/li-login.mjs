// Usage: node li-login.mjs <state.json>
// Opens a visible browser at LinkedIn sign-in, waits for the user to log in,
// and saves the session to <state.json>. Skips the prompt if the saved session still works.
import { chromium } from 'playwright';
import { existsSync } from 'node:fs';

const state = process.argv[2];
if (!state) { console.error('usage: node li-login.mjs <state.json>'); process.exit(2); }

const browser = await chromium.launch({ headless: false });
const context = await browser.newContext({
  viewport: { width: 1280, height: 900 },
  ...(existsSync(state) ? { storageState: state } : {}),
});
const page = await context.newPage();

// Signed in once the li_at session cookie exists.
const signedIn = async () => (await context.cookies('https://www.linkedin.com')).some(c => c.name === 'li_at');

if (!(await signedIn())) {
  await page.goto('https://www.linkedin.com/login');
  console.log('Waiting for sign-in...');
  const deadline = Date.now() + 10 * 60 * 1000;
  while (!(await signedIn())) {
    if (Date.now() > deadline) { console.log('Timed out waiting for sign-in'); process.exit(1); }
    await page.waitForTimeout(1500);
  }
  await page.waitForTimeout(3000);
}
await context.storageState({ path: state });
console.log('Session saved');
await browser.close();
