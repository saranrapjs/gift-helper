function doBuy(target) {
	var url = target.getAttribute('data-href'),
		text = target.getAttribute('data-text');

	document.querySelector('#purchase').textContent = text;

	document.querySelector('#url').value = url;

	document.querySelector('form').classList.add('open');

	document.querySelector('form a').href = url;

	document.querySelector('.url-preview').textContent = document.querySelector('form a').hostname;
}

function setupBuy(item) {
	return function() {
		doBuy(item.parentNode.querySelector('.purchase'))
	}
}

window.onload = function() {

	var items = document.querySelectorAll('.item');

	for (var i = 0; i < items.length; i++) {
		var item = items[i],
			url = item.href;

		if (item.getAttribute('data-bought') === 'true') {
			var span = document.createElement('span');
			span.textContent = '(fulfilled)';
			item.parentNode.insertBefore(span, item.nextSibling);
			continue;
		}

		var links = document.createElement('div');
		links.className = 'links';
		links.innerHTML =  '<a class="purchase" data-text="'+item.textContent+'" data-href="'+url+'" target="_blank">Gift this...</a>';
		item.parentNode.insertBefore(links, item.nextSibling);

		item.parentNode.addEventListener('click', setupBuy(item));

		item.removeAttribute('href');

		item.addEventListener('click', function() {
			this.parentNode.classList.toggle('open');
		});

	};

	document.querySelector('.cancel-btn').addEventListener('click', function() {
		document.querySelector('form').classList.remove('open');
	});
	document.querySelector('.close-btn').addEventListener('click', function() {
		document.querySelector('.message').classList.remove('shown');
	});

	document.querySelector('form').addEventListener('submit', function(e) {
		if (!document.querySelector('input[name="name"]').value.length || !document.querySelector('textarea').value.length) {
			e.preventDefault();
			alert("Please fill these fields out so we know who to thank :)")
		}
	});

}
