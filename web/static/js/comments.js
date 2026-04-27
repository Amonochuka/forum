// Escape user content before inserting into innerHTML
function fEsc(str) {
  if (!str) return '';
  return String(str)
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;');
}

function fClearAuthPrompt(scope) {
  if (!scope) return;
  const existing = scope.querySelector('.f-auth-prompt');
  if (existing) existing.remove();
}

function fShowAuthPrompt(scope, message) {
  if (!scope) return;
  fClearAuthPrompt(scope);

  const prompt = document.createElement('div');
  prompt.className = 'f-auth-prompt';
  prompt.innerHTML = `
    <p>${fEsc(message)}</p>
    <div class="f-auth-actions">
      <a href="/login">Log in</a>
      <a href="/register">Register</a>
    </div>
  `;

  scope.appendChild(prompt);
}

// Build a comment element
function fRenderComment(c) {
  const div = document.createElement('div');
  div.className = 'f-comment';
  div.id = 'f-comment-' + c.id;
  div.dataset.id = c.id;
 
  const replyCountLabel = c.replyCount === 1 ? '1 reply' : c.replyCount + ' replies';
 
  div.innerHTML = `
    <div class="f-comment-header">
      <span class="f-author">${fEsc(c.authorName)}</span>
      <span class="f-time">${fEsc(c.createdAt)}</span>
    </div>
    <p class="f-body">${fEsc(c.body)}</p>
    <div class="f-actions">

      <form action="/react" method="POST" style="display:inline;">
        <input type="hidden" name="comment_id" value="${c.id}">
        <input type="hidden" name="type" value="like">
        <button type="submit">👍 ${c.likes}</button>
      </form>

      <form action="/react" method="POST" style="display:inline;">
        <input type="hidden" name="comment_id" value="${c.id}">
        <input type="hidden" name="type" value="dislike">
        <button type="submit">👎 ${c.dislikes}</button>
      </form>

      <button onclick="fToggleReply(${c.id})">Reply</button>
      ${c.replyCount > 0 ? `
        <button id="f-view-replies-${c.id}" onclick="fLoadReplies(${c.id}, this)">
          ${replyCountLabel}
        </button>
      ` : ''}
    </div>
    <div id="f-reply-form-${c.id}" style="display:none;" class="f-reply-form">
      <textarea id="f-reply-text-${c.id}" placeholder="Write a reply..." rows="2"></textarea>
      <div>
        <button onclick="fSubmitReply(${c.id}, 'f-reply-text-${c.id}', 'f-replies-${c.id}')">Reply</button>
        <button onclick="fToggleReply(${c.id})">Cancel</button>
      </div>
    </div>
    <div id="f-replies-${c.id}" class="f-replies"></div>
  `;
 
  return div;
}

function fRenderReply(r) {
  const div = document.createElement('div');
  div.className = 'f-comment';
  div.id = 'f-comment-' + r.id;
 
  div.innerHTML = `
    <div class="f-comment-header">
      <span class="f-author">${fEsc(r.authorName)}</span>
      ${r.badge ? `<span class="f-badge">${fEsc(r.badge)}</span>` : ''}
      <span class="f-time">${fEsc(r.createdAt)}</span>
    </div>
    <p class="f-body">${fEsc(r.body)}</p>
    <div class="f-actions">

    <form action="/react" method="POST" style="display:inline;">
      <input type="hidden" name="comment_id" value="${r.id}">
      <input type="hidden" name="type" value="like">
      <button type="submit">👍 ${r.likes}</button>
    </form>

    <form action="/react" method="POST" style="display:inline;">
      <input type="hidden" name="comment_id" value="${r.id}">
      <input type="hidden" name="type" value="dislike">
      <button type="submit">👎 ${r.dislikes}</button>
    </form>
    </div>
  `;
 
  return div;
}

// Submit a top-level comment
function fSubmitComment(postID, textareaID, listID) {
  const textarea = document.getElementById(textareaID);
  const list = document.getElementById(listID);
  if (!textarea || !list) return;
 
  const content = textarea.value.trim();
  if (!content) return;

  const btn = textarea.closest('.f-compose-body').querySelector('button');
  const composeBody = textarea.closest('.f-compose-body');
  fClearAuthPrompt(composeBody);
  btn.disabled = true;
  btn.textContent = 'Posting...';
 
  const form = new FormData();
  form.append('content', content);
 
  fetch('/posts/' + postID + '/comments', { method: 'POST', body: form })
    .then(function(r) {
      if (r.status === 401) {
        throw new Error('auth');
      }
      if (!r.ok) throw new Error('failed');
      return r.json();
    })
    .then(function(data) {
      list.prepend(fRenderComment(data.comment));
      textarea.value = '';
 
      /* bump the count in the heading */
      const heading = document.querySelector('#f-comments h3');
      if (heading) {
        heading.textContent = heading.textContent.replace(/\d+/, function(n) {
          return parseInt(n, 10) + 1;
        });
      }
    })
    .catch(function(err) {
      if (err.message === 'auth') {
        fShowAuthPrompt(composeBody, 'You need to log in to post a comment.');
        return;
      }
      alert('Could not post comment. Please try again.');
    })
    .finally(function() {
      btn.disabled = false;
      btn.textContent = 'Post Comment';
    });
}

// Load next page of comments
function fLoadMore(postID, btn) {
  const list = document.getElementById('f-comment-list');
  const page = parseInt(btn.dataset.page || '2', 10);
 
  btn.disabled = true;
  btn.textContent = 'Loading...';
 
  fetch('/posts/' + postID + '/comments?page=' + page)
    .then(function(r) { 
      if (!r.ok) throw new Error('failed');
      return r.json(); 
    })
    .then(function(data) {
      (data.comments || []).forEach(function(c) {
        list.appendChild(fRenderComment(c));
      });
      btn.dataset.page = page + 1;
 
      const loaded = list.querySelectorAll('.f-comment').length;
      if (loaded >= data.total) {
        btn.style.display = 'none';
      } else {
        btn.disabled = false;
        btn.textContent = 'Load more';
      }
    })
    .catch(function() {
      btn.disabled = false;
      btn.textContent = 'Failed to load — try again';
    });
}

// Show/hide the reply form under a comment
function fToggleReply(commentID) {
  const form = document.getElementById('f-reply-form-' + commentID);
  if (!form) return;
  const opening = form.style.display === 'none';
  form.style.display = opening ? 'block' : 'none';
  if (opening) {
    document.getElementById('f-reply-text-' + commentID)?.focus();
  }
}

// Load replies for a comment on first click, toggle on subsequent clicks
function fLoadReplies(commentID, btn) {
  const container = document.getElementById('f-replies-' + commentID);
  if (!container) return;
 
  if (container.dataset.loaded) {
    const isHidden = container.style.display === 'none';
    container.style.display = isHidden ? 'flex' : 'none';
    btn.textContent = isHidden ? '▲ hide replies' : btn.dataset.label;
    return;
  }
 
  btn.textContent = 'Loading...';
  btn.disabled = true;
 
  fetch('/comments/' + commentID + '/replies')
    .then(function(r) { 
      if (!r.ok) throw new Error('failed');
      return r.json(); 
    })
    .then(function(data) {
      const replies = data.replies || [];
      replies.forEach(function(r) {
        container.appendChild(fRenderReply(r));
      });
 
      const label = replies.length === 1 ? '1 reply' : replies.length + ' replies';
      container.dataset.loaded = 'true';
      btn.dataset.label = label;
      container.style.display  = 'flex';
 
      btn.textContent = '▲ hide replies';
      btn.disabled = false;
    })
    .catch(function() {
      btn.textContent = 'Could not load replies';
      btn.disabled = false;
    });
}

// submit a reply to a comment
function fSubmitReply(parentID, textareaID, repliesContainerID) {
  const textarea  = document.getElementById(textareaID);
  const container = document.getElementById(repliesContainerID);
  if (!textarea || !container) return;
 
  const content = textarea.value.trim();
  if (!content) return;
 
  const postID = parseInt(
    document.getElementById('f-comments')?.dataset.postId || '0', 10
  );

  const replyForm = textarea.closest('.f-reply-form');
  const btn = replyForm.querySelector('button');
  fClearAuthPrompt(replyForm);
  btn.disabled = true;
  btn.textContent = 'Posting...';
 
  const form = new FormData();
  form.append('content', content);
  form.append('post_id', postID);
 
  fetch('/comments/' + parentID + '/replies', { method: 'POST', body: form })
    .then(function(r) {
      if (r.status === 401) {
        throw new Error('auth');
      }
      if (!r.ok) throw new Error('failed');
      return r.json();
    })
    .then(function(data) {
      container.style.display = 'flex';
      container.dataset.loaded = 'true';
      container.appendChild(fRenderReply(data.comment));
 
      // update or create the view-replies button
      const viewBtn = document.getElementById('f-view-replies-' + parentID);
      if (viewBtn) {
        const prev  = parseInt(viewBtn.dataset.label || '0', 10);
        const count = prev + 1;
        const label = count === 1 ? '1 reply' : count + ' replies';
        viewBtn.dataset.label = label;
        viewBtn.textContent   = '▲ hide replies';
      } else {
        const actions = document.querySelector('#f-comment-' + parentID + ' .f-actions');
        if (actions) {
          const newBtn = document.createElement('button');
          newBtn.id = 'f-view-replies-' + parentID;
          newBtn.dataset.label = '1 reply';
          newBtn.textContent   = '▲ hide replies';
          newBtn.onclick = function() { fLoadReplies(parentID, newBtn); };
          actions.appendChild(newBtn);
        }
      }
 
      textarea.value = '';
      fToggleReply(parentID);
    })
    .catch(function(err) {
      if (err.message === 'auth') {
        fShowAuthPrompt(replyForm, 'You need to log in to post a reply.');
        return;
      }
      alert('Could not post reply. Please try again.');
    })
    .finally(function() {
      btn.disabled = false;
      btn.textContent = 'Reply';
    });
}
