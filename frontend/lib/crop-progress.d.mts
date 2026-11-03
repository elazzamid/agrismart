export type CropGrowthStageRange = {
  name: string;
  sequence_no: number;
  min_days?: number | null;
  max_days?: number | null;
};

export function getCropAgeDays(plantingDate: string, today?: Date): number | null;
export function shouldLoadGrowthStageCatalog(status: string, ageDays: number | null): boolean;
export function findGrowthStage<T extends CropGrowthStageRange>(stages: T[], ageDays: number): T | null;
