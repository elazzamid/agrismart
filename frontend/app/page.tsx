"use client";

import { FormEvent, useCallback, useEffect, useState } from "react";
import {
  clearToken,
  createFarm,
  currentUser,
  Farm,
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

function Dashboard({ user, onSignOut }: { user: User; onSignOut: () => void }) {
  const [farms, setFarms] = useState<Farm[]>([]);
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
        <div className="dashboard-content">
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
          {loading ? <div className="loading-card"><span className="spinner" /> Memuat daftar kebun…</div> : farms.length === 0 ? <div className="empty-state"><div className="empty-icon">⌂</div><h3>Belum ada kebun tercatat</h3><p>Tambahkan kebun pertama Anda untuk mulai menyusun catatan budidaya.</p><button className="button primary" onClick={() => setShowForm(true)}>Tambahkan kebun pertama ↗</button></div> : <div className="farm-grid">{farms.map((farm, index) => <article className="farm-card" key={farm.id}><div className={`farm-card-art art-${index % 3}`}><span className="farm-badge">KEBUN {String(index + 1).padStart(2, "0")}</span><span className="farm-symbol">⌁</span></div><div className="farm-card-body"><div className="farm-card-title"><h3>{farm.name}</h3><span className="arrow-mark">↗</span></div><p className="farm-location"><span>⌖</span> {farm.location_name || "Lokasi belum ditambahkan"}</p><div className="farm-card-foot"><span><span className="status-dot" /> Terdaftar</span><span className="farm-id">ID · {farm.id.slice(0, 8)}</span></div></div></article>)}</div>}
          <section className="next-step"><div className="next-step-icon">✳</div><div><strong>Langkah selanjutnya: petak dan siklus tanam</strong><p>Setelah kebun dibuat, tahap berikutnya adalah mencatat petak, memilih tanaman, dan tanggal tanam.</p></div><span className="coming-soon">BERIKUTNYA</span></section>
          <footer className="dashboard-footer"><span>AgriSmart <span>·</span> Catatan tumbuh, keputusan lebih baik.</span><span>Data kebun hanya ditampilkan untuk akun Anda.</span></footer>
        </div>
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
