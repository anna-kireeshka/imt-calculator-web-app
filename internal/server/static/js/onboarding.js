// Анкета: два шага в одном документе, переключение без перезагрузки.
(function () {
  'use strict';

  const form = document.getElementById('onboarding');
  if (!form) return;

  const steps = document.getElementById('ob-steps');
  const bars = steps.querySelectorAll('.steps__bar');
  const panels = form.querySelectorAll('[data-step]');

  const year = document.getElementById('ob-year');
  const month = document.getElementById('ob-month');
  const day = document.getElementById('ob-day');
  const genderGroup = document.getElementById('ob-gender');
  const height = document.getElementById('ob-height');
  const weight = document.getElementById('ob-weight');
  const goalGroup = document.getElementById('ob-goal');
  const activity = document.getElementById('ob-activity');

  const error1 = document.getElementById('ob-error-1');
  const error2 = document.getElementById('ob-error-2');

  let gender = null;
  let goal = pressedValue(goalGroup);

  genderGroup.addEventListener('segment:change', function (event) {
    gender = event.detail.value;
    error1.textContent = '';
  });

  goalGroup.addEventListener('segment:change', function (event) {
    goal = event.detail.value;
  });

  function pressedValue(group) {
    const pressed = group.querySelector('[aria-pressed="true"]');
    return pressed ? pressed.dataset.value : null;
  }

  function showStep(number) {
    panels.forEach(function (panel) {
      panel.hidden = Number(panel.dataset.step) !== number;
    });
    bars.forEach(function (bar, index) {
      if (index < number) {
        bar.setAttribute('data-done', '1');
      } else {
        bar.removeAttribute('data-done');
      }
    });
    steps.setAttribute('aria-valuenow', String(number));
    document.querySelector('.app-main').scrollTop = 0;
  }

  // Telegram рисует свою кнопку «назад» в шапке — на втором шаге она
  // должна возвращать на первый, а не закрывать мини-апп.
  const backButton = window.fitcalc.tg && window.fitcalc.tg.BackButton;

  function goToStep(number) {
    showStep(number);
    if (!backButton) return;
    try {
      if (number === 2) backButton.show(); else backButton.hide();
    } catch (e) {}
  }

  if (backButton) {
    try {
      backButton.onClick(function () { goToStep(1); });
    } catch (e) {}
  }

  const nextButton = document.getElementById('ob-next');
  const submitButton = document.getElementById('ob-submit');

  nextButton.addEventListener('click', function () {
    if (!year.value || !month.value || !day.value) {
      error1.textContent = 'Заполните год, месяц и день рождения';
      return;
    }

    const born = birthDate();
    if (born === null) {
      error1.textContent = 'Такой даты не существует — проверьте день и месяц';
      return;
    }

    const age = ageAt(born);
    if (age < 10 || age > 120) {
      error1.textContent = 'Проверьте дату рождения: возраст должен быть от 10 до 120 лет';
      return;
    }
    if (!gender) {
      error1.textContent = 'Выберите пол';
      return;
    }

    error1.textContent = '';
    window.fitcalc.haptic('light');

    // Пол и дата рождения уходят на сервер сразу: расчёт на экране результата
    // считается там же, и без записи пользователя ему нечего будет соединять
    // с замерами. Не сохранилось — на второй шаг не пускаем.
    send(nextButton, error1, api('/api/v1/user', {
      dob: isoDateTime(born),
      gender: gender === '1'
    })).then(function () {
      goToStep(2);
    }).catch(function () {});
  });

  document.getElementById('ob-back').addEventListener('click', function () {
    goToStep(1);
  });

  form.addEventListener('submit', function (event) {
    event.preventDefault();

    const heightValue = parseFloat(height.value);
    const weightValue = parseFloat(weight.value);

    if (!heightValue || !weightValue) {
      error2.textContent = 'Заполните рост и вес';
      return;
    }
    if (heightValue < 100 || heightValue > 250) {
      error2.textContent = 'Рост укажите в сантиметрах, от 100 до 250';
      return;
    }
    if (weightValue < 20 || weightValue > 400) {
      error2.textContent = 'Вес укажите в килограммах, от 20 до 400';
      return;
    }

    const born = birthDate();
    if (born === null) {
      error2.textContent = 'Проверьте дату рождения на первом шаге';
      return;
    }

    error2.textContent = '';

    const measurements = {
      height: Math.round(heightValue),
      weight: weightValue,
      goal: goal,
      activity_type: activity.value
    };

    // localStorage остался только для возврата в анкету из «Изменить данные»:
    // сами числа на экране результата приходят с сервера.
    try {
      localStorage.setItem('fitcalc.profile', JSON.stringify({
        dob: isoDate(born),
        gender: gender,
        height: measurements.height,
        weight: measurements.weight,
        goal: measurements.goal,
        activity_type: measurements.activity_type
      }));
    } catch (e) {}

    window.fitcalc.haptic('medium');

    send(submitButton, error2, api('/api/v1/measurements', measurements))
      .then(function () {
        window.location.href = '/';
      })
      .catch(function () {});
  });

  // Без подписи Telegram сервер ответит 401, поэтому её отсутствие — сразу
  // ошибка, а не молчаливый пропуск: иначе анкета выглядела бы сохранённой.
  function api(path, body) {
    const tg = window.fitcalc.tg;
    if (!tg || !tg.initData) {
      return Promise.reject(new Error('Откройте приложение через Telegram — иначе анкету не сохранить'));
    }

    return fetch(path, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': 'tma ' + tg.initData
      },
      body: JSON.stringify(body)
    }).then(function (response) {
      if (response.ok) return response;

      return response.text().then(function (text) {
        throw new Error(text.trim() || ('Не удалось сохранить, сервер ответил ' + response.status));
      });
    }, function () {
      throw new Error('Нет связи с сервером. Проверьте соединение и попробуйте ещё раз');
    });
  }

  // send держит кнопку заблокированной на время запроса и кладёт текст ошибки
  // рядом с шагом. Промис остаётся отклонённым — следующий шаг не наступит.
  function send(button, errorBox, request) {
    const label = button.textContent;
    button.disabled = true;
    button.textContent = 'Сохраняем…';

    return request
      .catch(function (error) {
        errorBox.textContent = error.message;
        throw error;
      })
      .finally(function () {
        button.disabled = false;
        button.textContent = label;
      });
  }

  // Go разбирает dob в time.Time, а это RFC3339: голая дата «1994-03-07»
  // до сервера не доедет.
  function isoDateTime(born) {
    return isoDate(born) + 'T00:00:00Z';
  }


  // birthDate возвращает null для несуществующей даты: Date сам нормализует
  // 31 февраля в 3 марта, поэтому сверяем, что поля пережили конструктор.
  function birthDate() {
    const y = parseInt(year.value, 10);
    const m = parseInt(month.value, 10);
    const d = parseInt(day.value, 10);
    if (!y || !m || !d) return null;

    const date = new Date(Date.UTC(y, m - 1, d));
    if (date.getUTCFullYear() !== y || date.getUTCMonth() !== m - 1 || date.getUTCDate() !== d) {
      return null;
    }

    return date;
  }

  // API и localStorage продолжают получать ISO-строку — три поля живут
  // только в разметке. Строка собирается из уже разобранной даты, а не из
  // сырых value: «007» в поле дня иначе дало бы «1994-03-007».
  function isoDate(born) {
    return [
      String(born.getUTCFullYear()).padStart(4, '0'),
      String(born.getUTCMonth() + 1).padStart(2, '0'),
      String(born.getUTCDate()).padStart(2, '0')
    ].join('-');
  }

  function ageAt(born) {
    const now = new Date();
    const monthDiff = now.getMonth() - born.getUTCMonth();

    let age = now.getFullYear() - born.getUTCFullYear();
    if (monthDiff < 0 || (monthDiff === 0 && now.getDate() < born.getUTCDate())) age--;

    return age;
  }

  function press(group, value) {
    if (!value) return;
    group.querySelectorAll('[aria-pressed]').forEach(function (button) {
      button.setAttribute('aria-pressed', button.dataset.value === value ? 'true' : 'false');
    });
  }

  // Возврат в анкету из «Изменить данные»: показываем, что было введено.
  try {
    const saved = JSON.parse(localStorage.getItem('fitcalc.profile') || 'null');
    if (saved) {
      const parts = (saved.dob || '').split('-');
      if (parts.length === 3) {
        year.value = parts[0];
        month.value = parts[1];
        day.value = String(Number(parts[2]));
      }
      height.value = saved.height || '';
      weight.value = saved.weight || '';
      activity.value = saved.activity_type || activity.value;
      press(genderGroup, saved.gender);
      press(goalGroup, saved.goal);
      gender = saved.gender || null;
      goal = saved.goal || goal;
    }
  } catch (e) {}

  [year, month, day].forEach(function (field) {
    field.addEventListener('input', function () {
      error1.textContent = '';
    });
  });

  goToStep(1);
}());
