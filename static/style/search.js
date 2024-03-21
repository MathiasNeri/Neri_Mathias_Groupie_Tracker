
function searchAnime() {
    console.log('toto');
    const searchBox = document.querySelector('.searchInput');
    const query = searchBox.value;
    if (query.trim()) {
        window.location.href = `/result_search?q=${encodeURIComponent(query)}`;
    } else {
        alert('Veuillez entrer un terme de recherche.');
    }
}
