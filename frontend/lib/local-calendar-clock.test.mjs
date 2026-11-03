import assert from "node:assert/strict";
import test from "node:test";
import { millisecondsUntilNextLocalMidnight } from "./local-calendar-clock.mjs";

test("schedules the next local calendar midnight rather than a fixed 24-hour interval", () => {
  const now = new Date(2026, 2, 8, 0, 0, 0, 0);
  const nextMidnight = new Date(2026, 2, 9, 0, 0, 0, 0);
  assert.equal(millisecondsUntilNextLocalMidnight(now), nextMidnight.getTime() - now.getTime());
});

test("handles a local day that is longer than 24 hours when DST ends", () => {
  const now = new Date(2026, 10, 1, 0, 0, 0, 0);
  const nextMidnight = new Date(2026, 10, 2, 0, 0, 0, 0);
  assert.equal(millisecondsUntilNextLocalMidnight(now), nextMidnight.getTime() - now.getTime());
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
