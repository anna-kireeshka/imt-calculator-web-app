// Экран результата: ИМТ, норма калорий и БЖУ приходят готовыми с сервера —
// GET /api/v1/measurements/bmr. На устройстве больше ничего не считается:
// вторая реализация формул неизбежно разъезжается с той, по которой живёт база.
(function () {
  'use strict';

  const bmiValue = document.getElementById('bmi-value');
  if (!bmiValue) return;

  const bmiZone = document.getElementById('bmi-zone');
  const bmiPin = document.getElementById('bmi-pin');
  const errorText = document.getElementById('result-error');

  const out = {
    target: document.getElementById('cal-target'),
    protein: document.getElementById('cal-protein'),
    fat: document.getElementById('cal-fat'),
    carbs: document.getElementById('cal-carbs')
  };

  // Шкала нарисована от 0 до 40: доли сегментов 18,5 : 6,5 : 5 : 10.
  // Пороги зон — оформление, а не расчёт, поэтому остаются на клиенте.
  const SCALE_MAX = 40;

  const ZONES = [
    { max: 18.5, name: 'Недостаток', key: 'low' },
    { max: 25, name: 'Норма', key: 'ok' },
    { max: 30, name: 'Избыток', key: 'warn' },
    { max: Infinity, name: 'Ожирение', key: 'bad' }
  ];

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

    bmiValue.textContent = window.fitcalc.decimal(data.imt, 1);
    bmiZone.textContent = zone.name;
    bmiZone.setAttribute('data-zone', zone.key);

    // Значения за краями шкалы прижимают маркер к её концу.
    bmiPin.style.left = Math.max(0, Math.min(100, data.imt / SCALE_MAX * 100)) + '%';
    bmiPin.hidden = false;

    out.target.textContent = typeof data.ccal === 'number'
      ? window.fitcalc.number(data.ccal) + ' ккал'
      : '—';

    const macros = data.macros || {};
    out.protein.textContent = grams(macros.protein);
    out.fat.textContent = grams(macros.fat);
    out.carbs.textContent = grams(macros.carbohydrates);
  }

  function zoneOf(bmi) {
    for (let i = 0; i < ZONES.length; i++) {
      if (bmi < ZONES[i].max) return ZONES[i];
    }
    return ZONES[ZONES.length - 1];
  }

  // Сервер отдаёт нижнюю и верхнюю границы нормы: «90–120 г». Совпали —
  // печатаем одно число, а не «90–90 г».
  function grams(bounds) {
    if (!Array.isArray(bounds) || bounds.length !== 2) return '—';

    const low = window.fitcalc.number(bounds[0]);
    const high = window.fitcalc.number(bounds[1]);

    return (low === high ? low : low + '–' + high) + ' г';
  }
}());
