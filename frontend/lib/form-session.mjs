export function isCurrentFormSession(submittedSessionId, currentSessionId) {
  return Number.isSafeInteger(submittedSessionId)
    && Number.isSafeInteger(currentSessionId)
    && submittedSessionId === currentSessionId;
}

export function shouldShowFarmSaveNotice(submittedSessionId, currentSessionId, dashboardVisible) {
  return isCurrentFormSession(submittedSessionId, currentSessionId)
    && dashboardVisible === true;
}

export function shouldOpenFarmForm(isFormOpen, isSaving, forceNew = false) {
  return isSaving !== true && (isFormOpen !== true || forceNew === true);
}

export function shouldShowFarmSavePending(isFormOpen, isSaving) {
  return isFormOpen !== true && isSaving === true;
}

export function getDismissedFarmSaveOutcome(formWasDismissed, succeeded, dashboardVisible) {
  if (formWasDismissed !== true) return null;
  if (succeeded === true) return dashboardVisible === true ? "success" : null;
  return "error";
}
