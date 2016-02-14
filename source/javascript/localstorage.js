(function(root) {
  root.localStorageSetItem = function(key, value) {
    try {
      root.localStorage.setItem(key, value);
    } catch(e) {
    }
  };

  root.localStorageGetItem = function(key) {
    try {
      return root.localStorage.getItem(key);
    } catch(e) {
      return null;
    }
  };
})(window);
