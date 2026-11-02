"use client";

import { useCallback, useEffect, useRef, useState } from "react";
import { listCropGrowthStages, type CropGrowthStage } from "@/lib/api";
import { findGrowthStage, getCropAgeDays } from "@/lib/crop-progress.mjs";

type Props = { cropId: string; plantingDate: string; status: string };

export default function CropCycleProgress({ cropId, plantingDate, status }: Props) {
  const ageDays = getCropAgeDays(plantingDate);
  const [stages, setStages] = useState<CropGrowthStage[]>([]);
  const [loading, setLoading] = useState(false);
  const [failed, setFailed] = useState(false);
  const mountedRef = useRef(false);
  const requestId = useRef(0);

  const load = useCallback(async () => {
    const currentRequest = ++requestId.current;
    setLoading(true);
    setFailed(false);
    try {
      const result = await listCropGrowthStages(cropId);
      if (!mountedRef.current || requestId.current !== currentRequest) return;
      setStages(result);
    } catch {
      if (!mountedRef.current || requestId.current !== currentRequest) return;
      setFailed(true);
    } finally {
      if (mountedRef.current && requestId.current === currentRequest) setLoading(false);
    }
  }, [cropId]);

  useEffect(() => {
    mountedRef.current = true;
    if (status === "active" && ageDays !== null) void load();
    return () => {
      mountedRef.current = false;
      requestId.current += 1;
    };
  }, [ageDays, load, status]);

  if (status === "planned") return <p className="inline-hint">Siklus direncanakan; umur tanaman belum dimulai.</p>;
  if (status === "harvested") return <p className="inline-hint">Siklus selesai dipanen.</p>;
  if (status === "cancelled") return <p className="inline-hint">Siklus dibatalkan.</p>;
  if (status !== "active") return <p className="inline-hint">Status siklus belum dikenali; tahap tidak ditampilkan.</p>;
  if (ageDays === null) return <p className="inline-hint">Umur belum dihitung karena tanggal tanam berada di masa depan atau tidak valid.</p>;

  const stage = findGrowthStage(stages, ageDays);
  return (
    <div className="inline-hint" aria-live="polite">
      <p><strong>Umur tanaman: {ageDays} hari</strong></p>
      {loading ? <p>Memuat label tahap dari katalog…</p> : failed ? (
        <p role="alert">Katalog tahap belum dapat dimuat. <button type="button" className="text-button" onClick={() => void load()}>Coba lagi</button></p>
      ) : stage ? <p>Label tahap katalog contoh: <strong>{stage.name}</strong></p> : <p>Tahap belum dapat dipetakan dari rentang katalog contoh.</p>}
      <p>Label tahap memakai rentang metadata contoh; bukan rekomendasi budidaya, penilaian kondisi tanaman, atau prediksi panen.</p>
    </div>
  );
}
