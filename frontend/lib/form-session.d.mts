export function isCurrentFormSession(submittedSessionId: number, currentSessionId: number): boolean;
export function shouldShowFarmSaveNotice(submittedSessionId: number, currentSessionId: number, dashboardVisible: boolean): boolean;
export function shouldOpenFarmForm(isFormOpen: boolean, isSaving: boolean, forceNew?: boolean): boolean;
export function shouldResetFarmFormDraft(forceNew: boolean, isSaving: boolean): boolean;
export function shouldShowFarmSavePending(isFormOpen: boolean, isSaving: boolean): boolean;
export function getDismissedFarmSaveOutcome(formWasDismissed: boolean, succeeded: boolean, dashboardVisible: boolean): "success" | "error" | null;
