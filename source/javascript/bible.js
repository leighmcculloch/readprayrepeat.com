(function(root) {
  root.setBible = function(bible) {
    root.localStorageSetItem('bible-selected', bible);
  };

  root.getBible = function() {
    return root.localStorageGetItem('bible-selected');
  };
})(window);
