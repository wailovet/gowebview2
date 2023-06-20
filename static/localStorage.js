window.localStorage.getItem = function (key) {
  //sync ajax
  var xhr = new XMLHttpRequest();
  key = encodeURIComponent(key);
  xhr.open("GET", "http://127.0.0.1:%d/localStorage/getItem?key=" + key, false);
  xhr.send();
  var data = JSON.parse(xhr.responseText);
  return data?.data;
};

window.localStorage.setItem = function (key, value) {
  //sync ajax
  var xhr = new XMLHttpRequest();
  key = encodeURIComponent(key);
  value = encodeURIComponent(value);
  xhr.open(
    "GET",
    "http://127.0.0.1:%d/localStorage/setItem?key=" + key + "&value=" + value,
    false
  );
  xhr.send();
};

window.localStorage.removeItem = function (key) {
  //sync ajax
  var xhr = new XMLHttpRequest();
  key = encodeURIComponent(key);
  xhr.open(
    "GET",
    "http://127.0.0.1:%d/localStorage/removeItem?key=" + key,
    false
  );
  xhr.send();
};
