import { chromium } from "playwright";
import fs from "fs";

const outDir = "/opt/cursor/artifacts";
fs.mkdirSync(outDir, { recursive: true });

const browser = await chromium.launch({ headless: true });
const page = await browser.newPage({ viewport: { width: 1440, height: 900 } });
const errors = [];
page.on("pageerror", (e) => errors.push(String(e)));
page.on("console", (msg) => {
  if (msg.type() !== "error") return;
  const t = msg.text();
  if (/ERR_CONNECTION_REFUSED|7350|nakama/i.test(t)) return;
  errors.push(`console: ${t}`);
});

await page.goto("http://127.0.0.1:5173/", { waitUntil: "domcontentloaded", timeout: 30000 });
await page.waitForSelector(".landing__brand", { timeout: 20000 });
await page.screenshot({ path: `${outDir}/01-landing.png`, fullPage: true });

await page.fill('input[aria-label="Имя"]', "Тестер");
await page.click('button[type="submit"]');
await page.waitForSelector(".shell", { timeout: 20000 });
await page.waitForTimeout(1500);
await page.screenshot({ path: `${outDir}/02-game.png`, fullPage: true });

const buy = page.locator(".list button", { hasText: "Купить" }).first();
if (await buy.isEnabled()) {
  await buy.click();
  await page.waitForTimeout(800);
}
await page.click('nav.tabs button:has-text("Фирмы")');
await page.waitForTimeout(400);
await page.screenshot({ path: `${outDir}/03-firms.png`, fullPage: true });

const rootHTML = await page.locator("#root").innerHTML();
const ok = rootHTML.includes("404Eco") || rootHTML.includes("hud__brand");
console.log(JSON.stringify({ ok, errors, rootLen: rootHTML.length }, null, 2));
if (!ok || errors.length) process.exitCode = 1;
await browser.close();
