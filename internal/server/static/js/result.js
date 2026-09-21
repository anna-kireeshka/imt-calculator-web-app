// Экран результата: ИМТ, норма калорий и БЖУ приходят готовыми с сервера —
// GET /api/v1/measurements/bmr. На устройстве больше ничего не считается:
// вторая реализация формул неизбежно разъезжается с той, по которой живёт база.
(function () {
  'use strict';

  const bmiValue = document.getElementById('bmi-value');
  if (!bmiValue) return;

  const bmiZone = document.getElementById('bmi-zone');
  const bmiPin = document.getElementById('bmi-pin');
  const bmiBubble = document.getElementById('bmi-bubble');
  const errorText = document.getElementById('result-error');

  const out = {
    target: document.getElementById('cal-target'),
    protein: document.getElementById('cal-protein'),
    fat: document.getElementById('cal-fat'),
    carbs: document.getElementById('cal-carbs')
  };

  // Зоны ИМТ и их доли на шкале. from/to — границы значения, width — ширина
  // сегмента в разметке: доли взяты из макета (6,5 : 6,5 : 5 : 10) и шире
  // диапазонов не поровну, поэтому позиция считается внутри своего сегмента,
  // а не одной линейной формулой по всей шкале.
  const ZONES = [
    { from: 0,    to: 18.5, width: 6.5, name: 'Недостаток', key: 'low' },
    { from: 18.5, to: 25,   width: 6.5, name: 'Норма',      key: 'ok' },
    { from: 25,   to: 30,   width: 5,   name: 'Избыток',    key: 'warn' },
    { from: 30,   to: 40,   width: 10,  name: 'Ожирение',   key: 'bad' }
  ];

  const SCALE_WIDTH = ZONES.reduce(function (sum, zone) { return sum + zone.width; }, 0);

  load();

  function load() {
    const tg = window.fitcalc.tg;
    if (!tg || !tg.initData) {
      errorText.textContent = 'Откройте приложение через Telegram — показатели хранятся на сервере';
      return;
    }

    fetch('/api/v1/measurements/bmr', {
      headers: { 'Authorization': 'tma ' + tg.initData }
    }).then(function (response) {
      // Показывать нечего — значит человек ещё не дозаполнил анкету. Куда
      // именно его вести, решает сам онбординг: он спросит, заведён ли
      // пользователь, и откроет нужный шаг.
      if (response.status === 404) {
        if (toOnboarding()) return null;
        throw new Error('Пока нечего показывать — заполните анкету, по ней и считаются показатели');
      }
      if (!response.ok) {
        return response.text().then(function (text) {
          throw new Error(text.trim() || ('Не удалось получить расчёт, сервер ответил ' + response.status));
        });
      }
      return response.json();
    }, function () {
      throw new Error('Нет связи с сервером. Проверьте соединение и попробуйте ещё раз');
    }).then(function (data) {
      // null означает, что уже уходим на онбординг — рисовать нечего.
      if (data !== null) render(data);
    }).catch(function (error) {
      errorText.textContent = error.message;
    });
  }

  // Один переход на анкету за сессию: если бы /bmr продолжал отвечать 404 и
  // после её заполнения, экраны зациклились бы, перекидывая друг на друга.
  function toOnboarding() {
    try {
      if (sessionStorage.getItem('fitcalc.redirected') === '1') return false;
      sessionStorage.setItem('fitcalc.redirected', '1');
    } catch (e) {}

    window.location.href = '/onboarding';
    return true;
  }

  function render(data) {
    if (!data || typeof data.imt !== 'number') {
      throw new Error('Сервер вернул расчёт без показателей');
    }

    errorText.textContent = '';

    const zone = zoneOf(data.imt);
    const shown = window.fitcalc.decimal(data.imt, 1);

    bmiValue.textContent = shown;
    bmiZone.textContent = zone.name;
    bmiZone.setAttribute('data-zone', zone.key);

    bmiBubble.textContent = shown;
    bmiPin.style.left = position(data.imt) + '%';
    bmiPin.hidden = false;

    out.target.textContent = typeof data.ccal === 'number'
      ? window.fitcalc.number(data.ccal)
      : '—';

    const macros = data.macros || {};
    out.protein.textContent = grams(macros.protein);
    out.fat.textContent = grams(macros.fat);
    out.carbs.textContent = grams(macros.carbohydrates);
  }

  function zoneOf(bmi) {
    for (let i = 0; i < ZONES.length; i++) {
      if (bmi < ZONES[i].to) return ZONES[i];
    }
    return ZONES[ZONES.length - 1];
  }

  // Доля от ширины шкалы в процентах. Значения за краями прижимают указатель
  // к концу полосы, иначе он ушёл бы за неё.
  function position(bmi) {
    let offset = 0;

    for (let i = 0; i < ZONES.length; i++) {
      const zone = ZONES[i];

      if (bmi < zone.to || i === ZONES.length - 1) {
        const share = (bmi - zone.from) / (zone.to - zone.from);
        const clamped = Math.max(0, Math.min(1, share));
        return (offset + clamped * zone.width) / SCALE_WIDTH * 100;
      }

      offset += zone.width;
    }

    return 0;
  }

  // Сервер отдаёт нижнюю и верхнюю границы нормы: «90–120 г». Совпали —
  // печатаем одно число, а не «90–90 г».
  function grams(bounds) {
    if (!Array.isArray(bounds) || bounds.length !== 2) return '—';

    const low = window.fitcalc.number(bounds[0]);
    const high = window.fitcalc.number(bounds[1]);

    return low === high ? low : low + '–' + high;
  }
}());
