import assert from "node:assert/strict";
import test from "node:test";
import { isCurrentFormSession } from "./form-session.mjs";

test("the current journal form accepts its own save result", () => {
  assert.equal(isCurrentFormSession(4, 4), true);
});

test("a delayed journal success cannot reset or publish success feedback to a newer form", () => {
  assert.equal(isCurrentFormSession(4, 5), false);
});

test("a delayed crop-cycle success cannot publish success feedback to a newer form", () => {
  assert.equal(isCurrentFormSession(6, 7), false);
});

test("a delayed journal failure cannot replace newer-form feedback", () => {
  assert.equal(isCurrentFormSession(7, 8), false);
});

test("invalid or unsafe session identifiers are rejected", () => {
  assert.equal(isCurrentFormSession(Number.NaN, 5), false);
  assert.equal(isCurrentFormSession(4, Number.NaN), false);
  assert.equal(isCurrentFormSession(Number.MAX_SAFE_INTEGER + 1, Number.MAX_SAFE_INTEGER + 1), false);
});

test("the current farm form accepts its own save result", () => {
  assert.equal(isCurrentFormSession(9, 9), true);
});

test("a delayed farm save cannot reset or publish success feedback to a newer form", () => {
  assert.equal(isCurrentFormSession(9, 10), false);
});

test("the current plot form accepts its own save result", () => {
  assert.equal(isCurrentFormSession(11, 11), true);
});

test("a delayed plot save cannot replace newer-form feedback", () => {
  assert.equal(isCurrentFormSession(11, 12), false);
});
