import { useState } from "react";

type View = "site" | "manager";

interface Booking {
  id: number;
  name: string;
  phone: string;
  service: string;
  date: string;
  time: string;
  note: string;
  status: "new" | "called" | "confirmed" | "cancelled";
  createdAt: string;
}

const SERVICES = [
  "Кузовной ремонт после ДТП",
  "Локальная покраска элемента",
  "Полная перекраска автомобиля",
  "Полировка кузова (одно- / двухступенчатая)",
  "Нанокерамика / защитное покрытие",
  "Тонировка стёкол (плёнка / напыление)",
  "Атмосферная подсветка салона (RGB / RGBW)",
  "Антигравийная плёнка (PPF)",
  "Рихтовка без покраски (PDR)",
  "Восстановление геометрии кузова",
];

const TIMES = [
  "09:00", "09:30", "10:00", "10:30", "11:00", "11:30",
  "12:00", "13:00", "13:30", "14:00", "14:30", "15:00",
  "15:30", "16:00", "16:30", "17:00", "17:30", "18:00",
];

const initialBookings: Booking[] = [
  { id: 1, name: "Алексей Петров", phone: "+7 (912) 345-67-89", service: "Полировка кузова (одно- / двухступенчатая)", date: "2026-08-20", time: "10:00", note: "Toyota Camry 2019, белый перламутр", status: "called", createdAt: "2026-08-18 09:14" },
  { id: 2, name: "Мария Соколова", phone: "+7 (903) 211-55-44", service: "Атмосферная подсветка салона (RGB / RGBW)", date: "2026-08-21", time: "11:30", note: "BMW X5, хочет синий цвет", status: "new", createdAt: "2026-08-18 11:02" },
  { id: 3, name: "Дмитрий Захаров", phone: "+7 (925) 778-12-30", service: "Кузовной ремонт после ДТП", date: "2026-08-22", time: "09:00", note: "VW Polo, помята правая дверь и крыло", status: "confirmed", createdAt: "2026-08-17 16:45" },
];

const statusLabel: Record<Booking["status"], string> = {
  new: "Новая",
  called: "Позвонили",
  confirmed: "Подтверждена",
  cancelled: "Отменена",
};

const statusColor: Record<Booking["status"], string> = {
  new: "bg-amber-500 text-black",
  called: "bg-blue-600 text-white",
  confirmed: "bg-green-600 text-white",
  cancelled: "bg-zinc-600 text-white",
};

export default function App() {
  const [view, setView] = useState<View>("site");
  const [bookings, setBookings] = useState<Booking[]>(initialBookings);
  const [navOpen, setNavOpen] = useState(false);
  const [form, setForm] = useState({ name: "", phone: "", service: "", date: "", time: "", note: "" });
  const [formStatus, setFormStatus] = useState<"idle" | "success" | "error">("idle");
  const [managerPassword, setManagerPassword] = useState("");
  const [managerAuth, setManagerAuth] = useState(false);
  const [authError, setAuthError] = useState(false);
  const [filterStatus, setFilterStatus] = useState<"all" | Booking["status"]>("all");

  // modal: create client
  const [showCreate, setShowCreate] = useState(false);
  const emptyCreate = { name: "", phone: "", service: "", date: "", time: "", note: "" };
  const [createForm, setCreateForm] = useState(emptyCreate);
  const [createError, setCreateError] = useState(false);

  // modal: note view/edit
  const [noteModal, setNoteModal] = useState<{ id: number; note: string } | null>(null);

  const submitBooking = (e: React.FormEvent) => {
    e.preventDefault();
    if (!form.name || !form.phone || !form.service || !form.date || !form.time) {
      setFormStatus("error");
      return;
    }
    const newBooking: Booking = {
      id: Date.now(),
      ...form,
      status: "new",
      createdAt: new Date().toLocaleString("ru-RU").replace(",", ""),
    };
    setBookings((prev) => [newBooking, ...prev]);
    setForm({ name: "", phone: "", service: "", date: "", time: "", note: "" });
    setFormStatus("success");
    setTimeout(() => setFormStatus("idle"), 4000);
  };

  const updateStatus = (id: number, status: Booking["status"]) => {
    setBookings((prev) => prev.map((b) => (b.id === id ? { ...b, status } : b)));
  };

  const submitCreate = (e: React.FormEvent) => {
    e.preventDefault();
    if (!createForm.name || !createForm.phone || !createForm.service || !createForm.date || !createForm.time) {
      setCreateError(true);
      return;
    }
    const newBooking: Booking = {
      id: Date.now(),
      ...createForm,
      status: "new",
      createdAt: new Date().toLocaleString("ru-RU").replace(",", ""),
    };
    setBookings((prev) => [newBooking, ...prev]);
    setCreateForm(emptyCreate);
    setCreateError(false);
    setShowCreate(false);
  };

  const saveNote = () => {
    if (!noteModal) return;
    setBookings((prev) => prev.map((b) => (b.id === noteModal.id ? { ...b, note: noteModal.note } : b)));
    setNoteModal(null);
  };

  const handleManagerLogin = (e: React.FormEvent) => {
    e.preventDefault();
    if (managerPassword === "admin123") {
      setManagerAuth(true);
      setAuthError(false);
    } else {
      setAuthError(true);
    }
  };

  const navLinks = [
    { label: "Услуги", href: "#services" },
    { label: "О нас", href: "#about" },
    { label: "Запись", href: "#booking" },
    { label: "Контакты", href: "#contacts" },
  ];

  const filtered = filterStatus === "all" ? bookings : bookings.filter((b) => b.status === filterStatus);

  if (view === "manager") {
    if (!managerAuth) {
      return (
        <div className="min-h-screen flex items-center justify-center bg-[#0d0d0d] px-4">
          <div className="w-full max-w-sm">
            <div className="mb-8">
              <div className="text-[#e84a0e] font-['Oswald'] text-2xl font-bold tracking-widest uppercase mb-1">АвтоМастер</div>
              <h2 className="font-['Oswald'] text-3xl font-bold uppercase text-[#f0ede8]">Панель менеджера</h2>
            </div>
            <form onSubmit={handleManagerLogin} className="space-y-4">
              <input
                type="password"
                placeholder="Пароль"
                value={managerPassword}
                onChange={(e) => setManagerPassword(e.target.value)}
                className="w-full bg-[#161616] border border-[#2a2a2a] text-[#f0ede8] px-4 py-3 focus:outline-none focus:border-[#e84a0e] text-sm placeholder:text-[#4a4a4a]"
              />
              {authError && <p className="text-[#e84a0e] text-sm">Неверный пароль</p>}
              <button
                type="submit"
                className="w-full bg-[#e84a0e] text-white font-['Oswald'] font-semibold uppercase tracking-widest py-3 hover:bg-[#c73c0a] transition-colors"
              >
                Войти
              </button>
              <button
                type="button"
                onClick={() => setView("site")}
                className="w-full text-[#7a7670] text-sm hover:text-[#f0ede8] transition-colors mt-2"
              >
                ← На сайт
              </button>
            </form>
          </div>
        </div>
      );
    }

    return (
      <div className="min-h-screen bg-[#0d0d0d] text-[#f0ede8]">
        {/* Manager header */}
        <header className="border-b border-[#2a2a2a] px-6 py-4 flex items-center justify-between">
          <div className="flex items-center gap-4">
            <span className="text-[#e84a0e] font-['Oswald'] text-xl font-bold tracking-widest uppercase">АвтоМастер</span>
            <span className="text-[#7a7670] text-sm">/ Записи клиентов</span>
          </div>
          <div className="flex items-center gap-4">
            <button
              onClick={() => { setCreateForm(emptyCreate); setCreateError(false); setShowCreate(true); }}
              className="bg-[#e84a0e] text-white font-['Oswald'] uppercase tracking-widest text-xs px-4 py-2 hover:bg-[#c73c0a] transition-colors"
            >
              + Новый клиент
            </button>
            <button onClick={() => setView("site")} className="text-[#7a7670] text-sm hover:text-[#f0ede8] transition-colors">
              На сайт
            </button>
            <button onClick={() => { setManagerAuth(false); setManagerPassword(""); }} className="text-[#7a7670] text-sm hover:text-[#e84a0e] transition-colors">
              Выйти
            </button>
          </div>
        </header>

        <div className="max-w-7xl mx-auto px-4 md:px-8 py-8">
          {/* Stats */}
          <div className="grid grid-cols-2 md:grid-cols-4 gap-3 mb-8">
            {(["all", "new", "called", "confirmed"] as const).map((s) => (
              <button
                key={s}
                onClick={() => setFilterStatus(s)}
                className={`p-4 border text-left transition-all ${filterStatus === s ? "border-[#e84a0e]" : "border-[#2a2a2a] hover:border-[#e84a0e]/40"} bg-[#161616]`}
              >
                <div className="font-['Oswald'] text-2xl font-bold text-[#e84a0e]">
                  {s === "all" ? bookings.length : bookings.filter((b) => b.status === s).length}
                </div>
                <div className="text-[#7a7670] text-xs uppercase tracking-wider mt-1">
                  {s === "all" ? "Всего" : statusLabel[s]}
                </div>
              </button>
            ))}
          </div>

          {/* Bookings table */}
          <div className="border border-[#2a2a2a] overflow-x-auto">
            <table className="w-full text-sm">
              <thead>
                <tr className="border-b border-[#2a2a2a] bg-[#161616]">
                  {["Клиент", "Телефон", "Услуга", "Дата / Время", "Примечание", "Статус", "Действие"].map((h) => (
                    <th key={h} className="px-4 py-3 text-left font-['Oswald'] text-xs uppercase tracking-wider text-[#7a7670] font-medium whitespace-nowrap">{h}</th>
                  ))}
                </tr>
              </thead>
              <tbody>
                {filtered.length === 0 && (
                  <tr><td colSpan={7} className="px-4 py-8 text-center text-[#7a7670]">Записей нет</td></tr>
                )}
                {filtered.map((b) => (
                  <tr key={b.id} className="border-b border-[#1e1e1e] hover:bg-[#161616] transition-colors">
                    <td className="px-4 py-3 font-medium whitespace-nowrap">{b.name}</td>
                    <td className="px-4 py-3">
                      <a href={`tel:${b.phone}`} className="text-[#e84a0e] hover:underline whitespace-nowrap">{b.phone}</a>
                    </td>
                    <td className="px-4 py-3 text-[#a09c96] max-w-[160px]">{b.service}</td>
                    <td className="px-4 py-3 whitespace-nowrap">
                      <div>{new Date(b.date).toLocaleDateString("ru-RU", { day: "numeric", month: "short" })}</div>
                      <div className="text-[#7a7670]">{b.time}</div>
                    </td>
                    <td className="px-4 py-3 max-w-[180px]">
                      <button
                        onClick={() => setNoteModal({ id: b.id, note: b.note })}
                        className="text-left w-full group"
                        title="Нажмите для просмотра / редактирования"
                      >
                        <span className="text-[#a09c96] line-clamp-2 group-hover:text-[#f0ede8] transition-colors">
                          {b.note || <span className="text-[#3a3a3a] italic">нет</span>}
                        </span>
                        <span className="text-[#e84a0e] text-[10px] uppercase tracking-wider opacity-0 group-hover:opacity-100 transition-opacity">
                          изменить
                        </span>
                      </button>
                    </td>
                    <td className="px-4 py-3">
                      <span className={`px-2 py-1 text-xs font-['Oswald'] uppercase tracking-wider ${statusColor[b.status]}`}>
                        {statusLabel[b.status]}
                      </span>
                    </td>
                    <td className="px-4 py-3">
                      <select
                        value={b.status}
                        onChange={(e) => updateStatus(b.id, e.target.value as Booking["status"])}
                        className="bg-[#1e1e1e] border border-[#2a2a2a] text-[#f0ede8] text-xs px-2 py-1 focus:outline-none focus:border-[#e84a0e]"
                      >
                        <option value="new">Новая</option>
                        <option value="called">Позвонили</option>
                        <option value="confirmed">Подтверждена</option>
                        <option value="cancelled">Отменена</option>
                      </select>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
          <p className="text-[#4a4a4a] text-xs mt-4">Записей: {filtered.length}. Создано в текущей сессии.</p>
        </div>

        {/* ── Modal: Create client ─────────────────────────────────── */}
        {showCreate && (
          <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/70 px-4" onClick={() => setShowCreate(false)}>
            <div className="bg-[#161616] border border-[#2a2a2a] w-full max-w-lg p-6 md:p-8" onClick={(e) => e.stopPropagation()}>
              <div className="flex items-center justify-between mb-6">
                <h3 className="font-['Oswald'] text-xl font-bold uppercase text-[#f0ede8]">Новый клиент</h3>
                <button onClick={() => setShowCreate(false)} className="text-[#7a7670] hover:text-[#f0ede8] text-xl leading-none">&times;</button>
              </div>
              <form onSubmit={submitCreate} className="flex flex-col gap-4">
                <div className="grid grid-cols-2 gap-4">
                  <div>
                    <label className="block text-xs text-[#7a7670] uppercase tracking-wider mb-1.5">Имя *</label>
                    <input type="text" value={createForm.name} onChange={(e) => setCreateForm({ ...createForm, name: e.target.value })}
                      placeholder="Иван Петров"
                      className="w-full bg-[#0d0d0d] border border-[#2a2a2a] text-[#f0ede8] px-3 py-2.5 text-sm focus:outline-none focus:border-[#e84a0e] placeholder:text-[#3a3a3a]" />
                  </div>
                  <div>
                    <label className="block text-xs text-[#7a7670] uppercase tracking-wider mb-1.5">Телефон *</label>
                    <input type="tel" value={createForm.phone} onChange={(e) => setCreateForm({ ...createForm, phone: e.target.value })}
                      placeholder="+7 (___) ___-__-__"
                      className="w-full bg-[#0d0d0d] border border-[#2a2a2a] text-[#f0ede8] px-3 py-2.5 text-sm focus:outline-none focus:border-[#e84a0e] placeholder:text-[#3a3a3a]" />
                  </div>
                </div>
                <div>
                  <label className="block text-xs text-[#7a7670] uppercase tracking-wider mb-1.5">Услуга *</label>
                  <select value={createForm.service} onChange={(e) => setCreateForm({ ...createForm, service: e.target.value })}
                    className="w-full bg-[#0d0d0d] border border-[#2a2a2a] text-[#f0ede8] px-3 py-2.5 text-sm focus:outline-none focus:border-[#e84a0e]">
                    <option value="">Выберите услугу</option>
                    {SERVICES.map((s) => <option key={s} value={s}>{s}</option>)}
                  </select>
                </div>
                <div className="grid grid-cols-2 gap-4">
                  <div>
                    <label className="block text-xs text-[#7a7670] uppercase tracking-wider mb-1.5">Дата *</label>
                    <input type="date" value={createForm.date} min={new Date().toISOString().split("T")[0]}
                      onChange={(e) => setCreateForm({ ...createForm, date: e.target.value })}
                      className="w-full bg-[#0d0d0d] border border-[#2a2a2a] text-[#f0ede8] px-3 py-2.5 text-sm focus:outline-none focus:border-[#e84a0e]" />
                  </div>
                  <div>
                    <label className="block text-xs text-[#7a7670] uppercase tracking-wider mb-1.5">Время *</label>
                    <select value={createForm.time} onChange={(e) => setCreateForm({ ...createForm, time: e.target.value })}
                      className="w-full bg-[#0d0d0d] border border-[#2a2a2a] text-[#f0ede8] px-3 py-2.5 text-sm focus:outline-none focus:border-[#e84a0e]">
                      <option value="">— выберите —</option>
                      {TIMES.map((t) => <option key={t} value={t}>{t}</option>)}
                    </select>
                  </div>
                </div>
                <div>
                  <label className="block text-xs text-[#7a7670] uppercase tracking-wider mb-1.5">Примечание</label>
                  <textarea value={createForm.note} onChange={(e) => setCreateForm({ ...createForm, note: e.target.value })}
                    rows={3} placeholder="Марка авто, цвет, описание задачи..."
                    className="w-full bg-[#0d0d0d] border border-[#2a2a2a] text-[#f0ede8] px-3 py-2.5 text-sm focus:outline-none focus:border-[#e84a0e] placeholder:text-[#3a3a3a] resize-none" />
                </div>
                {createError && <p className="text-[#e84a0e] text-sm">Заполните все обязательные поля (*).</p>}
                <div className="flex gap-3 mt-1">
                  <button type="submit"
                    className="flex-1 bg-[#e84a0e] text-white font-['Oswald'] uppercase tracking-widest text-sm py-3 hover:bg-[#c73c0a] transition-colors">
                    Создать запись
                  </button>
                  <button type="button" onClick={() => setShowCreate(false)}
                    className="px-5 border border-[#2a2a2a] text-[#7a7670] text-sm hover:text-[#f0ede8] hover:border-[#f0ede8] transition-colors">
                    Отмена
                  </button>
                </div>
              </form>
            </div>
          </div>
        )}

        {/* ── Modal: Note view / edit ──────────────────────────────── */}
        {noteModal && (
          <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/70 px-4" onClick={() => setNoteModal(null)}>
            <div className="bg-[#161616] border border-[#2a2a2a] w-full max-w-md p-6 md:p-8" onClick={(e) => e.stopPropagation()}>
              <div className="flex items-center justify-between mb-4">
                <h3 className="font-['Oswald'] text-lg font-bold uppercase text-[#f0ede8]">Примечание</h3>
                <button onClick={() => setNoteModal(null)} className="text-[#7a7670] hover:text-[#f0ede8] text-xl leading-none">&times;</button>
              </div>
              <textarea
                value={noteModal.note}
                onChange={(e) => setNoteModal({ ...noteModal, note: e.target.value })}
                rows={6}
                placeholder="Добавьте примечание..."
                className="w-full bg-[#0d0d0d] border border-[#2a2a2a] text-[#f0ede8] px-3 py-2.5 text-sm focus:outline-none focus:border-[#e84a0e] placeholder:text-[#3a3a3a] resize-none mb-4"
              />
              <div className="flex gap-3">
                <button onClick={saveNote}
                  className="flex-1 bg-[#e84a0e] text-white font-['Oswald'] uppercase tracking-widest text-sm py-2.5 hover:bg-[#c73c0a] transition-colors">
                  Сохранить
                </button>
                <button onClick={() => setNoteModal(null)}
                  className="px-5 border border-[#2a2a2a] text-[#7a7670] text-sm hover:text-[#f0ede8] hover:border-[#f0ede8] transition-colors">
                  Отмена
                </button>
              </div>
            </div>
          </div>
        )}
      </div>
    );
  }

  // ── PUBLIC SITE ─────────────────────────────────────────────────────────────
  return (
    <div className="min-h-screen bg-[#0d0d0d] text-[#f0ede8]">

      {/* NAV */}
      <nav className="fixed top-0 left-0 right-0 z-50 border-b border-[#2a2a2a] bg-[#0d0d0d]/95 backdrop-blur-sm">
        <div className="max-w-7xl mx-auto px-4 md:px-8 flex items-center justify-between h-14">
          <a href="#" className="font-['Oswald'] text-xl font-bold text-[#e84a0e] tracking-widest uppercase">
            АвтоМастер
          </a>
          <div className="hidden md:flex items-center gap-8">
            {navLinks.map((l) => (
              <a key={l.label} href={l.href} className="text-sm text-[#a09c96] hover:text-[#f0ede8] transition-colors uppercase tracking-wider font-medium">
                {l.label}
              </a>
            ))}
            <a href="#booking" className="bg-[#e84a0e] text-white font-['Oswald'] uppercase tracking-widest text-sm px-4 py-2 hover:bg-[#c73c0a] transition-colors">
              Записаться
            </a>
          </div>
          <button onClick={() => setNavOpen(!navOpen)} className="md:hidden flex flex-col gap-1.5 p-1">
            <span className={`block h-0.5 w-5 bg-[#f0ede8] transition-all ${navOpen ? "rotate-45 translate-y-2" : ""}`} />
            <span className={`block h-0.5 w-5 bg-[#f0ede8] transition-all ${navOpen ? "opacity-0" : ""}`} />
            <span className={`block h-0.5 w-5 bg-[#f0ede8] transition-all ${navOpen ? "-rotate-45 -translate-y-2" : ""}`} />
          </button>
        </div>
        {navOpen && (
          <div className="md:hidden border-t border-[#2a2a2a] bg-[#0d0d0d] px-4 py-4 flex flex-col gap-4">
            {navLinks.map((l) => (
              <a key={l.label} href={l.href} onClick={() => setNavOpen(false)} className="text-[#a09c96] hover:text-[#f0ede8] uppercase tracking-wider text-sm">
                {l.label}
              </a>
            ))}
          </div>
        )}
      </nav>

      {/* HERO */}
      <section className="relative min-h-screen flex items-end pt-14 overflow-hidden">
        <div className="absolute inset-0">
          <img
            src="https://images.unsplash.com/photo-1632605185825-fd583793fa73?w=1600&h=900&fit=crop&auto=format"
            alt="Мастер в покрасочном боксе"
            className="w-full h-full object-cover object-center"
          />
          <div className="absolute inset-0 bg-gradient-to-t from-[#0d0d0d] via-[#0d0d0d]/60 to-[#0d0d0d]/20" />
          <div className="absolute inset-0 bg-gradient-to-r from-[#0d0d0d]/80 to-transparent" />
        </div>
        <div className="relative z-10 max-w-7xl mx-auto px-4 md:px-8 pb-20 md:pb-32 w-full">
          <div className="max-w-2xl">
            <div className="flex items-center gap-3 mb-6">
              <span className="w-10 h-0.5 bg-[#e84a0e]" />
              <span className="text-[#e84a0e] text-sm font-medium uppercase tracking-widest">Москва, с 2010 года</span>
            </div>
            <h1 className="font-['Oswald'] text-5xl md:text-7xl font-bold uppercase leading-none mb-6 text-[#f0ede8]">
              Авто<br />
              <span className="text-[#e84a0e]">Мастер</span>
            </h1>
            <p className="text-[#a09c96] text-lg md:text-xl leading-relaxed mb-10 max-w-lg">
              Кузовной ремонт, покраска, полировка, тонировка и атмосферная подсветка салона. Ваш автомобиль — наше искусство.
            </p>
            <div className="flex flex-col sm:flex-row gap-4">
              <a href="#booking" className="inline-flex items-center justify-center bg-[#e84a0e] text-white font-['Oswald'] uppercase tracking-widest text-base px-8 py-4 hover:bg-[#c73c0a] transition-colors">
                Записаться онлайн
              </a>
              <a href="tel:+74951234567" className="inline-flex items-center justify-center border border-[#f0ede8] text-[#f0ede8] font-['Oswald'] uppercase tracking-widest text-base px-8 py-4 hover:bg-[#f0ede8] hover:text-[#0d0d0d] transition-colors">
                +7 (495) 123-45-67
              </a>
            </div>
          </div>
        </div>
        {/* Scroll indicator */}
        <div className="absolute bottom-8 right-8 md:right-16 z-10 flex flex-col items-center gap-2">
          <span className="text-[#7a7670] text-xs uppercase tracking-widest rotate-90 origin-center translate-x-4">Scroll</span>
          <div className="w-px h-12 bg-gradient-to-b from-[#e84a0e] to-transparent" />
        </div>
      </section>

      {/* TRUST BAR */}
      <section className="border-y border-[#2a2a2a] bg-[#161616]">
        <div className="max-w-7xl mx-auto px-4 md:px-8">
          <div className="grid grid-cols-2 md:grid-cols-4 divide-x divide-[#2a2a2a]">
            {[
              { num: "15+", label: "Лет на рынке" },
              { num: "8 500+", label: "Кузовных работ" },
              { num: "98%", label: "Клиентов возвращаются" },
              { num: "2 года", label: "Гарантия на покраску" },
            ].map((item) => (
              <div key={item.label} className="px-6 md:px-10 py-8 text-center">
                <div className="font-['Oswald'] text-3xl md:text-4xl font-bold text-[#e84a0e] mb-1">{item.num}</div>
                <div className="text-[#7a7670] text-xs uppercase tracking-wider">{item.label}</div>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* SERVICES */}
      <section id="services" className="py-20 md:py-28 max-w-7xl mx-auto px-4 md:px-8">
        <div className="mb-14">
          <span className="accent-line" />
          <h2 className="font-['Oswald'] text-4xl md:text-5xl font-bold uppercase text-[#f0ede8]">Наши услуги</h2>
        </div>
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-px bg-[#2a2a2a]">
          {[
            { icon: "🔨", title: "Кузовной ремонт", desc: "Восстановление после ДТП: рихтовка, сварка, вытяжка геометрии на стапеле. Любые повреждения." },
            { icon: "🎨", title: "Покраска", desc: "Локальная покраска одного элемента или полная перекраска в камере. Подбор цвета по VIN-коду." },
            { icon: "✨", title: "Полировка", desc: "Одно- и двухступенчатая полировка, абразивное выравнивание, восстановление блеска лакокрасочного покрытия." },
            { icon: "🪟", title: "Тонировка стёкол", desc: "Профессиональная тонировочная плёнка — любой процент светопропускания. Напыление под заказ." },
            { icon: "💡", title: "Атмосферная подсветка", desc: "Монтаж RGB/RGBW лент и точечных диодов в салоне. Управление со смартфона, синхронизация с музыкой." },
            { icon: "🛡️", title: "Антигравийная плёнка (PPF)", desc: "Защита кузова от сколов, царапин и химии. Прозрачная самовосстанавливающаяся плёнка." },
            { icon: "💎", title: "Нанокерамика", desc: "Нанесение керамического защитного покрытия. Гидрофобный эффект, блеск и защита до 3 лет." },
            { icon: "🔧", title: "PDR — рихтовка без покраски", desc: "Удаление вмятин без нарушения лакокрасочного покрытия. Быстро, дёшево, незаметно." },
            { icon: "🚗", title: "Восстановление геометрии", desc: "Выправление лонжеронов и порогов на стапеле Blackhawk. Контроль в 5 точках по паспортным данным." },
          ].map((s) => (
            <div key={s.title} className="bg-[#0d0d0d] p-6 md:p-8 group hover:bg-[#161616] transition-colors cursor-default">
              <div className="text-2xl mb-4">{s.icon}</div>
              <h3 className="font-['Oswald'] text-lg font-semibold uppercase text-[#f0ede8] mb-2 group-hover:text-[#e84a0e] transition-colors">{s.title}</h3>
              <p className="text-[#7a7670] text-sm leading-relaxed">{s.desc}</p>
            </div>
          ))}
        </div>
      </section>

      {/* ABOUT */}
      <section id="about" className="bg-[#161616] border-y border-[#2a2a2a]">
        <div className="max-w-7xl mx-auto px-4 md:px-8 py-20 md:py-28">
          <div className="grid grid-cols-1 md:grid-cols-2 gap-12 md:gap-20 items-center">
            <div>
              <span className="accent-line" />
              <h2 className="font-['Oswald'] text-4xl md:text-5xl font-bold uppercase text-[#f0ede8] mb-6">
                О сервисе
              </h2>
              <p className="text-[#a09c96] leading-relaxed mb-6">
                Мы открылись в 2010 году и специализируемся исключительно на кузовных и малярных работах. Наши мастера — художники: они не просто красят металл, они возвращают автомобилю идеальный вид или создают уникальный образ.
              </p>
              <p className="text-[#a09c96] leading-relaxed mb-8">
                Покрасочная камера итальянского производства, стапель Blackhawk, профессиональные полировальные машинки Rupes и плёнки Llumar/STEK для PPF и тонировки. Каждый проект фотографируется до и после.
              </p>
              <div className="flex flex-col gap-3">
                {[
                  "Покрасочная камера с инфракрасной сушкой",
                  "Гарантия 2 года на малярные работы",
                  "Подбор цвета по VIN — точное совпадение",
                  "Фото- и видеоотчёт по каждому этапу",
                ].map((item) => (
                  <div key={item} className="flex items-center gap-3">
                    <span className="w-1.5 h-1.5 bg-[#e84a0e] rounded-full flex-shrink-0" />
                    <span className="text-[#f0ede8] text-sm">{item}</span>
                  </div>
                ))}
              </div>
            </div>
            <div className="relative">
              <div className="absolute -top-3 -left-3 w-full h-full border border-[#e84a0e]/30" />
              <img
                src="https://images.unsplash.com/photo-1632605192331-085fa2082575?w=700&h=500&fit=crop&auto=format"
                alt="Мастер красит автомобиль в покрасочной камере"
                className="w-full h-72 md:h-96 object-cover object-center relative z-10"
              />
            </div>
          </div>
        </div>
      </section>

      {/* BOOKING */}
      <section id="booking" className="py-20 md:py-28 max-w-7xl mx-auto px-4 md:px-8">
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-12 md:gap-20">
          <div>
            <span className="accent-line" />
            <h2 className="font-['Oswald'] text-4xl md:text-5xl font-bold uppercase text-[#f0ede8] mb-4">
              Запись на сервис
            </h2>
            <p className="text-[#7a7670] leading-relaxed mb-8">
              Заполните форму — менеджер позвонит в течение 30 минут, уточнит детали и назначит удобное время. Запись работает круглосуточно.
            </p>
            <div className="flex flex-col gap-4 text-sm text-[#a09c96]">
              <div className="flex items-center gap-3">
                <span className="w-8 h-8 bg-[#e84a0e] flex items-center justify-center font-['Oswald'] font-bold text-white text-xs">1</span>
                Заполните форму с данными
              </div>
              <div className="flex items-center gap-3">
                <span className="w-8 h-8 bg-[#2a2a2a] flex items-center justify-center font-['Oswald'] font-bold text-[#7a7670] text-xs">2</span>
                Менеджер позвонит и подтвердит время
              </div>
              <div className="flex items-center gap-3">
                <span className="w-8 h-8 bg-[#2a2a2a] flex items-center justify-center font-['Oswald'] font-bold text-[#7a7670] text-xs">3</span>
                Приезжайте — ваш мастер уже ждёт
              </div>
            </div>
          </div>

          <div className="bg-[#161616] border border-[#2a2a2a] p-6 md:p-8">
            {formStatus === "success" ? (
              <div className="flex flex-col items-start justify-center h-full min-h-[300px] gap-4">
                <div className="w-12 h-12 bg-green-600 flex items-center justify-center text-white text-2xl">✓</div>
                <h3 className="font-['Oswald'] text-2xl font-bold uppercase text-[#f0ede8]">Заявка принята!</h3>
                <p className="text-[#a09c96] text-sm leading-relaxed">
                  Менеджер свяжется с вами по указанному номеру в течение 30 минут для подтверждения записи.
                </p>
                <button onClick={() => setFormStatus("idle")} className="mt-2 text-[#e84a0e] text-sm hover:underline">
                  Записаться ещё раз
                </button>
              </div>
            ) : (
              <form onSubmit={submitBooking} className="flex flex-col gap-4">
                <h3 className="font-['Oswald'] text-xl font-semibold uppercase text-[#f0ede8] mb-2">Ваши данные</h3>

                <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                  <div>
                    <label className="block text-xs text-[#7a7670] uppercase tracking-wider mb-1.5">Имя *</label>
                    <input
                      type="text"
                      value={form.name}
                      onChange={(e) => setForm({ ...form, name: e.target.value })}
                      placeholder="Иван Петров"
                      className="w-full bg-[#0d0d0d] border border-[#2a2a2a] text-[#f0ede8] px-3 py-2.5 text-sm focus:outline-none focus:border-[#e84a0e] placeholder:text-[#3a3a3a] transition-colors"
                    />
                  </div>
                  <div>
                    <label className="block text-xs text-[#7a7670] uppercase tracking-wider mb-1.5">Телефон *</label>
                    <input
                      type="tel"
                      value={form.phone}
                      onChange={(e) => setForm({ ...form, phone: e.target.value })}
                      placeholder="+7 (___) ___-__-__"
                      className="w-full bg-[#0d0d0d] border border-[#2a2a2a] text-[#f0ede8] px-3 py-2.5 text-sm focus:outline-none focus:border-[#e84a0e] placeholder:text-[#3a3a3a] transition-colors"
                    />
                  </div>
                </div>

                <div>
                  <label className="block text-xs text-[#7a7670] uppercase tracking-wider mb-1.5">Услуга *</label>
                  <select
                    value={form.service}
                    onChange={(e) => setForm({ ...form, service: e.target.value })}
                    className="w-full bg-[#0d0d0d] border border-[#2a2a2a] text-[#f0ede8] px-3 py-2.5 text-sm focus:outline-none focus:border-[#e84a0e] transition-colors"
                  >
                    <option value="">Выберите услугу</option>
                    {SERVICES.map((s) => <option key={s} value={s}>{s}</option>)}
                  </select>
                </div>

                <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                  <div>
                    <label className="block text-xs text-[#7a7670] uppercase tracking-wider mb-1.5">Дата *</label>
                    <input
                      type="date"
                      value={form.date}
                      min={new Date().toISOString().split("T")[0]}
                      onChange={(e) => setForm({ ...form, date: e.target.value })}
                      className="w-full bg-[#0d0d0d] border border-[#2a2a2a] text-[#f0ede8] px-3 py-2.5 text-sm focus:outline-none focus:border-[#e84a0e] transition-colors"
                    />
                  </div>
                  <div>
                    <label className="block text-xs text-[#7a7670] uppercase tracking-wider mb-1.5">Время *</label>
                    <select
                      value={form.time}
                      onChange={(e) => setForm({ ...form, time: e.target.value })}
                      className="w-full bg-[#0d0d0d] border border-[#2a2a2a] text-[#f0ede8] px-3 py-2.5 text-sm focus:outline-none focus:border-[#e84a0e] transition-colors"
                    >
                      <option value="">— выберите —</option>
                      {TIMES.map((t) => <option key={t} value={t}>{t}</option>)}
                    </select>
                  </div>
                </div>

                <div>
                  <label className="block text-xs text-[#7a7670] uppercase tracking-wider mb-1.5">Марка и модель авто / Примечание</label>
                  <textarea
                    value={form.note}
                    onChange={(e) => setForm({ ...form, note: e.target.value })}
                    rows={2}
                    placeholder="Toyota Camry 2020, пробег 85 000 км..."
                    className="w-full bg-[#0d0d0d] border border-[#2a2a2a] text-[#f0ede8] px-3 py-2.5 text-sm focus:outline-none focus:border-[#e84a0e] placeholder:text-[#3a3a3a] resize-none transition-colors"
                  />
                </div>

                {formStatus === "error" && (
                  <p className="text-[#e84a0e] text-sm">Заполните все обязательные поля (*).</p>
                )}

                <button
                  type="submit"
                  className="bg-[#e84a0e] text-white font-['Oswald'] uppercase tracking-widest text-sm px-6 py-3.5 hover:bg-[#c73c0a] transition-colors mt-1"
                >
                  Отправить заявку
                </button>
                <p className="text-[#4a4a4a] text-xs">
                  Нажимая кнопку, вы соглашаетесь с обработкой персональных данных.
                </p>
              </form>
            )}
          </div>
        </div>
      </section>

      {/* CONTACTS */}
      <section id="contacts" className="bg-[#161616] border-t border-[#2a2a2a]">
        <div className="max-w-7xl mx-auto px-4 md:px-8 py-20 md:py-28">
          <div className="grid grid-cols-1 md:grid-cols-3 gap-10 md:gap-16">
            <div>
              <span className="accent-line" />
              <h2 className="font-['Oswald'] text-4xl md:text-5xl font-bold uppercase text-[#f0ede8] mb-6">Контакты</h2>
              <p className="text-[#7a7670] text-sm leading-relaxed">
                Приезжайте без записи на экспресс-услуги или звоните — ответим на любые вопросы.
              </p>
            </div>
            <div className="flex flex-col gap-6">
              {[
                { label: "Телефон", value: "+7 (495) 123-45-67", href: "tel:+74951234567" },
                { label: "WhatsApp", value: "+7 (912) 987-65-43", href: "https://wa.me/79129876543" },
                { label: "E-mail", value: "info@avtomaster.ru", href: "mailto:info@avtomaster.ru" },
              ].map((c) => (
                <div key={c.label}>
                  <div className="text-xs text-[#7a7670] uppercase tracking-wider mb-1">{c.label}</div>
                  <a href={c.href} className="text-[#f0ede8] hover:text-[#e84a0e] transition-colors font-medium">{c.value}</a>
                </div>
              ))}
            </div>
            <div className="flex flex-col gap-6">
              <div>
                <div className="text-xs text-[#7a7670] uppercase tracking-wider mb-1">Адрес</div>
                <div className="text-[#f0ede8] font-medium">г. Москва, ул. Авиамоторная, 12с3</div>
                <div className="text-[#7a7670] text-sm mt-1">м. Авиамоторная, 5 мин. пешком</div>
              </div>
              <div>
                <div className="text-xs text-[#7a7670] uppercase tracking-wider mb-1">Режим работы</div>
                <div className="text-[#f0ede8] font-medium">Пн–Пт: 09:00–20:00</div>
                <div className="text-[#f0ede8]">Сб–Вс: 10:00–18:00</div>
              </div>
            </div>
          </div>
        </div>
      </section>

      {/* FOOTER */}
      <footer className="border-t border-[#2a2a2a] bg-[#0d0d0d] px-4 md:px-8 py-6">
        <div className="max-w-7xl mx-auto flex flex-col md:flex-row items-center justify-between gap-4 text-xs text-[#4a4a4a]">
          <span className="font-['Oswald'] text-[#e84a0e] text-base font-bold tracking-widest uppercase">АвтоМастер</span>
          <span>© 2010–2026 АвтоМастер. Все права защищены.</span>
          <button
            onClick={() => setView("manager")}
            className="text-[#2a2a2a] hover:text-[#4a4a4a] transition-colors text-[10px] uppercase tracking-widest"
          >
            Панель менеджера
          </button>
        </div>
      </footer>
    </div>
  );
}
