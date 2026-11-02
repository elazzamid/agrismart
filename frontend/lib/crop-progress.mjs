const DAY_MS = 24 * 60 * 60 * 1000;

function parseDateOnly(value) {
  if (typeof value !== "string") return null;
  const match = /^(\d{4})-(\d{2})-(\d{2})$/.exec(value);
  if (!match) return null;

  const year = Number(match[1]);
  const month = Number(match[2]);
  const day = Number(match[3]);
  if (year < 1000 || month < 1 || month > 12 || day < 1 || day > 31) return null;

  const timestamp = Date.UTC(year, month - 1, day);
  const parsed = new Date(timestamp);
  if (parsed.getUTCFullYear() !== year || parsed.getUTCMonth() !== month - 1 || parsed.getUTCDate() !== day) return null;
  return timestamp;
}

// Uses the browser's local calendar date and UTC date arithmetic to avoid DST drift.
export function getCropAgeDays(plantingDate, today = new Date()) {
  if (!(today instanceof Date) || Number.isNaN(today.getTime())) return null;
  const plantedAt = parseDateOnly(plantingDate);
  if (plantedAt === null) return null;

  const todayAt = Date.UTC(today.getFullYear(), today.getMonth(), today.getDate());
  const elapsedDays = Math.floor((todayAt - plantedAt) / DAY_MS);
  return elapsedDays < 0 ? null : elapsedDays;
}

// Return a stage only when exactly one valid catalog range matches the age.
// Overlapping or unbounded sample ranges are deliberately treated as ambiguous.
export function findGrowthStage(stages, ageDays) {
  if (!Array.isArray(stages) || !Number.isInteger(ageDays) || ageDays < 0) return null;

  const matches = stages.filter(stage => {
    if (!stage || typeof stage !== "object") return false;
    const min = stage.min_days ?? null;
    const max = stage.max_days ?? null;
    if (min === null || max === null) return false;
    if (!Number.isInteger(min) || min < 0) return false;
    if (!Number.isInteger(max) || max < 0) return false;
    if (min > max) return false;
    return ageDays >= min && ageDays <= max;
  });

  return matches.length === 1 ? matches[0] : null;
}
