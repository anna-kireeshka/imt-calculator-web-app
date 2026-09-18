// Дневник: добавление приёма пищи и пересчёт остатка от нормы.
(function () {
  'use strict';

  var form = document.getElementById('diary-form');
  if (!form) return;

  var nameInput = document.getElementById('diary-name');
  var kcalInput = document.getElementById('diary-kcal');
  var errorText = document.getElementById('diary-error');
  var list = document.getElementById('diary-list');
  var totalOut = document.getElementById('diary-total');
  var leftOut = document.getElementById('diary-left');
  var bar = document.getElementById('diary-bar');

  var NORM = 2100;
  var total = 1640;

  form.addEventListener('submit', function (event) {
    event.preventDefault();

    var name = nameInput.value.trim();
    var kcal = parseInt(kcalInput.value, 10);

    if (!name || !kcal) {
      errorText.textContent = 'Укажите название и калории';
      return;
    }

    errorText.textContent = '';

    var row = document.createElement('div');
    row.className = 'row';

    var title = document.createElement('span');
    title.textContent = name;

    var amount = document.createElement('span');
    amount.textContent = kcal;

    row.append(title, amount);
    list.append(row);

    total += kcal;
    var over = total > NORM;

    totalOut.textContent = window.fitcalc.number(total);
    bar.style.width = Math.min(100, Math.round(total / NORM * 100)) + '%';
    bar.setAttribute('data-over', over ? '1' : '0');
    leftOut.textContent = over
      ? 'Превышение на ' + window.fitcalc.number(total - NORM) + ' ккал'
      : 'Осталось ' + window.fitcalc.number(NORM - total) + ' ккал';

    nameInput.value = '';
    kcalInput.value = '';
    window.fitcalc.haptic('light');
  });

  [nameInput, kcalInput].forEach(function (input) {
    input.addEventListener('input', function () {
      errorText.textContent = '';
    });
  });
}());
