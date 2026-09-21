// Экран истории: всё на экране строится из GET /api/v1/measurements/history.
(function () {
  'use strict';

  const button = document.getElementById('weight-button');
  const form = document.getElementById('weight-edit');
  if (!button || !form) return;

  const value = document.getElementById('weight-value');
  const input = document.getElementById('weight-input');
  const errorText = document.getElementById('weight-error');
  const cancel = document.getElementById('weight-cancel');

  const loadError = document.getElementById('history-error');
  const empty = document.getElementById('history-empty');
  const summary = document.getElementById('history-summary');
  const chartCard = document.getElementById('history-chart-card');
  const list = document.getElementById('history-list');

  const MIN = 20;
  const MAX = 400;

  // Поле viewBox графика и отступы, внутри которых живёт ломаная.
  const PLOT = { left: 10, right: 250, top: 18, bottom: 88 };

  load();

  // ── данные ────────────────────────────────────────────────────

  function load() {
    const tg = window.fitcalc.tg;
    if (!tg || !tg.initData) {
      loadError.textContent = 'Откройте приложение через Telegram — история хранится на сервере';
      return;
    }

    fetch('/api/v1/measurements/history', {
      headers: { 'Authorization': 'tma ' + tg.initData }
    }).then(function (response) {
      if (!response.ok) {
        return response.text().then(function (text) {
          throw new Error(text.trim() || ('Не удалось загрузить историю, сервер ответил ' + response.status));
        });
      }
      return response.json();
    }, function () {
      throw new Error('Нет связи с сервером. Проверьте соединение и попробуйте ещё раз');
    }).then(render).catch(function (error) {
      loadError.textContent = error.message;
    });
  }

  function render(data) {
    // Пустая выборка приходит как null, а не как [].
    // Отрисовка может идти после неудачной попытки — снимаем прежнее сообщение.
    loadError.textContent = '';

    const entries = Array.isArray(data) ? data.slice() : [];
    if (entries.length === 0) {
      empty.hidden = false;
      return;
    }

    // Сервер отдаёт от новых к старым; графику и списку удобнее наоборот.
    entries.reverse();

    const latest = entries[entries.length - 1];
    const oldest = entries[0];

    value.textContent = window.fitcalc.decimal(latest.weight, 1);

    fillSummary(latest, oldest, entries.length);
    fillList(entries);
    drawChart(entries);
  }

  function fillSummary(latest, oldest, count) {
    document.getElementById('summary-weight').textContent = window.fitcalc.decimal(latest.weight, 1);
    document.getElementById('summary-imt').textContent = window.fitcalc.decimal(latest.imt, 1);

    // По одной записи сравнивать не с чем — дельту не показываем.
    if (count > 1) {
      fillDelta(document.getElementById('summary-weight-delta'), latest.weight - oldest.weight, ' кг');
      fillDelta(document.getElementById('summary-imt-delta'), latest.imt - oldest.imt, '');
    }

    summary.hidden = false;
  }

  // Знак рисуем сами: минус U+2212 вместо дефиса, как в макете.
  function fillDelta(node, diff, unit) {
    const rounded = Math.round(diff * 10) / 10;
    if (rounded === 0) {
      node.textContent = '';
      return;
    }

    node.className = rounded < 0 ? 'delta--down' : 'delta--up';
    node.textContent = (rounded < 0 ? '−' : '+') + window.fitcalc.decimal(Math.abs(rounded), 1) + unit;
  }

  function fillList(entries) {
    list.textContent = '';

    // Сверху последние взвешивания.
    entries.slice().reverse().forEach(function (entry) {
      const row = document.createElement('div');
      row.className = 'row';

      const date = document.createElement('span');
      date.textContent = shortDate(entry.created_at);

      const numbers = document.createElement('span');
      numbers.textContent = window.fitcalc.decimal(entry.weight, 1) + ' кг ';

      const imt = document.createElement('span');
      imt.className = 'muted';
      imt.textContent = window.fitcalc.decimal(entry.imt, 1);

      numbers.appendChild(imt);
      row.appendChild(date);
      row.appendChild(numbers);
      list.appendChild(row);
    });

    list.hidden = false;
  }

  // ── график ────────────────────────────────────────────────────

  // По оси x — номер записи, а не дата: взвешиваются неравномерно, и на
  // временной оси несколько подряд идущих дней сжались бы в одну точку.
  function drawChart(entries) {
    const weights = entries.map(function (entry) { return entry.weight; });
    const min = Math.min.apply(null, weights);
    const max = Math.max.apply(null, weights);

    const points = entries.map(function (entry, index) {
      return x(index, entries.length) + ',' + y(entry.weight, min, max);
    });

    document.getElementById('chart-line').setAttribute('points', points.join(' '));

    const last = points[points.length - 1].split(',');
    const dot = document.getElementById('chart-dot');
    dot.setAttribute('cx', last[0]);
    dot.setAttribute('cy', last[1]);

    document.getElementById('chart-from').textContent = shortDate(entries[0].created_at);
    document.getElementById('chart-to').textContent = shortDate(entries[entries.length - 1].created_at);

    document.getElementById('history-chart').setAttribute('aria-label',
      'Вес с ' + window.fitcalc.decimal(entries[0].weight, 1) + ' кг ' + shortDate(entries[0].created_at) +
      ' до ' + window.fitcalc.decimal(entries[entries.length - 1].weight, 1) + ' кг ' +
      shortDate(entries[entries.length - 1].created_at));

    chartCard.hidden = false;
  }

  function x(index, count) {
    if (count === 1) return (PLOT.left + PLOT.right) / 2;
    return PLOT.left + index / (count - 1) * (PLOT.right - PLOT.left);
  }

  // Ровная линия при одинаковом весе легла бы на верхний край, поэтому
  // вырожденный диапазон рисуем по середине поля.
  function y(weight, min, max) {
    if (max === min) return (PLOT.top + PLOT.bottom) / 2;
    return PLOT.bottom - (weight - min) / (max - min) * (PLOT.bottom - PLOT.top);
  }

  function shortDate(value) {
    const date = new Date(value);
    if (isNaN(date.getTime())) return '';

    // Короткий месяц в русской локали идёт с точкой — «21 сент.». Точку
    // оставляем: без неё получается не сокращение, а обрубок.
    return date.toLocaleDateString('ru-RU', { day: 'numeric', month: 'short' });
  }

  // ── ввод веса ─────────────────────────────────────────────────

  button.addEventListener('click', open);
  cancel.addEventListener('click', close);

  input.addEventListener('input', function () {
    errorText.textContent = '';
  });

  form.addEventListener('submit', function (event) {
    event.preventDefault();
    save();
  });

  form.addEventListener('keydown', function (event) {
    if (event.key === 'Escape') close();
  });

  function open() {
    input.value = value.textContent.trim().replace(',', '.');
    errorText.textContent = '';

    button.hidden = true;
    button.setAttribute('aria-expanded', 'true');
    form.hidden = false;

    input.focus();
    input.select();
  }

  function close() {
    form.hidden = true;
    button.hidden = false;
    button.setAttribute('aria-expanded', 'false');
    button.focus();
  }

  // Новый вес пока никуда не отправляется — по отдельному решению бэкенд
  // здесь не трогаем. Когда дойдёт очередь: POST /api/v1/measurements
  // с ростом, целью и активностью из последнего замера, а затем load().
  function save() {
    const weight = parseFloat(input.value);

    if (!weight) {
      errorText.textContent = 'Введите вес';
      return;
    }
    if (weight < MIN || weight > MAX) {
      errorText.textContent = 'Вес укажите в килограммах, от ' + MIN + ' до ' + MAX;
      return;
    }

    value.textContent = window.fitcalc.decimal(Math.round(weight * 10) / 10, 1);

    window.fitcalc.haptic('medium');
    close();
  }
}());
