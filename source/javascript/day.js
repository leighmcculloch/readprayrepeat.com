(function() {
  var $textControlVerses;

  function init() {
    $textControlVerses = document.querySelectorAll('.text-controls-verses');

    for (var i = 0; i < $textControlVerses.length; ++i) {
      $textControlVerses[i].addEventListener("click", onToggleVerses);
    }
  }

  function onToggleVerses() {
    var $text = document.querySelectorAll('.text');
    for (var i = 0; i < $text.length; ++i) {
      if ($text[i].className.indexOf('show-verses') == -1) {
        $text[i].className += ' show-verses';
      } else {
        $text[i].className = $text[i].className.replace('show-verses', '');
      }
    }
  }

  window.addEventListener('load', init);
})();
