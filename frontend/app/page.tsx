"use client";

import { FormEvent, useCallback, useEffect, useRef, useState } from "react";
import CycleRecords from "@/components/CycleRecords";
import CropCycleProgress from "@/components/CropCycleProgress";
import {
  clearToken,
  createFarm,
  currentUser,
  Farm,
  Crop,
  CropVariety,
  CropCycle,
  Plot,
  createCropCycle,
  createPlot,
  listCropCycles,
  listCrops,
  listCropVarieties,
  listPlots,
  getToken,
  listFarms,
  login,
  register,
  saveToken,
  User,
} from "@/lib/api";
import { getDismissedFarmSaveOutcome, isCurrentFormSession, shouldOpenFarmForm, shouldShowFarmSaveNotice, shouldShowFarmSavePending } from "@/lib/form-session.mjs";

type AuthMode = "login" | "register";

function Brand() {
  return (
    <div className="brand">
      <div className="brand-mark" aria-hidden="true">a<span>+</span></div>
      <span>agri<span className="brand-accent">smart</span></span>
    </div>
  );
}

function Field({ label, ...props }: React.InputHTMLAttributes<HTMLInputElement> & { label: string }) {
  return (
    <label className="field">
      <span>{label}</span>
      <input {...props} />
    </label>
  );
}

function AuthScreen({ onAuthenticated }: { onAuthenticated: (user: User, token: string) => void }) {
  const [mode, setMode] = useState<AuthMode>("login");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [fullName, setFullName] = useState("");
  const [phone, setPhone] = useState("");
  const [busy, setBusy] = useState(false);
  const [error, setError] = useState("");
  const mountedRef = useRef(true);

  useEffect(() => {
    mountedRef.current = true;
    return () => { mountedRef.current = false; };
  }, []);

  async function submit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setBusy(true);
    setError("");
    try {
      if (mode === "register") {
        await register({ email: email.trim(), password, full_name: fullName.trim(), phone: phone.trim() });
        if (!mountedRef.current) return;
      }
      const response = await login(email.trim(), password);
      if (!mountedRef.current) return;
      saveToken(response.access_token);
      onAuthenticated(response.user, response.access_token);
    } catch (err) {
      if (!mountedRef.current) return;
      setError(err instanceof Error ? err.message : "Terjadi kesalahan. Coba lagi.");
    } finally {
      if (mountedRef.current) setBusy(false);
    }
  }

  return (
    <main className="auth-layout">
      <section className="auth-story">
        <Brand />
        <div className="story-content">
          <span className="eyebrow"><span className="status-dot" /> RUANG KERJA PETANI</span>
          <h1>Rawat kebun dengan lebih <em>terarah.</em></h1>
          <p>Satukan catatan lahan, kegiatan budidaya, biaya, dan hasil panen. Sedikit demi sedikit, keputusan di kebun jadi lebih mudah.</p>
          <div className="story-tags"><span>01 · Catat</span><span>02 · Pantau</span><span>03 · Evaluasi</span></div>
        </div>
        <p className="story-foot">Dibuat untuk perjalanan bertani di Indonesia.</p>
        <div className="decor-leaf leaf-one" /><div className="decor-leaf leaf-two" />
      </section>
      <section className="auth-panel">
        <div className="auth-card">
          <div className="mobile-brand"><Brand /></div>
          <span className="eyebrow muted">{mode === "login" ? "SELAMAT DATANG KEMBALI" : "MULAI PERJALANAN ANDA"}</span>
          <h2>{mode === "login" ? "Masuk ke AgriSmart" : "Buat akun petani"}</h2>
          <p className="subcopy">{mode === "login" ? "Lanjutkan mencatat perkembangan kebun Anda." : "Daftar untuk mulai mengelola catatan pertanian."}</p>
          <form onSubmit={submit} className="form-stack">
            {mode === "register" && <>
              <Field label="Nama lengkap" autoComplete="name" value={fullName} onChange={e => setFullName(e.target.value)} required maxLength={120} />
              <Field label="Nomor telepon (opsional)" autoComplete="tel" inputMode="tel" value={phone} onChange={e => setPhone(e.target.value)} />
            </>}
            <Field label="Email" type="email" autoComplete="email" value={email} onChange={e => setEmail(e.target.value)} required />
            <Field label="Kata sandi" type="password" autoComplete={mode === "login" ? "current-password" : "new-password"} value={password} onChange={e => setPassword(e.target.value)} required minLength={8} />
            {error && <p className="form-error" role="alert">{error}</p>}
            <button className="button primary full" disabled={busy}>{busy ? "Memproses…" : mode === "login" ? "Masuk" : "Buat akun"} <span aria-hidden="true">↗</span></button>
          </form>
          <p className="switch-copy">
            {mode === "login" ? "Belum punya akun?" : "Sudah punya akun?"}{" "}
            <button className="text-button" onClick={() => { setMode(mode === "login" ? "register" : "login"); setError(""); }}>
              {mode === "login" ? "Daftar sekarang" : "Masuk"}
            </button>
          </p>
          <p className="privacy-note">Akun Anda digunakan untuk menjaga catatan kebun tetap pribadi.</p>
        </div>
      </section>
    </main>
  );
}


function FarmWorkspace({ farm, onBack }: { farm: Farm; onBack: () => void }) {
  const [plots, setPlots] = useState<Plot[]>([]);
  const [crops, setCrops] = useState<Crop[]>([]);
  const [cycles, setCycles] = useState<Record<string, CropCycle[]>>({});
  const cycleRequestIds = useRef<Record<string, number>>({});
  const plotIdsRef = useRef<string[]>([]);
  const varietyRequestId = useRef(0);
  const invalidateCycleRequests = useCallback(() => {
    for (const plotId of Object.keys(cycleRequestIds.current)) {
      cycleRequestIds.current[plotId] = (cycleRequestIds.current[plotId] ?? 0) + 1;
    }
  }, []);
  const [cycleLoadErrors, setCycleLoadErrors] = useState<string[]>([]);
  const [cycleLoadPending, setCycleLoadPending] = useState<string[]>([]);
  const [cycleLoadedIds, setCycleLoadedIds] = useState<string[]>([]);
  const [cropsLoadFailed, setCropsLoadFailed] = useState(false);
  const [plotsLoaded, setPlotsLoaded] = useState(false);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [showPlotForm, setShowPlotForm] = useState(false);
  const [plotName, setPlotName] = useState("");
  const [plotArea, setPlotArea] = useState("");
  const [savingPlot, setSavingPlot] = useState(false);
  const plotFormSessionId = useRef(0);
  const [cycleFormPlot, setCycleFormPlot] = useState<string | null>(null);
  const cycleFormSessionId = useRef(0);
  const [recordCycleId, setRecordCycleId] = useState<string | null>(null);
  const [cropId, setCropId] = useState("");
  const [varietyId, setVarietyId] = useState("");
  const [varieties, setVarieties] = useState<CropVariety[]>([]);
  const [varietiesLoading, setVarietiesLoading] = useState(false);
  const [varietiesLoadFailed, setVarietiesLoadFailed] = useState(false);
  const [plantingDate, setPlantingDate] = useState("");
  const [savingCycle, setSavingCycle] = useState(false);
  const [notice, setNotice] = useState("");

  const loadCycles = useCallback(async (plotId: string) => {
    const requestId = (cycleRequestIds.current[plotId] ?? 0) + 1;
    cycleRequestIds.current[plotId] = requestId;
    setCycleLoadPending(current => current.includes(plotId) ? current : [...current, plotId]);
    setCycleLoadErrors(current => current.filter(id => id !== plotId));
    try {
      const nextCycles = await listCropCycles(plotId);
      if (cycleRequestIds.current[plotId] !== requestId) return;
      setCycles(previous => ({ ...previous, [plotId]: nextCycles }));
      setCycleLoadedIds(current => current.includes(plotId) ? current : [...current, plotId]);
    } catch {
      if (cycleRequestIds.current[plotId] !== requestId) return;
      setCycleLoadErrors(current => current.includes(plotId) ? current : [...current, plotId]);
    } finally {
      if (cycleRequestIds.current[plotId] === requestId) {
        setCycleLoadPending(current => current.filter(id => id !== plotId));
      }
    }
  }, []);

  const mountedRef = useRef(false);
  const refreshRequestId = useRef(0);
  const refresh = useCallback(async () => {
    const requestId = ++refreshRequestId.current;
    invalidateCycleRequests();
    const cycleRequestIdsAtStart = { ...cycleRequestIds.current };
    setCycleLoadPending(plotIdsRef.current);
    setCycleLoadErrors([]);
    setLoading(true);
    setError("");
    try {
      const plotsRequest = listPlots(farm.id).then(nextPlots => {
        if (refreshRequestId.current !== requestId) return;
        plotIdsRef.current = nextPlots.map(plot => plot.id);
        setPlots(nextPlots);
        setPlotsLoaded(true);
        const nextPlotIds = new Set(nextPlots.map(plot => plot.id));
        const plotsToLoad = nextPlots.filter(plot =>
          (cycleRequestIds.current[plot.id] ?? 0) === (cycleRequestIdsAtStart[plot.id] ?? 0)
        );
        const plotsToLoadIds = new Set(plotsToLoad.map(plot => plot.id));
        setCycleLoadErrors(current => current.filter(plotId => nextPlotIds.has(plotId)));
        setCycleLoadedIds(current => current.filter(plotId => nextPlotIds.has(plotId)));
        setCycleLoadPending(current => {
          const newerPending = current.filter(plotId => nextPlotIds.has(plotId) && !plotsToLoadIds.has(plotId));
          return [...new Set([...newerPending, ...plotsToLoadIds])];
        });
        void Promise.all(plotsToLoad.map(plot => loadCycles(plot.id)));
      }).catch(err => {
        if (refreshRequestId.current !== requestId) return;
        const plotsToLoad = plotIdsRef.current.filter(plotId =>
          (cycleRequestIds.current[plotId] ?? 0) === (cycleRequestIdsAtStart[plotId] ?? 0)
        );
        void Promise.all(plotsToLoad.map(plotId => loadCycles(plotId)));
        setError(err instanceof Error ? err.message : "Daftar petak belum dapat dimuat.");
      });
      const cropsRequest = listCrops().then(nextCrops => {
        if (refreshRequestId.current !== requestId) return;
        setCrops(nextCrops);
        setCropsLoadFailed(false);
      }).catch(() => {
        if (refreshRequestId.current !== requestId) return;
        setCropsLoadFailed(true);
      });
      await Promise.allSettled([plotsRequest, cropsRequest]);
    } catch (err) {
      if (refreshRequestId.current !== requestId) return;
      setError(err instanceof Error ? err.message : "Data kebun belum dapat dimuat.");
    } finally {
      if (refreshRequestId.current === requestId) setLoading(false);
    }
  }, [farm.id, invalidateCycleRequests, loadCycles]);

  const retryCycles = useCallback(async (plotId: string) => {
    await loadCycles(plotId);
  }, [loadCycles]);

  const loadVarieties = useCallback(async (selectedCropId: string) => {
    const requestId = ++varietyRequestId.current;
    setVarieties([]);
    setVarietiesLoadFailed(false);
    if (!selectedCropId) {
      setVarietiesLoading(false);
      return;
    }
    setVarietiesLoading(true);
    try {
      const nextVarieties = await listCropVarieties(selectedCropId);
      if (!mountedRef.current || varietyRequestId.current !== requestId) return;
      setVarieties(nextVarieties);
    } catch {
      if (!mountedRef.current || varietyRequestId.current !== requestId) return;
      setVarietiesLoadFailed(true);
    } finally {
      if (mountedRef.current && varietyRequestId.current === requestId) setVarietiesLoading(false);
    }
  }, []);

  useEffect(() => { mountedRef.current = true; void refresh(); return () => { mountedRef.current = false; refreshRequestId.current += 1; varietyRequestId.current += 1; invalidateCycleRequests(); }; }, [refresh, invalidateCycleRequests]);

  async function submitPlot(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    const submittedFormSessionId = plotFormSessionId.current;
    setSavingPlot(true);
    setError("");
    setNotice("");
    try {
      await createPlot(farm.id, { name: plotName.trim(), area_m2: Number(plotArea) });
      if (!mountedRef.current) return;
      if (isCurrentFormSession(submittedFormSessionId, plotFormSessionId.current)) {
        plotFormSessionId.current += 1;
        setPlotName("");
        setPlotArea("");
        setShowPlotForm(false);
        setNotice("Petak berhasil ditambahkan.");
      }
      await refresh();
    } catch (err) {
      if (!mountedRef.current) return;
      if (isCurrentFormSession(submittedFormSessionId, plotFormSessionId.current)) {
        setError(err instanceof Error ? err.message : "Petak belum berhasil disimpan.");
      }
    } finally {
      if (mountedRef.current) setSavingPlot(false);
    }
  }

  async function submitCycle(event: FormEvent<HTMLFormElement>, plotId: string) {
    event.preventDefault();
    if (cropsLoadFailed) {
      setError("Katalog tanaman belum berhasil diperbarui. Coba muat ulang sebelum membuat siklus.");
      return;
    }
    if (varietiesLoading || varietiesLoadFailed) {
      setError("Katalog varietas belum berhasil dimuat. Coba muat ulang sebelum membuat siklus.");
      return;
    }
    const submittedFormSessionId = cycleFormSessionId.current;
    setSavingCycle(true);
    setError("");
    setNotice("");
    try {
      await createCropCycle(plotId, { crop_id: cropId, variety_id: varietyId || null, planting_date: plantingDate });
      if (!mountedRef.current) return;
      if (isCurrentFormSession(submittedFormSessionId, cycleFormSessionId.current)) {
        cycleFormSessionId.current += 1;
        setCycleFormPlot(null);
        setCropId("");
        setVarietyId("");
        setVarieties([]);
        setVarietiesLoading(false);
        setVarietiesLoadFailed(false);
        varietyRequestId.current += 1;
        setPlantingDate("");
        setNotice("Siklus tanam berhasil ditambahkan.");
      }
      await refresh();
    } catch (err) {
      if (!mountedRef.current) return;
      if (isCurrentFormSession(submittedFormSessionId, cycleFormSessionId.current)) {
        setError(err instanceof Error ? err.message : "Siklus tanam belum berhasil disimpan.");
      }
    } finally {
      if (mountedRef.current) setSavingCycle(false);
    }
  }

  return (
    <div className="dashboard-content workspace-content">
      <button className="back-link" onClick={onBack}>← Kembali ke daftar kebun</button>
      <section className="welcome-row workspace-welcome">
        <div><span className="eyebrow muted">KEBUN PILIHAN</span><h1>{farm.name}<span className="period">.</span></h1><p>{farm.location_name || "Lokasi belum ditambahkan"} · Kelola petak dan siklus tanam di kebun ini.</p></div>
        <button className="button primary" onClick={() => { plotFormSessionId.current += 1; setShowPlotForm(v => !v); }}>{showPlotForm ? "Tutup formulir" : "+ Tambah petak"}</button>
      </section>
      <div className="workspace-banner"><div className="workspace-banner-icon">⌂</div><div><strong>Setiap petak punya ceritanya sendiri.</strong><p>Tambahkan petak, pilih tanaman, lalu masukkan tanggal tanam untuk memulai pencatatan budidaya.</p></div></div>
      {notice && <p className="success-note" role="status">{notice}</p>}
      {error && <div className="error-banner" role="alert"><span>{error}</span><button className="text-button" onClick={() => void refresh()}>Coba lagi</button></div>}
      {cropsLoadFailed && <div className="error-banner" role="alert"><span>Katalog tanaman belum dapat dimuat. Siklus tanam tidak dapat dibuat sampai katalog berhasil diperbarui.</span><button className="text-button" onClick={() => void refresh()}>Coba lagi</button></div>}
      {showPlotForm && <form className="farm-form" onSubmit={submitPlot}><div className="form-title"><div><strong>Petak baru</strong><p>Gunakan nama dan luas petak sesuai catatan Anda.</p></div><span className="form-step">PETAK</span></div><div className="form-grid"><Field label="Nama petak" placeholder="Contoh: Petak Utara" value={plotName} onChange={e => { plotFormSessionId.current += 1; setPlotName(e.target.value); }} required maxLength={120} /><Field label="Luas (m²)" type="number" inputMode="decimal" min="0.01" step="any" placeholder="Contoh: 250" value={plotArea} onChange={e => { plotFormSessionId.current += 1; setPlotArea(e.target.value); }} required /></div><div className="form-actions"><button type="button" className="button ghost" onClick={() => { plotFormSessionId.current += 1; setShowPlotForm(false); }}>Batal</button><button className="button primary" disabled={savingPlot}>{savingPlot ? "Menyimpan…" : "Simpan petak ↗"}</button></div></form>}
      <section className="section-heading"><div><span className="eyebrow muted">PEMBAGIAN LAHAN</span><h2>Petak kebun <span className="count">{plots.length}</span></h2></div></section>
      {loading && !plotsLoaded ? <div className="loading-card"><span className="spinner" /> Memuat petak dan siklus tanam…</div> : !plotsLoaded ? null : plots.length === 0 ? <div className="empty-state"><div className="empty-icon">▦</div><h3>Belum ada petak</h3><p>Buat petak pertama untuk mulai mencatat tanaman dan tanggal tanam.</p><button className="button primary" onClick={() => { plotFormSessionId.current += 1; setShowPlotForm(true); }}>Tambah petak pertama ↗</button></div> : <div className="plot-list">{plots.map((plot, index) => <article className="plot-card" key={plot.id}><div className="plot-card-head"><div className="plot-index">{String(index + 1).padStart(2, "0")}</div><div className="plot-main"><h3>{plot.name}</h3><p>{plot.area_m2.toLocaleString("id-ID")} m² · Petak budidaya</p></div><button className="button outline small-button" onClick={() => { cycleFormSessionId.current += 1; setCycleFormPlot(cycleFormPlot === plot.id ? null : plot.id); setCropId(""); setVarietyId(""); void loadVarieties(""); setPlantingDate(""); }}>{cycleFormPlot === plot.id ? "Batal" : "+ Siklus tanam"}</button></div>
        {cycleFormPlot === plot.id && <form className="cycle-form" onSubmit={event => void submitCycle(event, plot.id)}><Field label="Tanggal tanam" type="date" value={plantingDate} onChange={e => { cycleFormSessionId.current += 1; setPlantingDate(e.target.value); }} required /><label className="field"><span>Tanaman</span><select value={cropId} onChange={e => { cycleFormSessionId.current += 1; setCropId(e.target.value); setVarietyId(""); void loadVarieties(e.target.value); }} required><option value="">Pilih tanaman</option>{crops.map(crop => <option key={crop.id} value={crop.id}>{crop.name}</option>)}</select></label>{cropId && (varietiesLoading ? <p className="inline-hint" role="status">Memuat pilihan varietas…</p> : varietiesLoadFailed ? <p className="inline-hint" role="alert">Varietas belum dapat dimuat. <button type="button" className="text-button" onClick={() => void loadVarieties(cropId)}>Coba lagi</button></p> : varieties.length === 0 ? <p className="inline-hint">Tidak ada varietas aktif untuk tanaman ini. Anda dapat melanjutkan tanpa varietas.</p> : <label className="field"><span>Varietas (opsional)</span><select value={varietyId} onChange={e => { cycleFormSessionId.current += 1; setVarietyId(e.target.value); }}><option value="">Tanpa varietas spesifik</option>{varieties.map(variety => <option key={variety.id} value={variety.id}>{variety.name}</option>)}</select></label>)}<button className="button primary" disabled={savingCycle || crops.length === 0 || cropsLoadFailed || varietiesLoading || varietiesLoadFailed}>{savingCycle ? "Menyimpan…" : "Simpan siklus ↗"}</button>{cropsLoadFailed ? <p className="inline-hint">Katalog tanaman belum berhasil diperbarui. Coba muat ulang sebelum membuat siklus.</p> : crops.length === 0 && <p className="inline-hint">Katalog tanaman aktif belum tersedia.</p>}</form>}
        <div className="cycle-list">{cycleLoadPending.includes(plot.id) && !cycleLoadedIds.includes(plot.id) ? <p className="no-cycle" role="status">Memuat siklus tanam…</p> : cycleLoadErrors.includes(plot.id) && !cycleLoadedIds.includes(plot.id) ? <p className="no-cycle" role="status">Siklus tanam belum berhasil dimuat. <button type="button" className="text-button" onClick={() => void retryCycles(plot.id)}>Coba lagi</button></p> : <>{cycleLoadPending.includes(plot.id) && <p className="no-cycle" role="status">Memuat pembaruan; menampilkan siklus terakhir yang berhasil dimuat.</p>}{cycleLoadErrors.includes(plot.id) && <p className="no-cycle" role="status">Pembaruan gagal; menampilkan siklus terakhir yang berhasil dimuat. <button type="button" className="text-button" onClick={() => void retryCycles(plot.id)}>Coba lagi</button></p>}{cycleLoadedIds.includes(plot.id) && (cycles[plot.id] || []).length === 0 ? <p className="no-cycle">Belum ada siklus tanam pada petak ini.</p> : (cycles[plot.id] || []).map(cycle => <div className="cycle-entry" key={cycle.id}><div className="cycle-row"><span className="cycle-plant-icon">✳</span><div className="cycle-copy"><strong>{crops.find(crop => crop.id === cycle.crop_id)?.name || "Tanaman"}</strong><span>Mulai tanam {new Date(cycle.planting_date + "T00:00:00").toLocaleDateString("id-ID", { day: "numeric", month: "long", year: "numeric" })}</span></div><span className="cycle-status">{cycle.status}</span><button type="button" className="text-button cycle-record-toggle" aria-expanded={recordCycleId === cycle.id} onClick={() => setRecordCycleId(recordCycleId === cycle.id ? null : cycle.id)}>{recordCycleId === cycle.id ? "Tutup catatan" : "Catatan"}</button></div><CropCycleProgress cropId={cycle.crop_id} plantingDate={cycle.planting_date} status={cycle.status} />{recordCycleId === cycle.id && <CycleRecords cycleId={cycle.id} />}</div>)}</>}</div>
      </article>)}</div>}
      <footer className="dashboard-footer"><span>AgriSmart <span>·</span> {farm.name}</span><span>Petak dan siklus tanam mengikuti data API akun Anda.</span></footer>
    </div>
  );
}

function Dashboard({ user, onSignOut }: { user: User; onSignOut: () => void }) {
  const [farms, setFarms] = useState<Farm[]>([]);
  const [selectedFarm, setSelectedFarm] = useState<Farm | null>(null);
  const [loading, setLoading] = useState(true);
  const [farmsLoaded, setFarmsLoaded] = useState(false);
  const [error, setError] = useState("");
  const [farmListError, setFarmListError] = useState("");
  const [showForm, setShowForm] = useState(false);
  const [farmName, setFarmName] = useState("");
  const [locationName, setLocationName] = useState("");
  const [saving, setSaving] = useState(false);
  const farmFormSessionId = useRef(0);
  const farmFormDismissedDuringSaveRef = useRef(false);
  const [notice, setNotice] = useState("");

  const mountedRef = useRef(false);
  const refreshRequestId = useRef(0);
  const refreshFarms = useCallback(async () => {
    const requestId = ++refreshRequestId.current;
    setLoading(true);
    setFarmListError("");
    try {
      const nextFarms = await listFarms();
      if (refreshRequestId.current !== requestId) return;
      setFarms(nextFarms);
      setFarmsLoaded(true);
    } catch (err) {
      if (refreshRequestId.current !== requestId) return;
      setFarmListError(err instanceof Error ? err.message : "Data kebun belum dapat dimuat.");
    } finally {
      if (refreshRequestId.current === requestId) setLoading(false);
    }
  }, []);

  useEffect(() => { mountedRef.current = true; void refreshFarms(); return () => { mountedRef.current = false; refreshRequestId.current += 1; }; }, [refreshFarms]);

  const dashboardVisibleRef = useRef(true);

  function openFarm(farm: Farm) {
    dashboardVisibleRef.current = false;
    setSelectedFarm(farm);
  }

  function openFarmForm(forceNew = false) {
    dashboardVisibleRef.current = true;
    if (!shouldOpenFarmForm(showForm, saving, forceNew)) return;
    farmFormDismissedDuringSaveRef.current = false;
    farmFormSessionId.current += 1;
    setError("");
    setNotice("");
    setShowForm(true);
  }

  function closeFarmForm() {
    farmFormDismissedDuringSaveRef.current = saving;
    farmFormSessionId.current += 1;
    setShowForm(false);
    setFarmName("");
    setLocationName("");
    setError("");
  }

  async function submitFarm(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    const submittedFormSessionId = farmFormSessionId.current;
    let postSettled = false;
    let formWasDismissed = false;
    setSaving(true);
    setError("");
    setNotice("");
    try {
      await createFarm({ name: farmName.trim(), location_name: locationName.trim() });
      postSettled = true;
      if (!mountedRef.current) return;
      setSaving(false);
      formWasDismissed = farmFormDismissedDuringSaveRef.current;
      farmFormDismissedDuringSaveRef.current = false;
      const formSessionIsCurrent = isCurrentFormSession(submittedFormSessionId, farmFormSessionId.current);
      const shouldNotifySuccess = shouldShowFarmSaveNotice(submittedFormSessionId, farmFormSessionId.current, dashboardVisibleRef.current);
      const dismissedSaveOutcome = getDismissedFarmSaveOutcome(formWasDismissed, true, dashboardVisibleRef.current);
      if (formSessionIsCurrent) {
        farmFormSessionId.current += 1;
        setFarmName("");
        setLocationName("");
        setShowForm(false);
        if (shouldNotifySuccess) setNotice("Kebun berhasil ditambahkan.");
      } else if (dismissedSaveOutcome === "success") {
        setNotice("Kebun berhasil disimpan setelah formulir ditutup.");
      }
      await refreshFarms();
    } catch (err) {
      if (!mountedRef.current) return;
      formWasDismissed = farmFormDismissedDuringSaveRef.current;
      farmFormDismissedDuringSaveRef.current = false;
      const dismissedSaveOutcome = getDismissedFarmSaveOutcome(formWasDismissed, false, dashboardVisibleRef.current);
      if (isCurrentFormSession(submittedFormSessionId, farmFormSessionId.current)) {
        setError(err instanceof Error ? err.message : "Kebun belum berhasil disimpan.");
      } else if (dismissedSaveOutcome === "error") {
        setError("Penyimpanan kebun gagal setelah formulir ditutup. Buka formulir kembali untuk mencoba lagi.");
      }
    } finally {
      if (mountedRef.current && !postSettled) setSaving(false);
    }
  }

  const firstName = user.email.split("@")[0];
  return (
    <div className="dashboard-shell">
      <aside className="sidebar">
        <Brand />
        <div className="side-label">MENU UTAMA</div>
        <div className="nav-item active"><span className="nav-icon">⌂</span> Ringkasan</div>
        <div className="nav-item muted-nav"><span className="nav-icon">▦</span> Kebun saya</div>
        <div className="nav-item muted-nav"><span className="nav-icon">↗</span> Catatan budidaya</div>
        <div className="side-bottom">
          <div className="help-card"><span className="help-icon">✳</span><strong>Langkah kecil, hasil berarti.</strong><p>Catatan yang rapi membantu Anda belajar dari setiap musim tanam.</p></div>
          <div className="profile-row"><div className="avatar">{firstName.slice(0, 1).toUpperCase()}</div><div className="profile-copy"><strong>{user.email}</strong><span>Petani AgriSmart</span></div><button className="icon-button" aria-label="Keluar dari akun" title="Keluar" onClick={onSignOut}>↪</button></div>
        </div>
      </aside>
      <main className="dashboard-main">
        <header className="topbar">
          <div className="breadcrumb">Ruang kerja <span>/</span> <strong>Ringkasan</strong></div>
          <div className="topbar-right"><span className="online-pill"><span className="status-dot" /> Data tersinkron saat dibuka</span><div className="avatar small">{firstName.slice(0, 1).toUpperCase()}</div></div>
        </header>
        {selectedFarm ? <FarmWorkspace key={selectedFarm.id} farm={selectedFarm} onBack={() => { dashboardVisibleRef.current = true; setSelectedFarm(null); }} /> : <div className="dashboard-content">
          <section className="welcome-row">
            <div><span className="eyebrow muted">RINGKASAN PERTANIAN</span><h1>Selamat datang kembali<span className="period">.</span></h1><p>Mulai dari mana Anda ingin merawat kebun hari ini?</p></div>
            <div className="date-chip"><span>◷</span> Ruang kerja pribadi</div>
          </section>
          <section className="hero-card">
            <div className="hero-copy"><span className="hero-kicker">LANGKAH PERTAMA</span><h2>Semua cerita kebun<br />dimulai dari <em>catatan.</em></h2><p>Tambahkan kebun Anda, lalu lanjutkan dengan petak dan siklus tanam agar aktivitas, biaya, serta panen tercatat rapi.</p><button className="button light" onClick={() => openFarmForm()}>+ Tambah kebun</button></div>
            <div className="hero-art" aria-hidden="true"><div className="sun" /><div className="hill hill-back" /><div className="hill hill-front" /><div className="plant plant-a"><i /><i /><i /><b /></div><div className="plant plant-b"><i /><i /><i /><b /></div><div className="plant plant-c"><i /><i /><i /><b /></div><div className="art-ground" /></div>
          </section>
          <section className="section-heading"><div><span className="eyebrow muted">DATA ANDA</span><h2>Kebun saya <span className="count">{farms.length}</span></h2></div><button className="button outline" onClick={() => showForm ? closeFarmForm() : openFarmForm()}>{showForm ? "Tutup formulir" : "+ Kebun baru"}</button></section>
          {notice && <p className="success-note" role="status">{notice}</p>}
          {farmListError && <div className="error-banner" role="alert"><span>{farmListError}</span><button className="text-button" onClick={() => void refreshFarms()}>Coba lagi</button></div>}
          {error && <div className="error-banner" role="alert"><span>{error}</span></div>}
          {shouldShowFarmSavePending(showForm, saving) && <section className="next-step" role="status"><div className="next-step-icon"><span className="spinner" /></div><div><strong>Penyimpanan kebun masih diproses</strong><p>Form baru dapat dibuka setelah proses sebelumnya selesai.</p></div></section>}
          {showForm && <form className="farm-form" onSubmit={submitFarm}><div className="form-title"><div><strong>Tambahkan kebun</strong><p>Mulai dengan nama dan lokasi yang mudah Anda kenali.</p></div><span className="form-step">01 / 01</span></div><div className="form-grid"><Field label="Nama kebun" placeholder="Contoh: Kebun Cabai Belakang" value={farmName} onChange={e => { farmFormSessionId.current += 1; setFarmName(e.target.value); }} required maxLength={120} /><Field label="Lokasi (opsional)" placeholder="Desa, kecamatan, kabupaten" value={locationName} onChange={e => { farmFormSessionId.current += 1; setLocationName(e.target.value); }} maxLength={160} /></div><div className="form-actions"><button type="button" className="button ghost" onClick={closeFarmForm}>Batal</button><button className="button primary" disabled={saving}>{saving ? "Menyimpan…" : "Simpan kebun ↗"}</button></div></form>}
          {loading && !farmsLoaded ? <div className="loading-card"><span className="spinner" /> Memuat daftar kebun…</div> : !farmsLoaded ? null : farms.length === 0 ? <div className="empty-state"><div className="empty-icon">⌂</div><h3>Belum ada kebun tercatat</h3><p>Tambahkan kebun pertama Anda untuk mulai menyusun catatan budidaya.</p><button className="button primary" onClick={() => openFarmForm()}>Tambahkan kebun pertama ↗</button></div> : <div className="farm-grid">{farms.map((farm, index) => <article className="farm-card" key={farm.id} role="button" tabIndex={0} onClick={() => openFarm(farm)} onKeyDown={event => { if (event.key === "Enter" || event.key === " ") { event.preventDefault(); openFarm(farm); } }} aria-label={`Buka kebun ${farm.name}`}><div className={`farm-card-art art-${index % 3}`}><span className="farm-badge">KEBUN {String(index + 1).padStart(2, "0")}</span><span className="farm-symbol">⌁</span></div><div className="farm-card-body"><div className="farm-card-title"><h3>{farm.name}</h3><span className="arrow-mark">↗</span></div><p className="farm-location"><span>⌖</span> {farm.location_name || "Lokasi belum ditambahkan"}</p><div className="farm-card-foot"><span><span className="status-dot" /> Terdaftar</span><span className="farm-id">ID · {farm.id.slice(0, 8)}</span></div></div></article>)}</div>}
          <section className="next-step"><div className="next-step-icon">✳</div><div><strong>Langkah selanjutnya: petak dan siklus tanam</strong><p>Setelah kebun dibuat, tahap berikutnya adalah mencatat petak, memilih tanaman, dan tanggal tanam.</p></div><span className="coming-soon">BERIKUTNYA</span></section>
          <footer className="dashboard-footer"><span>AgriSmart <span>·</span> Catatan tumbuh, keputusan lebih baik.</span><span>Data kebun hanya ditampilkan untuk akun Anda.</span></footer>
        </div>}
      </main>
      <nav className="mobile-nav" aria-label="Navigasi utama"><span className="mobile-nav-active">⌂<small>Ringkasan</small></span><button onClick={() => { setSelectedFarm(null); openFarmForm(true); window.scrollTo({ top: 0, behavior: "smooth" }); }}>＋<small>Tambah kebun</small></button><button onClick={onSignOut}>↪<small>Keluar</small></button></nav>
    </div>
  );
}

export default function HomePage() {
  const [user, setUser] = useState<User | null>(null);
  const [checking, setChecking] = useState(true);

  const signOut = useCallback(() => {
    clearToken();
    setUser(null);
  }, []);

  useEffect(() => {
    let active = true;
    const unauthorized = () => {
      if (active) setUser(null);
    };
    window.addEventListener("agrismart:unauthorized", unauthorized);
    const token = getToken();
    if (!token) {
      if (active) setChecking(false);
    } else {
      currentUser()
        .then((nextUser) => {
          if (active) setUser(nextUser);
        })
        .catch(() => {
          if (active) setUser(null);
        })
        .finally(() => {
          if (active) setChecking(false);
        });
    }
    return () => {
      active = false;
      window.removeEventListener("agrismart:unauthorized", unauthorized);
    };
  }, []);

  if (checking) return <main className="splash"><Brand /><span className="spinner" /><p>Menyiapkan ruang kerja Anda…</p></main>;
  if (!user) return <AuthScreen onAuthenticated={(nextUser) => { setUser(nextUser); setChecking(false); }} />;
  return <Dashboard user={user} onSignOut={signOut} />;
}
