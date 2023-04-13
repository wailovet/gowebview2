window.gowebview2 = true;

window.chrome.storage = {
  sync: {
    get: function (keys, callback) {
      var result = {};
      var count = 0;
      for (var i = 0; i < keys.length; i++) {
        storageGetItem(keys[i]).then(
          (function (value) {
            return function (data) {
              result[value] = data;
              count++;
              if (count === keys.length) {
                callback(result);
              }
            };
          })(keys[i])
        );
      }
    },
    set: function (result, callback) {
      var count = 0;

      result = JSON.parse(JSON.stringify(result));

      for (var key in result) {
        storageSetItem(key, JSON.stringify(result[key])).then(() => {
          count++;
          if (count === result.length) {
            callback();
          }
        });
      }
    },
    remove: function (key, callback) {
      storageRemoveItem(key);
      setTimeout(() => {
        callback();
      }, 500);
    },
    clear: function (callback) {
      storageClear();
      setTimeout(() => {
        callback();
      }, 500);
    },
  },
};
