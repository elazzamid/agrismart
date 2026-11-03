// Schedule refresh against the next local calendar midnight, not a fixed 24 hours.
// Local days can be shorter or longer when daylight-saving time changes.
export function millisecondsUntilNextLocalMidnight(now = new Date()) {
  if (!(now instanceof Date) || Number.isNaN(now.getTime())) return null;
  const nextMidnight = new Date(now.getFullYear(), now.getMonth(), now.getDate() + 1);
  return Math.max(1, nextMidnight.getTime() - now.getTime());
}
