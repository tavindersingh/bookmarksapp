$('table tbody tr.item').hover(function () {
    $(this).find('td.address>p>span').toggleClass('visible');
});