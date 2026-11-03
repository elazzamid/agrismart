import assert from "node:assert/strict";
import test from "node:test";
import { shouldInvalidateSession } from "./auth-session.mjs";

test("a 401 invalidates the session when it belongs to the current token", () => {
  assert.equal(shouldInvalidateSession("token-a", "token-a"), true);
});

test("a stale 401 cannot clear a newer login", () => {
  assert.equal(shouldInvalidateSession("token-a", "token-b"), false);
});

test("a stale 401 after sign-out cannot invalidate a session", () => {
  assert.equal(shouldInvalidateSession("token-a", null), false);
});

test("an unauthenticated login or registration 401 cannot clear a session", () => {
  assert.equal(shouldInvalidateSession(null, "token-b"), false);
  assert.equal(shouldInvalidateSession(null, null), false);
});

test("empty request tokens are not treated as authenticated sessions", () => {
  assert.equal(shouldInvalidateSession("", ""), false);
});
