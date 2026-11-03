import assert from "node:assert/strict";
import test from "node:test";
import { findGrowthStage, getCropAgeDays, shouldLoadGrowthStageCatalog } from "./crop-progress.mjs";
import { isCurrentCycleFormSession } from "./cycle-form-session.mjs";

const localDate = (year, month, day) => new Date(year, month - 1, day, 12);

test("crop age is zero on the planting calendar date", () => {
  assert.equal(getCropAgeDays("2026-11-03", localDate(2026, 11, 3)), 0);
});

test("crop age counts elapsed calendar days without timezone or DST arithmetic", () => {
  assert.equal(getCropAgeDays("2026-11-02", localDate(2026, 11, 3)), 1);
  assert.equal(getCropAgeDays("2026-10-31", localDate(2026, 11, 3)), 3);
});

test("does not calculate crop age before a browser-local date is available", () => {
  assert.equal(getCropAgeDays("2026-11-02", null), null);
});

test("future and malformed planting dates do not produce an age", () => {
  assert.equal(getCropAgeDays("2026-11-04", localDate(2026, 11, 3)), null);
  assert.equal(getCropAgeDays("2026-02-29", localDate(2026, 3, 1)), null);
  assert.equal(getCropAgeDays("03-11-2026", localDate(2026, 11, 3)), null);
});

test("leap-day planting dates are accepted when valid", () => {
  assert.equal(getCropAgeDays("2024-02-29", localDate(2024, 3, 1)), 1);
});

test("growth-stage catalog eligibility is stable while an active crop age advances", () => {
  assert.equal(shouldLoadGrowthStageCatalog("active", 0), true);
  assert.equal(shouldLoadGrowthStageCatalog("active", 1), true);
  assert.equal(shouldLoadGrowthStageCatalog("active", 365), true);
  assert.equal(shouldLoadGrowthStageCatalog("planned", 0), false);
  assert.equal(shouldLoadGrowthStageCatalog("harvested", 1), false);
  assert.equal(shouldLoadGrowthStageCatalog("cancelled", 1), false);
  assert.equal(shouldLoadGrowthStageCatalog("active", null), false);
  assert.equal(shouldLoadGrowthStageCatalog("active", -1), false);
});

test("growth-stage range boundaries are inclusive", () => {
  const stages = [
    { name: "Awal", sequence_no: 1, min_days: 0, max_days: 7 },
    { name: "Lanjut", sequence_no: 2, min_days: 8, max_days: 20 },
  ];
  assert.equal(findGrowthStage(stages, 0)?.name, "Awal");
  assert.equal(findGrowthStage(stages, 7)?.name, "Awal");
  assert.equal(findGrowthStage(stages, 8)?.name, "Lanjut");
  assert.equal(findGrowthStage(stages, 21), null);
});

test("overlapping, unbounded, and invalid ranges are not presented as a stage", () => {
  assert.equal(findGrowthStage([
    { name: "A", sequence_no: 1, min_days: 0, max_days: 10 },
    { name: "B", sequence_no: 2, min_days: 8, max_days: 20 },
  ], 9), null);
  assert.equal(findGrowthStage([{ name: "Unbounded", sequence_no: 1 }], 4), null);
  assert.equal(findGrowthStage([{ name: "Open lower", sequence_no: 1, max_days: 10 }], 4), null);
  assert.equal(findGrowthStage([{ name: "Open upper", sequence_no: 1, min_days: 4 }], 4), null);
  assert.equal(findGrowthStage([
    { name: "Boundary A", sequence_no: 1, min_days: 0, max_days: 8 },
    { name: "Boundary B", sequence_no: 2, min_days: 8, max_days: 20 },
  ], 8), null);
  assert.equal(findGrowthStage([{ name: "Invalid", sequence_no: 1, min_days: 9, max_days: 2 }], 4), null);
});

test("late crop-cycle submissions only reset the form session that submitted them", () => {
  assert.equal(isCurrentCycleFormSession(4, 4), true);
  assert.equal(isCurrentCycleFormSession(4, 5), false);
  assert.equal(isCurrentCycleFormSession(Number.NaN, 5), false);
  assert.equal(isCurrentCycleFormSession(4, Number.NaN), false);
});

test("late crop-cycle failures are ignored after the form session changes", () => {
  assert.equal(isCurrentCycleFormSession(7, 7), true);
  assert.equal(isCurrentCycleFormSession(7, 8), false);
});
