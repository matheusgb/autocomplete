const socket = new WebSocket('ws://localhost:8080/ws');

socket.addEventListener('open', function (event) {
  console.log('WebSocket is open now.');
  socket.send(JSON.stringify({ message: "Page loaded" }))
});

socket.addEventListener('message', function (event) {
  const data = JSON.parse(event.data);
  console.log(data)
  const suggestions = document.getElementById('suggestions');
  const frequent = document.getElementById('frequent-list');
  suggestions.innerHTML = '';
  if (data.suggestions !== null) {
    data.suggestions.forEach(item => {
      const li = document.createElement('li');
      li.textContent = item;
      suggestions.appendChild(li);
    });
  }
  while (frequent.firstChild) {
    frequent.removeChild(frequent.firstChild);
  }
  data.frequent_values.forEach(item => {
    const li = document.createElement('li');
    li.textContent = `${item.key}:${item.doc_count}`;
    frequent.appendChild(li);
  });
});

document.getElementById('autocomplete').addEventListener('input', function () {
  const query = this.value;
  if (query.length > 0) {
    socket.send(query);
  } else {
    document.getElementById('suggestions').innerHTML = '';
  }
});

document.getElementById('send-button').addEventListener('click', function () {
  const word = document.getElementById('autocomplete').value;

  const data = JSON.stringify({ word: word });

  fetch('http://localhost:8080/send', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: data,
  })
    .then(response => response.json())
    .then(data => {
      console.log('Success:', data);
    })
    .catch((error) => {
      console.error('Error:', error);
    });
  socket.send("Word send")
});
