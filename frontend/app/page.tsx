"use client";

import { FormEvent, useCallback, useEffect, useState } from "react";
import CycleRecords from "@/components/CycleRecords";
import {
  clearToken,
  createFarm,
  currentUser,
  Farm,
  Crop,
  CropCycle,
  Plot,
  createCropCycle,
  createPlot,
  listCropCycles,
  listCrops,
  listPlots,
  getToken,
  listFarms,
  login,
  register,
  saveToken,
  User,
} from "@/lib/api";

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

  async function submit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setBusy(true);
    setError("");
    try {
      if (mode === "register") {
        await register({ email: email.trim(), password, full_name: fullName.trim(), phone: phone.trim() });
      }
      const response = await login(email.trim(), password);
      saveToken(response.access_token);
      onAuthenticated(response.user, response.access_token);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Terjadi kesalahan. Coba lagi.");
    } finally {
      setBusy(false);
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
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [showPlotForm, setShowPlotForm] = useState(false);
  const [plotName, setPlotName] = useState("");
  const [plotArea, setPlotArea] = useState("");
  const [savingPlot, setSavingPlot] = useState(false);
  const [cycleFormPlot, setCycleFormPlot] = useState<string | null>(null);
  const [cropId, setCropId] = useState("");
  const [plantingDate, setPlantingDate] = useState("");
  const [savingCycle, setSavingCycle] = useState(false);
  const [notice, setNotice] = useState("");

  const refresh = useCallback(async () => {
    setLoading(true);
    setError("");
    try {
      const [nextPlots, nextCrops] = await Promise.all([listPlots(farm.id), listCrops()]);
      const pairs = await Promise.all(nextPlots.map(async plot => [plot.id, await listCropCycles(plot.id)] as const));
      setPlots(nextPlots);
      setCrops(nextCrops);
      setCycles(Object.fromEntries(pairs));
    } catch (err) {
      setError(err instanceof Error ? err.message : "Data kebun belum dapat dimuat.");
    } finally {
      setLoading(false);
    }
  }, [farm.id]);

  useEffect(() => { void refresh(); }, [refresh]);

  async function submitPlot(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setSavingPlot(true);
    setError("");
    setNotice("");
    try {
      await createPlot(farm.id, { name: plotName.trim(), area_m2: Number(plotArea) });
      setPlotName("");
      setPlotArea("");
      setShowPlotForm(false);
      setNotice("Petak berhasil ditambahkan.");
      await refresh();
    } catch (err) {
      setError(err instanceof Error ? err.message : "Petak belum berhasil disimpan.");
    } finally {
      setSavingPlot(false);
    }
  }

  async function submitCycle(event: FormEvent<HTMLFormElement>, plotId: string) {
    event.preventDefault();
    setSavingCycle(true);
    setError("");
    setNotice("");
    try {
      await createCropCycle(plotId, { crop_id: cropId, planting_date: plantingDate });
      setCycleFormPlot(null);
      setCropId("");
      setPlantingDate("");
      setNotice("Siklus tanam berhasil ditambahkan.");
      await refresh();
    } catch (err) {
      setError(err instanceof Error ? err.message : "Siklus tanam belum berhasil disimpan.");
    } finally {
      setSavingCycle(false);
    }
  }

  return (
    <div className="dashboard-content workspace-content">
      <button className="back-link" onClick={onBack}>← Kembali ke daftar kebun</button>
      <section className="welcome-row workspace-welcome">
        <div><span className="eyebrow muted">KEBUN PILIHAN</span><h1>{farm.name}<span className="period">.</span></h1><p>{farm.location_name || "Lokasi belum ditambahkan"} · Kelola petak dan siklus tanam di kebun ini.</p></div>
        <button className="button primary" onClick={() => setShowPlotForm(v => !v)}>{showPlotForm ? "Tutup formulir" : "+ Tambah petak"}</button>
      </section>
      <div className="workspace-banner"><div className="workspace-banner-icon">⌂</div><div><strong>Setiap petak punya ceritanya sendiri.</strong><p>Tambahkan petak, pilih tanaman, lalu masukkan tanggal tanam untuk memulai pencatatan budidaya.</p></div></div>
      {notice && <p className="success-note" role="status">{notice}</p>}
      {error && <div className="error-banner" role="alert"><span>{error}</span><button className="text-button" onClick={() => void refresh()}>Coba lagi</button></div>}
      {showPlotForm && <form className="farm-form" onSubmit={submitPlot}><div className="form-title"><div><strong>Petak baru</strong><p>Gunakan nama dan luas petak sesuai catatan Anda.</p></div><span className="form-step">PETAK</span></div><div className="form-grid"><Field label="Nama petak" placeholder="Contoh: Petak Utara" value={plotName} onChange={e => setPlotName(e.target.value)} required maxLength={120} /><Field label="Luas (m²)" type="number" inputMode="decimal" min="0.01" step="any" placeholder="Contoh: 250" value={plotArea} onChange={e => setPlotArea(e.target.value)} required /></div><div className="form-actions"><button type="button" className="button ghost" onClick={() => setShowPlotForm(false)}>Batal</button><button className="button primary" disabled={savingPlot}>{savingPlot ? "Menyimpan…" : "Simpan petak ↗"}</button></div></form>}
      <section className="section-heading"><div><span className="eyebrow muted">PEMBAGIAN LAHAN</span><h2>Petak kebun <span className="count">{plots.length}</span></h2></div></section>
      {loading ? <div className="loading-card"><span className="spinner" /> Memuat petak dan siklus tanam…</div> : plots.length === 0 ? <div className="empty-state"><div className="empty-icon">▦</div><h3>Belum ada petak</h3><p>Buat petak pertama untuk mulai mencatat tanaman dan tanggal tanam.</p><button className="button primary" onClick={() => setShowPlotForm(true)}>Tambah petak pertama ↗</button></div> : <div className="plot-list">{plots.map((plot, index) => <article className="plot-card" key={plot.id}><div className="plot-card-head"><div className="plot-index">{String(index + 1).padStart(2, "0")}</div><div className="plot-main"><h3>{plot.name}</h3><p>{plot.area_m2.toLocaleString("id-ID")} m² · Petak budidaya</p></div><button className="button outline small-button" onClick={() => { setCycleFormPlot(cycleFormPlot === plot.id ? null : plot.id); setCropId(""); setPlantingDate(""); }}>{cycleFormPlot === plot.id ? "Batal" : "+ Siklus tanam"}</button></div>
        {cycleFormPlot === plot.id && <form className="cycle-form" onSubmit={event => void submitCycle(event, plot.id)}><Field label="Tanggal tanam" type="date" value={plantingDate} onChange={e => setPlantingDate(e.target.value)} required /><label className="field"><span>Tanaman</span><select value={cropId} onChange={e => setCropId(e.target.value)} required><option value="">Pilih tanaman</option>{crops.map(crop => <option key={crop.id} value={crop.id}>{crop.name}</option>)}</select></label><button className="button primary" disabled={savingCycle || crops.length === 0}>{savingCycle ? "Menyimpan…" : "Simpan siklus ↗"}</button>{crops.length === 0 && <p className="inline-hint">Katalog tanaman aktif belum tersedia.</p>}</form>}
        <div className="cycle-list">{(cycles[plot.id] || []).length === 0 ? <p className="no-cycle">Belum ada siklus tanam pada petak ini.</p> : (cycles[plot.id] || []).map(cycle => <div className="cycle-entry" key={cycle.id}><div className="cycle-row"><span className="cycle-plant-icon">✳</span><div className="cycle-copy"><strong>{crops.find(crop => crop.id === cycle.crop_id)?.name || "Tanaman"}</strong><span>Mulai tanam {new Date(cycle.planting_date + "T00:00:00").toLocaleDateString("id-ID", { day: "numeric", month: "long", year: "numeric" })}</span></div><span className="cycle-status">{cycle.status}</span></div><CycleRecords cycleId={cycle.id} /></div>)}</div>
      </article>)}</div>}
      <footer className="dashboard-footer"><span>AgriSmart <span>·</span> {farm.name}</span><span>Petak dan siklus tanam mengikuti data API akun Anda.</span></footer>
    </div>
  );
}

function Dashboard({ user, onSignOut }: { user: User; onSignOut: () => void }) {
  const [farms, setFarms] = useState<Farm[]>([]);
  const [selectedFarm, setSelectedFarm] = useState<Farm | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [showForm, setShowForm] = useState(false);
  const [farmName, setFarmName] = useState("");
  const [locationName, setLocationName] = useState("");
  const [saving, setSaving] = useState(false);
  const [notice, setNotice] = useState("");

  const refreshFarms = useCallback(async () => {
    setLoading(true);
    setError("");
    try {
      setFarms(await listFarms());
    } catch (err) {
      setError(err instanceof Error ? err.message : "Data kebun belum dapat dimuat.");
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => { void refreshFarms(); }, [refreshFarms]);

  async function submitFarm(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setSaving(true);
    setError("");
    setNotice("");
    try {
      await createFarm({ name: farmName.trim(), location_name: locationName.trim() });
      setFarmName("");
      setLocationName("");
      setShowForm(false);
      setNotice("Kebun berhasil ditambahkan.");
      await refreshFarms();
    } catch (err) {
      setError(err instanceof Error ? err.message : "Kebun belum berhasil disimpan.");
    } finally {
      setSaving(false);
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
        {selectedFarm ? <FarmWorkspace farm={selectedFarm} onBack={() => setSelectedFarm(null)} /> : <div className="dashboard-content">
          <section className="welcome-row">
            <div><span className="eyebrow muted">RINGKASAN PERTANIAN</span><h1>Selamat datang kembali<span className="period">.</span></h1><p>Mulai dari mana Anda ingin merawat kebun hari ini?</p></div>
            <div className="date-chip"><span>◷</span> Ruang kerja pribadi</div>
          </section>
          <section className="hero-card">
            <div className="hero-copy"><span className="hero-kicker">LANGKAH PERTAMA</span><h2>Semua cerita kebun<br />dimulai dari <em>catatan.</em></h2><p>Tambahkan kebun Anda, lalu lanjutkan dengan petak dan siklus tanam agar aktivitas, biaya, serta panen tercatat rapi.</p><button className="button light" onClick={() => setShowForm(true)}>+ Tambah kebun</button></div>
            <div className="hero-art" aria-hidden="true"><div className="sun" /><div className="hill hill-back" /><div className="hill hill-front" /><div className="plant plant-a"><i /><i /><i /><b /></div><div className="plant plant-b"><i /><i /><i /><b /></div><div className="plant plant-c"><i /><i /><i /><b /></div><div className="art-ground" /></div>
          </section>
          <section className="section-heading"><div><span className="eyebrow muted">DATA ANDA</span><h2>Kebun saya <span className="count">{farms.length}</span></h2></div><button className="button outline" onClick={() => setShowForm(!showForm)}>{showForm ? "Tutup formulir" : "+ Kebun baru"}</button></section>
          {notice && <p className="success-note" role="status">{notice}</p>}
          {error && <div className="error-banner" role="alert"><span>{error}</span><button className="text-button" onClick={() => void refreshFarms()}>Coba lagi</button></div>}
          {showForm && <form className="farm-form" onSubmit={submitFarm}><div className="form-title"><div><strong>Tambahkan kebun</strong><p>Mulai dengan nama dan lokasi yang mudah Anda kenali.</p></div><span className="form-step">01 / 01</span></div><div className="form-grid"><Field label="Nama kebun" placeholder="Contoh: Kebun Cabai Belakang" value={farmName} onChange={e => setFarmName(e.target.value)} required maxLength={120} /><Field label="Lokasi (opsional)" placeholder="Desa, kecamatan, kabupaten" value={locationName} onChange={e => setLocationName(e.target.value)} maxLength={160} /></div><div className="form-actions"><button type="button" className="button ghost" onClick={() => setShowForm(false)}>Batal</button><button className="button primary" disabled={saving}>{saving ? "Menyimpan…" : "Simpan kebun ↗"}</button></div></form>}
          {loading ? <div className="loading-card"><span className="spinner" /> Memuat daftar kebun…</div> : farms.length === 0 ? <div className="empty-state"><div className="empty-icon">⌂</div><h3>Belum ada kebun tercatat</h3><p>Tambahkan kebun pertama Anda untuk mulai menyusun catatan budidaya.</p><button className="button primary" onClick={() => setShowForm(true)}>Tambahkan kebun pertama ↗</button></div> : <div className="farm-grid">{farms.map((farm, index) => <article className="farm-card" key={farm.id} role="button" tabIndex={0} onClick={() => setSelectedFarm(farm)} onKeyDown={event => { if (event.key === "Enter" || event.key === " ") { event.preventDefault(); setSelectedFarm(farm); } }} aria-label={`Buka kebun ${farm.name}`}><div className={`farm-card-art art-${index % 3}`}><span className="farm-badge">KEBUN {String(index + 1).padStart(2, "0")}</span><span className="farm-symbol">⌁</span></div><div className="farm-card-body"><div className="farm-card-title"><h3>{farm.name}</h3><span className="arrow-mark">↗</span></div><p className="farm-location"><span>⌖</span> {farm.location_name || "Lokasi belum ditambahkan"}</p><div className="farm-card-foot"><span><span className="status-dot" /> Terdaftar</span><span className="farm-id">ID · {farm.id.slice(0, 8)}</span></div></div></article>)}</div>}
          <section className="next-step"><div className="next-step-icon">✳</div><div><strong>Langkah selanjutnya: petak dan siklus tanam</strong><p>Setelah kebun dibuat, tahap berikutnya adalah mencatat petak, memilih tanaman, dan tanggal tanam.</p></div><span className="coming-soon">BERIKUTNYA</span></section>
          <footer className="dashboard-footer"><span>AgriSmart <span>·</span> Catatan tumbuh, keputusan lebih baik.</span><span>Data kebun hanya ditampilkan untuk akun Anda.</span></footer>
        </div>}
      </main>
      <nav className="mobile-nav" aria-label="Navigasi utama"><span className="mobile-nav-active">⌂<small>Ringkasan</small></span><button onClick={() => { setShowForm(true); window.scrollTo({ top: 0, behavior: "smooth" }); }}>＋<small>Tambah kebun</small></button><button onClick={onSignOut}>↪<small>Keluar</small></button></nav>
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
    const unauthorized = () => setUser(null);
    window.addEventListener("agrismart:unauthorized", unauthorized);
    const token = getToken();
    if (!token) {
      setChecking(false);
    } else {
      currentUser().then(setUser).catch(() => setUser(null)).finally(() => setChecking(false));
    }
    return () => window.removeEventListener("agrismart:unauthorized", unauthorized);
  }, []);

  if (checking) return <main className="splash"><Brand /><span className="spinner" /><p>Menyiapkan ruang kerja Anda…</p></main>;
  if (!user) return <AuthScreen onAuthenticated={(nextUser) => { setUser(nextUser); setChecking(false); }} />;
  return <Dashboard user={user} onSignOut={signOut} />;
}
