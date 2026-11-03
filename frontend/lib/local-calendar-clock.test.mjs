import assert from "node:assert/strict";
import { spawnSync } from "node:child_process";
import test from "node:test";
import { millisecondsUntilNextLocalMidnight } from "./local-calendar-clock.mjs";

test("schedules the next local calendar midnight rather than a fixed 24-hour interval", () => {
  const now = new Date(2026, 2, 8, 0, 0, 0, 0);
  const nextMidnight = new Date(2026, 2, 9, 0, 0, 0, 0);
  assert.equal(millisecondsUntilNextLocalMidnight(now), nextMidnight.getTime() - now.getTime());
});

test("verifies 23-hour and 25-hour local days across New York DST transitions", () => {
  const moduleUrl = new URL("./local-calendar-clock.mjs", import.meta.url).href;
  const source = `
    import { millisecondsUntilNextLocalMidnight } from ${JSON.stringify(moduleUrl)};
    const spring = new Date(2026, 2, 8, 0, 0, 0, 0);
    const fall = new Date(2026, 10, 1, 0, 0, 0, 0);
    console.log(JSON.stringify({
      springHours: millisecondsUntilNextLocalMidnight(spring) / 3600000,
      fallHours: millisecondsUntilNextLocalMidnight(fall) / 3600000
    }));
  `;
  const result = spawnSync(process.execPath, ["--input-type=module", "-e", source], {
    encoding: "utf8",
    env: { ...process.env, TZ: "America/New_York" },
  });
  assert.equal(result.status, 0, result.stderr);
  assert.deepEqual(JSON.parse(result.stdout), { springHours: 23, fallHours: 25 });
});

test("returns the remaining local-day interval from the current time", () => {
  const now = new Date(2026, 10, 3, 18, 30, 0, 0);
  const nextMidnight = new Date(2026, 10, 4, 0, 0, 0, 0);
  assert.equal(millisecondsUntilNextLocalMidnight(now), nextMidnight.getTime() - now.getTime());
});

test("keeps the timer positive immediately before midnight", () => {
  const now = new Date(2026, 10, 3, 23, 59, 59, 999);
  assert.equal(millisecondsUntilNextLocalMidnight(now), 1);
});

test("returns null for an invalid date", () => {
  assert.equal(millisecondsUntilNextLocalMidnight(new Date(Number.NaN)), null);
});
