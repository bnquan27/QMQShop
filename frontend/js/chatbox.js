(function () {
  'use strict';

  var STORAGE_KEY = 'qmq_chat_history';
  var MAX_HISTORY = 20;
  var chatHistory = [];
  var isOpen = false;

  function escapeHtml(text) {
    var d = document.createElement('div');
    d.textContent = text;
    return d.innerHTML;
  }

  function applyInline(text) {
    text = text.replace(/`([^`]+)`/g, '<code>$1</code>');
    text = text.replace(/\*\*([^*]+)\*\*/g, '<strong>$1</strong>');
    text = text.replace(/\*([^*]+)\*/g, '<em>$1</em>');
    text = text.replace(/\[([^\]]+)\]\(([^)]+)\)/g, '<a href="$2" target="_blank" rel="noopener">$1</a>');
    return text;
  }

  function renderMarkdown(text) {
    var safe = escapeHtml(text);
    var lines = safe.split('\n');
    var out = '';
    var inUl = false;
    var inOl = false;

    function closeList() {
      if (inUl) { out += '</ul>'; inUl = false; }
      if (inOl) { out += '</ol>'; inOl = false; }
    }

    for (var i = 0; i < lines.length; i++) {
      var line = lines[i];
      var t = line.trim();

      if (/^---+\s*$/.test(t)) { closeList(); out += '<hr>'; continue; }

      var h = t.match(/^(#{1,3})\s+(.+)$/);
      if (h) { closeList(); out += '<h' + h[1].length + '>' + applyInline(h[2]) + '</h' + h[1].length + '>'; continue; }

      var ul = t.match(/^[*\-+]\s+(.+)$/);
      if (ul) {
        if (inOl) { out += '</ol>'; inOl = false; }
        if (!inUl) { out += '<ul>'; inUl = true; }
        out += '<li>' + applyInline(ul[1]) + '</li>';
        continue;
      }

      var ol = t.match(/^\d+[.)]\s+(.+)$/);
      if (ol) {
        if (inUl) { out += '</ul>'; inUl = false; }
        if (!inOl) { out += '<ol>'; inOl = true; }
        out += '<li>' + applyInline(ol[1]) + '</li>';
        continue;
      }

      closeList();

      if (t === '') { out += '<br>'; continue; }

      out += '<p>' + applyInline(line) + '</p>';
    }

    closeList();
    return out;
  }

  function scrollBottom() {
    var area = document.getElementById('chat-msg-area');
    requestAnimationFrame(function () { area.scrollTop = area.scrollHeight; });
  }

  function getToken() {
    try {
      if (typeof Auth !== 'undefined' && Auth.getToken) return Auth.getToken();
      return localStorage.getItem('token');
    } catch (e) { return null; }
  }

  function saveHistory() {
    try {
      var trimmed = chatHistory.slice(-MAX_HISTORY);
      localStorage.setItem(STORAGE_KEY, JSON.stringify(trimmed));
    } catch (e) {}
  }

  function renderHistory() {
    var area = document.getElementById('chat-msg-area');
    var welcome = area.querySelector('.chat-welcome');
    if (chatHistory.length > 0) {
      if (welcome) welcome.style.display = 'none';
    } else {
      if (welcome) welcome.style.display = '';
    }

    var existing = area.querySelectorAll('.chat-msg');
    for (var i = 0; i < existing.length; i++) existing[i].remove();

    for (var i = 0; i < chatHistory.length; i++) {
      var m = chatHistory[i];
      var div = document.createElement('div');
      div.className = 'chat-msg ' + m.role;
      div.innerHTML = '<div class="chat-bubble chat-bubble-markdown">' + renderMarkdown(m.text) + '</div>';
      area.appendChild(div);
    }
    scrollBottom();
  }

  function addMessage(role, text) {
    chatHistory.push({ role: role, text: text });
    saveHistory();

    var area = document.getElementById('chat-msg-area');
    var typing = area.querySelector('.typing');
    if (typing) typing.remove();

    var div = document.createElement('div');
    div.className = 'chat-msg ' + role;
    div.innerHTML = '<div class="chat-bubble chat-bubble-markdown">' + renderMarkdown(text) + '</div>';
    area.appendChild(div);
    scrollBottom();
  }

  function toggle() {
    isOpen = !isOpen;
    var panel = document.getElementById('chat-panel');
    var fab = document.getElementById('chat-fab');
    panel.classList.toggle('hidden', !isOpen);
    fab.classList.toggle('active', isOpen);
    if (isOpen) {
      setTimeout(function () {
        document.getElementById('chat-input').focus();
      }, 200);
    }
  }

  function send() {
    var input = document.getElementById('chat-input');
    var text = input.value.trim();
    if (!text) return;

    input.value = '';
    document.getElementById('chat-send').disabled = true;

    addMessage('user', text);

    var area = document.getElementById('chat-msg-area');
    var typingDiv = document.createElement('div');
    typingDiv.className = 'chat-msg assistant typing';
    typingDiv.innerHTML =
      '<div class="chat-bubble typing-dots"><span></span><span></span><span></span></div>';
    area.appendChild(typingDiv);
    scrollBottom();

    var token = getToken();
    var headers = { 'Content-Type': 'application/json' };
    if (token) headers['Authorization'] = 'Bearer ' + token;

    fetch('/api/chat', {
      method: 'POST',
      headers: headers,
      body: JSON.stringify({ message: text, history: chatHistory })
    })
      .then(function (res) { return res.json(); })
      .then(function (data) {
        typingDiv.remove();
        if (data.reply) {
          addMessage('assistant', data.reply);
        } else {
          addMessage('assistant', 'Xin lỗi, tôi không thể trả lời ngay lúc này. Vui lòng thử lại sau.');
        }
      })
      .catch(function () {
        typingDiv.remove();
        addMessage('assistant', 'Xin lỗi, đã xảy ra lỗi kết nối. Vui lòng thử lại sau.');
      })
      .finally(function () {
        document.getElementById('chat-send').disabled = false;
        document.getElementById('chat-input').focus();
      });
  }

  function handleInput() {
    var input = document.getElementById('chat-input');
    var sendBtn = document.getElementById('chat-send');
    sendBtn.disabled = !input.value.trim();
  }

  function handleKey(e) {
    if (e.key === 'Enter') send();
  }

  function createDOM() {
    var c = document.createElement('div');
    c.id = 'qmq-chatbox';
    c.innerHTML =
      '<button id="chat-fab" onclick="QChat.toggle()" aria-label="Chat v\u1EDBi AI">' +
        '<i class="fa-solid fa-comment-dots"></i>' +
      '</button>' +
      '<div id="chat-panel" class="chat-panel hidden">' +
        '<div class="chat-header">' +
          '<span class="chat-title"><i class="fa-solid fa-robot"></i> QMQ AI</span>' +
          '<button class="chat-close" onclick="QChat.toggle()"><i class="fa-solid fa-xmark"></i></button>' +
        '</div>' +
        '<div class="chat-msg-area" id="chat-msg-area">' +
          '<div class="chat-welcome">' +
            '<div class="chat-welcome-icon"><i class="fa-solid fa-robot"></i></div>' +
            '<p>Xin ch\u00E0o! T\u00F4i l\u00E0 tr\u1EE3 l\u00FD AI c\u1EE7a <strong>QMQ Shop</strong>.</p>' +
            '<p>T\u00F4i c\u00F3 th\u1EC3 gi\u00FAp b\u1EA1n:</p>' +
            '<ul>' +
              '<li>\uD83D\uDD0D T\u01B0 v\u1EA5n ch\u1ECDn m\u00E1y t\u00EDnh / linh ki\u1EC7n</li>' +
              '<li>\uD83D\uDCB0 G\u1EE3i \u00FD s\u1EA3n ph\u1EA9m ph\u00F9 h\u1EE3p ng\u00E2n s\u00E1ch</li>' +
              '<li>\uD83C\uDFF7\uFE0F Gi\u1EDBi thi\u1EC7u s\u1EA3n ph\u1EA9m \u0111ang sale</li>' +
              '<li>\uD83D\uDCDE Gi\u1EA3i \u0111\u00E1p th\u1EAFc m\u1EAFc v\u1EC1 shop</li>' +
            '</ul>' +
            '<p class="chat-welcome-hint"><em>H\u00E3y \u0111\u1EB7t c\u00E2u h\u1ECFi cho t\u00F4i nh\u00E9!</em></p>' +
          '</div>' +
        '</div>' +
        '<div class="chat-input-row">' +
          '<input type="text" id="chat-input" placeholder="Nh\u1EADp tin nh\u1EAFn..." autocomplete="off">' +
          '<button id="chat-send" disabled><i class="fa-solid fa-paper-plane"></i></button>' +
        '</div>' +
      '</div>';
    document.body.appendChild(c);

    document.getElementById('chat-input').addEventListener('input', handleInput);
    document.getElementById('chat-input').addEventListener('keydown', handleKey);
    document.getElementById('chat-send').addEventListener('click', send);
    document.addEventListener('keydown', function (e) {
      if (e.key === 'Escape' && isOpen) toggle();
    });
  }

  function loadHistory() {
    try {
      var saved = localStorage.getItem(STORAGE_KEY);
      if (saved) {
        var parsed = JSON.parse(saved);
        if (Array.isArray(parsed) && parsed.length > 0) {
          chatHistory = parsed;
          renderHistory();
        }
      }
    } catch (e) {
      chatHistory = [];
    }
  }

  function init() {
    createDOM();
    loadHistory();
  }

  window.QChat = { toggle: toggle, send: send };

  if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', init);
  } else {
    init();
  }
})();
