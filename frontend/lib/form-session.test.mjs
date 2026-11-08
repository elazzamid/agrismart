import assert from "node:assert/strict";
import test from "node:test";
import { getDismissedFarmSaveOutcome, isCurrentFormSession, shouldOpenFarmForm, shouldShowFarmSaveNotice, shouldShowFarmSavePending } from "./form-session.mjs";

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

test("a delayed farm failure cannot replace feedback for a newer form", () => {
  assert.equal(isCurrentFormSession(9, 10), false);
});

test("a current farm save shows success feedback while the dashboard is visible", () => {
  assert.equal(shouldShowFarmSaveNotice(12, 12, true), true);
});

test("a current farm save suppresses success feedback while the dashboard is hidden", () => {
  assert.equal(shouldShowFarmSaveNotice(12, 12, false), false);
});

test("a stale farm save suppresses success feedback even when the dashboard is visible", () => {
  assert.equal(shouldShowFarmSaveNotice(12, 13, true), false);
});

test("an unsafe farm form session cannot publish success feedback", () => {
  assert.equal(shouldShowFarmSaveNotice(Number.NaN, 13, true), false);
});

test("a closed farm form can start a new session", () => {
  assert.equal(shouldOpenFarmForm(false, false), true);
});

test("reopening an already-open form cannot invalidate a pending farm save", () => {
  assert.equal(shouldOpenFarmForm(true, true, true), false);
});

test("an explicit mobile add action can start a new session after a failed save", () => {
  assert.equal(shouldOpenFarmForm(true, false, true), true);
});

test("a repeated non-forced open action preserves the current farm form session", () => {
  assert.equal(shouldOpenFarmForm(true, false, false), false);
});

test("a closed farm form cannot start a new session while its previous save is pending", () => {
  assert.equal(shouldOpenFarmForm(false, true), false);
});

test("mobile force-new cannot invalidate a pending save after the form was canceled", () => {
  assert.equal(shouldOpenFarmForm(false, true, true), false);
});

test("a canceled form can be reopened after the previous save settles", () => {
  assert.equal(shouldOpenFarmForm(false, false), true);
});

test("a dismissed farm save reports successful completion", () => {
  assert.equal(getDismissedFarmSaveOutcome(true, true, true), "success");
});

test("a dismissed farm save reports failed completion", () => {
  assert.equal(getDismissedFarmSaveOutcome(true, false, true), "error");
  assert.equal(getDismissedFarmSaveOutcome(true, false, false), "error");
});

test("a dismissed farm save does not queue success feedback while the dashboard is hidden", () => {
  assert.equal(getDismissedFarmSaveOutcome(true, true, false), null);
});

test("an active farm form does not receive a dismissed-save outcome", () => {
  assert.equal(getDismissedFarmSaveOutcome(false, true, true), null);
  assert.equal(getDismissedFarmSaveOutcome(false, false, true), null);
});
