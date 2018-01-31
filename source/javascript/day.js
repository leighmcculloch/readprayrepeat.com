(function() {
  var AUDIO_PLAY = '🔊';
  var AUDIO_STOP = '🔇';

  var $textControlVerses;
  var $textControlBibles;
  var $metaDay;

  var voice;

  function init() {
    $textControlAudios = document.querySelectorAll('.text-controls-audio');
    $textControlVerses = document.querySelectorAll('.text-controls-verses');
    $metaDay = document.querySelector('meta[name="readprayrepeat:day"]');

    if ('speechSynthesis' in window) {
      var voiceInit = function() {
        var voices = window.speechSynthesis.getVoices();
        var preferredVoices = voices.filter(function(voice) {
          return voice.localService && voice.lang == 'en-GB';
        });
        voice = preferredVoices[0] || allVoices[0];
        for (var i = 0; i < $textControlAudios.length; ++i) {
          $textControlAudios[i].style.display = 'block';
          $textControlAudios[i].innerText = AUDIO_PLAY;
          $textControlAudios[i].addEventListener("click", onTextControlAudiosClick);
        }
      };
      window.speechSynthesis.cancel();
      window.speechSynthesis.onvoiceschanged = voiceInit;
      var voices = window.speechSynthesis.getVoices();
      if (voices.length > 0) {
        voiceInit();
      }
    }

    for (var i = 0; i < $textControlVerses.length; ++i) {
      $textControlVerses[i].addEventListener("click", onTextControlVersesClick);
    }

    for (var i = 0; i < $textControlBibles.length; ++i) {
      $textControlBibles[i].addEventListener("change", onTextControlBiblesChange);
    }

    renderCurrentDate();
  }

  function onTextControlAudiosClick(e) {
    e.preventDefault();
    var $target = e.target || e.srcElement;
    toggleAudio($target)
    return false;
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

  function toggleAudio($domElement) {
    if ($domElement.innerText == AUDIO_PLAY) {
      var text = $domElement.nextElementSibling.innerText;
      var utterance = new SpeechSynthesisUtterance(text);
      utterance.voice = voice;
      console.log('Voice: ' + utterance.voice.name + ' Lang: ' + utterance.voice.lang);
      utterance.onend = function() {
        $domElement.innerText = AUDIO_PLAY;
      };
      window.speechSynthesis.speak(utterance);
      for (var i = 0; i < $textControlAudios.length; ++i) {
        $textControlAudios[i].innerText = AUDIO_PLAY;
      }
      $domElement.innerText = AUDIO_STOP;
    } else if ($domElement.innerText == AUDIO_STOP) {
      window.speechSynthesis.cancel();
      $domElement.innerText = AUDIO_PLAY;
    }
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
