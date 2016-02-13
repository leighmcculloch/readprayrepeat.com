(function() {
  var $textControlVerses;
  var $textControlBibles;
  var $metaDay;

  function init() {
    $textControlVerses = document.querySelectorAll('.text-controls-verses');
    $textControlBibles = document.querySelectorAll('.text-controls-bibles');
    $metaDay = document.querySelector('meta[name="readprayrepeat:day"]');

    for (var i = 0; i < $textControlVerses.length; ++i) {
      $textControlVerses[i].addEventListener("click", onTextControlVersesClick);
    }

    for (var i = 0; i < $textControlBibles.length; ++i) {
      $textControlBibles[i].addEventListener("change", onTextControlBiblesChange);
    }

    renderCurrentDate();
  }

  function onTextControlVersesClick(e) {
    e.preventDefault();
    toggleVerses();
    return false;
  }

  function onTextControlBiblesChange(e) {
    e.preventDefault();
    changeToSelectedBible(this);
    return false;
  }

  function toggleVerses() {
    var $text = document.querySelectorAll('.text');
    for (var i = 0; i < $text.length; ++i) {
      if ($text[i].className.indexOf('show-verses') == -1) {
        $text[i].className += ' show-verses';
      } else {
        $text[i].className = $text[i].className.replace('show-verses', '');
      }
    }
  }

  function changeToSelectedBible($select) {
    var selectedBiblePath = $select.options[$select.selectedIndex].value;
    window.location.href = selectedBiblePath;
  }

  function renderCurrentDate() {
    var day = $metaDay.getAttribute("content");
    var date = getDateForDay(day);
    var $navCurrentDate = document.querySelectorAll('.nav-current-date');
    for (var i = 0; i < $navCurrentDate.length; ++i) {
      $navCurrentDate[i].innerText = date.toDateString();
    }
  }

  function getDateForDay(day) {
    var start = getStartDate();
    var oneDay = 1000 * 60 * 60 * 24;
    var daysInTime = oneDay * day;
    return new Date(start.getTime() + daysInTime);
  }

  function getStartDate() {
    var now = new Date();
    return new Date(now.getFullYear(), 0, 0);
  }

  window.addEventListener('load', init);
})();
