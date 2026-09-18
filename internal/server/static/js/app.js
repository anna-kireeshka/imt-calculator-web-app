// Общее для всех экранов: инициализация Telegram, переключатели, повтор запроса.
(function () {
  'use strict';

  var webApp = window.Telegram && window.Telegram.WebApp;
  var inTelegram = !!(webApp && ((webApp.initData && webApp.initData.length > 0) ||
    (webApp.platform && webApp.platform !== 'unknown')));

  window.fitcalc = {
    tg: inTelegram ? webApp : null,

    // Отклик на смену зоны или успешное действие — только внутри Telegram.
    haptic: function (style) {
      var tg = window.fitcalc.tg;
      if (!tg || !tg.HapticFeedback) return;
      try {
        tg.HapticFeedback.impactOccurred(style || 'soft');
      } catch (e) {}
    },

    number: function (value) {
      return Math.round(value).toLocaleString('ru-RU');
    },

    decimal: function (value, digits) {
      return value.toFixed(digits).replace('.', ',');
    }
  };

  if (inTelegram) {
    try {
      webApp.ready();
      webApp.expand();
    } catch (e) {}
    document.documentElement.setAttribute('data-tg', '1');
  }

  // Сегментированные переключатели: один активный на группу.
  document.querySelectorAll('.segmented').forEach(function (group) {
    var buttons = group.querySelectorAll('button[aria-pressed]');
    if (!buttons.length) return;

    buttons.forEach(function (button) {
      button.addEventListener('click', function () {
        buttons.forEach(function (other) {
          other.setAttribute('aria-pressed', other === button ? 'true' : 'false');
        });
        group.dispatchEvent(new CustomEvent('segment:change', {
          detail: { value: button.dataset.value || button.textContent.trim() }
        }));
      });
    });
  });

  document.querySelectorAll('.switch[role="switch"]').forEach(function (control) {
    control.addEventListener('click', function () {
      var on = control.getAttribute('aria-checked') === 'true';
      control.setAttribute('aria-checked', on ? 'false' : 'true');
      window.fitcalc.haptic('light');
    });
  });

  document.querySelectorAll('[data-retry]').forEach(function (button) {
    button.addEventListener('click', function () {
      var label = button.innerHTML;
      button.textContent = 'Подключаемся…';
      button.disabled = true;
      setTimeout(function () {
        button.innerHTML = label;
        button.disabled = false;
        window.location.reload();
      }, 900);
    });
  });
}());
