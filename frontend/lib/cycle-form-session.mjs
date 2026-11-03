export function isCurrentCycleFormSession(submittedSessionId, currentSessionId) {
  return Number.isSafeInteger(submittedSessionId)
    && Number.isSafeInteger(currentSessionId)
    && submittedSessionId === currentSessionId;
}
