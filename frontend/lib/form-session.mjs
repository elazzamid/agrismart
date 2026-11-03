export function isCurrentFormSession(submittedSessionId, currentSessionId) {
  return Number.isSafeInteger(submittedSessionId)
    && Number.isSafeInteger(currentSessionId)
    && submittedSessionId === currentSessionId;
}
