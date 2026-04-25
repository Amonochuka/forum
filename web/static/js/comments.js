// Submit a top-level comment
function fSubmitComment(postID, textareaID, listID) {
  const textarea = document.getElementById(textareaID);
  const list = document.getElementById(listID);
  if (!textarea || !list) return;
 
  const content = textarea.value.trim();
  if (!content) return;
 
  const btn = textarea.closest('.f-compose-body').querySelector('button');
  btn.disabled = true;
  btn.textContent = 'Posting...';
 
  const form = new FormData();
  form.append('content', content);
 
  fetch('/posts/' + postID + '/comments', { method: 'POST', body: form })
    .then(function(r) {
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
    .catch(function() {
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
    .then(function(r) { return r.json(); })
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
      btn.textContent = 'Load more';
    });
}

