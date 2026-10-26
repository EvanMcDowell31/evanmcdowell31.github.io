// Site-wide niceties: "/" jumps to search, and code blocks get a copy button.
(() => {
  const searchPath = '/search/';

  document.addEventListener('keydown', (e) => {
    if (e.key !== '/' || e.metaKey || e.ctrlKey || e.altKey) return;
    const t = e.target;
    if (t.isContentEditable || /^(INPUT|TEXTAREA|SELECT)$/.test(t.tagName)) return;
    e.preventDefault();
    const q = document.getElementById('q');
    if (q) q.focus();
    else location.href = searchPath;
  });

  // Keep the active library sub-tab visible when the row scrolls sideways on phones.
  const activeSub = document.querySelector('.subtabs [aria-current="page"]');
  if (activeSub) {
    const row = activeSub.parentElement;
    row.scrollLeft = activeSub.offsetLeft - (row.clientWidth - activeSub.offsetWidth) / 2;
  }

  const copyIcon = '<svg viewBox="0 0 24 24" width="16" height="16" fill="none" stroke="currentColor" stroke-width="1.75" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><path d="M9 9h11v11H9zM5 15H4V4h11v1"/></svg>';

  document.querySelectorAll('.prose pre').forEach((pre) => {
    const wrap = pre.closest('.highlight') || pre;
    if (wrap.querySelector('.copy')) return;
    const btn = document.createElement('button');
    btn.type = 'button';
    btn.className = 'copy';
    btn.innerHTML = `${copyIcon}<span>Copy</span>`;
    btn.addEventListener('click', async () => {
      const label = btn.querySelector('span');
      try {
        await navigator.clipboard.writeText(pre.innerText.replace(/\n$/, ''));
        label.textContent = 'Copied';
      } catch {
        label.textContent = 'Press ⌘C';
        const range = document.createRange();
        range.selectNodeContents(pre);
        const sel = getSelection();
        sel.removeAllRanges();
        sel.addRange(range);
      }
      setTimeout(() => { label.textContent = 'Copy'; }, 1600);
    });
    wrap.classList.add('has-copy');
    wrap.appendChild(btn);
  });
})();
