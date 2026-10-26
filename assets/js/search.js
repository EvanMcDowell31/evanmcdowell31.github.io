// Client-side search over the build-time index at /search-index.json.
// Every query word must match somewhere; titles and tags weigh more than body text.
(() => {
  const script = document.currentScript;
  const input = document.getElementById('q');
  const results = document.getElementById('results');
  const status = document.getElementById('status');
  const browse = document.getElementById('browse');
  const form = input.form;
  const filters = form.querySelectorAll('input[name="type"]');

  let index = null;

  const norm = (s) => (s || '').toLowerCase().normalize('NFD').replace(/[\u0300-\u036f]/g, '');
  const esc = (s) => s.replace(/[&<>"']/g, (c) => ({ '&': '&amp;', '<': '&lt;', '>': '&gt;', '"': '&quot;', "'": '&#39;' })[c]);
  const reEsc = (s) => s.replace(/[.*+?^${}()|[\]\\]/g, '\\$&');

  const kindClass = { Story: 'story', Documents: 'document', Bookmarks: 'bookmark', Scripts: 'terminal' };

  function score(item, words) {
    let total = 0;
    for (const w of words) {
      let s = 0;
      if (item._title.includes(w)) s += item._title.startsWith(w) ? 12 : 8;
      if (item._tags.some((t) => t.startsWith(w))) s += 6;
      if (item._summary.includes(w)) s += 3;
      if (item._body.includes(w)) s += 1;
      if (!s) return 0;
      total += s;
    }
    return total;
  }

  function highlight(text, words) {
    let out = esc(text);
    for (const w of words) {
      if (w.length < 2) continue;
      out = out.replace(new RegExp(`(${reEsc(esc(w))})`, 'gi'), '<mark>$1</mark>');
    }
    return out;
  }

  // Short excerpt of the body centred on the first matching word.
  function excerpt(item, words) {
    const body = item.body || '';
    const lower = norm(body);
    const hit = words.map((w) => lower.indexOf(w)).filter((i) => i >= 0).sort((a, b) => a - b)[0];
    if (hit === undefined) return item.summary;
    const start = Math.max(0, hit - 60);
    return (start ? '…' : '') + body.slice(start, start + 160).trim() + '…';
  }

  function render() {
    const raw = input.value.trim();
    const type = form.querySelector('input[name="type"]:checked').value;
    const url = new URL(location);
    raw ? url.searchParams.set('q', raw) : url.searchParams.delete('q');
    type !== 'All' ? url.searchParams.set('type', type) : url.searchParams.delete('type');
    history.replaceState(null, '', url);

    if (!index) return;
    const pool = type === 'All' ? index : index.filter((i) => i.kind === type);

    if (!raw) {
      results.innerHTML = '';
      browse.hidden = false;
      status.textContent = type === 'All' ? '' : `${pool.length} in ${type}`;
      if (type !== 'All') renderList(pool.map((item) => ({ item, s: 1 })), []);
      return;
    }

    const words = norm(raw).split(/\s+/).filter(Boolean);
    const hits = pool.map((item) => ({ item, s: score(item, words) })).filter((h) => h.s).sort((a, b) => b.s - a.s);
    browse.hidden = hits.length > 0;
    status.textContent = hits.length
      ? `${hits.length} ${hits.length === 1 ? 'result' : 'results'} for “${raw}”`
      : `Nothing for “${raw}” yet. Try another word, or wander by tag below.`;
    renderList(hits, words);
  }

  function renderList(hits, words) {
    results.innerHTML = hits.slice(0, 30).map(({ item }) => {
      const href = item.link || item.url;
      const ext = item.link ? ' rel="noopener"' : '';
      const text = words.length ? excerpt(item, words) : item.summary;
      return `<li><a class="result kind-${kindClass[item.kind] || 'document'}" href="${esc(href)}"${ext}>
        <span class="result-kind">${esc(item.kind)}</span>
        <strong>${highlight(item.title, words)}</strong>
        <span class="result-text">${highlight(text || '', words)}</span>
        ${item.link ? `<span class="result-url">${esc(new URL(item.link).host.replace(/^www\./, ''))}</span>` : ''}
      </a></li>`;
    }).join('');
  }

  const params = new URLSearchParams(location.search);
  input.value = params.get('q') || '';
  const startType = params.get('type');
  filters.forEach((f) => { if (f.value === startType) f.checked = true; });

  input.addEventListener('input', render);
  filters.forEach((f) => f.addEventListener('change', render));
  form.addEventListener('submit', (e) => { e.preventDefault(); render(); });
  browse.addEventListener('click', (e) => {
    const tag = e.target.closest('[data-tag]');
    if (!tag) return;
    input.value = tag.dataset.tag;
    input.focus();
    render();
  });

  status.textContent = 'Warming up the index…';
  fetch(script.dataset.index)
    .then((r) => r.json())
    .then((data) => {
      index = data.map((d) => ({
        ...d,
        _title: norm(d.title),
        _summary: norm(d.summary),
        _body: norm(d.body),
        _tags: (d.tags || []).map(norm),
      }));
      status.textContent = '';
      render();
    })
    .catch(() => { status.textContent = 'Search is taking a nap. Try refreshing the page.'; });

  // Skip autofocus on touch screens so the keyboard doesn't cover the page.
  if (!input.value && matchMedia('(hover: hover)').matches) input.focus({ preventScroll: true });
})();
