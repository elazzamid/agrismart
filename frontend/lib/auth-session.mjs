// A 401 may invalidate only the exact session token used by that request.
// Stale responses must not clear a newer login or a session already signed out.
export function shouldInvalidateSession(requestToken, currentToken) {
  return typeof requestToken === "string" && requestToken.length > 0 && requestToken === currentToken;
}
