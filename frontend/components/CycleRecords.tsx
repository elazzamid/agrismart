"use client";

import { FormEvent, useCallback, useEffect, useRef, useState } from "react";
import {
  createActivity, createExpense, createHarvest, createInput,
  FarmActivity, FarmExpense, FarmHarvest, FarmInput,
  listActivities, listExpenses, listHarvests, listInputs,
} from "@/lib/api";
import { isCurrentFormSession } from "@/lib/form-session.mjs";

type Kind = "activity" | "input" | "expense" | "harvest";
type RecordRow = FarmActivity | FarmInput | FarmExpense | FarmHarvest;
const today = () => {
  const now = new Date();
  const month = String(now.getMonth() + 1).padStart(2, "0");
  const day = String(now.getDate()).padStart(2, "0");
  return `${now.getFullYear()}-${month}-${day}`;
};
const labels: Record<Kind, string> = { activity: "Kegiatan", input: "Sarana produksi", expense: "Biaya", harvest: "Panen" };
const activityNames: Record<string, string> = { watering: "Penyiraman", weeding: "Penyiangan", pruning: "Pemangkasan", planting: "Penanaman", inspection: "Pemeriksaan", other: "Lainnya" };
const inputNames: Record<string, string> = { fertilizer: "Pupuk", pesticide: "Pestisida", seed: "Benih", other: "Lainnya" };
const expenseNames: Record<string, string> = { seed: "Benih", fertilizer: "Pupuk", pesticide: "Pestisida", labor: "Tenaga kerja", irrigation: "Irigasi", equipment: "Peralatan", transport: "Transportasi", other: "Lainnya" };


function formatDate(value: string) {
  return new Date(value + "T00:00:00").toLocaleDateString("id-ID", { day: "numeric", month: "short", year: "numeric" });
}
function money(value: number) {
  return new Intl.NumberFormat("id-ID", { style: "currency", currency: "IDR", maximumFractionDigits: 0 }).format(value);
}

export default function CycleRecords({ cycleId }: { cycleId: string }) {
  const [kind, setKind] = useState<Kind>("activity");
  const [rows, setRows] = useState<Record<Kind, RecordRow[]>>({ activity: [], input: [], expense: [], harvest: [] });
  const [loading, setLoading] = useState(true);
  const [loadedSuccessfully, setLoadedSuccessfully] = useState(false);
  const [loadError, setLoadError] = useState("");
  const mountedRef = useRef(false);
  const refreshRequestId = useRef(0);
  const formSessionId = useRef(0);
  const [saving, setSaving] = useState(false);
  const [showForm, setShowForm] = useState(false);
  const [error, setError] = useState("");
  const [notice, setNotice] = useState("");
  const [date, setDate] = useState("");
  const [type, setType] = useState("watering");
  const [name, setName] = useState("");
  const [quantity, setQuantity] = useState("");
  const [unit, setUnit] = useState("kg");
  const [amount, setAmount] = useState("");
  const [grade, setGrade] = useState("");
  const [notes, setNotes] = useState("");

  const refresh = useCallback(async () => {
    const requestId = ++refreshRequestId.current;
    setLoading(true);
    setError("");
    setLoadError("");
    try {
      const [activities, inputs, expenses, harvests] = await Promise.all([
        listActivities(cycleId), listInputs(cycleId), listExpenses(cycleId), listHarvests(cycleId),
      ]);
      if (refreshRequestId.current !== requestId) return;
      setRows({ activity: activities, input: inputs, expense: expenses, harvest: harvests });
      setLoadedSuccessfully(true);
    } catch (err) {
      if (refreshRequestId.current !== requestId) return;
      setLoadError(err instanceof Error ? err.message : "Catatan budidaya gagal dimuat.");
    } finally {
      if (refreshRequestId.current === requestId) setLoading(false);
    }
  }, [cycleId]);

  useEffect(() => { mountedRef.current = true; void refresh(); return () => { mountedRef.current = false; refreshRequestId.current += 1; }; }, [refresh]);

  useEffect(() => {
    setDate(today());
  }, []);

  function changeKind(next: Kind) {
    formSessionId.current += 1;
    setKind(next);
    setShowForm(false);
    setError("");
    setNotice("");
    setType(next === "activity" ? "watering" : next === "input" ? "fertilizer" : next === "expense" ? "labor" : "");
    setDate(today()); setName(""); setQuantity(""); setUnit(next === "harvest" ? "kg" : ""); setAmount(""); setGrade(""); setNotes("");
  }

  async function submit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    const submittedFormSessionId = formSessionId.current;
    setSaving(true); setError(""); setNotice("");
    try {
      if (kind === "activity") {
        await createActivity(cycleId, { activity_type: type, activity_date: date, notes });
      } else if (kind === "input") {
        await createInput(cycleId, { input_type: type, product_name: name.trim(), quantity: Number(quantity), unit: unit.trim(), applied_date: date, notes });
      } else if (kind === "expense") {
        await createExpense(cycleId, { expense_type: type, description: name.trim(), amount: Number(amount), expense_date: date, notes });
      } else {
        await createHarvest(cycleId, { quantity: Number(quantity), unit: unit.trim(), harvest_date: date, quality_grade: grade.trim(), notes });
      }
      if (!mountedRef.current) return;
      if (isCurrentFormSession(submittedFormSessionId, formSessionId.current)) {
        formSessionId.current += 1;
        setShowForm(false); setName(""); setQuantity(""); setAmount(""); setGrade(""); setNotes("");
      }
      setNotice("Catatan berhasil disimpan.");
      await refresh();
    } catch (err) {
      if (!mountedRef.current) return;
      if (isCurrentFormSession(submittedFormSessionId, formSessionId.current)) {
        setError(err instanceof Error ? err.message : "Catatan belum berhasil disimpan.");
      }
    } finally {
      if (mountedRef.current) setSaving(false);
    }
  }

  const items = rows[kind];
  const totalExpenses = rows.expense.reduce((sum, item) => sum + (item as FarmExpense).amount, 0);
  const summary = [
    { label: "Kegiatan", value: loadedSuccessfully ? String(rows.activity.length) : "—" },
    { label: "Sarana dipakai", value: loadedSuccessfully ? String(rows.input.length) : "—" },
    { label: "Catatan panen", value: loadedSuccessfully ? String(rows.harvest.length) : "—" },
    { label: "Total biaya", value: loadedSuccessfully ? money(totalExpenses) : "—" },
  ];
  return <section className="records-panel">
    <div className="records-head"><div><span className="eyebrow muted">JURNAL BUDIDAYA</span><h4>Catatan siklus</h4></div><button type="button" className="button outline small-button" onClick={() => { formSessionId.current += 1; setShowForm(v => !v); }}>{showForm ? "Tutup" : "+ Catat"}</button></div>
    <div className="journal-summary">{summary.map(item => <div className="journal-metric" key={item.label}><span>{item.label}</span><strong>{item.value}</strong></div>)}</div>
    <div className="record-tabs" role="tablist" aria-label="Jenis catatan">
      {(Object.keys(labels) as Kind[]).map(key => <button type="button" role="tab" aria-selected={kind === key} className={kind === key ? "record-tab active" : "record-tab"} key={key} onClick={() => changeKind(key)}>{labels[key]} <span>{loadedSuccessfully ? rows[key].length : "—"}</span></button>)}
    </div>
    {notice && <p className="success-note" role="status">{notice}</p>}
    {loadError && <div className="form-error" role="alert"><span>{loadError}</span> <button type="button" className="text-button" onClick={() => void refresh()}>Coba lagi</button></div>}
    {error && <p className="form-error" role="alert">{error}</p>}
    {showForm && <form className="record-form" onSubmit={submit}>
      {kind === "activity" && <label className="field"><span>Jenis kegiatan</span><select value={type} onChange={e => { formSessionId.current += 1; setType(e.target.value); }} required><option value="watering">Penyiraman</option><option value="weeding">Penyiangan</option><option value="pruning">Pemangkasan</option><option value="planting">Penanaman</option><option value="inspection">Pemeriksaan</option><option value="other">Lainnya</option></select></label>}
      {kind === "input" && <><label className="field"><span>Jenis sarana</span><select value={type} onChange={e => { formSessionId.current += 1; setType(e.target.value); }}><option value="fertilizer">Pupuk</option><option value="pesticide">Pestisida</option><option value="seed">Benih</option><option value="other">Lainnya</option></select></label><label className="field"><span>Nama produk</span><input value={name} onChange={e => { formSessionId.current += 1; setName(e.target.value); }} required maxLength={120} /></label></>}
      {kind === "expense" && <><label className="field"><span>Kategori biaya</span><select value={type} onChange={e => { formSessionId.current += 1; setType(e.target.value); }}>{["seed","fertilizer","pesticide","labor","irrigation","equipment","transport","other"].map(v => <option key={v} value={v}>{({seed:"Benih",fertilizer:"Pupuk",pesticide:"Pestisida",labor:"Tenaga kerja",irrigation:"Irigasi",equipment:"Peralatan",transport:"Transportasi",other:"Lainnya"} as Record<string,string>)[v]}</option>)}</select></label><label className="field"><span>Deskripsi biaya</span><input value={name} onChange={e => { formSessionId.current += 1; setName(e.target.value); }} required maxLength={160} /></label><label className="field"><span>Jumlah biaya (Rp)</span><input type="number" min="0.01" step="any" value={amount} onChange={e => { formSessionId.current += 1; setAmount(e.target.value); }} required /></label></>}
      {kind === "input" && <div className="record-grid"><label className="field"><span>Kuantitas</span><input type="number" min="0.001" step="any" value={quantity} onChange={e => { formSessionId.current += 1; setQuantity(e.target.value); }} required /></label><label className="field"><span>Satuan</span><input value={unit} onChange={e => { formSessionId.current += 1; setUnit(e.target.value); }} required maxLength={20} /></label></div>}
      {kind === "harvest" && <><div className="record-grid"><label className="field"><span>Jumlah panen</span><input type="number" min="0.001" step="any" value={quantity} onChange={e => { formSessionId.current += 1; setQuantity(e.target.value); }} required /></label><label className="field"><span>Satuan</span><input value={unit} onChange={e => { formSessionId.current += 1; setUnit(e.target.value); }} required maxLength={20} /></label></div><label className="field"><span>Kualitas (opsional)</span><input value={grade} onChange={e => { formSessionId.current += 1; setGrade(e.target.value); }} maxLength={40} placeholder="Contoh: Grade A" /></label></>}
      <label className="field"><span>Tanggal {kind === "activity" ? "kegiatan" : kind === "input" ? "aplikasi" : kind === "expense" ? "biaya" : "panen"}</span><input type="date" value={date} onChange={e => { formSessionId.current += 1; setDate(e.target.value); }} required /></label>
      <label className="field"><span>Catatan (opsional)</span><textarea value={notes} onChange={e => { formSessionId.current += 1; setNotes(e.target.value); }} maxLength={2000} rows={2} /></label>
      <div className="form-actions"><button type="button" className="button ghost" onClick={() => { formSessionId.current += 1; setShowForm(false); }}>Batal</button><button className="button primary" disabled={saving}>{saving ? "Menyimpan…" : "Simpan catatan ↗"}</button></div>
      {kind === "input" && <p className="inline-hint">Catat penggunaan aktual; formulir ini tidak menghitung atau merekomendasikan dosis.</p>}
    </form>}
    {loading && loadedSuccessfully && <p className="no-cycle" role="status">Memuat pembaruan; menampilkan catatan terakhir yang berhasil dimuat.</p>}
    {loading && !loadedSuccessfully ? <p className="no-cycle">Memuat catatan…</p> : !loadedSuccessfully ? <p className="no-cycle" role="status">Catatan belum berhasil dimuat.</p> : items.length === 0 ? <p className="no-cycle">Belum ada catatan {labels[kind].toLowerCase()} untuk siklus ini.</p> : <div className="record-list">{items.map(item => <article className="record-item" key={item.id}>
      <div className="record-item-top"><strong>{kind === "activity" ? activityNames[(item as FarmActivity).activity_type] || (item as FarmActivity).activity_type : kind === "input" ? (item as FarmInput).product_name : kind === "expense" ? (item as FarmExpense).description : "Hasil panen"}</strong><time>{formatDate(kind === "activity" ? (item as FarmActivity).activity_date : kind === "input" ? (item as FarmInput).applied_date : kind === "expense" ? (item as FarmExpense).expense_date : (item as FarmHarvest).harvest_date)}</time></div>
      <p>{kind === "input" ? `${(item as FarmInput).quantity} ${(item as FarmInput).unit} · ${(inputNames[(item as FarmInput).input_type] || (item as FarmInput).input_type)}` : kind === "expense" ? money((item as FarmExpense).amount) + " · " + (expenseNames[(item as FarmExpense).expense_type] || (item as FarmExpense).expense_type) : kind === "harvest" ? `${(item as FarmHarvest).quantity} ${(item as FarmHarvest).unit}${(item as FarmHarvest).quality_grade ? " · " + (item as FarmHarvest).quality_grade : ""}` : (item as FarmActivity).notes || "Kegiatan tercatat"}</p>
      {kind !== "activity" && item.notes && <p className="record-note">{item.notes}</p>}
    </article>)}</div>}
  </section>;
}
